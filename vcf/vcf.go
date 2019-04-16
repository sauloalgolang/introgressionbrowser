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

var DEBUG bool = false
var ONLYFIRST bool = false
var BREAKAT_THREAD int64 = 0
var BREAKAT_CHROM int64 = 0

//
//
// Chrosmosome Callback
//
//

type ChromosomeCallbackRegister struct {
	registerCallBack interfaces.VCFCallBack
	chromosomeNames  []string
	wg               *sizedwaitgroup.SizedWaitGroup
	// wg             *sync.WaitGroup
}

func (cc *ChromosomeCallbackRegister) ChromosomeCallback(r io.Reader, callBackParameters interfaces.CallBackParameters) {
	defer cc.wg.Done()

	cc.ChromosomeCallbackSingleThreaded(r, callBackParameters)
}

func (cc *ChromosomeCallbackRegister) ChromosomeCallbackSingleThreaded(r io.Reader, callBackParameters interfaces.CallBackParameters) {
	bufreader := bufio.NewReader(r)

	ProcessVcfRaw(bufreader, callBackParameters, cc.registerCallBack, cc.chromosomeNames)

	fmt.Println("Finished reading chromosomes   :", cc.chromosomeNames)
}

//
//
// File checker
//
//
type VcfFormat struct {
	isTar bool
	isGz  bool
}

func CheckVcfFormat(sourceFile string) VcfFormat {
	vf := VcfFormat{
		isTar: false,
		isGz:  false,
	}

	sourceFileLower := strings.ToLower(sourceFile)

	if strings.HasSuffix(sourceFileLower, ".vcf.tar.gz") {
		fmt.Println(" .tar.gz format")
		vf.isTar = true
		vf.isGz = true
	} else if strings.HasSuffix(sourceFileLower, ".vcf.gz") {
		fmt.Println(" .gz format")
		vf.isTar = false
		vf.isGz = true
	} else if strings.HasSuffix(sourceFileLower, ".vcf") {
		fmt.Println(" .vcf format")
		vf.isTar = false
		vf.isGz = false
	} else {
		fmt.Println("unknown file suffix!")
		os.Exit(1)
	}

	return vf
}

//
//
// Open VCF
//
//

// func OpenVcfFile(sourceFile string, continueOnError bool, numThreads int, registerCallBack interfaces.VCFMaskedReaderChromosomeType) {
func OpenVcfFile(sourceFile string, callBackParameters interfaces.CallBackParameters, registerCallBack interfaces.VCFCallBack) {
	fmt.Println("OpenVcfFile :: ",
		"sourceFile", sourceFile,
		"numBits", callBackParameters.NumBits,
		"continueOnError", callBackParameters.ContinueOnError,
		"numThreads", callBackParameters.NumThreads)

	vcfFormat := CheckVcfFormat(sourceFile)

	chromosomeNames := GatherChromosomeNames(sourceFile, vcfFormat.isTar, vcfFormat.isGz, callBackParameters)

	p := message.NewPrinter(language.English)
	p.Print("Gathered Chromosome Names:\n")
	p.Printf(" NumChromosomes : %12d\n", chromosomeNames.NumChromosomes)
	p.Printf(" StartPosition  : %12d\n", chromosomeNames.StartPosition)
	p.Printf(" EndPosition    : %12d\n", chromosomeNames.EndPosition)
	p.Printf(" NumRegisters   : %12d\n", chromosomeNames.NumRegisters)

	if callBackParameters.NumThreads == 1 {
		fmt.Println("Running single threaded")

		chromosomeGroup := make([]string, chromosomeNames.NumChromosomes, chromosomeNames.NumChromosomes)

		for _, chromosomeInfo := range chromosomeNames.Infos {
			chromosomeGroup = append(chromosomeGroup, chromosomeInfo.ChromosomeName)
		}

		ccr := ChromosomeCallbackRegister{
			registerCallBack: registerCallBack,
			chromosomeNames:  chromosomeGroup,
		}

		openfile.OpenFile(sourceFile, vcfFormat.isTar, vcfFormat.isGz, callBackParameters, ccr.ChromosomeCallbackSingleThreaded)

		fmt.Println("Finished reading file")

	} else {
		chromosomeGroups := SpreadChromosomes(chromosomeNames, callBackParameters.NumThreads)

		// wg := sync.WaitGroup
		wg := sizedwaitgroup.New(callBackParameters.NumThreads)
		for _, chromosomeGroup := range chromosomeGroups {
			ccr := ChromosomeCallbackRegister{
				registerCallBack: registerCallBack,
				chromosomeNames:  chromosomeGroup,
				wg:               &wg,
			}

			// wg.Add(1)
			wg.Add()

			go openfile.OpenFile(sourceFile, vcfFormat.isTar, vcfFormat.isGz, callBackParameters, ccr.ChromosomeCallback)

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
