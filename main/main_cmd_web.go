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

type WebCommand struct {
	Host      string `long:"host" description:"Hostname" default:"127.0.0.1"`
	Port      int    `long:"port" description:"Port" default:"8000"`
	Dir       string `long:"dir" description:"Directory to be served" default:"res/"`
	Verbose   []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	verbosity int
}

func (w WebCommand) String() (res string) {
	res += fmt.Sprintf(" Host                   : %s\n", w.Host)
	res += fmt.Sprintf(" Port                   : %d\n", w.Port)
	res += fmt.Sprintf(" Dir                    : %s\n", w.Dir)
	return res
}

var webCommand WebCommand

func (x *WebCommand) Execute(args []string) error {
	x.verbosity = len(x.Verbose)

	var verbosityLevel log.Level
	switch x.verbosity {
	case 0:
		log.Println("LOG LEVEL ", x.verbosity, "Info")
		verbosityLevel = log.InfoLevel
	case 1:
		log.Println("LOG LEVEL ", x.verbosity, "Debug")
		verbosityLevel = log.DebugLevel
	case 2:
		log.Println("LOG LEVEL ", x.verbosity, "Trace")
		verbosityLevel = log.TraceLevel
	default:
		if x.verbosity > 2 {
			log.Println("LOG LEVEL ", x.verbosity, "Trace")
			verbosityLevel = log.TraceLevel
		} else {
			log.Println("LOG LEVEL ", x.verbosity, "Warn")
			verbosityLevel = log.WarnLevel
			// log.SetLevel(log.ErrorLevel)
			// log.SetLevel(log.FatalLevel)
			// log.SetLevel(log.PanicLevel)
		}
	}

	log.SetLevel(verbosityLevel)
	// Only log the warning severity or above.

	log.Info("Web")
	log.Info("\n", x)

	fi, err := os.Stat(x.Dir)
	if err != nil {
		log.Fatal(err)
	}

	if !fi.Mode().IsDir() {
		log.Fatal("input folder ", x.Dir, " is not a folder")
	}

	web.NewWeb(x.Dir, x.Host, x.Port, verbosityLevel)

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
