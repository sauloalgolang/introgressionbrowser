package interfaces

type DistanceRow []uint64
type DistanceTable []uint64
type DistanceMatrix2D [][]uint64
type DistanceMatrix1D struct {
	Data      DistanceRow
	Dimension uint64
	Size      uint64
}

type DistanceMatrix = DistanceMatrix1D

var NewDistanceMatrix = NewDistanceMatrix1D
