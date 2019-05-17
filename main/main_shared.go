package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"runtime/pprof"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

//
// SaveLoadOptions
//

// SaveLoadOptions holds the sahred parameters for the commandline save and load commands
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

// Process processes save/load options
func (s *SaveLoadOptions) Process() {
	if _, ok := save.Formats[s.Format]; !ok {
		log.Println("Unknown format: ", s.Format, ". valid formats are:")
		for k := range save.Formats {
			log.Println(" ", k)
		}
		os.Exit(1)
	}

	if _, ok := save.Compressors[s.Compression]; !ok {
		log.Println("Unknown compression: ", s.Compression, ". valid formats are:")
		for k := range save.Compressors {
			log.Println(" ", k)
		}
		os.Exit(1)
	}
}

// ProcessParameters processes parameters
func (s *SaveLoadOptions) ProcessParameters(parameters *Parameters) {
	parameters.Compression = s.Compression
	parameters.Format = s.Format
}

//
// ProfileOptions
//

// ProfileOptions holds the shared parameters for the commandline profile options
type ProfileOptions struct {
	CPUProfile     string `long:"CPUProfile" description:"Write cpu profile to file" default:""`
	MemProfile     string `long:"memProfile" description:"Write memory profile to file" default:""`
	cpuFileHandler *os.File
}

func (p ProfileOptions) String() (res string) {
	res += fmt.Sprintf("Profile:\n")
	res += fmt.Sprintf(" CPUProfile             : %s\n", p.CPUProfile)
	res += fmt.Sprintf(" MemProfile             : %s\n", p.MemProfile)
	return res
}

// CPUStart starts cpu profiling
func (p *ProfileOptions) CPUStart() {
	if p.CPUProfile != "" {
		err := *new(error)
		p.cpuFileHandler, err = os.Create(p.CPUProfile)
		if err != nil {
			log.Fatal("could not create CPU profile", err)
		}
		if err := pprof.StartCPUProfile(p.cpuFileHandler); err != nil {
			log.Fatal("could not start CPU profile", err)
		}
	}
}

// CPUEnd ends cpu profiling
func (p *ProfileOptions) CPUEnd() {
	if p.CPUProfile != "" {
		pprof.StopCPUProfile()
		p.cpuFileHandler.Close()
	}
}

// MemStart starts memory profiling
func (p ProfileOptions) MemStart() {

}

// MemEnd ends memory profiling
func (p ProfileOptions) MemEnd() {
	runtime.GC() // get up-to-date statistics

	if p.MemProfile != "" {
		f, err := os.Create(p.MemProfile)
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

// Process profile options
func (p ProfileOptions) Process() func() {
	p.CPUStart()
	p.MemStart()

	return func() {
		p.CPUEnd()
		p.MemEnd()
	}
}

//
// Debug
//

// DebugOptions commandline debug options
type DebugOptions struct {
	Debug                  bool  `long:"debug" description:"Print debug information"`
	DebugFirstOnly         bool  `long:"debugFirstOnly" description:"Read only fist chromosome from each thread"`
	DebugMaxRegisterThread int64 `long:"debugMaxRegisterThread" description:"Maximum number of registers to read per thread" default:"0"`
	DebugMaxRegisterChrom  int64 `long:"debugMaxRegisterChrom" description:"Maximum number of registers to read per chromosome" default:"0"`
}

func (d DebugOptions) String() (res string) {
	res += fmt.Sprintf("Debug:\n")
	res += fmt.Sprintf(" Debug                  : %#v\n", d.Debug)
	res += fmt.Sprintf(" DebugFirstOnly         : %#v\n", d.DebugFirstOnly)
	res += fmt.Sprintf(" DebugMaxRegisterThread : %d\n", d.DebugMaxRegisterThread)
	res += fmt.Sprintf(" DebugMaxRegisterChrom  : %d\n", d.DebugMaxRegisterChrom)
	return res
}

// Process debug options
func (d *DebugOptions) Process() {
	vcf.Debug = d.Debug
	vcf.OnlyFirst = d.DebugFirstOnly
	vcf.BreakAtThread = d.DebugMaxRegisterThread
	vcf.BreakAtChrom = d.DebugMaxRegisterChrom
}

// ProcessParameters processes parameters
func (d *DebugOptions) ProcessParameters(parameters *Parameters) {
	parameters.DebugFirstOnly = d.DebugFirstOnly
	parameters.DebugMaxRegisterThread = d.DebugMaxRegisterThread
	parameters.DebugMaxRegisterChrom = d.DebugMaxRegisterChrom
}


//
// Args
// 

func processArgs(args []string) (sourceFile string) {
	if len(args) == 0 {
		log.Println("no database prefix given")
		os.Exit(1)
	}

	if len(args) > 1 {
		log.Println("more than one database prefix given")
		os.Exit(1)
	}

	sourceFile = args[0]

	return sourceFile
}
