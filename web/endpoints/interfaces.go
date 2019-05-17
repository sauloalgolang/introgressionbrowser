package endpoints

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

// DatabaseDir is the folder to search for database
var DatabaseDir = "database/"

// DataEndpoint is the folder where data should be found
var DataEndpoint = "dataep/"

// Verbosity is the verbosity level
// var Verbosity = log.WarnLevel

// SliceIndex is the function to search for the index of a value in a list
var SliceIndex = tools.SliceIndex

// GuessPrefixFormat is the function to guess the format of the file given a file prefix
var GuessPrefixFormat = save.GuessPrefixFormat

// GuessFormat is the function to guess the format of the file given a file name
var GuessFormat = save.GuessFormat

// NewIBrowser is the function to create a new ibrowser instance
var NewIBrowser = ibrowser.NewIBrowser

// Parameters are the command line parameter
type Parameters = ibrowser.Parameters

// IBrowser is the ibrowser main struct
type IBrowser = ibrowser.IBrowser

// IBChromosome is the ibrowser chromosome struct
type IBChromosome = ibrowser.IBChromosome

// IBBlock is the ibrowser block struct
type IBBlock = ibrowser.IBBlock

// IBMatrix is the ibrowser matrix struct
type IBMatrix = ibrowser.IBDistanceMatrix

// IBDistanceTable is the ibrowser matrix table struct
type IBDistanceTable = ibrowser.IBDistanceTable

//
// DbDb
//

// DbDb is the struct holding all databases
type DbDb struct {
	Databases map[string]*DatabaseInfo
}

// NewDbDb creates a new DbDb instance
func NewDbDb() (db *DbDb) {
	db = &DbDb{
		Databases: make(map[string]*DatabaseInfo, 0),
	}
	return db
}

// Register registers a new database
func (d *DbDb) Register(fileName string, path string) (err error) {
	err = nil

	if _, ok := d.Databases[fileName]; ok {
		log.Debugf("Registering db :: filename: '%s' path: '%s' - Exists", fileName, path)
		return err

	}

	log.Infof("Registering db :: filename: '%s' path: '%s' - loading\n", fileName, path)

	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			log.Warningf("Registering db :: filename: '%s' path: '%s' - Error registering", fileName, path)
			err = fmt.Errorf("Error registering :: fileName '%s' path '%s'", fileName, path)
			return
		}
	}()

	ib := NewIBrowser(Parameters{})
	ib.EasyLoadFile(path, true)

	dbi := NewDatabaseInfo(fileName, path, ib)

	d.Databases[fileName] = dbi

	return err
}

//
// Get database
//

func (d *DbDb) getDatabase(fileName string) (*DatabaseInfo, *IBrowser, bool) {
	dbi, hasDb := d.Databases[fileName]

	if !hasDb {
		return nil, nil, hasDb
	}

	ib := dbi.ib

	return dbi, ib, true
}

//
// Get database summary
//

func (d *DbDb) getDatabaseSummaryBlock(fileName string) (*DatabaseInfo, *IBrowser, *IBBlock, bool) {
	dbi, ib, hasDatabase := d.getDatabase(fileName)

	if !hasDatabase {
		return nil, nil, nil, hasDatabase
	}

	block, hasBlock := ib.GetSummaryBlock()

	if !hasBlock {
		return nil, nil, nil, hasBlock
	}

	return dbi, ib, block, true
}

func (d *DbDb) getDatabaseSummaryBlockMatrix(fileName string) (*DatabaseInfo, *IBrowser, *IBBlock, *IBMatrix, bool) {
	dbi, ib, block, hasBlock := d.getDatabaseSummaryBlock(fileName)

	if !hasBlock {
		return nil, nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, nil, hasMatrix
	}

	return dbi, ib, block, matrix, true
}

func (d *DbDb) getDatabaseSummaryBlockMatrixData(fileName string) (*DatabaseInfo, *IBrowser, *IBBlock, *IBMatrix, *IBDistanceTable, bool) {
	dbi, ib, block, matrix, hasMatrix := d.getDatabaseSummaryBlockMatrix(fileName)

	if !hasMatrix {
		return nil, nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := block.GetMatrixTable()

	if !hasTable {
		return nil, nil, nil, nil, nil, hasTable
	}

	return dbi, ib, block, matrix, table, true
}

//
// Get chromosome
//

func (d *DbDb) getChromosomeNames(fileName string) (*DatabaseInfo, []string, bool) {
	dbi, ib, ok := d.getDatabase(fileName)

	if !ok {
		return nil, nil, ok
	}

	chromosomes := ib.GetChromosomeNames()

	return dbi, chromosomes, true
}

func (d *DbDb) getChromosomes(fileName string) (*DatabaseInfo, *IBrowser, []*IBChromosome, bool) {
	dbi, ib, ok := d.getDatabase(fileName)

	if !ok {
		return nil, nil, nil, ok
	}

	chromosomes := ib.GetChromosomes()

	return dbi, ib, chromosomes, true
}

func (d *DbDb) getChromosome(fileName string, chromosome string) (*DatabaseInfo, *IBrowser, *IBChromosome, bool) {
	dbi, ib, hasDb := d.getDatabase(fileName)

	if !hasDb {
		return nil, nil, nil, hasDb
	}

	chrom, hasChrom := ib.GetChromosome(chromosome)

	if !hasChrom {
		return nil, nil, nil, hasChrom
	}

	return dbi, ib, chrom, true
}

func (d *DbDb) getChromosomeSummaryBlock(fileName string, chromosome string) (*DatabaseInfo, *IBrowser, *IBChromosome, *IBBlock, bool) {
	dbi, ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, nil, nil, nil, hasChrom
	}

	block, hasBlock := chrom.GetSummaryBlock()

	if !hasBlock {
		return nil, nil, nil, nil, hasBlock
	}

	return dbi, ib, chrom, block, true
}

func (d *DbDb) getChromosomeSummaryBlockMatrix(fileName string, chromosome string) (*DatabaseInfo, *IBrowser, *IBChromosome, *IBBlock, *IBMatrix, bool) {
	dbi, ib, chrom, block, hasBlock := d.getChromosomeSummaryBlock(fileName, chromosome)

	if !hasBlock {
		return nil, nil, nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, nil, nil, hasMatrix
	}

	return dbi, ib, chrom, block, matrix, true
}

func (d *DbDb) getChromosomeSummaryBlockMatrixTable(fileName string, chromosome string) (*DatabaseInfo, *IBrowser, *IBChromosome, *IBBlock, *IBMatrix, *IBDistanceTable, bool) {
	dbi, ib, chrom, block, matrix, hasMatrix := d.getChromosomeSummaryBlockMatrix(fileName, chromosome)

	if !hasMatrix {
		return nil, nil, nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, nil, nil, nil, nil, nil, hasTable
	}

	return dbi, ib, chrom, block, matrix, table, true
}

//
// Get blocks
//
func (d *DbDb) getChromosomeBlocks(fileName string, chromosome string) (*DatabaseInfo, *IBrowser, *IBChromosome, []*IBBlock, bool) {
	dbi, ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, nil, nil, nil, hasChrom
	}

	blocks, hasBlock := chrom.GetBlocks()

	if !hasBlock {
		return nil, nil, nil, nil, hasBlock
	}

	return dbi, ib, chrom, blocks, true
}

func (d *DbDb) getChromosomeBlock(fileName string, chromosome string, blockNum uint64) (*DatabaseInfo, *IBrowser, *IBChromosome, *IBBlock, bool) {
	dbi, ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, nil, nil, nil, hasChrom
	}

	block, hasBlock := chrom.GetBlock(blockNum)

	if !hasBlock {
		return nil, nil, nil, nil, hasBlock
	}

	return dbi, ib, chrom, block, true
}

func (d *DbDb) getChromosomeBlockMatrix(fileName string, chromosome string, blockNum uint64) (*DatabaseInfo, *IBrowser, *IBChromosome, *IBBlock, *IBMatrix, bool) {
	dbi, ib, chrom, block, hasBlock := d.getChromosomeBlock(fileName, chromosome, blockNum)

	if !hasBlock {
		return nil, nil, nil, nil, nil, hasBlock
	}

	matrix, hasMatrix := block.GetMatrix()

	if !hasMatrix {
		return nil, nil, nil, nil, nil, hasMatrix
	}

	return dbi, ib, chrom, block, matrix, true
}

func (d *DbDb) getChromosomeBlockMatrixTable(fileName string, chromosome string, blockNum uint64) (*DatabaseInfo, *IBrowser, *IBChromosome, *IBBlock, *IBMatrix, *IBDistanceTable, bool) {
	dbi, ib, chrom, block, matrix, hasMatrix := d.getChromosomeBlockMatrix(fileName, chromosome, blockNum)

	if !hasMatrix {
		return nil, nil, nil, nil, nil, nil, hasMatrix
	}

	table, hasTable := matrix.GetTable()

	if !hasTable {
		return nil, nil, nil, nil, nil, nil, hasTable
	}

	return dbi, ib, chrom, block, matrix, table, true
}

//
// Get web versions of data
//

//
// Databases

// GetDatabases returns a list of all databases
func (d *DbDb) GetDatabases() (files []*DatabaseInfo) {
	files = make([]*DatabaseInfo, 0, len(d.Databases))

	for _, value := range d.Databases {
		files = append(files, value)
	}

	return files
}

// GetDatabase returns a database
func (d *DbDb) GetDatabase(fileName string) (*DatabaseInfo, bool) {
	if dbi, ok := d.Databases[fileName]; ok {
		return dbi, ok
	}

	return nil, false
}

//
// Database summary block

// GetDatabaseSummaryBlock returns the summary block of a database
func (d *DbDb) GetDatabaseSummaryBlock(fileName string) (*BlockInfo, bool) {
	dbi, ib, block, ok := d.getDatabaseSummaryBlock(fileName)

	if !ok {
		return nil, ok
	}

	bi := NewBlockInfo(dbi, ib, nil, block)

	return bi, true
}

// GetDatabaseSummaryBlockMatrix returns the summary block matrix of a database
func (d *DbDb) GetDatabaseSummaryBlockMatrix(fileName string) (*MatrixInfo, bool) {
	dbi, ib, block, matrix, ok := d.getDatabaseSummaryBlockMatrix(fileName)

	if !ok {
		return nil, ok
	}

	mi := NewMatrixInfo(dbi, ib, nil, block, matrix)

	return mi, true
}

// GetDatabaseSummaryBlockMatrixTable returns the summary block matrix table of a database
func (d *DbDb) GetDatabaseSummaryBlockMatrixTable(fileName string) (*TableInfo, bool) {
	dbi, ib, block, matrix, table, ok := d.getDatabaseSummaryBlockMatrixData(fileName)

	if !ok {
		return nil, ok
	}

	ti := NewTableInfo(dbi, ib, nil, block, matrix, table, true)

	return ti, true
}

//
// Chromosomes

// GetChromosomes returns a list of the choromosomes of a database
func (d *DbDb) GetChromosomes(fileName string) ([]*ChromosomeInfo, bool) {
	dbi, ib, chromosomes, ok := d.getChromosomes(fileName)

	if !ok {
		return nil, ok
	}

	numChromosomes := len(chromosomes)
	chromosomesi := make([]*ChromosomeInfo, numChromosomes, numChromosomes)

	for cl, chromosome := range chromosomes {
		ci := NewChromosomeInfo(dbi, ib, chromosome)
		chromosomesi[cl] = ci
	}

	return chromosomesi, true
}

// GetChromosome returns a chromosome of a database
func (d *DbDb) GetChromosome(fileName string, chromosome string) (*ChromosomeInfo, bool) {
	dbi, ib, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	ci := NewChromosomeInfo(dbi, ib, chrom)

	return ci, true
}

// GetChromosomeSummaryBlock returns a chromosome summary block of a database
func (d *DbDb) GetChromosomeSummaryBlock(fileName string, chromosome string) (*BlockInfo, bool) {
	dbi, ib, chrom, block, ok := d.getChromosomeSummaryBlock(fileName, chromosome)

	if !ok {
		return nil, ok
	}

	bl := NewBlockInfo(dbi, ib, chrom, block)

	return bl, true
}

// GetChromosomeSummaryBlockMatrix returns a chromosome summary block matrix of a database
func (d *DbDb) GetChromosomeSummaryBlockMatrix(fileName string, chromosome string) (*MatrixInfo, bool) {
	dbi, ib, chrom, block, matrix, ok := d.getChromosomeSummaryBlockMatrix(fileName, chromosome)

	if !ok {
		return nil, ok
	}

	bl := NewMatrixInfo(dbi, ib, chrom, block, matrix)

	return bl, true
}

// GetChromosomeSummaryBlockMatrixTable returns a chromosome summary block matrix table of a database
func (d *DbDb) GetChromosomeSummaryBlockMatrixTable(fileName string, chromosome string) (*TableInfo, bool) {
	dbi, ib, chrom, block, matrix, table, ok := d.getChromosomeSummaryBlockMatrixTable(fileName, chromosome)

	if !ok {
		return nil, ok
	}

	bl := NewTableInfo(dbi, ib, chrom, block, matrix, table, true)

	return bl, true
}

//
// Blocks

// GetBlocks returns a list of blocks for a chromosome in a database
func (d *DbDb) GetBlocks(fileName string, chromosome string) ([]*BlockInfo, bool) {
	dbi, ib, chrom, blocks, hasChrom := d.getChromosomeBlocks(fileName, chromosome)

	if !hasChrom {
		return nil, hasChrom
	}

	numBlocks := len(blocks)
	blocksi := make([]*BlockInfo, numBlocks, numBlocks)

	for bl, block := range blocks {
		blocksi[bl] = NewBlockInfo(dbi, ib, chrom, block)
	}

	return blocksi, true
}

// GetBlock returns a block for a chromosome in a database
func (d *DbDb) GetBlock(fileName string, chromosome string, blockNum uint64) (*BlockInfo, bool) {
	dbi, ib, chrom, block, hasBlock := d.getChromosomeBlock(fileName, chromosome, blockNum)

	if !hasBlock {
		return nil, hasBlock
	}

	bi := NewBlockInfo(dbi, ib, chrom, block)

	return bi, true
}

// GetBlockMatrix returns a block matrix for a chromosome in a database
func (d *DbDb) GetBlockMatrix(fileName string, chromosome string, blockNum uint64) (*MatrixInfo, bool) {
	dbi, ib, chrom, block, matrix, ok := d.getChromosomeBlockMatrix(fileName, chromosome, blockNum)

	if !ok {
		return nil, ok
	}

	mi := NewMatrixInfo(dbi, ib, chrom, block, matrix)

	return mi, true
}

// GetBlockMatrixTable returns a block matrix table for a chromosome in a database
func (d *DbDb) GetBlockMatrixTable(fileName string, chromosome string, blockNum uint64) (*TableInfo, bool) {
	dbi, ib, chrom, block, matrix, table, ok := d.getChromosomeBlockMatrixTable(fileName, chromosome, blockNum)

	if !ok {
		return nil, ok
	}

	ti := NewTableInfo(dbi, ib, chrom, block, matrix, table, false)

	return ti, true
}

//
// Plots
//

// router.HandleFunc(PLOTS_ENDPOINT+"/{database}/{chromosome}/{referenceName}", endpoints.Plots).Methods("GET").Name("plots")

func (d *DbDb) referenceName2referenceNumber(fileName string, referenceName string) (referenceNumber int, ok bool) {
	referenceNumber = 0
	ok = false

	_, ib, hasDb := d.getDatabase(fileName)

	if !hasDb {
		return
	}

	referenceNumber, ok = ib.GetSampleID(referenceName)

	return
}

func (d *DbDb) referenceNumber2referenceName(fileName string, referenceNumber int) (referenceName string, ok bool) {
	referenceName = ""
	ok = false

	_, ib, hasDb := d.getDatabase(fileName)

	if !hasDb {
		return
	}

	referenceName, ok = ib.GetSampleName(referenceNumber)

	return
}

// GetPlotTable returns a table ready to be plotted
func (d *DbDb) GetPlotTable(fileName string, chromosome string, referenceName string) (*PlotInfo, bool) {
	referenceNumber, hasRfn := d.referenceName2referenceNumber(fileName, referenceName)

	log.Printf("GetPlotTable :: fileName %s chromosome %s referenceName %s referenceNumber %d\n",
		fileName,
		chromosome,
		referenceName,
		referenceNumber)

	if !hasRfn {
		return nil, false
	}

	_, _, chrom, hasChrom := d.getChromosome(fileName, chromosome)

	if !hasChrom {
		return nil, false
	}

	distanceTable, hasTable := chrom.GetColumn(referenceNumber)

	if !hasTable {
		return nil, false
	}

	pi := NewPlotInfo(distanceTable)

	return pi, true
}

// PlotInfo contains the information to generate a plot
type PlotInfo struct {
	DistanceTable *[]*IBDistanceTable
}

// NewPlotInfo generates a new instance of PlotInfo
func NewPlotInfo(d *[]*IBDistanceTable) (pi *PlotInfo) {
	pi = &PlotInfo{
		DistanceTable: d,
	}
	return pi
}

//
//
//
// STRUCT
//
//
//

//
// DatabaseInfo
//

// DatabaseInfo contains the information of a database
type DatabaseInfo struct {
	DatabaseName     string
	FilePath         string
	Parameters       Parameters
	Samples          []string
	NumSamples       uint64
	BlockSize        uint64
	KeepEmptyBlock   bool
	NumRegisters     uint64
	NumSNPS          uint64
	NumBlocks        uint64
	CounterBits      uint64
	ChromosomesNames []string
	ib               *IBrowser
}

// NewDatabaseInfo creates a new DatabaseInfo instance
func NewDatabaseInfo(databaseName string, filePath string, ib *IBrowser) (di *DatabaseInfo) {
	di = &DatabaseInfo{
		DatabaseName:   databaseName,
		FilePath:       filePath,
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
	res += fmt.Sprintf(" DatabaseName     %s\n", d.DatabaseName)
	res += fmt.Sprintf(" FilePath         %s\n", d.FilePath)
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

// ChromosomeInfo contains the information of a chromosome
type ChromosomeInfo struct {
	DatabaseName string
	Name         string
	Number       int
	MinPosition  uint64
	MaxPosition  uint64
	NumBlocks    uint64
	NumSNPS      uint64
	chromosome   *IBChromosome
	ib           *IBrowser
	dbi          *DatabaseInfo
}

// NewChromosomeInfo creates a new ChromosomeInfo instance
func NewChromosomeInfo(dbi *DatabaseInfo, ib *IBrowser, chromosome *IBChromosome) (c *ChromosomeInfo) {
	c = &ChromosomeInfo{
		DatabaseName: dbi.DatabaseName,
		Name:         chromosome.ChromosomeName,
		Number:       chromosome.ChromosomeNumber,
		MinPosition:  chromosome.MinPosition,
		MaxPosition:  chromosome.MaxPosition,
		NumSNPS:      chromosome.NumSNPS,
		chromosome:   chromosome,
		ib:           ib,
		dbi:          dbi,
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

// BlockInfo contains the information of a block
type BlockInfo struct {
	DatabaseName  string
	MinPosition   uint64
	MaxPosition   uint64
	NumSNPS       uint64
	NumSamples    uint64
	BlockPosition uint64
	BlockNumber   uint64
	Serial        uint64
	block         *IBBlock
	chromosome    *IBChromosome
	ib            *IBrowser
	dbi           *DatabaseInfo
}

// NewBlockInfo creates a new BlockInfo instance
func NewBlockInfo(dbi *DatabaseInfo, ib *IBrowser, chromosome *IBChromosome, block *IBBlock) (m *BlockInfo) {
	m = &BlockInfo{
		DatabaseName:  dbi.DatabaseName,
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
		dbi:           dbi,
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

// MatrixInfo contains the information of a block matrix
type MatrixInfo struct {
	DatabaseName  string
	Dimension     uint64
	Size          uint64
	BlockPosition uint64
	BlockNumber   uint64
	Serial        uint64
	matrix        *IBMatrix
	block         *IBBlock
	chromosome    *IBChromosome
	ib            *IBrowser
	dbi           *DatabaseInfo
}

// NewMatrixInfo creates a new MatrixInfo instance
func NewMatrixInfo(dbi *DatabaseInfo, ib *IBrowser, chromosome *IBChromosome, block *IBBlock, matrix *IBMatrix) (m *MatrixInfo) {
	m = &MatrixInfo{
		DatabaseName:  dbi.DatabaseName,
		Dimension:     matrix.Dimension,
		Size:          matrix.Size,
		BlockPosition: matrix.BlockPosition,
		BlockNumber:   matrix.BlockNumber,
		Serial:        matrix.Serial,
		matrix:        matrix,
		block:         block,
		chromosome:    chromosome,
		ib:            ib,
		dbi:           dbi,
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

// TableInfo contains the information of a block matrix table
type TableInfo struct {
	DatabaseName     string
	FileName         string
	RegisterPosition uint64
	RegisterSize     uint64
	Serial           uint64
	matrix           *IBMatrix
	block            *IBBlock
	chromosome       *IBChromosome
	ib               *IBrowser
	dbi              *DatabaseInfo
}

// NewTableInfo creates a new TableInfo instance
func NewTableInfo(dbi *DatabaseInfo, ib *IBrowser, chromosome *IBChromosome, block *IBBlock, matrix *IBMatrix, table *IBDistanceTable, isSummary bool) (m *TableInfo) {
	fileName := ""
	
	if chromosome == nil {
		fileName = ib.GenMatrixDumpFileName(dbi.FilePath)
	} else {
		fileName = ib.GenMatrixChromosomeDumpFileName(dbi.FilePath, chromosome.ChromosomeName)
	}

	if DatabaseDir[len(DatabaseDir)-1] == '/' {
		fileName = strings.TrimPrefix(fileName, DatabaseDir)
	} else {
		fileName = strings.TrimPrefix(fileName, DatabaseDir+"/")
	}
	fileName = strings.Join([]string{strings.TrimSuffix(DataEndpoint, "/"), fileName}, "/")


	RegisterPosition := ib.RegisterSize * uint64(matrix.Serial)


	m = &TableInfo{
		DatabaseName:     dbi.DatabaseName,
		FileName:         fileName,
		RegisterPosition: RegisterPosition,
		RegisterSize:     ib.RegisterSize,
		Serial:           uint64(matrix.Serial),
		matrix:           matrix,
		block:            block,
		chromosome:       chromosome,
		ib:               ib,
		dbi:              dbi,
	}

	return
}

func (t TableInfo) String() (res string) {
	res += fmt.Sprintf(" FileName         %s\n", t.FileName)
	res += fmt.Sprintf(" RegisterPosition %d\n", t.RegisterPosition)
	res += fmt.Sprintf(" RegisterSize     %d\n", t.RegisterSize)
	res += fmt.Sprintf(" Serial           %d\n", t.Serial)
	return res
}

//
// List new databases
//

// ListDatabases list all databases in a folder
func ListDatabases() {
	log.Tracef("ListDatabases")

	err := filepath.Walk(DatabaseDir, func(path string, info os.FileInfo, err error) error {
		found, _, _, prefix := GuessFormat(path)

		if found {
			log.Tracef("ListDatabases :: path '%s' valid database", path)

			fi, err := os.Stat(path)

			if err != nil {
				log.Fatal(err)
			}

			if fi.Mode().IsRegular() {
				log.Tracef("ListDatabases :: path '%s' is file", path)

				fn := strings.TrimPrefix(prefix, DatabaseDir+"/")
				parts := filepath.SplitList(fn)
				fn = strings.Join(parts, " - ")

				log.Tracef("ListDatabases :: path '%s' prefix '%s'", path, fn)

				regErr := databases.Register(fn, path)

				if regErr == nil {
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
