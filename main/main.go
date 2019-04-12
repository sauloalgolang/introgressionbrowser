package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	// "strconv"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "Write memory profile to `file`")
var outfile = flag.String("outfile", "output", "Output file prefix")
var format = flag.String("format", "yaml", "File format: yaml, bson, binary")
var chromosomes = flag.String("chromosomes", "", "Comma separated list of chromomomes to read")
var continueOnError = flag.Bool("continueonerror", true, "Continue reading the file on parsing error")
var blockSize = flag.Uint64("blocksize", 100000, "Block size")
var minSnpPerBlock = flag.Uint64("minsnpperblock", 10, "Minimum number of SNPs per block")
var maxSnpPerBlock = flag.Uint64("maxsnpperblock", math.MaxUint64, "Maximum number of SNPs per block")
var keepEmptyBlock = flag.Bool("keepemptyblocks", true, "Keep empty blocks")
var numThreads = flag.Int("numberthreads", 4, "Number of threads")
var debug = flag.Bool("debug", false, "Print debug information")
var debug_first_only = flag.Bool("debug_first_only", false, "Read only fist chromosome from each thread")
var debug_maxregister = flag.Int64("debug_maxregister", 0, "Maximum number of registers to read per thread")
var version = flag.Bool("version", false, "Print version and exit")

// var BREAKAT string

func main() {
	// get the arguments from the command line
	flag.Parse()

	fmt.Println("cpuprofile       :", *cpuprofile)
	fmt.Println("memprofile       :", *memprofile)
	fmt.Println("outfile          :", *outfile)
	fmt.Println("format           :", *format)
	fmt.Println("chromosomes      :", *chromosomes) // TODO: implement
	fmt.Println("continueonerror  :", *continueOnError)
	fmt.Println("blocksize        :", *blockSize)
	fmt.Println("minsnpperblock   :", *minSnpPerBlock) // TODO: implement
	fmt.Println("maxsnpperblock   :", *maxSnpPerBlock) // TODO: implement
	fmt.Println("keepemptyblock   :", *keepEmptyBlock)
	fmt.Println("numthreads       :", *numThreads)
	fmt.Println("debug            :", *debug)
	fmt.Println("debug_first_only :", *debug_first_only)
	fmt.Println("debug_maxregister:", *debug_maxregister)
	fmt.Println("version          :", *version)

	vcf.DEBUG = *debug
	vcf.ONLYFIRST = *debug_first_only

	if *debug_maxregister != 0 {
		vcf.BREAKAT = *debug_maxregister
	}

	if *version {
		fmt.Println("IBROWSER_GIT_COMMIT_HASH    :", IBROWSER_GIT_COMMIT_HASH)
		fmt.Println("IBROWSER_GIT_COMMIT_AUTHOR  :", IBROWSER_GIT_COMMIT_AUTHOR)
		fmt.Println("IBROWSER_GIT_COMMIT_COMMITER:", IBROWSER_GIT_COMMIT_COMMITER)
		fmt.Println("IBROWSER_GIT_COMMIT_NOTES   :", IBROWSER_GIT_COMMIT_NOTES)
		fmt.Println("IBROWSER_GIT_COMMIT_TITLE   :", IBROWSER_GIT_COMMIT_TITLE)
		fmt.Println("IBROWSER_GIT_STATUS         :", IBROWSER_GIT_STATUS)
		fmt.Println("IBROWSER_GIT_DIFF           :", IBROWSER_GIT_DIFF)
		os.Exit(0)
	}

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	if _, ok := save.Formats[*format]; !ok {
		fmt.Println("Unknown format: ", *format, ". valid formats are:")
		for k, _ := range save.Formats {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}

	sourceFile := flag.Arg(0)

	if sourceFile == "" {
		fmt.Println("Dude, you didn't pass a input file!")
		os.Exit(1)
	} else {
		fmt.Println("Openning", sourceFile)
	}

	ibrowser := ibrowser.NewIBrowser(vcf.ProcessVcf, *blockSize, *keepEmptyBlock)

	vcf.OpenVcfFile(sourceFile, *continueOnError, *numThreads, ibrowser.ReaderCallBack)

	ibrowser.Save(*outfile, *format)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
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
