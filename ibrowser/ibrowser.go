package ibrowser

import (
	// "encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"sync/atomic"
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
	atomic.StoreUint64(&ib.NumSamples, uint64(numSamples))
	ib.Block = NewIBBlock(0, atomic.LoadUint64(&ib.NumSamples))

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

func (ib *IBrowser) ReaderCallBack(r io.Reader, continueOnError bool, chromosomeName string) {
	ib.reader(r, ib.RegisterCallBack, continueOnError, chromosomeName)
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

	// if reg.Chromosome != ib.lastChrom {
	// 	if ib.lastChrom != "" {
	// 		ib.GetOrCreateChromosome(ib.lastChrom).Save("output")
	// 	}

	// 	ib.lastChrom = reg.Chromosome
	// 	ib.lastPosition = 0
	// 	fmt.Println("New chromosome: ", reg.Chromosome)

	// } else {
	// 	if !(reg.Position > ib.lastPosition) {
	// 		fmt.Println("Coordinate mismatch")
	// 		fmt.Println(ib.lastPosition, ">=", reg.Position)
	// 		os.Exit(1)
	// 	}
	// }

	atomic.AddUint64(&ib.NumRegisters, 1)

	// TODO
	// FILTERING
	//

	//
	// Adding distance
	//

	atomic.AddUint64(&ib.NumSNPs, 1)

	ib.Block.AddAtomic(0, reg.Distance)

	position := reg.Position
	blockNum := position / ib.BlockSize
	chromosome := ib.GetOrCreateChromosome(reg.Chromosome)

	if _, hasBlock := chromosome.GetBlock(blockNum); !hasBlock {
		atomic.AddUint64(&ib.NumBlocks, 1)
	}

	chromosome.Add(blockNum, position, reg.Distance)
}

func (ib *IBrowser) SaveChromosomes(outPrefix string) {
	for chromosomeName, chromosome := range ib.Chromosomes {

		fmt.Print("saving chromosome: ", chromosomeName)

		outfile := outPrefix + "." + chromosomeName + ".yaml"

		if _, err := os.Stat(outfile); err == nil {
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
		chromosome.Save(outPrefix)
	}
}

func (ib *IBrowser) Save(outPrefix string) {
	// ibB, _ := json.Marshal(ib)
	// fmt.Println(string(ibB))

	// fmt.Println("ib.Block", ib.Block)
	// fmt.Println()

	// chromosome := ib.Chromosomes[ib.ChromosomesNames[0]]

	// fmt.Println("chromosome.Block", chromosome.Block)
	// fmt.Println()

	// block := chromosome.Blocks[0]

	// fmt.Println("chromosome.Blocks[0]", block)
	// fmt.Println()

	d, err := yaml.Marshal(ib)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	// fmt.Printf("--- dump:\n%s\n\n", d)
	outfile := outPrefix + ".yaml"
	fmt.Println("saving ibrowser to ", outfile)
	err = ioutil.WriteFile(outfile, d, 0644)
	fmt.Println("done")
}
