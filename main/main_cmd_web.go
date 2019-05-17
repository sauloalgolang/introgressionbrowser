package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/web"
)

// WebCommand holds the parameter for the commandline web command
type WebCommand struct {
	Host          string `long:"host" description:"Hostname" default:"127.0.0.1"`
	Port          int    `long:"port" description:"Port" default:"8000"`
	DatabaseDir   string `long:"DatabaseDir" description:"Databases folder" default:"res/"`
	HTTPDir       string `long:"HTTPDir" description:"Web page folder to be served folder" default:"http/"`
	LoggerOptions LoggerOptions
}

func (w WebCommand) String() (res string) {
	res += fmt.Sprintf(" Host                   : %s\n", w.Host)
	res += fmt.Sprintf(" Port                   : %d\n", w.Port)
	res += fmt.Sprintf(" DatabaseDir            : %s\n", w.DatabaseDir)
	res += fmt.Sprintf(" HTTPDir                : %s\n", w.HTTPDir)
	return res
}

var webCommand WebCommand

// Execute runs the processing of the commandline parameters
func (w *WebCommand) Execute(args []string) error {
	w.LoggerOptions.Process()

	log.Info("Web")
	log.Info("\n", w)

	fi, err := os.Stat(w.DatabaseDir)
	if err != nil {
		log.Fatal(err)
	}

	if !fi.Mode().IsDir() {
		log.Fatal("input folder ", w.DatabaseDir, " is not a folder")
	}

	web.NewWeb(w.DatabaseDir, w.HTTPDir, w.Host, w.Port)

	return nil
}

func init() {
	parser.AddCommand(
		"web",
		"Start web interface",
		"Start web interface",
		&webCommand,
	)
}
