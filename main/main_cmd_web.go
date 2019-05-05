package main

import (
	"fmt"
	// "log"
	log "github.com/sirupsen/logrus"
	"os"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/web"
)

// WebCommand holds the parameter for the commandline web command
type WebCommand struct {
	Host        string `long:"host" description:"Hostname" default:"127.0.0.1"`
	Port        int    `long:"port" description:"Port" default:"8000"`
	DatabaseDir string `long:"DatabaseDir" description:"Databases folder" default:"res/"`
	HTTPDir     string `long:"HTTPDir" description:"Web page folder to be served folder" default:"http/"`
	Verbose     []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	verbosity   int
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
	w.verbosity = len(w.Verbose)

	var verbosityLevel log.Level
	switch w.verbosity {
	case 0:
		log.Println("LOG LEVEL ", w.verbosity, "Info")
		verbosityLevel = log.InfoLevel
	case 1:
		log.Println("LOG LEVEL ", w.verbosity, "Debug")
		verbosityLevel = log.DebugLevel
	case 2:
		log.Println("LOG LEVEL ", w.verbosity, "Trace")
		verbosityLevel = log.TraceLevel
	default:
		if w.verbosity > 2 {
			log.Println("LOG LEVEL ", w.verbosity, "Trace")
			verbosityLevel = log.TraceLevel
		} else {
			log.Println("LOG LEVEL ", w.verbosity, "Warn")
			verbosityLevel = log.WarnLevel
			// log.SetLevel(log.ErrorLevel)
			// log.SetLevel(log.FatalLevel)
			// log.SetLevel(log.PanicLevel)
		}
	}

	log.SetLevel(verbosityLevel)
	// Only log the warning severity or above.

	log.Info("Web")
	log.Info("\n", w)

	fi, err := os.Stat(w.DatabaseDir)
	if err != nil {
		log.Fatal(err)
	}

	if !fi.Mode().IsDir() {
		log.Fatal("input folder ", w.DatabaseDir, " is not a folder")
	}

	web.NewWeb(w.DatabaseDir, w.HTTPDir, w.Host, w.Port, verbosityLevel)

	return nil
}

func init() {
	parser.AddCommand(
		"web",
		"Start web interface",
		"Start web interface",
		&webCommand,
	)

	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}
