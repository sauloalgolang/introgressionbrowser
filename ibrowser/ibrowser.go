package ibrowser

import (
	// "encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"
import "github.com/sauloalgolang/introgressionbrowser/tools"

//
//
// IBROWSER SECTION
//
//

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
	Chromosomes      map[string]*IBChromosome
	ChromosomesNames []string
	//
	Block *IBBlock
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
		Chromosomes:      make(map[string]*IBChromosome, 100),
		ChromosomesNames: make([]string, 0, 100),
		//
		Block: NewIBBlock(0, 0),
	}

	return &ib
}

func (ib *IBrowser) SetSamples(samples *interfaces.VCFSamples) {
	numSamples := len(*samples)
	ib.Samples = make(interfaces.VCFSamples, numSamples, numSamples)
	ib.NumSamples = uint64(numSamples)
	ib.Block = NewIBBlock(0, ib.NumSamples)

	for samplePos, sampleName := range *samples {
		// fmt.Println(samplePos, sampleName)
		ib.Samples[samplePos] = sampleName
	}
}

func (ib *IBrowser) GetSamples() interfaces.VCFSamples {
	return ib.Samples
}

func (ib *IBrowser) GetChromosome(chromosomeName string) (*IBChromosome, bool) {
	if chromosome, ok := ib.Chromosomes[chromosomeName]; ok {
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

	ib.Chromosomes[chromosomeName] = NewIBChromosome(chromosomeName, ib.NumSamples, ib.KeepEmptyBlock)
	ib.ChromosomesNames = append(ib.ChromosomesNames, chromosomeName)

	return ib.Chromosomes[chromosomeName]
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
	if ib.NumSamples == 0 {
		ib.SetSamples(samples)
	} else {
		if len(ib.Samples) != len(*samples) {
			fmt.Println("Sample mismatch")
			fmt.Println(len(ib.Samples), "!=", len(*samples))
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

	ib.NumRegisters++

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
	position := reg.Pos
	blockNum := position / ib.BlockSize
	distance := tools.CalculateDistance(ib.NumSamples, reg)

	if _, hasBlock := chromosome.GetBlock(blockNum); !hasBlock {
		ib.NumBlocks++
	}

	ib.NumSNPs++

	chromosome.Add(blockNum, position, distance)

	ib.Block.Add(position, distance)

	// ib.Save()

	// os.Exit(0)
}

func (ib *IBrowser) Save() {
	// ibB, _ := json.Marshal(ib)
	// fmt.Println(string(ibB))

	fmt.Println(ib.Block)

	d, err := yaml.Marshal(ib)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("--- dump:\n%s\n\n", d)
}
