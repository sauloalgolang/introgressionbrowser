package vcf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

import (
	"github.com/remeh/sizedwaitgroup"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
)

const DEBUG = false
const BREAKAT = 100

// https://github.com/brentp/vcfgo/blob/master/examples/main.go

//
//
// Chromosome Gather
//
//

func GatherChromosomeNames(sourceFile string, isTar bool, isGz bool, continueOnError bool) (chromosomeNames interfaces.ChromosomeNamesType) {
	indexfile := chromosomeNames.IndexFileName(sourceFile)

	if _, err := os.Stat(indexfile); err == nil {
		// path/to/whatever exists
		fmt.Println(" exists")
		chromosomeNames.Load(sourceFile)
		// fmt.Println(chromosomeNames)

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println(" creating")

		addToNames := func(SampleNames *interfaces.VCFSamples, register *interfaces.VCFRegister) {
			fmt.Println("adding chromosome ", register.Chromosome, chromosomeNames)
			chromosomeNames.Add(register.Chromosome, register.LineNumber)
		}

		getNames := func(r io.Reader, continueOnError bool) {
			ProcessVcfRaw(r, addToNames, continueOnError, []string{""})
		}

		openfile.OpenFile(sourceFile, isTar, isGz, continueOnError, getNames)

		chromosomeNames.Save(sourceFile)
	}

	return chromosomeNames
}

//
//
// Chrosmosome Callback
//
//

type ChromosomeCallbackRegister struct {
	callBack        interfaces.VCFMaskedReaderChromosomeType
	chromosomeNames []string
	wg              *sizedwaitgroup.SizedWaitGroup
	// wg             *sync.WaitGroup
}

func (cc *ChromosomeCallbackRegister) ChromosomeCallback(r io.Reader, continueOnError bool) {
	defer cc.wg.Done()

	cc.callBack(r, continueOnError, cc.chromosomeNames)

	fmt.Println("Finished reading chromosomes   :", cc.chromosomeNames)
}

//
//
// Open VCF
//
//

func OpenVcfFile(sourceFile string, continueOnError bool, numThreads int, callBack interfaces.VCFMaskedReaderChromosomeType) {
	fmt.Println("OpenVcfFile :: ",
		"sourceFile", sourceFile,
		"continueOnError", continueOnError,
		"numThreads", numThreads)

	isTar := false
	isGz := false

	if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf.tar.gz") {
		fmt.Println(" .tar.gz format")
		isTar = true
		isGz = true
	} else if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf.gz") {
		fmt.Println(" .gz format")
		isTar = false
		isGz = true
	} else if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf") {
		fmt.Println(" .vcf format")
		isTar = false
		isGz = false
	} else {
		fmt.Println("unknown file suffix!")
		os.Exit(1)
	}

	chromosomeNames := GatherChromosomeNames(sourceFile, isTar, isGz, continueOnError)
	for _, chromosomeInfo := range chromosomeNames.Infos {
		fmt.Println("Gathered Chromosome Names :: ", chromosomeInfo.ChromosomeName, " ", chromosomeInfo.NumRegisters)
	}

	chromosomeGroups := make([][]string, numThreads, numThreads)

	// wg := sync.WaitGroup
	wg := sizedwaitgroup.New(numThreads)
	for _, chromosomeGroup := range chromosomeGroups {
		ccr := ChromosomeCallbackRegister{
			callBack:        callBack,
			chromosomeNames: chromosomeGroup,
			wg:              &wg,
		}

		// wg.Add(1)
		wg.Add()

		go openfile.OpenFile(sourceFile, isTar, isGz, continueOnError, ccr.ChromosomeCallback)
	}

	fmt.Println("Waiting for all chromosomes to complete")
	wg.Wait()
	fmt.Println("All chromosomes completed")
}

func ProcessVcf(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool, chromosomeNames []string) {
	bufreader := bufio.NewReader(r)
	ProcessVcfRaw(bufreader, callback, continueOnError, chromosomeNames)
	// ProcessVcfVcfGo(bufreader, callback, continueOnError, chromosomeNames)
}
