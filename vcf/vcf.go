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

//
//
// Chrosmosome Callback
//
//

// ChromosomeCallbackRegister struct holding callback to be used for a list of chromosomes
type ChromosomeCallbackRegister struct {
	registerCallBack RegisterCallBack
	chromosomeNames  []string
	wg               *SizedWaitGroup
	// wg             *sync.WaitGroup
}

// ChromosomeCallback callback to be used to read chromosomes
func (cc *ChromosomeCallbackRegister) ChromosomeCallback(r io.Reader, callBackParameters CallBackParameters) {
	defer cc.wg.Done()

	cc.ChromosomeCallbackSingleThreaded(r, callBackParameters)
}

// ChromosomeCallbackSingleThreaded callback to be used to read chromosomes
func (cc *ChromosomeCallbackRegister) ChromosomeCallbackSingleThreaded(r io.Reader, callBackParameters CallBackParameters) {
	bufreader := bufio.NewReader(r)

	ProcessVcfRaw(bufreader, callBackParameters, cc.registerCallBack, cc.chromosomeNames)

	fmt.Println("Finished reading chromosomes   :", cc.chromosomeNames)
}

//
//
// File checker
//
//

// Format description of the vcf formats
type Format struct {
	isTar bool
	isGz  bool
}

// CheckFormat identifies the file format
func CheckFormat(sourceFile string) Format {
	vf := Format{
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

// OpenFile opens file, call callback when a register is found
func OpenFile(sourceFile string, callBackParameters CallBackParameters, registerCallBack RegisterCallBack) {
	fmt.Println("OpenFile :: ",
		"sourceFile", sourceFile,
		"numBits", callBackParameters.NumBits,
		"continueOnError", callBackParameters.ContinueOnError,
		"numThreads", callBackParameters.NumThreads)

	vcfFormat := CheckFormat(sourceFile)

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

		OpenAnyFile(sourceFile, vcfFormat.isTar, vcfFormat.isGz, callBackParameters, ccr.ChromosomeCallbackSingleThreaded)

		fmt.Println("Finished reading file")

	} else {
		threads := callBackParameters.NumThreads

		if threads > int(chromosomeNames.NumChromosomes) {
			threads = int(chromosomeNames.NumChromosomes)
		}

		chromosomeGroups := SpreadChromosomes(chromosomeNames, threads)

		// wg := sync.WaitGroup
		wg := sizedwaitgroup.New(threads)
		for _, chromosomeGroup := range chromosomeGroups {
			ccr := ChromosomeCallbackRegister{
				registerCallBack: registerCallBack,
				chromosomeNames:  chromosomeGroup,
				wg:               &wg,
			}

			// wg.Add(1)
			wg.Add()

			go OpenAnyFile(
				sourceFile,
				vcfFormat.isTar,
				vcfFormat.isGz,
				callBackParameters,
				ccr.ChromosomeCallback,
			)

			if OnlyFirst {
				fmt.Println("Only sending first")
				break
			}
		}

		fmt.Println("Waiting for all chromosomes to complete")
		wg.Wait()
		fmt.Println("All chromosomes completed")
	}
}
