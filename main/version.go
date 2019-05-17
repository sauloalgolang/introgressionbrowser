package main

import log "github.com/sirupsen/logrus"

// IBrowserVersion - Current ibrowser version
const IBrowserVersion = 0.5

// PrintVersion - prints the version
func PrintVersion() {
	log.Println("IBrowser Git Commit Author  :", IBrowserGitCommitAuthor)
	log.Println("IBrowser Git Commit Commiter:", IBrowserGitCommitCommiter)
	log.Println("IBrowser Git Commit Hash    :", IBrowserGitCommitHash)
	log.Println("IBrowser Git Commit Notes   :", IBrowserGitCommitNotes)
	log.Println("IBrowser Git Commit Title   :", IBrowserGitCommitTitle)
	log.Println("IBrowser Git Status         :", IBrowserGitStatus)
	log.Println("IBrowser Git Diff           :", IBrowserGitDiff)
	log.Println("IBrowser Go Version         :", IBrowserGoVersion)
	log.Println("IBrowser Version            :", IBrowserVersion)
}
