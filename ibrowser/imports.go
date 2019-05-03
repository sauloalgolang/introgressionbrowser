package ibrowser

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

//
// tools

// Min64 returns the minimal value between two values
var Min64 = tools.Min64

// Max64 returns the maximal value between two values
var Max64 = tools.Max64

// SliceIndex finds the first index of a given value in a slice
var SliceIndex = tools.SliceIndex

//
// save

// NewSaverCompressed creates a new compressor for saving files
var NewSaverCompressed = save.NewSaverCompressed

// NewMultiArrayFile creates a new binary dumper
var NewMultiArrayFile = save.NewMultiArrayFile

// MultiArrayFile defines the file descriptor for binary files
type MultiArrayFile = save.MultiArrayFile

// Parameters holds the program parameters
type Parameters = interfaces.Parameters

// type DistanceRow16 = imports.DistanceRow16
// type DistanceRow32 = imports.DistanceRow32
// type DistanceRow64 = imports.DistanceRow64

// type IBDistanceTable = imports.IBDistanceTable
// type IBDistanceMatrix = imports.IBDistanceMatrix

// var NewDistanceMatrix = imports.NewDistanceMatrix

//
// Types
//

// VCFSamples holds the sample names as read in the VCF file
type VCFSamples = vcf.VCFSamples

// VCFRegister holds a single SNP position
type VCFRegister = vcf.VCFRegister

// VCFDistanceMatrix holds the distance matrix used to calculate distance between
// SNP calls and holds the summary of all distances
type VCFDistanceMatrix = vcf.DistanceMatrix

// NamePosPair holds a name/position pair to keep chromosomes in order
type NamePosPair struct {
	Name string
	Pos  int
}

// NamePosPairList holds all chromosomes and their order
type NamePosPairList []NamePosPair

func (s NamePosPairList) Len() int {
	return len(s)
}

func (s NamePosPairList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s NamePosPairList) Less(i, j int) bool {
	return s[i].Pos < s[j].Pos
}
