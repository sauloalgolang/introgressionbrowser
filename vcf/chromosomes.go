package vcf

import (
	log "github.com/sirupsen/logrus"
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

// GatherChromosomeNames Gets chromosome name from a file and saves in a index file. if index file exists, reads it directly
func GatherChromosomeNames(sourceFile string, isTar bool, isGz bool, callBackParameters interfaces.CallBackParameters) (chromosomeNames interfaces.ChromosomeNamesType) {
	exists, _ := chromosomeNames.Exists(sourceFile)

	if exists {
		log.Println(" exists")
		chromosomeNames.Load(sourceFile)

	} else {
		log.Println(" creating")

		addToNames := func(SampleNames *Samples, register *Register) {
			log.Println("adding chromosome ", register.Chromosome)
			chromosomeNames.Add(register.Chromosome, register.LineNumber)
		}

		getNames := func(r io.Reader, callBackParameters interfaces.CallBackParameters) {
			ProcessVcfRaw(r,
				callBackParameters,
				addToNames,
				[]string{""})
		}

		openfile.OpenFile(sourceFile, isTar, isGz, callBackParameters, getNames)

		chromosomeNames.Save(sourceFile)
	}

	return chromosomeNames
}

// SpreadChromosomes spreads the available chromosomes among the threads
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

	if len(chromosomeGroups[len(chromosomeGroups)-1]) == 0 {
		for idx, cNames := range chromosomeGroups {
			if idx != (len(chromosomeGroups) - 1) {
				lastChrom := cNames[len(cNames)-1]
				chromosomeGroups[idx+1] = append([]string{lastChrom}, chromosomeGroups[idx+1]...) //prepend
				chromosomeGroups[idx] = cNames[:len(cNames)-1]
			}
		}
	}

	return chromosomeGroups
}
