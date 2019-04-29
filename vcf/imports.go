package vcf

import (
	"github.com/remeh/sizedwaitgroup"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/matrix"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
	"github.com/sauloalgolang/introgressionbrowser/tools"
)

var DEBUG bool = false
var ONLYFIRST bool = false
var BREAKAT_THREAD int64 = 0
var BREAKAT_CHROM int64 = 0

var SliceIndex = tools.SliceIndex

type CallBackParameters = interfaces.CallBackParameters

type SizedWaitGroup = sizedwaitgroup.SizedWaitGroup

type VCFRegisterRaw = interfaces.VCFRegisterRaw
type VCFCallBack = interfaces.VCFCallBack
type VCFGT = interfaces.VCFGT
type VCFGTVal = interfaces.VCFGTVal
type VCFSamplesGT = interfaces.VCFSamplesGT
type VCFRegister = interfaces.VCFRegister

type DistanceMatrix = matrix.DistanceMatrix
type DistanceTable = matrix.DistanceTable

var OpenFile = openfile.OpenFile

var NewDistanceMatrix = matrix.NewDistanceMatrix
