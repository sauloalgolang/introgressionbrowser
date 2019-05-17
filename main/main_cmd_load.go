package main

import (
	log "github.com/sirupsen/logrus"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
)

// LoadCommand commandline load parameters
type LoadCommand struct {
	Infile          LoadArgsOptions `long:"indb" description:"Input database prefix" positional-args:"true" positional-arg-name:"Input Database Prefix" hidden:"true"`
	IsNotSoft       bool            `long:"notSoftLoad" description:"Force load matrices"`
	ProfileOptions  ProfileOptions
	SaveLoadOptions SaveLoadOptions
	LoggerOptions   LoggerOptions
}

// LoadArgsOptions commandline load parameters - options
type LoadArgsOptions struct {
	DbPrefix string `long:"indb" description:"Input database prefix" required:"true" positional-arg-name:"Input Database Prefix"`
}

// LoadCommand instance
var loadCommand LoadCommand

// Execute runs the processing of the commandline parameters
func (x *LoadCommand) Execute(args []string) error {
	x.LoggerOptions.Process()

	log.Printf("Load\n")

	// sourceFile := processArgs(args)
	sourceFile := x.Infile.DbPrefix

	parameters := Parameters{}
	x.SaveLoadOptions.NoCheck = false

	x.SaveLoadOptions.ProcessParameters(&parameters)

	log.Printf(" sourceFile             : %s\n", sourceFile)
	log.Println(parameters)
	log.Println(x.LoggerOptions)
	log.Println(x.ProfileOptions)

	x.SaveLoadOptions.Process()
	profileCloser := x.ProfileOptions.Process()

	log.Println("Openning", sourceFile)

	ibrowser := ibrowser.NewIBrowser(parameters)
	isSoft := !x.IsNotSoft
	ibrowser.Load(sourceFile, x.SaveLoadOptions.Format, x.SaveLoadOptions.Compression, isSoft)

	if !x.SaveLoadOptions.NoCheck {
		checkRes := ibrowser.Check()

		if checkRes {
			log.Println("Passed all tests")
		} else {
			log.Println("Failed tests")
		}
	}

	profileCloser()

	return nil
}

func init() {
	parser.AddCommand("load",
		"Load database",
		"Load previously created database",
		&loadCommand)
}
