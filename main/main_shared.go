package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/save"
)

type SaveLoadOptions struct {
	NoCheck     bool   `long:"check" description:"Check for self consistency"`
	Compression string `long:"compression" description:"Compression format: none, snappy, gzip" choice:"none" choice:"snappy" choice:"gzip" default:"none"`
	Format      string `long:"format" description:"File format: yaml" choice:"yaml" default:"yaml"`
	NumThreads  int    `long:"threads" description:"Number of threads" default:"4"`
}

func (s SaveLoadOptions) String() (res string) {
	res += fmt.Sprintf("Debug:\n")
	res += fmt.Sprintf(" Check                  : %#v\n", !s.NoCheck)
	res += fmt.Sprintf(" Compression            : %s\n", s.Compression)
	res += fmt.Sprintf(" Format                 : %s\n", s.Format)
	res += fmt.Sprintf(" NumThreads             : %d\n", s.NumThreads)
	return res
}

type ProfileOptions struct {
	CpuProfile     string `long:"cpuProfile" description:"Write cpu profile to file" default:""`
	MemProfile     string `long:"memProfile" description:"Write memory profile to file" default:""`
	cpuFileHandler *os.File
}

func (p ProfileOptions) String() (res string) {
	res += fmt.Sprintf("Profile:\n")
	res += fmt.Sprintf(" CpuProfile             : %s\n", p.CpuProfile)
	res += fmt.Sprintf(" MemProfile             : %s\n", p.MemProfile)
	return res
}

func processSaveLoad(opts SaveLoadOptions) {
	if _, ok := save.Formats[opts.Format]; !ok {
		fmt.Println("Unknown format: ", opts.Format, ". valid formats are:")
		for k, _ := range save.Formats {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}

	if _, ok := save.Compressors[opts.Compression]; !ok {
		fmt.Println("Unknown compression: ", opts.Compression, ". valid formats are:")
		for k, _ := range save.Compressors {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}
}

func profileCpuStart(opts ProfileOptions) {
	if opts.CpuProfile != "" {
		err := *new(error)
		opts.cpuFileHandler, err = os.Create(opts.CpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile", err)
		}
		if err := pprof.StartCPUProfile(opts.cpuFileHandler); err != nil {
			log.Fatal("could not start CPU profile", err)
		}
	}
}

func profileCpuEnd(opts ProfileOptions) {
	if opts.CpuProfile != "" {
		pprof.StopCPUProfile()
		opts.cpuFileHandler.Close()
	}
}

func profileMemStart(opts ProfileOptions) {

}

func profileMemEnd(opts ProfileOptions) {
	runtime.GC() // get up-to-date statistics

	if opts.MemProfile != "" {
		f, err := os.Create(opts.MemProfile)
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

func processProfile(opts ProfileOptions) func() {
	profileCpuStart(opts)
	profileMemStart(opts)

	return func() {
		profileCpuEnd(opts)
		profileMemEnd(opts)
	}
}

func processArgs(args []string) (sourceFile string) {
	if len(args) == 0 {
		fmt.Println("no database prefix given")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("more than one database prefix given")
		os.Exit(1)
	}

	sourceFile = args[0]

	return sourceFile
}

func processSaveLoadParameters(parameters *Parameters, saveLoadOptions SaveLoadOptions) {
	parameters.Compression = saveLoadOptions.Compression
	parameters.Format = saveLoadOptions.Format
}
