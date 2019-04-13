package interfaces

import (
	"math"
)

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
