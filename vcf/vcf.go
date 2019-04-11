package vcf

import (
	"bufio"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

const BREAKAT = 0
const ONLYFIRST = false

// const BREAKAT = int64(1000000)
// const ONLYFIRST = true

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
			fmt.Println("adding chromosome ", register.Chromosome)
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

func (cc *ChromosomeCallbackRegister) ChromosomeCallbackSingleThreaded(r io.Reader, continueOnError bool) {
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

	p := message.NewPrinter(language.English)
	p.Print("Gathered Chromosome Names:\n")
	p.Printf(" NumChromosomes : %12d\n", chromosomeNames.NumChromosomes)
	p.Printf(" StartPosition  : %12d\n", chromosomeNames.StartPosition)
	p.Printf(" EndPosition    : %12d\n", chromosomeNames.EndPosition)
	p.Printf(" NumRegisters   : %12d\n", chromosomeNames.NumRegisters)

	if numThreads == 1 {
		fmt.Println("Running single threaded")

		chromosomeGroup := make([]string, chromosomeNames.NumChromosomes, chromosomeNames.NumChromosomes)

		for _, chromosomeInfo := range chromosomeNames.Infos {
			chromosomeGroup = append(chromosomeGroup, chromosomeInfo.ChromosomeName)
		}

		ccr := ChromosomeCallbackRegister{
			callBack:        callBack,
			chromosomeNames: chromosomeGroup,
		}

		openfile.OpenFile(sourceFile, isTar, isGz, continueOnError, ccr.ChromosomeCallbackSingleThreaded)

		fmt.Println("Finished reading file")

	} else {
		chromosomeGroups := make([][]string, numThreads, numThreads)
		chromosomeGroupsSizes := make([]int64, numThreads, numThreads)
		fraq := chromosomeNames.NumRegisters / int64(numThreads) / int64(numThreads)

		p.Printf(" Fraction       : %12d\n", fraq)
		p.Println()

		cummChromSize := int64(0)

		for _, chromosomeInfo := range chromosomeNames.Infos {
			idx := cummChromSize / fraq / int64(numThreads+(numThreads/3))

			if idx >= int64(numThreads) {
				p.Printf("%12d\n", idx)
				idx = int64(numThreads) - 1
			}

			cummChromSize += chromosomeInfo.NumRegisters

			p.Printf(" Chromosome Name: %s\n", chromosomeInfo.ChromosomeName)
			p.Printf("  Start Position: %12d\n", chromosomeInfo.StartPosition)
			p.Printf("  Registers     : %12d\n", chromosomeInfo.NumRegisters)
			p.Printf("  Cumm Registers: %12d\n", cummChromSize)
			p.Printf("  Thread        : %12d\n", idx)

			if len(chromosomeGroups[idx]) == 0 {
				chromosomeGroups[idx] = make([]string, 0, 0)
			}

			chromosomeGroupsSizes[idx] += chromosomeInfo.NumRegisters
			chromosomeGroups[idx] = append(chromosomeGroups[idx], chromosomeInfo.ChromosomeName)

			p.Printf("  Group Size    : %12d\n", chromosomeGroupsSizes[idx])
			p.Printf("  Group         : %v\n", chromosomeGroups[idx])
			p.Println()
		}

		p.Println()

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

			if ONLYFIRST {
				fmt.Println("Only sending first")
				break
			}
		}

		fmt.Println("Waiting for all chromosomes to complete")
		wg.Wait()
		fmt.Println("All chromosomes completed")
	}
}

func ProcessVcf(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool, chromosomeNames []string) {
	bufreader := bufio.NewReader(r)
	ProcessVcfRaw(bufreader, callback, continueOnError, chromosomeNames)
	// ProcessVcfVcfGo(bufreader, callback, continueOnError, chromosomeNames)
}
