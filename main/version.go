package main

import "fmt"

// IBrowserVersion - Current ibrowser version
const IBrowserVersion = 0.5

// PrintVersion - prints the version
func PrintVersion() {
	fmt.Println("IBrowser Git Commit Author  :", IBrowserGitCommitAuthor)
	fmt.Println("IBrowser Git Commit Commiter:", IBrowserGitCommitCommiter)
	fmt.Println("IBrowser Git Commit Hash    :", IBrowserGitCommitHash)
	fmt.Println("IBrowser Git Commit Notes   :", IBrowserGitCommitNotes)
	fmt.Println("IBrowser Git Commit Title   :", IBrowserGitCommitTitle)
	fmt.Println("IBrowser Git Status         :", IBrowserGitStatus)
	fmt.Println("IBrowser Git Diff           :", IBrowserGitDiff)
	fmt.Println("IBrowser Go Version         :", IBrowserGoVersion)
	fmt.Println("IBrowser Version            :", IBrowserVersion)
}