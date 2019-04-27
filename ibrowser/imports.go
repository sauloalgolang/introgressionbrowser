package ibrowser

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/matrix"
	"github.com/sauloalgolang/introgressionbrowser/save"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

type DistanceMatrix = interfaces.DistanceMatrix
type DistanceTable = interfaces.DistanceTable
type MultiArrayFile = save.MultiArrayFile
type IBDistanceMatrix = DistanceMatrix
type IBDistanceTable = DistanceTable

var NewDistanceMatrix = matrix.NewDistanceMatrix

var Min64 = tools.Min64
var Max64 = tools.Max64
var NewSaverCompressed = save.NewSaverCompressed
var NewMultiArrayFile = save.NewMultiArrayFile

//
// Types
//

// type VCFReaderType = interfaces.VCFReaderType
type VCFSamples = interfaces.VCFSamples
type VCFRegister = interfaces.VCFRegister
type Parameters = interfaces.Parameters

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
