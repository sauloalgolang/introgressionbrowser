package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

func main() {
	// get the arguments from the command line

	// numPtr := flag.Int("n", 4, "an integer")

	continueOnError := *flag.Bool("continueOnError", true, "continue reading the file on error")

	flag.Parse()

	sourceFile := flag.Arg(0)

	ibrowser := ibrowser.NewIBrowser(vcf.ProcessVcf)

	if sourceFile == "" {
		fmt.Println("Dude, you didn't pass a input file!")
		os.Exit(1)
	} else {
		fmt.Println("Openning", sourceFile)
	}

	if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf.tar.gz") {
		fmt.Println(" .tar.gz format")
		openfile.OpenFile(sourceFile, true, true, continueOnError, ibrowser.ReaderCallBack)
	} else if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf.gz") {
		fmt.Println(" .gz format")
		openfile.OpenFile(sourceFile, false, true, continueOnError, ibrowser.ReaderCallBack)
	} else if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf") {
		fmt.Println(" .vcf format")
		openfile.OpenFile(sourceFile, false, false, continueOnError, ibrowser.ReaderCallBack)
	} else {
		fmt.Println("unknown file suffix!")
		os.Exit(1)
	}
}
