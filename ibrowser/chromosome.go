package ibrowser

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

import "runtime/debug"

import "github.com/sauloalgolang/introgressionbrowser/tools"

//
//
// CHROMOSOME SECTION
//
//

type IBChromosome struct {
	Chromosome     string
	MinPosition    uint64
	MaxPosition    uint64
	NumBlocks      uint64
	NumSNPS        uint64
	NumSamples     uint64
	KeepEmptyBlock bool
	Block          *IBBlock
	Blocks         []*IBBlock
	BlockNames     map[uint64]uint64
}

func NewIBChromosome(chromosome string, numSamples uint64, keepEmptyBlock bool) *IBChromosome {
	ibc := IBChromosome{
		Chromosome:     chromosome,
		NumSamples:     numSamples,
		MinPosition:    0,
		MaxPosition:    0,
		NumBlocks:      0,
		NumSNPS:        0,
		KeepEmptyBlock: keepEmptyBlock,
		Block:          NewIBBlock(0, numSamples),
		Blocks:         make([]*IBBlock, 0, 100),
		BlockNames:     make(map[uint64]uint64, 100),
	}

	return &ibc
}

func (ibc *IBChromosome) AppendBlock(blockNum uint64) {
	ibc.Blocks = append(ibc.Blocks, NewIBBlock(blockNum, ibc.NumSamples))
	ibc.BlockNames[blockNum] = uint64(len(ibc.Blocks)) - uint64(1)
	ibc.NumBlocks++
}

func (ibc *IBChromosome) GetBlock(blockNum uint64) (*IBBlock, bool) {
	if blockPos, ok := ibc.BlockNames[blockNum]; ok {
		if blockPos >= uint64(len(ibc.Blocks)) {
			fmt.Println(&ibc, "Index out of range. block num:", blockNum, "block pos:", blockPos, "len:", len(ibc.Blocks), "NumBlocks:", ibc.NumBlocks)
			fmt.Println(&ibc, "BlockNames", ibc.BlockNames)
			fmt.Println(&ibc, "Blocks", ibc.Blocks)
			debug.PrintStack()
			os.Exit(1)
		}

		return ibc.Blocks[blockPos], ok
	} else {
		return &IBBlock{}, ok
	}
}

func (ibc *IBChromosome) normalizeBlocks(blockNum uint64) {
	if _, hasBlock := ibc.GetBlock(blockNum); !hasBlock {
		if ibc.KeepEmptyBlock {
			lastBlockPos := uint64(0)
			NumBlocks := uint64(len(ibc.Blocks))

			if NumBlocks != 0 {
				lastBlockPos = NumBlocks - 1
			}

			for currBlockPos := lastBlockPos; currBlockPos < blockNum; currBlockPos++ {
				ibc.AppendBlock(currBlockPos)
			}
		}
		ibc.AppendBlock(blockNum)
	}
}

func (ibc *IBChromosome) Add(blockNum uint64, position uint64, distance *tools.DistanceMatrix) {
	ibc.normalizeBlocks(blockNum)

	if block, success := ibc.GetBlock(blockNum); success {
		block.Add(position, distance)
		ibc.Block.Add(position, distance)
		ibc.NumSNPS++
		ibc.MinPosition = tools.Min64(ibc.MinPosition, block.MinPosition)
		ibc.MaxPosition = tools.Max64(ibc.MaxPosition, block.MaxPosition)
	} else {
		fmt.Println("Failure getting block", blockNum)
		os.Exit(1)
	}
}

func (ibc *IBChromosome) Save(outfile string) {
	// ibB, _ := json.Marshal(ib)
	// fmt.Println(string(ibB))

	d, err := yaml.Marshal(ibc)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	// fmt.Printf("--- dump:\n%s\n\n", d)
	fmt.Println("saving chromosome ", ibc.Chromosome, " to ", outfile)
	err = ioutil.WriteFile(outfile+"yaml", d, 0644)
	fmt.Println("done")
}
