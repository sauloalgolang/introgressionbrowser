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
	samples    []string
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
	chromosomes      map[string]IBChromosome
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
		samples:    make([]string, 0, 100),
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
		chromosomes:      make(map[string]IBChromosome, 100),
		chromosomesNames: make([]string, 0, 100),
	}

	return ib
}

func (ib *IBrowser) SetSamples(samples *[]string) {
	ib.samples = *samples
	ib.numSamples = uint64(len(ib.samples))
}

func (ib *IBrowser) GetSamples() []string {
	return ib.samples
}

func (ib *IBrowser) HasChromosome(chromosome string) bool {
	if _, ok := ib.chromosomes[chromosome]; ok {
		return true
	} else {
		return false
	}
}

func (ib *IBrowser) AddChromosome(chromosome string) IBChromosome {
	if ib.HasChromosome(chromosome) {
		fmt.Println("Failed to add chromosome", chromosome, ". Already exists")
		os.Exit(1)
	}

	ib.chromosomes[chromosome] = NewIBChromosome(chromosome, ib.numSamples, ib.keepEmptyBlock)
	ib.chromosomesNames = append(ib.chromosomesNames, chromosome)

	return ib.chromosomes[chromosome]
}

func (ib *IBrowser) GetChromosome(chromosome string) (IBChromosome, bool) {
	if chromosome, ok := ib.chromosomes[chromosome]; ok {
		return chromosome, true
	} else {
		return IBChromosome{}, false
	}
}

func (ib *IBrowser) GetOrCreateChromosome(chromosomeName string) IBChromosome {
	if chromosome, ok := ib.GetChromosome(chromosomeName); ok {
		return chromosome
	} else {
		return ib.AddChromosome(chromosomeName)
	}
}

func (ib *IBrowser) ReaderCallBack(r io.Reader, continueOnError bool) {
	ib.reader(r, ib.RegisterCallBack, continueOnError)
}

func (ib *IBrowser) RegisterCallBack(reg *interfaces.VCFRegister) {
	fmt.Println("got register", reg)

	if ib.numSamples == 0 {
		ib.SetSamples(reg.Samples)
	} else {
		if len(ib.samples) != len(*reg.Samples) {
			fmt.Println("Sample mismatch")
			fmt.Println(len(ib.samples), "!=", len(*reg.Samples))
			os.Exit(1)
		}
	}

	if reg.Chromosome != ib.lastChrom {
		ib.lastChrom = reg.Chromosome
		ib.lastPosition = 0
	} else {
		if !(reg.Position > ib.lastPosition) {
			fmt.Println("Coordinate mismatch")
			fmt.Println(ib.lastPosition, ">=", reg.Position)
			os.Exit(1)
		}
	}

	ib.numRegisters++

	//TODO: FILTERING
	//
	// type VCFRegister struct {
	//  Samples      *[]string
	// 	IsHomozygous bool
	// 	IsIndel      bool
	// 	IsMNP        bool
	//  Row          uint64
	// 	Chromosome   string
	// 	Position     uint64
	// 	Quality      float32
	// 	Info         map[string]interface{}
	// 	Filter       string
	// 	NumAlt       uint64
	// 	Phased       bool
	// 	GT           [][]int
	// 	Fields       map[string]string
	// }

	chromosome := ib.GetOrCreateChromosome(reg.Chromosome)

	blockNum := reg.Position / ib.blockSize

	if !chromosome.HasBlock(blockNum) {
		ib.numBlocks++
	}

	ib.numSNPs++

	chromosome.Add(blockNum, reg)
}
