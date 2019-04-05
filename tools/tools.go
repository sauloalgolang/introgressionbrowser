package tools

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

func Min64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	} else if a == b {
		return a
	} else {
		return b
	}
}

func Max64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	} else if a == b {
		return a
	} else {
		return b
	}
}

type DistanceRow []uint64
type DistanceMatrix [][]uint64

func (d *DistanceMatrix) Add(e *DistanceMatrix) {

}

func NewDistanceMatrix(dimention uint64) *DistanceMatrix {
	r := make(DistanceMatrix, dimention, dimention)

	for i := range r {
		r[i] = make(DistanceRow, dimention, dimention)
		for j := range r[i] {
			r[i][j] = uint64(0)
		}
	}

	return &r
}

func CalculateDistance(numSamples uint64, reg *interfaces.VCFRegister) *DistanceMatrix {
	r := NewDistanceMatrix(numSamples)

	return r
}
