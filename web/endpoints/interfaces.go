package endpoints

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/save"
)

var DATABASE_DIR = "res/"
var VERBOSITY = log.WarnLevel

type Parameters = ibrowser.Parameters
type IBrowser = ibrowser.IBrowser
type IBChromosome = ibrowser.IBChromosome
type IBBlock = ibrowser.IBBlock
type IBMatrix = ibrowser.DistanceMatrix

var GuessPrefixFormat = save.GuessPrefixFormat
var GuessFormat = save.GuessFormat
var NewIBrowser = ibrowser.NewIBrowser

//
// DbDb
//
type DbDb struct {
	Databases map[string]*DatabaseInfo
}

func NewDbDb() (db *DbDb) {
	db = &DbDb{
		Databases: make(map[string]*DatabaseInfo, 0),
	}
	return db
}

func (d *DbDb) Register(fileName string, path string) (err error) {
	err = nil

	if _, ok := d.Databases[fileName]; ok {
		log.Debug("Registering db :: filename: '%s' path: '%s' - Exists", fileName, path)
		return err

	} else {
		log.Infof("Registering db :: filename: '%s' path: '%s' - loading\n", fileName, path)

		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			if recover() != nil {
				log.Warning("Error registering", fileName, path)
				err = errors.New(fmt.Sprintf("Error registering :: fileName '%s' path '%s'", fileName, path))
				// return err
				return
			}
		}()

		ib := NewIBrowser(Parameters{})
		ib.EasyLoadFile(path, true)

		dbi := NewDatabaseInfo(ib)
		dbi.Name = fileName
		dbi.Path = path

		d.Databases[fileName] = dbi

		return err
	}
}

//
// Get data from ibrowser
//

func (d *DbDb) getDatabase(fileName string) (*IBrowser, bool) {
	if dbi, ok := d.Databases[fileName]; ok {
		return dbi.ib, ok
	} else {
		return nil, ok
	}
}

func (d *DbDb) getChromosome(fileName string, chromosome string) (*IBChromosome, bool) {
	db, hasDb := d.getDatabase(fileName)

	if !hasDb {
		return nil, hasDb
	}

	chrom, hasChrom := db.GetChromosome(chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	return chrom, hasChrom
}

func (d *DbDb) getBlock(fileName string, chromosome string, blockNum uint64) (*IBBlock, bool) {
	chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	block, hasBlock := chrom.GetBlock(blockNum)

	if !hasBlock {
		return nil, hasBlock
	}

	return block, hasBlock
}

func (d *DbDb) getMatrix(fileName string, chromosome string, blockNum uint64) (*IBMatrix, bool) {
	block, hasBlock := d.getBlock(fileName, chromosome, blockNum)

	if !hasBlock {
		return nil, hasBlock
	}

	matrix := block.Matrix

	return matrix, true
}

//
// Get web versions of data
//

func (d *DbDb) GetDatabases() (files []*DatabaseInfo) {
	files = make([]*DatabaseInfo, 0, len(d.Databases))

	for _, value := range d.Databases {
		files = append(files, value)
	}

	return files
}

func (d *DbDb) GetDatabase(fileName string) (*DatabaseInfo, bool) {
	if dbi, ok := d.Databases[fileName]; ok {
		return dbi, ok
	} else {
		return nil, ok
	}
}

func (d *DbDb) GetDatabaseBlock(fileName string) (*BlockInfo, bool) {
	if dbi, ok := d.GetDatabase(fileName); ok {
		bl := NewBlockInfo(dbi.ib.Block)
		return bl, ok
	} else {
		return nil, ok
	}
}

func (d *DbDb) GetChromosomes(fileName string) ([]*ChromosomeInfo, bool) {
	db, hasDb := d.GetDatabase(fileName)

	if !hasDb {
		return nil, hasDb
	}

	dbib := db.ib

	numChromosomes := len(dbib.ChromosomesNames)
	chromosomes := make([]*ChromosomeInfo, numChromosomes, numChromosomes)

	for cl, chromNamePosPair := range dbib.ChromosomesNames {
		chromName := chromNamePosPair.Name
		// chromPos := chromNamePosPair.Pos
		chromosome := dbib.Chromosomes[chromName]

		ci := NewChromosomeInfo(chromosome)

		chromosomes[cl] = ci
	}

	return chromosomes, true
}

func (d *DbDb) GetChromosome(fileName string, chromosome string) (*ChromosomeInfo, bool) {
	chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	ci := NewChromosomeInfo(chrom)

	return ci, true
}

func (d *DbDb) GetChromosomeBlock(fileName string, chromosome string) (*BlockInfo, bool) {
	if dbi, ok := d.GetChromosome(fileName, chromosome); ok {
		bl := NewBlockInfo(dbi.chromosome.Block)
		return bl, ok
	} else {
		return nil, ok
	}
}

func (d *DbDb) GetBlocks(fileName string, chromosome string) ([]*BlockInfo, bool) {
	chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	numBlocks := chrom.NumBlocks
	blocks := make([]*BlockInfo, numBlocks, numBlocks)

	for bl, block := range chrom.Blocks {
		blocks[bl] = NewBlockInfo(block)
	}

	return blocks, true
}

func (d *DbDb) GetBlock(fileName string, chromosome string, blockNum uint64) (*BlockInfo, bool) {
	block, hasBlock := d.getBlock(fileName, chromosome, blockNum)

	if !hasBlock {
		return nil, hasBlock
	}

	bi := NewBlockInfo(block)

	return bi, true
}

func (d *DbDb) GetBlockMatrix(fileName string, chromosome string, blockNum uint64) (*MatrixInfo, bool) {
	matrix, hasMatrix := d.getMatrix(fileName, chromosome, blockNum)

	if !hasMatrix {
		return nil, hasMatrix
	}

	mi := NewMatrixInfo(matrix)

	return mi, true
}

//
// DatabaseInfo
//

type DatabaseInfo struct {
	Name             string
	Path             string
	Parameters       Parameters
	Samples          []string
	NumSamples       uint64
	BlockSize        uint64
	KeepEmptyBlock   bool
	NumRegisters     uint64
	NumSNPS          uint64
	NumBlocks        uint64
	CounterBits      int
	ChromosomesNames []string
	ib               *IBrowser
}

func NewDatabaseInfo(ib *IBrowser) (di *DatabaseInfo) {
	di = &DatabaseInfo{
		// Name: ib.Name,
		// Path:           Path
		Parameters:     ib.Parameters,
		Samples:        ib.Samples,
		NumSamples:     ib.NumSamples,
		BlockSize:      ib.BlockSize,
		KeepEmptyBlock: ib.KeepEmptyBlock,
		NumRegisters:   ib.NumRegisters,
		NumSNPS:        ib.NumSNPS,
		NumBlocks:      ib.NumBlocks,
		CounterBits:    ib.CounterBits,
		ib:             ib,
	}

	chromosomesNames := ib.ChromosomesNames
	di.ChromosomesNames = make([]string, len(chromosomesNames), len(chromosomesNames))

	for p, k := range chromosomesNames {
		di.ChromosomesNames[p] = k.Name
	}

	return
}

func (d DatabaseInfo) String() (res string) {
	res += fmt.Sprintf(" Name             %s\n", d.Name)
	res += fmt.Sprintf(" Path             %s\n", d.Path)
	res += fmt.Sprintf(" NumSamples       %d\n", d.NumSamples)
	res += fmt.Sprintf(" BlockSize        %d\n", d.BlockSize)
	res += fmt.Sprintf(" KeepEmptyBlock   %#v\n", d.KeepEmptyBlock)
	res += fmt.Sprintf(" NumRegisters     %d\n", d.NumRegisters)
	res += fmt.Sprintf(" NumSNPS          %d\n", d.NumSNPS)
	res += fmt.Sprintf(" NumBlocks        %d\n", d.NumBlocks)
	res += fmt.Sprintf(" CounterBits      %d\n", d.CounterBits)
	res += fmt.Sprintf(" ChromosomesNames %s\n", strings.Join(d.ChromosomesNames, ", "))
	res += fmt.Sprintf(" Samples          %s\n", strings.Join(d.Samples, ", "))
	res += fmt.Sprintf("%s\n", d.Parameters)
	return
}

//
// ChromosomeInfo
//

type ChromosomeInfo struct {
	Name        string
	Number      int
	MinPosition uint64
	MaxPosition uint64
	NumBlocks   uint64
	NumSNPS     uint64
	chromosome  *IBChromosome
}

func NewChromosomeInfo(chromosome *IBChromosome) (c *ChromosomeInfo) {
	c = &ChromosomeInfo{
		Name:        chromosome.ChromosomeName,
		Number:      chromosome.ChromosomeNumber,
		MinPosition: chromosome.MinPosition,
		MaxPosition: chromosome.MaxPosition,
		NumBlocks:   chromosome.NumBlocks,
		NumSNPS:     chromosome.NumSNPS,
		chromosome:  chromosome,
	}

	return
}

func (c ChromosomeInfo) String() (res string) {
	res += fmt.Sprintf(" Name             %s\n", c.Name)
	res += fmt.Sprintf(" Number           %d\n", c.Number)
	res += fmt.Sprintf(" MinPosition      %d\n", c.MinPosition)
	res += fmt.Sprintf(" MaxPosition      %d\n", c.MaxPosition)
	res += fmt.Sprintf(" NumBlocks        %d\n", c.NumBlocks)
	res += fmt.Sprintf(" NumSNPS          %d\n", c.NumSNPS)
	return
}

//
// BlockInfo
//

type BlockInfo struct {
	MinPosition   uint64
	MaxPosition   uint64
	NumSNPS       uint64
	NumSamples    uint64
	BlockPosition uint64
	BlockNumber   uint64
	Serial        int64
	block         *IBBlock
}

func NewBlockInfo(block *IBBlock) (m *BlockInfo) {
	m = &BlockInfo{
		MinPosition:   block.MinPosition,
		MaxPosition:   block.MaxPosition,
		NumSNPS:       block.NumSNPS,
		NumSamples:    block.NumSamples,
		BlockPosition: block.BlockPosition,
		BlockNumber:   block.BlockNumber,
		Serial:        block.Serial,
		block:         block,
	}
	return
}

func (b BlockInfo) String() (res string) {
	res += fmt.Sprintf(" MinPosition    %d\n", b.MinPosition)
	res += fmt.Sprintf(" MaxPosition    %d\n", b.MaxPosition)
	res += fmt.Sprintf(" NumSNPS        %d\n", b.NumSNPS)
	res += fmt.Sprintf(" BlockPosition  %d\n", b.BlockPosition)
	res += fmt.Sprintf(" BlockNumber    %d\n", b.BlockNumber)
	res += fmt.Sprintf(" Serial         %d\n", b.Serial)
	return res
}

//
// MatrixInfo
//

type MatrixInfo struct {
	Dimension     uint64
	Size          uint64
	BlockPosition uint64
	BlockNumber   uint64
	Serial        int64
	matrix        *IBMatrix
}

func NewMatrixInfo(matrix *IBMatrix) (m *MatrixInfo) {
	m = &MatrixInfo{
		Dimension:     matrix.Dimension,
		Size:          matrix.Size,
		BlockPosition: matrix.BlockPosition,
		BlockNumber:   matrix.BlockNumber,
		Serial:        matrix.Serial,
		matrix:        matrix,
	}
	return
}

func (m MatrixInfo) String() (res string) {
	res += fmt.Sprintf(" Dimension     %d\n", m.Dimension)
	res += fmt.Sprintf(" Size          %d\n", m.Size)
	res += fmt.Sprintf(" BlockPosition %d\n", m.BlockPosition)
	res += fmt.Sprintf(" BlockNumber   %d\n", m.BlockNumber)
	res += fmt.Sprintf(" Serial        %d\n", m.Serial)
	return res
}

//
// List new databases
//

func ListDatabases() {
	log.Tracef("ListDatabases")

	err := filepath.Walk(DATABASE_DIR, func(path string, info os.FileInfo, err error) error {
		found, _, _, prefix := GuessFormat(path)

		if found {
			log.Tracef("ListDatabases :: path '%s' valid database", path)

			fi, err := os.Stat(path)

			if err != nil {
				log.Fatal(err)
			}

			if fi.Mode().IsRegular() {
				log.Tracef("ListDatabases :: path '%s' is file", path)

				fn := strings.TrimPrefix(prefix, DATABASE_DIR)
				parts := filepath.SplitList(fn)
				fn = strings.Join(parts, " - ")

				log.Tracef("ListDatabases :: path '%s' prefix '%s'", path, fn)

				reg_err := databases.Register(fn, path)

				if reg_err == nil {
					log.Tracef("ListDatabases :: path '%s' prefix '%s' - success registering", path, fn)
				} else {
					log.Tracef("ListDatabases :: path '%s' prefix '%s' - failed registering", path, fn)
				}
			} else {
				log.Tracef("ListDatabases :: path '%s' is folder", path)
			}
		} else {
			log.Tracef("ListDatabases :: path '%s' invalid database", path)
		}
		return nil
	})

	if err != nil {
		log.Panic(err)
	}
}
