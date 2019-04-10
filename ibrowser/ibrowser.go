package ibrowser

import (
	// "encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/save"
)

//
//
// IBROWSER SECTION
//
//

var mutex = &sync.Mutex{}

type IBrowser struct {
	reader interfaces.VCFReaderType
	//
	Samples    interfaces.VCFSamples
	NumSamples uint64
	//
	BlockSize      uint64
	KeepEmptyBlock bool
	//
	NumRegisters uint64
	NumSNPs      uint64
	NumBlocks    uint64
	//
	lastChrom    string
	lastPosition uint64
	//
	ChromosomesNames []string
	chromosomes      map[string]*IBChromosome
	//
	block *IBBlock
	//
	// Parameters string
	// Header string
	//
	// TODO: per sample stats
}

func NewIBrowser(reader interfaces.VCFReaderType, blockSize uint64, keepEmptyBlock bool) *IBrowser {
	ib := IBrowser{
		reader: reader,
		//
		Samples:    make(interfaces.VCFSamples, 0, 100),
		NumSamples: 0,
		//
		BlockSize:      blockSize,
		KeepEmptyBlock: keepEmptyBlock,
		//
		NumRegisters: 0,
		NumSNPs:      0,
		NumBlocks:    0,
		//
		lastChrom:    "",
		lastPosition: 0,
		//
		chromosomes:      make(map[string]*IBChromosome, 100),
		ChromosomesNames: make([]string, 0, 100),
		//
		// Block: NewIBBlock(0, 0),
	}

	return &ib
}

func (ib *IBrowser) SetSamples(samples *interfaces.VCFSamples) {
	numSamples := len(*samples)
	ib.Samples = make(interfaces.VCFSamples, numSamples, numSamples)

	ib.NumSamples = uint64(numSamples)
	ib.block = NewIBBlock("_whole_genome", ib.BlockSize, 0, 0, ib.NumSamples)

	for samplePos, sampleName := range *samples {
		// fmt.Println(samplePos, sampleName)
		ib.Samples[samplePos] = sampleName
	}
}

func (ib *IBrowser) GetSamples() interfaces.VCFSamples {
	return ib.Samples
}

func (ib *IBrowser) GetChromosome(chromosomeName string) (*IBChromosome, bool) {
	if chromosome, ok := ib.chromosomes[chromosomeName]; ok {
		// fmt.Println("GetChromosome", chromosomeName, "exists", &chromosome)
		return chromosome, ok
	} else {
		// fmt.Println("GetChromosome", chromosomeName, "DOES NOT exists")
		return &IBChromosome{}, ok
	}
}

func (ib *IBrowser) AddChromosome(chromosomeName string) *IBChromosome {
	if chromosome, hasChromosome := ib.GetChromosome(chromosomeName); hasChromosome {
		fmt.Println("Failed to add chromosome", chromosomeName, ". Already exists", &chromosome)
		os.Exit(1)
	}

	ib.chromosomes[chromosomeName] = NewIBChromosome(chromosomeName, ib.BlockSize, ib.NumSamples, ib.KeepEmptyBlock)
	ib.ChromosomesNames = append(ib.ChromosomesNames, chromosomeName)

	return ib.chromosomes[chromosomeName]
}

func (ib *IBrowser) GetOrCreateChromosome(chromosomeName string) *IBChromosome {
	if chromosome, ok := ib.GetChromosome(chromosomeName); ok {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "exists", &chromosome)
		return chromosome
	} else {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "creating")
		return ib.AddChromosome(chromosomeName)
	}
}

func (ib *IBrowser) ReaderCallBack(r io.Reader, continueOnError bool, chromosomeNames []string) {
	ib.reader(r, ib.RegisterCallBack, continueOnError, chromosomeNames)
}

func (ib *IBrowser) RegisterCallBack(samples *interfaces.VCFSamples, reg *interfaces.VCFRegister) {
	if atomic.LoadUint64(&ib.NumSamples) == 0 {
		ib.SetSamples(samples)
	} else {
		if len(ib.Samples) != len(*samples) {
			fmt.Println("Sample mismatch")
			fmt.Println(len(ib.Samples), "!=", len(*samples))
			os.Exit(1)
		}
	}

	atomic.AddUint64(&ib.NumRegisters, 1)

	//
	// Adding distance
	//

	atomic.AddUint64(&ib.NumSNPs, 1)

	// ib.Block.AddAtomic(0, reg.Distance) // did not work

	// mutex.Lock() // did not work
	ib.block.Add(0, reg.Distance)
	// mutex.Unlock()

	chromosome := ib.GetOrCreateChromosome(reg.Chromosome)

	_, isNew := chromosome.Add(reg)

	if !isNew {
		atomic.AddUint64(&ib.NumBlocks, 1)
	}
}

func (ib *IBrowser) GenFilename(outPrefix string, format string) (baseName string, fileName string) {
	baseName = outPrefix
	fileName = save.GenFilename(baseName, format)
	return baseName, fileName
}

func (ib *IBrowser) Save(outPrefix string, format string) {
	baseName, _ := ib.GenFilename(outPrefix, format)
	save.Save(baseName, format, ib)
	ib.saveBlock(baseName, format)
	ib.saveChromosomes(baseName, format)
}

func (ib *IBrowser) saveBlock(outPrefix string, format string) {
	ib.block.Save(outPrefix+"_block", format)
}

func (ib *IBrowser) saveChromosomes(outPrefix string, format string) {
	for chromosomeName, chromosome := range ib.chromosomes {
		_, fileName := chromosome.GenFilename(outPrefix, format)

		fmt.Print("saving chromosome: ", chromosomeName, " to: ", fileName)

		if _, err := os.Stat(fileName); err == nil {
			// path/to/whatever exists
			fmt.Println(" exists")
			continue

		} else if os.IsNotExist(err) {
			fmt.Println(" creating")
			// path/to/whatever does *not* exist

		} else {
			// Schrodinger: file may or may not exist. See err for details.

			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		}

		chromosome.Save(outPrefix, format)
	}
}
