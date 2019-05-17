package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// import "github.com/jessevdk/go-flags"
import "github.com/sauloalgolang/go-flags"

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
)

// CallBackParameters - alias to interfaces.CallBackParameters
type CallBackParameters = interfaces.CallBackParameters

// Parameters - alias to interfaces.Parameters
type Parameters = interfaces.Parameters

// Options holds all commandline options
type Options struct {
	// LoggerOptions LoggerOptions
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// get the arguments from the command line
	_, argErr := parser.Parse()

	// options.LoggerOptions.Process()

	if argErr != nil {
		if flagsErr, ok := argErr.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			log.Println("errHelp")
			// log.Println(argErr)
			os.Exit(0)
		} else {
			log.Println("error")
			// log.Println(argErr)
			os.Exit(1)
		}
	}
}
