package ibrowser

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
	"github.com/sauloalgolang/introgressionbrowser/vcf"
)

// tools
var Min64 = tools.Min64
var Max64 = tools.Max64
var SliceIndex = tools.SliceIndex

// save
var NewSaverCompressed = save.NewSaverCompressed
var NewMultiArrayFile = save.NewMultiArrayFile

type MultiArrayFile = save.MultiArrayFile

// interfaces
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

type VCFSamples = vcf.VCFSamples
type VCFRegister = vcf.VCFRegister
type VCFDistanceMatrix = vcf.DistanceMatrix

type NamePosPair struct {
	Name string
	Pos  int
}

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
