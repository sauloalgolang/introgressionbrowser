package main

import (
	"fmt"
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
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {
	// get the arguments from the command line
	_, argErr := parser.Parse()
	if argErr != nil {
		if flagsErr, ok := argErr.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			fmt.Println("errHelp")
			// fmt.Println(argErr)
			os.Exit(0)
		} else {
			fmt.Println("error")
			// fmt.Println(argErr)
			os.Exit(1)
		}
	}
}
