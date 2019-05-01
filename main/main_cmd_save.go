package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

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
}

type SaveArgsOptions struct {
	VCF string `long:"infile" description:"Input VCF file" required:"true" positional-arg-name:"Input VCF file"`
}

type DebugOptions struct {
	Debug                  bool  `long:"debug" description:"Print debug information"`
	DebugFirstOnly         bool  `long:"debugFirstOnly" description:"Read only fist chromosome from each thread"`
	DebugMaxRegisterThread int64 `long:"debugMaxRegisterThread" description:"Maximum number of registers to read per thread" default:"0"`
	DebugMaxRegisterChrom  int64 `long:"debugMaxRegisterChrom" description:"Maximum number of registers to read per chromosome" default:"0"`
}

func (d DebugOptions) String() (res string) {
	res += fmt.Sprintf("Debug:\n")
	res += fmt.Sprintf(" Debug                  : %#v\n", d.Debug)
	res += fmt.Sprintf(" DebugFirstOnly         : %#v\n", d.Debug)
	res += fmt.Sprintf(" DebugMaxRegisterThread : %d\n", d.Debug)
	res += fmt.Sprintf(" DebugMaxRegisterChrom  : %d\n", d.Debug)
	return res
}

var saveCommand SaveCommand

func (x *SaveCommand) Execute(args []string) error {
	fmt.Printf("Save\n")

	// sourceFile := processArgs(args)
	sourceFile := x.Infile.VCF

	fi, err := os.Stat(sourceFile)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if !fi.Mode().IsRegular() {
		fmt.Println("input file ", sourceFile, " is not a file")
		os.Exit(1)
	}

	parameters := Parameters{
		SourceFile: sourceFile,
	}

	if x.Description == "" {
		x.Description = filepath.Base(sourceFile)
	}

	processDebugParameters(&parameters, x.DebugOptions)
	processSaveLoadParameters(&parameters, x.SaveLoadOptions)
	processSaveParameters(&parameters, *x)

	fmt.Printf(" sourceFile             : %s\n", sourceFile)
	fmt.Println(parameters)
	fmt.Println(x.ProfileOptions)
	fmt.Println(x.DebugOptions)

	processDebug(x.DebugOptions)
	processSaveLoad(x.SaveLoadOptions)
	profileCloser := processProfile(x.ProfileOptions)

	log.Println("Openning", sourceFile)

	ibrowser := ibrowser.NewIBrowser(parameters)

	callBackParameters := CallBackParameters{
		ContinueOnError: !x.NoContinueOnError,
		NumBits:         x.CounterBits,
		NumThreads:      x.SaveLoadOptions.NumThreads,
	}

	vcf.OpenVcfFile(sourceFile, callBackParameters, ibrowser.RegisterCallBack)

	if !x.SaveLoadOptions.NoCheck {
		checkRes := ibrowser.Check()

		if checkRes {
			log.Println("Passed all tests")
		} else {
			log.Println("Failed tests")
		}
	}

	ibrowser.Save(x.Outfile, x.SaveLoadOptions.Format, x.SaveLoadOptions.Compression)

	profileCloser()

	return nil
}

func processDebug(opts DebugOptions) {
	vcf.DEBUG = opts.Debug
	vcf.ONLYFIRST = opts.DebugFirstOnly
	vcf.BREAKAT_THREAD = opts.DebugMaxRegisterThread
	vcf.BREAKAT_CHROM = opts.DebugMaxRegisterChrom
}

func processSaveParameters(parameters *Parameters, saveCommand SaveCommand) {
	parameters.BlockSize = saveCommand.BlockSize
	parameters.Chromosomes = saveCommand.Chromosomes
	parameters.ContinueOnError = !saveCommand.NoContinueOnError
	parameters.CounterBits = saveCommand.CounterBits
	parameters.Description = saveCommand.Description
	parameters.KeepEmptyBlock = !saveCommand.NoKeepEmptyBlock
	parameters.MaxSnpPerBlock = saveCommand.MaxSnpPerBlock
	parameters.MinSnpPerBlock = saveCommand.MinSnpPerBlock
}

func processDebugParameters(parameters *Parameters, debugOptions DebugOptions) {
	parameters.DebugFirstOnly = debugOptions.DebugFirstOnly
	parameters.DebugMaxRegisterThread = debugOptions.DebugMaxRegisterThread
	parameters.DebugMaxRegisterChrom = debugOptions.DebugMaxRegisterChrom
}

func init() {
	parser.AddCommand("save",
		"Read VCF and save database",
		"Read VCF and save database",
		&saveCommand)
}
