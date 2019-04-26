package main

import (
	"fmt"
	"log"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
)

type LoadCommand struct {
	Infile          LoadArgsOptions `long:"indb" description:"Input database prefix" positional-args:"true" positional-arg-name:"Input Database Prefix" hidden:"true"`
	Soft            bool            `long:"softLoad" description:"Do not load matrices"`
	ProfileOptions  ProfileOptions
	SaveLoadOptions SaveLoadOptions
}

type LoadArgsOptions struct {
	DbPrefix string `long:"indb" description:"Input database prefix" required:"true" positional-arg-name:"Input Database Prefix"`
}

var loadCommand LoadCommand

func (x *LoadCommand) Execute(args []string) error {
	fmt.Printf("Load\n")

	// sourceFile := processArgs(args)
	sourceFile := x.Infile.DbPrefix

	parameters := Parameters{}
	x.SaveLoadOptions.NoCheck = false

	processSaveLoadParameters(&parameters, x.SaveLoadOptions)

	fmt.Printf(" sourceFile             : %s\n", sourceFile)
	fmt.Println(parameters)
	fmt.Println(x.ProfileOptions)

	processSaveLoad(x.SaveLoadOptions)
	profileCloser := processProfile(x.ProfileOptions)

	log.Println("Openning", sourceFile)

	ibrowser := ibrowser.NewIBrowser(parameters)

	ibrowser.Load(sourceFile, x.SaveLoadOptions.Format, x.SaveLoadOptions.Compression, x.Soft)

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
