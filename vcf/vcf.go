package vcf

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
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

type chromosomeNamesType struct {
	Names []string
}

func (cn *chromosomeNamesType) IndexFileName(outPrefix string) (indexFile string) {
	indexFile = outPrefix + ".ibindex"
	return indexFile
}

func (cn *chromosomeNamesType) Save(outPrefix string) {
	d, err := yaml.Marshal(cn)

	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	// fmt.Printf("--- dump:\n%s\n\n", d)
	outfile := cn.IndexFileName(outPrefix)
	fmt.Println(" saving index to ", outfile)
	err = ioutil.WriteFile(outfile, d, 0644)
	fmt.Println(" done")
}

func (cn *chromosomeNamesType) Load(outPrefix string) {
	outfile := cn.IndexFileName(outPrefix)

	data, err := ioutil.ReadFile(outfile)

	if err != nil {
		fmt.Printf("yamlFile. Get err   #%v ", err)
	}

	err = yaml.Unmarshal(data, &cn)

	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
	}
}

func (cn *chromosomeNamesType) Add(chromosomeName string) {
	cn.Names = append(cn.Names, chromosomeName)
}

func GatherChromosomeNames(sourceFile string, isTar bool, isGz bool, continueOnError bool) (chromosomeNames chromosomeNamesType) {
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
			chromosomeNames.Add(register.Chromosome)
		}

		getNames := func(r io.Reader, continueOnError bool) {
			ProcessVcfRaw(r, addToNames, continueOnError, "")
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
	callBack       interfaces.VCFMaskedReaderChromosomeType
	chromosomeName string
	wg             *sizedwaitgroup.SizedWaitGroup
	// wg             *sync.WaitGroup
}

func (cc *ChromosomeCallbackRegister) ChromosomeCallback(r io.Reader, continueOnError bool) {
	defer cc.wg.Done()
	cc.callBack(r, continueOnError, cc.chromosomeName)

	fmt.Println("Finished reading chromosome    :", cc.chromosomeName)
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
	for _, chromosomeName := range chromosomeNames.Names {
		fmt.Println("Gathered Chromosome Names :: ", chromosomeName)
	}

	// wg := sync.WaitGroup
	wg := sizedwaitgroup.New(numThreads)
	for _, chromosomeName := range chromosomeNames.Names {
		ccr := ChromosomeCallbackRegister{
			callBack:       callBack,
			chromosomeName: chromosomeName,
			wg:             &wg,
		}

		// wg.Add(1)
		wg.Add()

		go openfile.OpenFile(sourceFile, isTar, isGz, continueOnError, ccr.ChromosomeCallback)
	}

	fmt.Println("Waiting for all chromosomes to complete")
	wg.Wait()
	fmt.Println("All chromosomes completed")
}
