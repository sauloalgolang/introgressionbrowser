package main

import (
	"fmt"
	"os"
)

type VersionOptions struct {
	Version bool `short:"v" long:"version" description:"Print version and exit" command:"version"`
}

type VersionCommand struct {
}

var versionCommand VersionCommand

func (x *VersionCommand) Execute(args []string) error {
	fmt.Println("IBROWSER_GIT_COMMIT_AUTHOR  :", IBROWSER_GIT_COMMIT_AUTHOR)
	fmt.Println("IBROWSER_GIT_COMMIT_COMMITER:", IBROWSER_GIT_COMMIT_COMMITER)
	fmt.Println("IBROWSER_GIT_COMMIT_HASH    :", IBROWSER_GIT_COMMIT_HASH)
	fmt.Println("IBROWSER_GIT_COMMIT_NOTES   :", IBROWSER_GIT_COMMIT_NOTES)
	fmt.Println("IBROWSER_GIT_COMMIT_TITLE   :", IBROWSER_GIT_COMMIT_TITLE)
	fmt.Println("IBROWSER_GIT_STATUS         :", IBROWSER_GIT_STATUS)
	fmt.Println("IBROWSER_GIT_DIFF           :", IBROWSER_GIT_DIFF)
	fmt.Println("IBROWSER_GO_VERSION         :", IBROWSER_GO_VERSION)
	fmt.Println("IBROWSER_VERSION            :", IBROWSER_VERSION)
	os.Exit(0)
	return nil
}

func init() {
	parser.AddCommand("version",
		"Print version and exit",
		"Print version and exit",
		&versionCommand)
}
