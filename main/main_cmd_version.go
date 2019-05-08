package main

import (
	"os"
)

// VersionCommand holds the parameter for the commandline version command
type VersionCommand struct {
}

// VersionOptions holds the parameter for the commandline version - options
type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version and exit" command:"version"`
}

// VersionCommand instance
var versionCommand VersionCommand

// Execute runs the processing of the commandline parameters
func (x *VersionCommand) Execute(args []string) error {
	PrintVersion()
	os.Exit(0)
	return nil
}

func init() {
	parser.AddCommand("version",
		"Print version and exit",
		"Print version and exit",
		&versionCommand)
}
