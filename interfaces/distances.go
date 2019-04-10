package interfaces

type DistanceRow []uint64
type DistanceTable []uint64
type DistanceMatrix2D [][]uint64
type DistanceMatrix1D struct {
	ChromosomeName string
	BlockSize      uint64
	BlockPosition  uint64
	BlockNumber    uint64
	Dimension      uint64
	Size           uint64
	Data           DistanceRow
}

type DistanceMatrix = DistanceMatrix1D

var NewDistanceMatrix = NewDistanceMatrix1D
