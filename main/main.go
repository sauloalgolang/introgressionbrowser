package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

var cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "Write memory profile to `file`")
var outfile = flag.String("outfile", "output", "Outfile prefix")
var format = flag.String("format", "yaml", "File format: yaml, bson")
var continueOnError = flag.Bool("continueonerror", true, "Continue reading the file on error")
var blockSize = flag.Uint64("blocksize", 1000000, "Block size")
var keepEmptyBlock = flag.Bool("keepemptyblocks", true, "Keep empty blocks")
var numThreads = flag.Int("numberthreads", 4, "Number of threads")

func main() {
	// get the arguments from the command line
	flag.Parse()

	fmt.Println("cpuprofile     :", *cpuprofile)
	fmt.Println("memprofile     :", *memprofile)
	fmt.Println("outfile        :", *outfile)
	fmt.Println("format         :", *format)
	fmt.Println("continueonerror:", *continueOnError)
	fmt.Println("blocksize      :", *blockSize)
	fmt.Println("keepemptyblock :", *keepEmptyBlock)
	fmt.Println("numthreads     :", *numThreads)

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

	sourceFile := flag.Arg(0)

	if sourceFile == "" {
		fmt.Println("Dude, you didn't pass a input file!")
		os.Exit(1)
	} else {
		fmt.Println("Openning", sourceFile)
	}

	ibrowser := ibrowser.NewIBrowser(vcf.ProcessVcf, *blockSize, *keepEmptyBlock)

	vcf.OpenVcfFile(sourceFile, *continueOnError, *numThreads, ibrowser.ReaderCallBack)

	// ibrowser.SaveChromosomes(*outfile, *format)
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
