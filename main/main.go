package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strings"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

var DEFAULT_BLOCK_SIZE = uint64(100000)
var DEFAULT_OUTFILE = "output"
var DEFAULT_COUTNER_BITS = 32

var blockSize = flag.Uint64("blockSize", DEFAULT_BLOCK_SIZE, "Block size")
var chromosomes = flag.String("chromosomes", "", "Comma separated list of chromomomes to read")
var compression = flag.String("compression", save.DefaultCompressor, "Compression format: "+strings.Join(save.CompressorNames, ", "))
var continueOnError = flag.Bool("continueOnError", true, "Continue reading the file on parsing error")
var counterBits = flag.Int("counterbits", DEFAULT_COUTNER_BITS, "Number of bits")
var cpuProfile = flag.String("cpuProfile", "", "Write cpu profile to `file`")
var debug = flag.Bool("debug", false, "Print debug information")
var debugFirstOnly = flag.Bool("debugFirstOnly", false, "Read only fist chromosome from each thread")
var debugMaxRegisterThread = flag.Int64("debugMaxRegisterThread", 0, "Maximum number of registers to read per thread")
var debugMaxRegisterChrom = flag.Int64("debugMaxRegisterChrom", 0, "Maximum number of registers to read per chromosome")
var format = flag.String("format", save.DefaultFormat, "File format: "+strings.Join(save.FormatNames, ", "))
var keepEmptyBlock = flag.Bool("keepemptyblocks", true, "Keep empty blocks")
var maxSnpPerBlock = flag.Uint64("maxsnpperblock", math.MaxUint64, "Maximum number of SNPs per block")
var minSnpPerBlock = flag.Uint64("minsnpperblock", 10, "Minimum number of SNPs per block")
var memProfile = flag.String("memProfile", "", "Write memory profile to `file`")
var numThreads = flag.Int("threads", 4, "Number of threads")
var load = flag.Bool("load", false, "Load project")
var outfile = flag.String("outfile", "output", "Output file prefix")
var version = flag.Bool("version", false, "Print version and exit")

func main() {
	// get the arguments from the command line
	flag.Parse()

	fmt.Println("blockSize               :", *blockSize)
	fmt.Println("chromosomes             :", *chromosomes) // TODO: implement
	fmt.Println("compression             :", *compression)
	fmt.Println("continueOnError         :", *continueOnError)
	fmt.Println("counterbits             :", *counterBits)
	fmt.Println("cpuProfile              :", *cpuProfile)
	fmt.Println("debug                   :", *debug)
	fmt.Println("debugFirstOnly          :", *debugFirstOnly)
	fmt.Println("debugMaxRegisterThread  :", *debugMaxRegisterThread)
	fmt.Println("debugMaxRegisterChrom   :", *debugMaxRegisterChrom)
	fmt.Println("format                  :", *format)
	fmt.Println("keepemptyblock          :", *keepEmptyBlock)
	fmt.Println("maxsnpperblock          :", *maxSnpPerBlock) // TODO: implement
	fmt.Println("minsnpperblock          :", *minSnpPerBlock) // TODO: implement
	fmt.Println("memProfile              :", *memProfile)
	fmt.Println("numthreads              :", *numThreads)
	fmt.Println("load                    :", *load)
	fmt.Println("outfile                 :", *outfile)
	fmt.Println("version                 :", *version)

	vcf.DEBUG = *debug
	vcf.ONLYFIRST = *debugFirstOnly
	vcf.BREAKAT_THREAD = *debugMaxRegisterThread
	vcf.BREAKAT_CHROM = *debugMaxRegisterChrom

	if *version {
		fmt.Println("IBROWSER_GIT_COMMIT_AUTHOR  :", IBROWSER_GIT_COMMIT_AUTHOR)
		fmt.Println("IBROWSER_GIT_COMMIT_COMMITER:", IBROWSER_GIT_COMMIT_COMMITER)
		fmt.Println("IBROWSER_GIT_COMMIT_HASH    :", IBROWSER_GIT_COMMIT_HASH)
		fmt.Println("IBROWSER_GIT_COMMIT_NOTES   :", IBROWSER_GIT_COMMIT_NOTES)
		fmt.Println("IBROWSER_GIT_COMMIT_TITLE   :", IBROWSER_GIT_COMMIT_TITLE)
		fmt.Println("IBROWSER_GIT_STATUS         :", IBROWSER_GIT_STATUS)
		fmt.Println("IBROWSER_GIT_DIFF           :", IBROWSER_GIT_DIFF)
		fmt.Println("IBROWSER_GO_VERSION         :", IBROWSER_GO_VERSION)
		fmt.Println("IBROWSER_VERSION            :", IBROWSER_VERSION)
		os.Exit(0)
	}

	if _, ok := save.Formats[*format]; !ok {
		fmt.Println("Unknown format: ", *format, ". valid formats are:")
		for k, _ := range save.Formats {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}

	if _, ok := save.Compressors[*compression]; !ok {
		fmt.Println("Unknown compression: ", *compression, ". valid formats are:")
		for k, _ := range save.Compressors {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile", err)
		}
		defer pprof.StopCPUProfile()
	}

	//
	// Load or Save
	//

	sourceFile := flag.Arg(0)

	log.Println("Openning", sourceFile)

	if sourceFile == "" {
		log.Fatal("Dude, you didn't pass a input file!")
	}

	ibrowser := ibrowser.NewIBrowser(*blockSize, *counterBits, *keepEmptyBlock)

	if *load {
		if *blockSize != DEFAULT_BLOCK_SIZE {
			log.Fatal("cannot use blockSize and load at the same time")
		}

		if *chromosomes != "" {
			log.Fatal("cannot use chromosomes and load at the same time")
		}

		if *outfile != DEFAULT_OUTFILE {
			log.Fatal("cannot use outfile and load at the same time")
		}

		ibrowser.Load(sourceFile, *format, *compression)

	} else {
		callBackParameters := interfaces.CallBackParameters{
			ContinueOnError: *continueOnError,
			NumBits:         *counterBits,
			NumThreads:      *numThreads,
		}

		vcf.OpenVcfFile(sourceFile, callBackParameters, ibrowser.RegisterCallBack)

		ibrowser.Save(*outfile, *format, *compression)

	}

	runtime.GC() // get up-to-date statistics

	if *memProfile != "" {
		f, err := os.Create(*memProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
