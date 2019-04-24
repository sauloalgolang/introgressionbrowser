package main

import (
	"fmt"
	"log"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
)

type LoadCommand struct {
	ProfileOptions  ProfileOptions
	SaveLoadOptions SaveLoadOptions
}

var loadCommand LoadCommand

func (x *LoadCommand) Execute(args []string) error {
	fmt.Printf("Load\n")

	sourceFile := processArgs(args)

	parameters := Parameters{}

	processSaveLoadParameters(&parameters, x.SaveLoadOptions)

	fmt.Printf(" sourceFile             : %s\n", sourceFile)
	fmt.Println(parameters)
	fmt.Println(x.ProfileOptions)

	processSaveLoad(x.SaveLoadOptions)
	profileCloser := processProfile(x.ProfileOptions)

	log.Println("Openning", sourceFile)

	ibrowser := ibrowser.NewIBrowser(parameters)

	ibrowser.Load(sourceFile, x.SaveLoadOptions.Format, x.SaveLoadOptions.Compression)

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
