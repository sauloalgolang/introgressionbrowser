package ibrowser

import (
	"math"
)

// DistanceRow16 1d array of uint16
type DistanceRow16 = []uint16

// DistanceRow32 1d array of uint32
type DistanceRow32 = []uint32

// DistanceRow64 1d array of uint64
type DistanceRow64 = []uint64

// IBDistanceTable is the default distance table for ibrowser
type IBDistanceTable = DistanceRow64

// IBDistanceMatrix is the default distance matrix for ibrowser
type IBDistanceMatrix = DistanceMatrix1Dg

// NewDistanceMatrix creates a new instance of the default distance matrix
var NewDistanceMatrix = NewDistanceMatrix1Dg

// DistanceMatrix1DType is the interface for a distance matrix
type DistanceMatrix1DType interface {
	// Exported Methods
	Add(*DistanceMatrix1DType)
	AddVcfMatrix(*VCFDistanceMatrix)
	AddAtomic(*DistanceMatrix1DType)
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

	if fi > fj {
		fi, fj = fj, fi
	}

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
