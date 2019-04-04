package ibrowser

import (
	"fmt"
	"io"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

type IBBlock struct {
	MinPosition uint64
	MaxPosition uint64
	BlockNumber uint64
	NumSNPS     uint64
	NumSamples  uint64
	Matrix      [][]uint64
}

func NewIBBlock(numSamples uint64) IBBlock {
	// TODO: MAYBE UNSAFE
	ibb := *new(IBBlock)
	ibb.NumSamples = numSamples
	ibb.MinPosition = 0
	ibb.MaxPosition = 0
	ibb.BlockNumber = 0
	ibb.NumSNPS = 0
	ibb.Matrix = make([][]uint64, 0, 0)
	return ibb
}

func (ib *IBBlock) Add(reg *interfaces.VCFRegister) {
}

type IBChromosome struct {
	Chromosome     string
	MinPosition    uint64
	MaxPosition    uint64
	NumBlocks      uint64
	NumSNPS        uint64
	NumSamples     uint64
	KeepEmptyBlock bool
	Blocks         []IBBlock
	BlockNames     map[uint64]uint64
}

func NewIBChromosome(chromosome string, numSamples uint64, keepEmptyBlock bool) IBChromosome {
	// TODO: MAYBE UNSAFE
	ibc := *new(IBChromosome)
	ibc.Chromosome = chromosome
	ibc.NumSamples = numSamples
	ibc.MinPosition = 0
	ibc.MaxPosition = 0
	ibc.NumBlocks = 0
	ibc.NumSNPS = 0
	ibc.KeepEmptyBlock = keepEmptyBlock
	ibc.Blocks = make([]IBBlock, 0, 0)
	return ibc
}

func (ibc *IBChromosome) Add(blockNum uint64, reg *interfaces.VCFRegister) {
	if blockPos, ok := ibc.BlockNames[blockNum]; ok {
		if !ok {
			if ibc.KeepEmptyBlock {

			} else {
				ibc.BlockNames[blockNum] = uint64(len(ibc.Blocks)) + 1
				ibc.Blocks = append(ibc.Blocks, NewIBBlock(ibc.NumSamples))
			}
		}

		// TODO: MAKE SURE TO GET REFERENCE
		blockPos = ibc.BlockNames[blockNum]
		block := ibc.Blocks[blockPos]

		block.Add(reg)
	}
}

type IBrowser struct {
	reader     interfaces.VCFReaderType
	samples    []string
	numSamples uint64
	//
	blockSize      uint64
	keepEmptyBlock bool
	//
	numSNPs   uint64
	numBlocks uint64
	//
	lastChrom    string
	lastPosition uint64
	//
	// Parameters string
	// Header string
	//
	chromosomes      map[string]IBChromosome
	chromosomesNames []string
}

func NewIBrowser(reader interfaces.VCFReaderType, blockSize uint64, keepEmptyBlock bool) *IBrowser {
	ib := new(IBrowser)
	ib.reader = reader
	ib.samples = make([]string, 0, 0)
	ib.numSamples = 0
	ib.blockSize = blockSize
	ib.keepEmptyBlock = keepEmptyBlock
	return ib
}

func (ib *IBrowser) SetSamples(samples *[]string) {
	ib.samples = *samples
	ib.numSamples = uint64(len(ib.samples))
}

func (ib *IBrowser) GetSamples() []string {
	return ib.samples
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
		if ib.lastPosition >= reg.Position {
			fmt.Println("Coordinate mismatch")
			fmt.Println(ib.lastPosition, ">=", reg.Position)
			os.Exit(1)
		}
	}

	//TODO: FILTERING

	if val, ok := ib.chromosomes[reg.Chromosome]; ok {
		if !ok {
			ib.chromosomes[reg.Chromosome] = NewIBChromosome(reg.Chromosome, ib.numSamples, ib.keepEmptyBlock)
			ib.chromosomesNames = append(ib.chromosomesNames, reg.Chromosome)
		}

		blockNum := reg.Position / ib.blockSize

		val.Add(blockNum, reg)
	}
}
