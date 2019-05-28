package main

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

// SaveCommand commandline save parameters
type SaveCommand struct {
	BlockSize         uint64          `long:"blockSize" description:"Block size" default:"100000"`
	Chromosomes       string          `long:"chromosomes" description:"Comma separated list of chromomomes to read" default:""`
	NoContinueOnError bool            `long:"continueOnError" description:"Continue reading the file on parsing error"`
	CounterBits       uint64          `long:"counterBits" description:"Number of bits" default:"32"`
	NoKeepEmptyBlock  bool            `long:"keepEmptyBlocks" description:"Keep empty blocks"`
	MaxSnpPerBlock    uint64          `long:"maxSnpPerBlock" description:"Maximum number of SNPs per block" default:"18446744073709551615"`
	MinSnpPerBlock    uint64          `long:"minSnpPerBlock" description:"Minimum number of SNPs per block" default:"10"`
	Outfile           string          `long:"outfile" description:"Output file prefix" default:"res/output"`
	Description       string          `long:"description" description:"Description of the database" default:""`
	Infile            SaveArgsOptions `long:"infile" description:"Input VCF file" positional-args:"true" positional-arg-name:"Input VCF file" hidden:"true"`
	ProfileOptions    ProfileOptions
	SaveLoadOptions   SaveLoadOptions
	DebugOptions      DebugOptions
	LoggerOptions     LoggerOptions
}

// SaveArgsOptions commandline save parameters - options
type SaveArgsOptions struct {
	VCF string `long:"infile" description:"Input VCF file" required:"true" positional-arg-name:"Input VCF file"`
}

// SaveCommand instance
var saveCommand SaveCommand

// Execute runs the processing of the commandline parameters
func (x *SaveCommand) Execute(args []string) error {
	x.LoggerOptions.Process()

	log.Printf("Save\n")

	// sourceFile := processArgs(args)
	sourceFile := x.Infile.VCF
	outfile := x.Outfile

	fi, err := os.Stat(sourceFile)
	if err != nil {
		log.Println(err)
		return nil
	}

	if !fi.Mode().IsRegular() {
		log.Println("input file ", sourceFile, " is not a file")
		os.Exit(1)
	}

	parameters := Parameters{
		SourceFile: sourceFile,
		Outfile:    outfile,
	}

	if x.Description == "" {
		x.Description = filepath.Base(sourceFile)
	}

	x.DebugOptions.ProcessParameters(&parameters)
	x.SaveLoadOptions.ProcessParameters(&parameters)
	x.ProcessParameters(&parameters)

	log.Printf(" sourceFile             : %s\n", sourceFile)
	log.Println(parameters)
	// log.Println(x.LoggerOptions)
	log.Println(x.ProfileOptions)
	log.Println(x.DebugOptions)

	x.DebugOptions.Process()
	x.SaveLoadOptions.Process()
	profileCloser := x.ProfileOptions.Process()

	log.Println("Openning", sourceFile)

	ibrowser := ibrowser.NewIBrowser(parameters)

	callBackParameters := CallBackParameters{
		ContinueOnError: !x.NoContinueOnError,
		NumBits:         x.CounterBits,
		NumThreads:      x.SaveLoadOptions.NumThreads,
	}

	vcf.OpenFile(sourceFile, callBackParameters, ibrowser.RegisterCallBack)

	if !x.SaveLoadOptions.NoCheck {
		checkRes := ibrowser.Check()

		if checkRes {
			log.Println("Passed all tests")
		} else {
			log.Println("Failed tests")
		}
	}

	ibrowser.Save(
		x.SaveLoadOptions.Format,
		x.SaveLoadOptions.Compression,
	)

	profileCloser()

	return nil
}

// ProcessParameters processes parameters
func (x *SaveCommand) ProcessParameters(parameters *Parameters) {
	parameters.BlockSize = x.BlockSize
	parameters.Chromosomes = x.Chromosomes
	parameters.ContinueOnError = !x.NoContinueOnError
	parameters.CounterBits = x.CounterBits
	parameters.Description = x.Description
	parameters.KeepEmptyBlock = !x.NoKeepEmptyBlock
	parameters.MaxSnpPerBlock = x.MaxSnpPerBlock
	parameters.MinSnpPerBlock = x.MinSnpPerBlock
}

func init() {
	parser.AddCommand("save",
		"Read VCF and save database",
		"Read VCF and save database",
		&saveCommand)
}
