package ibrowser

import (
	"fmt"
	"math"
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
	// reader interfaces.VCFReaderType
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
	NumBits      int
	//
	ChromosomesNames []string
	chromosomes      map[string]*IBChromosome
	//
	block *IBBlock
	//
	lastChrom    string
	lastPosition uint64
	//
	//
	// Parameters string
	// Header string
	//
	// TODO: per sample stats
}

func NewIBrowser(blockSize uint64, numBits int, keepEmptyBlock bool) *IBrowser {
	if blockSize > uint64((math.MaxUint32/3)-1) {
		fmt.Println("block size too large")
		os.Exit(1)
	}

	ib := IBrowser{
		Samples:    make(interfaces.VCFSamples, 0, 100),
		NumSamples: 0,
		//
		BlockSize:      blockSize,
		KeepEmptyBlock: keepEmptyBlock,
		//
		NumRegisters: 0,
		NumSNPs:      0,
		NumBlocks:    0,
		NumBits:      numBits,
		//
		lastChrom:    "",
		lastPosition: 0,
		//
		ChromosomesNames: make([]string, 0, 100),
		chromosomes:      make(map[string]*IBChromosome, 100),
		//
		// block: NewIBBlock("_whole_genome", blockSize, numBits, 0, 0, 0),
	}

	return &ib
}

func (ib *IBrowser) SetSamples(samples *interfaces.VCFSamples) {
	numSamples := len(*samples)
	ib.Samples = make(interfaces.VCFSamples, numSamples, numSamples)

	ib.NumSamples = uint64(numSamples)
	ib.block = NewIBBlock("_whole_genome", ib.BlockSize, ib.NumBits, ib.NumSamples, 0, 0)

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

func (ib *IBrowser) GetOrCreateChromosome(chromosomeName string) *IBChromosome {
	if chromosome, ok := ib.GetChromosome(chromosomeName); ok {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "exists", &chromosome)
		return chromosome
	} else {
		// fmt.Println("GetOrCreateChromosome", chromosomeName, "creating")
		return ib.AddChromosome(chromosomeName)
	}
}

func (ib *IBrowser) AddChromosome(chromosomeName string) *IBChromosome {
	if chromosome, hasChromosome := ib.GetChromosome(chromosomeName); hasChromosome {
		fmt.Println("Failed to add chromosome", chromosomeName, ". Already exists", &chromosome)
		os.Exit(1)
	}

	ib.chromosomes[chromosomeName] = NewIBChromosome(chromosomeName, ib.BlockSize, ib.NumBits, ib.NumSamples, ib.KeepEmptyBlock)
	ib.ChromosomesNames = append(ib.ChromosomesNames, chromosomeName)

	return ib.chromosomes[chromosomeName]
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

	chromosome := ib.GetOrCreateChromosome(reg.Chromosome)

	_, isNew, numBlocksAdded := chromosome.Add(reg)

	mutex.Lock()
	{
		if isNew {
			ib.NumBlocks += numBlocksAdded
		}

		ib.NumRegisters++

		ib.NumSNPs++

		ib.block.Add(0, reg.Distance)
	}
	mutex.Unlock()
}

func (ib *IBrowser) GenFilename(outPrefix string, format string, compression string) (baseName string, fileName string) {
	baseName = outPrefix

	saver := save.NewSaverCompressed(baseName, format, compression)

	fileName = saver.GenFilename()

	return baseName, fileName
}

//
// Save
//
func (ib *IBrowser) Save(outPrefix string, format string, compression string) {
	ib.saveLoad(true, outPrefix, format, compression)
}

//
// Load
//

func (ib *IBrowser) Load(outPrefix string, format string, compression string) {
	ib.saveLoad(false, outPrefix, format, compression)
}

//
// SaveLoad
//

func (ib *IBrowser) saveLoad(isSave bool, outPrefix string, format string, compression string) {
	baseName, _ := ib.GenFilename(outPrefix, format, compression)
	saver := save.NewSaverCompressed(baseName, format, compression)

	if isSave {
		fmt.Println("saving global ibrowser status")
		saver.Save(ib)
	} else {
		fmt.Println("loading global ibrowser status")
		saver.Load(ib)
	}

	ib.saveLoadBlock(isSave, baseName, format, compression)
	ib.saveLoadChromosomes(isSave, baseName, format, compression)
}

func (ib *IBrowser) saveLoadBlock(isSave bool, outPrefix string, format string, compression string) {
	newPrefix := outPrefix + "_block"

	if isSave {
		fmt.Println("saving global ibrowser block")
		ib.block.Save(newPrefix, format, compression)
	} else {
		fmt.Println("loading global ibrowser block")
		ib.block.Load(newPrefix, format, compression)
	}
}

func (ib *IBrowser) saveLoadChromosomes(isSave bool, outPrefix string, format string, compression string) {
	for chromosomePos := 0; chromosomePos < len(ib.ChromosomesNames); chromosomePos++ {
		chromosomeName := ib.ChromosomesNames[chromosomePos]

		if isSave {
			fmt.Println("saving chromosome        : ", chromosomeName)
			chromosome := ib.chromosomes[chromosomeName]
			chromosome.Save(outPrefix, format, compression)

		} else {
			fmt.Println("loading chromosome       : ", chromosomeName)
			ib.chromosomes[chromosomeName] = NewIBChromosome(chromosomeName, ib.BlockSize, ib.NumBits, ib.NumSamples, ib.KeepEmptyBlock)
			chromosome := ib.chromosomes[chromosomeName]
			chromosome.Load(outPrefix, format, compression)
		}
	}
}
