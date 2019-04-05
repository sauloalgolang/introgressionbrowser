package ibrowser

import (
	"fmt"
	"io"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

//
//
// IBROWSER SECTION
//
//

type IBrowser struct {
	reader interfaces.VCFReaderType
	//
	samples    interfaces.VCFSamples
	numSamples uint64
	//
	blockSize      uint64
	keepEmptyBlock bool
	//
	numRegisters uint64
	numSNPs      uint64
	numBlocks    uint64
	//
	lastChrom    string
	lastPosition uint64
	//
	chromosomes      map[string]*IBChromosome
	chromosomesNames []string
	//
	// Parameters string
	// Header string
	//
	// TODO: per sample stats
}

func NewIBrowser(reader interfaces.VCFReaderType, blockSize uint64, keepEmptyBlock bool) IBrowser {
	ib := IBrowser{
		reader: reader,
		//
		samples:    make(interfaces.VCFSamples, 0, 100),
		numSamples: 0,
		//
		blockSize:      blockSize,
		keepEmptyBlock: keepEmptyBlock,
		//
		numRegisters: 0,
		numSNPs:      0,
		numBlocks:    0,
		//
		lastChrom:    "",
		lastPosition: 0,
		//
		chromosomes:      make(map[string]*IBChromosome, 100),
		chromosomesNames: make([]string, 0, 100),
	}

	return ib
}

func (ib *IBrowser) SetSamples(samples *interfaces.VCFSamples) {
	numSamples := len(*samples)
	ib.samples = make(interfaces.VCFSamples, numSamples, numSamples)
	ib.numSamples = uint64(numSamples)

	for samplePos, sampleName := range *samples {
		// fmt.Println(samplePos, sampleName)
		ib.samples[samplePos] = sampleName
	}
}

func (ib *IBrowser) GetSamples() interfaces.VCFSamples {
	return ib.samples
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

	nchr := NewIBChromosome(chromosomeName, ib.numSamples, ib.keepEmptyBlock)
	ib.chromosomes[chromosomeName] = &nchr
	ib.chromosomesNames = append(ib.chromosomesNames, chromosomeName)

	return &nchr
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

func (ib *IBrowser) ReaderCallBack(r io.Reader, continueOnError bool) {
	ib.reader(r, ib.RegisterCallBack, continueOnError)
}

func (ib *IBrowser) RegisterCallBack(samples *interfaces.VCFSamples, reg *interfaces.VCFRegister) {
	if ib.numSamples == 0 {
		ib.SetSamples(samples)
	} else {
		if len(ib.samples) != len(*samples) {
			fmt.Println("Sample mismatch")
			fmt.Println(len(ib.samples), "!=", len(*samples))
			os.Exit(1)
		}
	}

	if reg.Chromosome != ib.lastChrom {
		ib.lastChrom = reg.Chromosome
		ib.lastPosition = 0
	} else {
		if !(reg.Pos > ib.lastPosition) {
			fmt.Println("Coordinate mismatch")
			fmt.Println(ib.lastPosition, ">=", reg.Pos)
			os.Exit(1)
		}
	}

	ib.numRegisters++

	//TODO: FILTERING
	//
	// type Variant struct {
	// 	Chromosome      string
	// 	Pos        		uint64
	// 	Id         		string
	// 	Ref        		string
	// 	Alt        		[]string
	// 	Quality    		float32
	// 	Filter     		string
	// 	Info       		InfoMap
	// 	Format     		[]string
	// 	Samples    		[]*SampleGenotype
	// 	Header     		*Header
	// 	LineNumber 		int64
	// }

	chromosome := ib.GetOrCreateChromosome(reg.Chromosome)

	blockNum := reg.Pos / ib.blockSize

	if _, hasBlock := chromosome.GetBlock(blockNum); !hasBlock {
		ib.numBlocks++
	}

	ib.numSNPs++

	chromosome.Add(blockNum, reg)
}
