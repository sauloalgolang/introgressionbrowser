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
type IBDistanceTable = ibrowser.IBDistanceTable

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
// Get database
//

func (d *DbDb) getDatabase(fileName string) (*IBrowser, bool) {
	dbi, hasDb := d.Databases[fileName]

	if !hasDb {
		return nil, hasDb
	}

	ib := dbi.ib

	return ib, true
}

//
// Get database summary
//

func (d *DbDb) getDatabaseSummaryBlock(fileName string) (*IBrowser, *IBBlock, bool) {
	ib, hasDatabase := d.getDatabase(fileName)

	if !hasDatabase {
		return nil, nil, hasDatabase
	}

	block, hasBlock := ib.GetSummaryBlock()

	if !hasBlock {
		return nil, nil, hasBlock
	}

	return ib, block, true
}

func (d *DbDb) getDatabaseSummaryBlockMatrix(fileName string) (*IBrowser, *IBBlock, *IBMatrix, bool) {
	ib, block, hasBlock := d.getDatabaseSummaryBlock(fileName)

	if !hasBlock {
		return nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, hasMatrix
	}

	return ib, block, matrix, true
}

func (d *DbDb) getDatabaseSummaryBlockMatrixData(fileName string) (*IBrowser, *IBBlock, *IBMatrix, *IBDistanceTable, bool) {
	ib, block, matrix, hasMatrix := d.getDatabaseSummaryBlockMatrix(fileName)

	if !hasMatrix {
		return nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := block.GetMatrixData()

	if !hasTable {
		return nil, nil, nil, nil, hasTable
	}

	return ib, block, matrix, table, true
}

//
// Get chromosome
//

func (d *DbDb) getChromosomeNames(fileName string) ([]string, bool) {
	ib, ok := d.getDatabase(fileName)

	if !ok {
		return nil, ok
	}

	chromosomes := ib.GetChromosomeNames()

	return chromosomes, true
}

func (d *DbDb) getChromosomes(fileName string) (*IBrowser, []*IBChromosome, bool) {
	ib, ok := d.getDatabase(fileName)

	if !ok {
		return nil, nil, ok
	}

	chromosomes := ib.GetChromosomes()

	return ib, chromosomes, true
}

func (d *DbDb) getChromosome(fileName string, chromosome string) (*IBrowser, *IBChromosome, bool) {
	ib, hasDb := d.getDatabase(fileName)

	if !hasDb {
		return nil, nil, hasDb
	}

	chrom, hasChrom := ib.GetChromosome(chromosome)

	if !hasChrom {
		return nil, nil, hasChrom
	}

	return ib, chrom, true
}

func (d *DbDb) getChromosomeSummaryBlock(fileName string, chromosome string) (*IBrowser, *IBChromosome, *IBBlock, bool) {
	ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, nil, nil, hasChrom
	}

	block, hasBlock := chrom.GetSummaryBlock()

	if !hasBlock {
		return nil, nil, nil, hasBlock
	}

	return ib, chrom, block, true
}

func (d *DbDb) getChromosomeSummaryBlockMatrix(fileName string, chromosome string) (*IBrowser, *IBChromosome, *IBBlock, *IBMatrix, bool) {
	ib, chrom, block, hasBlock := d.getChromosomeSummaryBlock(fileName, chromosome)

	if !hasBlock {
		return nil, nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, nil, hasMatrix
	}

	return ib, chrom, block, matrix, true
}

func (d *DbDb) getChromosomeSummaryBlockMatrixTable(fileName string, chromosome string) (*IBrowser, *IBChromosome, *IBBlock, *IBMatrix, *IBDistanceTable, bool) {
	ib, chrom, block, matrix, hasMatrix := d.getChromosomeSummaryBlockMatrix(fileName, chromosome)

	if !hasMatrix {
		return nil, nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := matrix.GetMatrix()

	if !hasTable {
		return nil, nil, nil, nil, nil, hasTable
	}

	return ib, chrom, block, matrix, table, true
}

//
// Get blocks
//
func (d *DbDb) getChromosomeBlocks(fileName string, chromosome string) (*IBrowser, *IBChromosome, []*IBBlock, bool) {
	ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, nil, nil, hasChrom
	}

	blocks, hasBlock := chrom.GetBlocks()

	if !hasBlock {
		return nil, nil, nil, hasBlock
	}

	return ib, chrom, blocks, true
}

func (d *DbDb) getChromosomeBlock(fileName string, chromosome string, blockNum uint64) (*IBrowser, *IBChromosome, *IBBlock, bool) {
	ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, nil, nil, hasChrom
	}

	block, hasBlock := chrom.GetBlock(blockNum)

	if !hasBlock {
		return nil, nil, nil, hasBlock
	}

	return ib, chrom, block, true
}

func (d *DbDb) getChromosomeBlockMatrix(fileName string, chromosome string, blockNum uint64) (*IBrowser, *IBChromosome, *IBBlock, *IBMatrix, bool) {
	ib, chrom, block, hasBlock := d.getChromosomeBlock(fileName, chromosome, blockNum)

	if !hasBlock {
		return nil, nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, nil, hasMatrix
	}

	return ib, chrom, block, matrix, true
}

func (d *DbDb) getChromosomeBlockMatrixTable(fileName string, chromosome string, blockNum uint64) (*IBrowser, *IBChromosome, *IBBlock, *IBMatrix, *IBDistanceTable, bool) {
	ib, chrom, block, matrix, hasMatrix := d.getChromosomeBlockMatrix(fileName, chromosome, blockNum)

	if !hasMatrix {
		return nil, nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := matrix.GetMatrix()

	if !hasTable {
		return nil, nil, nil, nil, nil, hasTable
	}

	return ib, chrom, block, matrix, table, true
}

//
// Get web versions of data
//

//
// Database
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

//
// Database summary block
func (d *DbDb) GetDatabaseSummaryBlock(fileName string) (*BlockInfo, bool) {
	ib, block, ok := d.getDatabaseSummaryBlock(fileName)

	if !ok {
		return nil, ok
	}

	bi := NewBlockInfo(ib, nil, block)

	return bi, true
}

func (d *DbDb) GetDatabaseSummaryBlockMatrix(fileName string) (*MatrixInfo, bool) {
	ib, block, matrix, ok := d.getDatabaseSummaryBlockMatrix(fileName)

	if !ok {
		return nil, ok
	}

	mi := NewMatrixInfo(ib, nil, block, matrix)

	return mi, true
}

func (d *DbDb) GetDatabaseSummaryBlockMatrixTable(fileName string) (*TableInfo, bool) {
	ib, block, matrix, table, ok := d.getDatabaseSummaryBlockMatrixData(fileName)

	if !ok {
		return nil, ok
	}

	ti := NewTableInfo(ib, nil, block, matrix, table)

	return ti, true
}

//
// Chromosomes
func (d *DbDb) GetChromosomes(fileName string) ([]*ChromosomeInfo, bool) {
	ib, chromosomes, ok := d.getChromosomes(fileName)

	if !ok {
		return nil, ok
	}

	numChromosomes := len(chromosomes)
	chromosomesi := make([]*ChromosomeInfo, numChromosomes, numChromosomes)

	for cl, chromosome := range chromosomes {
		ci := NewChromosomeInfo(ib, chromosome)
		chromosomesi[cl] = ci
	}

	return chromosomesi, true
}

func (d *DbDb) GetChromosome(fileName string, chromosome string) (*ChromosomeInfo, bool) {
	ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	ci := NewChromosomeInfo(ib, chrom)

	return ci, true
}

func (d *DbDb) GetChromosomeSummaryBlock(fileName string, chromosome string) (*BlockInfo, bool) {
	ib, chrom, block, ok := d.getChromosomeSummaryBlock(fileName, chromosome)

	if !ok {
		return nil, ok
	}

	bl := NewBlockInfo(ib, chrom, block)

	return bl, true
}

func (d *DbDb) GetChromosomeSummaryBlockMatrix(fileName string, chromosome string) (*MatrixInfo, bool) {
	ib, chrom, block, matrix, ok := d.getChromosomeSummaryBlockMatrix(fileName, chromosome)

	if !ok {
		return nil, ok
	}

	bl := NewMatrixInfo(ib, chrom, block, matrix)

	return bl, true
}

func (d *DbDb) GetChromosomeSummaryBlockMatrixTable(fileName string, chromosome string) (*TableInfo, bool) {
	ib, chrom, block, matrix, table, ok := d.getChromosomeSummaryBlockMatrixTable(fileName, chromosome)

	if !ok {
		return nil, ok
	}

	bl := NewTableInfo(ib, chrom, block, matrix, table)

	return bl, true
}

//
// Blocks
func (d *DbDb) GetBlocks(fileName string, chromosome string) ([]*BlockInfo, bool) {
	ib, chrom, blocks, hasChrom := d.getChromosomeBlocks(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	numBlocks := len(blocks)
	blocksi := make([]*BlockInfo, numBlocks, numBlocks)

	for bl, block := range blocks {
		blocksi[bl] = NewBlockInfo(ib, chrom, block)
	}

	return blocksi, true
}

func (d *DbDb) GetBlock(fileName string, chromosome string, blockNum uint64) (*BlockInfo, bool) {
	ib, chrom, block, hasBlock := d.getChromosomeBlock(fileName, chromosome, blockNum)

	if !hasBlock {
		return nil, hasBlock
	}

	bi := NewBlockInfo(ib, chrom, block)

	return bi, true
}

func (d *DbDb) GetBlockMatrix(fileName string, chromosome string, blockNum uint64) (*MatrixInfo, bool) {
	ib, chrom, block, matrix, ok := d.getChromosomeBlockMatrix(fileName, chromosome, blockNum)

	if !ok {
		return nil, ok
	}

	mi := NewMatrixInfo(ib, chrom, block, matrix)

	return mi, true
}

func (d *DbDb) GetBlockMatrixTable(fileName string, chromosome string, blockNum uint64) (*TableInfo, bool) {
	ib, chrom, block, matrix, table, ok := d.getChromosomeBlockMatrixTable(fileName, chromosome, blockNum)

	if !ok {
		return nil, ok
	}

	ti := NewTableInfo(ib, chrom, block, matrix, table)

	return ti, true
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
	ib          *IBrowser
}

func NewChromosomeInfo(ib *IBrowser, chromosome *IBChromosome) (c *ChromosomeInfo) {
	c = &ChromosomeInfo{
		Name:        chromosome.ChromosomeName,
		Number:      chromosome.ChromosomeNumber,
		MinPosition: chromosome.MinPosition,
		MaxPosition: chromosome.MaxPosition,
		NumBlocks:   chromosome.NumBlocks,
		NumSNPS:     chromosome.NumSNPS,
		chromosome:  chromosome,
		ib:          ib,
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
	chromosome    *IBChromosome
	ib            *IBrowser
}

func NewBlockInfo(ib *IBrowser, chromosome *IBChromosome, block *IBBlock) (m *BlockInfo) {
	m = &BlockInfo{
		MinPosition:   block.MinPosition,
		MaxPosition:   block.MaxPosition,
		NumSNPS:       block.NumSNPS,
		NumSamples:    block.NumSamples,
		BlockPosition: block.BlockPosition,
		BlockNumber:   block.BlockNumber,
		Serial:        block.Serial,
		block:         block,
		chromosome:    chromosome,
		ib:            ib,
	}

	// output_360_merged_2.50.vcf.gz_chromosomes_SL2.50ch00.bin

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
	block         *IBBlock
	chromosome    *IBChromosome
	ib            *IBrowser
}

func NewMatrixInfo(ib *IBrowser, chromosome *IBChromosome, block *IBBlock, matrix *IBMatrix) (m *MatrixInfo) {
	m = &MatrixInfo{
		Dimension:     matrix.Dimension,
		Size:          matrix.Size,
		BlockPosition: matrix.BlockPosition,
		BlockNumber:   matrix.BlockNumber,
		Serial:        matrix.Serial,
		matrix:        matrix,
		block:         block,
		chromosome:    chromosome,
		ib:            ib,
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
// TableInfo
//

type TableInfo struct {
	matrix     *IBMatrix
	block      *IBBlock
	chromosome *IBChromosome
	ib         *IBrowser
}

func NewTableInfo(ib *IBrowser, chromosome *IBChromosome, block *IBBlock, matrix *IBMatrix, table *IBDistanceTable) (m *TableInfo) {
	m = &TableInfo{
		matrix:     matrix,
		block:      block,
		chromosome: chromosome,
		ib:         ib,
	}
	return
}

func (m TableInfo) String() (res string) {
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
