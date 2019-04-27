package matrix

import (
	"math"
)

type DistanceTable = DistanceRow64

// type DistanceMatrix = DistanceMatrix1D64

// var NewDistanceMatrix = NewDistanceMatrix1D64

type DistanceMatrix = DistanceMatrix1Dg

var NewDistanceMatrix = NewDistanceMatrix1Dg

type DistanceMatrix1D_T interface {
	// Exported Methods
	Add(*DistanceMatrix1D_T)
	AddAtomic(*DistanceMatrix1D_T)
	Clean()
	Check() bool
	Set(uint64, uint64, uint64)
	Get(uint64, uint64, uint64) uint64
	GenFilename(string, string, string) (string, string)
	Save(string, string, string)
	Load(string, string, string)
	// Unexported Methods
	ijToK(uint64, uint64) uint64
	kToIJ(uint64) (uint64, uint64)
	saveLoad(bool, string, string, string)
}

//
// Calc
//

// # https://stackoverflow.com/questions/27086195/linear-index-upper-triangular-matrix

func ijToK(dimension uint64, i uint64, j uint64) uint64 {
	dim := float64(dimension)
	fi := float64(i)
	fj := float64(j)

	fk := (dim * (dim - 1) / 2) - (dim-fi)*((dim-fi)-1)/2 + fj - fi - 1

	return uint64(fk)
}

func kToIJ(dimension uint64, k uint64) (uint64, uint64) {
	dim := float64(dimension)
	idx := float64(k)

	fi := dim - 2 - math.Floor(math.Sqrt(-8*idx+4*dim*(dim-1)-7)/2.0-0.5)
	fj := idx + fi + 1 - dim*(dim-1)/2 + (dim-fi)*((dim-fi)-1)/2

	return uint64(fi), uint64(fj)
}
