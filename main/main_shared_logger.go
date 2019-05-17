package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// LoggerOptions holds the shared parameters for the logger options
type LoggerOptions struct {
	Verbose                []bool `short:"v" long:"verbose"      description:"Show verbose debug information"`
	DisableColors          bool   `long:"DisableColors"          description:"Remove colors"`
	DisableTimestamp       bool   `long:"DisableTimestamp"       description:"Remove timestamp from test output"`
	FullTimestamp          bool   `long:"FullTimestamp"          description:"Show full timestamp instead of elapsed time"`
	EnableSorting          bool   `long:"EnableSorting"          description:"Enable sorting"`
	DisableLevelTruncation bool   `long:"DisableLevelTruncation" description:"Disables the truncation of the level text to 4 characters"`
	QuoteEmptyFields       bool   `long:"QuoteEmptyFields"       description:"QuoteEmptyFields will wrap empty fields in quotes if true"`
	ReportCaller           bool   `long:"ReportCaller"           description:"SetReportCaller sets whether the standard logger will include the calling method as a field"`
	verbosity              int
	verbosityLevel         log.Level
	// ElapsetTimeStamp bool `long:"ElapsetTimeStamp"       description:"Show elapsed time instead of full timestamp"`
	// DisableSorting bool         `long:"DisableSorting"         description:"Disable sorting"`
}

func (l LoggerOptions) String() (res string) {
	res += fmt.Sprintf("Logger:\n")
	res += fmt.Sprintf(" Verbose                : %#v\n", l.Verbose)
	res += fmt.Sprintf(" DisableColors          : %#v\n", l.DisableColors)
	res += fmt.Sprintf(" DisableTimestamp       : %#v\n", l.DisableTimestamp)
	res += fmt.Sprintf(" FullTimestamp          : %#v\n", l.FullTimestamp)
	res += fmt.Sprintf(" EnableSorting          : %#v\n", l.EnableSorting)
	res += fmt.Sprintf(" DisableLevelTruncation : %#v\n", l.DisableLevelTruncation)
	res += fmt.Sprintf(" QuoteEmptyFields       : %#v\n", l.QuoteEmptyFields)
	res += fmt.Sprintf(" ReportCaller           : %#v\n", l.ReportCaller)
	res += fmt.Sprintf(" verbosity              : %d\n", l.verbosity)
	// res += fmt.Sprintf(" ElapsetTimeStamp       : %#v\n", l.ElapsetTimeStamp)
	// res += fmt.Sprintf(" DisableSorting         : %#v\n", l.DisableSorting)
	return res
}

// Process sets log level
func (l *LoggerOptions) Process() {
	l.verbosity = len(l.Verbose)

	switch l.verbosity {
	case 0:
		// log.Println("LOG LEVEL ", l.verbosity, "Info")
		l.verbosityLevel = log.InfoLevel
	case 1:
		// log.Println("LOG LEVEL ", l.verbosity, "Debug")
		l.verbosityLevel = log.DebugLevel
	case 2:
		// log.Println("LOG LEVEL ", l.verbosity, "Trace")
		l.verbosityLevel = log.TraceLevel
	default:
		if l.verbosity > 2 {
			// log.Println("LOG LEVEL ", l.verbosity, "Trace")
			l.verbosityLevel = log.TraceLevel
		} else {
			// log.Println("LOG LEVEL ", l.verbosity, "Warn")
			l.verbosityLevel = log.WarnLevel
			// log.SetLevel(log.ErrorLevel)
			// log.SetLevel(log.FatalLevel)
			// log.SetLevel(log.PanicLevel)
		}
	}

	formatter := new(log.TextFormatter)
	formatter.DisableColors = l.DisableColors       // remove colors
	formatter.DisableTimestamp = l.DisableTimestamp // remove timestamp from test output
	formatter.FullTimestamp = l.FullTimestamp       // Enable logging the full timestamp when a TTY is attached instead of just the time passed since beginning of execution.
	formatter.DisableSorting = !l.EnableSorting
	formatter.DisableLevelTruncation = l.DisableLevelTruncation
	formatter.QuoteEmptyFields = l.QuoteEmptyFields
	// formatter.FullTimestamp = !l.ElapsetTimeStamp // Enable logging the full timestamp when a TTY is attached instead of just the time passed since beginning of execution.

	log.SetFormatter(formatter) //default
	log.SetReportCaller(l.ReportCaller)
	log.SetLevel(l.verbosityLevel)
	// Only log the warning severity or above.

	log.Println("LOG LEVEL ", l.verbosity, "Level", l.verbosityLevel)

}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetReportCaller(true)

	// log.SetFormatter(new(log.JSONFormatter))

	formatter := new(log.TextFormatter)
	formatter.DisableColors = false    // remove colors
	formatter.DisableTimestamp = false // remove timestamp from test output
	formatter.FullTimestamp = true     // Enable logging the full timestamp when a TTY is attached instead of just the time passed since beginning of execution.
	formatter.DisableSorting = true
	formatter.DisableLevelTruncation = true
	formatter.QuoteEmptyFields = false

	log.SetFormatter(formatter) //default
	log.SetLevel(log.TraceLevel)
}
