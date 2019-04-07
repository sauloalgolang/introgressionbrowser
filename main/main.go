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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var outfile = *flag.String("outfile prefix", "output", "write memory profile to `file`")
var continueOnError = *flag.Bool("continueOnError", true, "continue reading the file on error")
var blockSize = *flag.Uint64("BlockSize", 1000000, "Block size")
var keepEmptyBlock = *flag.Bool("KeepEmptyBlocks", true, "Keep empty blocks")
var numThreads = *flag.Int("NumberThreads", 4, "Number of threads")

func main() {
	// get the arguments from the command line
	flag.Parse()

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

	ibrowser := ibrowser.NewIBrowser(vcf.ProcessVcf, blockSize, keepEmptyBlock)

	if sourceFile == "" {
		fmt.Println("Dude, you didn't pass a input file!")
		os.Exit(1)
	} else {
		fmt.Println("Openning", sourceFile)
	}

	vcf.OpenVcfFile(sourceFile, continueOnError, numThreads, ibrowser.ReaderCallBack)

	// ibrowser.SaveChromosomes(outfile)
	ibrowser.Save(outfile)

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
