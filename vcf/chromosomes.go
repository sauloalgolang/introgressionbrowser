package vcf

import (
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
)

//
//
// Chromosome Gather
//
//

func GatherChromosomeNames(sourceFile string, isTar bool, isGz bool, continueOnError bool) (chromosomeNames interfaces.ChromosomeNamesType) {
	exists, _ := chromosomeNames.Exists(sourceFile)

	if exists {
		fmt.Println(" exists")
		chromosomeNames.Load(sourceFile)

	} else {
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

func SpreadChromosomes(chromosomeNames interfaces.ChromosomeNamesType, numThreads int) [][]string {
	chromosomeGroups := make([][]string, numThreads, numThreads)
	chromosomeGroupsSizes := make([]int64, numThreads, numThreads)
	fraq := chromosomeNames.NumRegisters / int64(numThreads) / int64(numThreads)

	p := message.NewPrinter(language.English)

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

	return chromosomeGroups
}
