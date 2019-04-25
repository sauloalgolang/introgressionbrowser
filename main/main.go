package main

import (
	// "flag"
	"fmt"
	// "log"
	// "math"
	"os"
	// "strings"
)

// import "github.com/jessevdk/go-flags"
import "github.com/sauloalgolang/go-flags"

import (
	// "github.com/sauloalgolang/introgressionbrowser/api"
	// "github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	// "github.com/sauloalgolang/introgressionbrowser/save"
	// "github.com/sauloalgolang/introgressionbrowser/vcf"
)

type CallBackParameters = interfaces.CallBackParameters
type Parameters = interfaces.Parameters

var DEFAULT_BLOCK_SIZE = uint64(100000)
var DEFAULT_OUTFILE = "output"
var DEFAULT_COUNTER_BITS = 32

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
