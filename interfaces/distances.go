package interfaces

import (
	"sync/atomic"
)

type DistanceRow []uint64
type DistanceTable []uint64
type DistanceMatrix2D [][]uint64
type DistanceMatrix1D struct {
	Data      DistanceRow
	Dimension uint64
}

type DistanceMatrix = DistanceMatrix2D

var NewDistanceMatrix = NewDistanceMatrix2D

//
//
// Matrix 2D
//
//

func NewDistanceMatrix2D(dimension uint64) *DistanceMatrix2D {
	r := make(DistanceMatrix2D, dimension, dimension)

	for i := range r {
		r[i] = make(DistanceRow, dimension, dimension)
		ri := &r[i]
		for j := range *ri {
			(*ri)[j] = uint64(0)
		}
	}

	return &r
}

func (d *DistanceMatrix2D) add(e *DistanceMatrix2D, isAtomic bool) {
	for i := range *d {
		di := &(*d)[i]
		ei := &(*e)[i]

		for j := i + 1; j < len(*d); j++ {
			if isAtomic {
				atomic.AddUint64(&(*di)[j], atomic.LoadUint64(&(*ei)[j]))

			} else {
				(*di)[j] += (*ei)[j]
			}
		}
	}
}

func (d *DistanceMatrix2D) Add(e *DistanceMatrix2D) {
	d.add(e, false)
}
func (d *DistanceMatrix2D) AddAtomic(e *DistanceMatrix2D) {
	d.add(e, true)
}

func (d *DistanceMatrix2D) Clean() {
	for i := range *d {
		ti := &(*d)[i]
		for j := i + 1; j < len(*d); j++ {
			(*ti)[j] = uint64(0)
		}
	}
}

func (d *DistanceMatrix2D) Set(p1 uint64, p2 uint64, val uint64) {
	(*d)[p1][p2] += val
}

func (d *DistanceMatrix2D) Get(p1 uint64, p2 uint64) uint64 {
	return (*d)[p1][p2]
}

//
//
// Matrix 1D
//
//

// https://stackoverflow.com/questions/3187957/how-to-store-a-symmetric-matrix
//
// Here is a good method to store a symmetric matrix, it requires only N(N+1)/2 memory:
//
// int fromMatrixToVector(int i, int j, int N)
// {
//    if (i <= j)
//       return i * N - (i - 1) * i / 2 + j - i;
//    else
//       return j * N - (j - 1) * j / 2 + i - j;
// }
// For some triangular matrix
//
// 0 1 2 3
//   4 5 6
//     7 8
//       9
// 1D representation (stored in std::vector, for example) looks like as follows:
//
// [0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
// And call fromMatrixToVector(1, 2, 4) returns 5, so the matrix data is vector[5] -> 5.
//
// The first expression can be rewritten as (2*N - i - 1)*i/2 + j

func fromMatrixToVector(i uint64, j uint64, N uint64) uint64 {
	if i <= j {
		return i*N - (i-1)*i/2 + j - i
	} else {
		return j*N - (j-1)*j/2 + i - j
	}
}

func NewDistanceMatrix1D(dimension uint64) *DistanceMatrix1D {
	r := DistanceMatrix1D{
		Data:      make(DistanceRow, dimension, dimension),
		Dimension: dimension,
	}

	for i := range r.Data {
		r.Data[i] = uint64(0)
	}

	return &r
}

func (d *DistanceMatrix1D) add(e *DistanceMatrix1D, isAtomic bool) {
	// TODO

	// for i := range *d {
	// 	di := &(*d)[i]
	// 	ei := &(*e)[i]

	// 	for j := i + 1; j < len(*d); j++ {
	// 		if isAtomic {
	// 			atomic.AddUint64(&(*di)[j], atomic.LoadUint64(&(*ei)[j]))

	// 		} else {
	// 			(*di)[j] += (*ei)[j]
	// 		}
	// 	}
	// }
}

func (d *DistanceMatrix1D) Add(e *DistanceMatrix1D) {
	d.add(e, false)
}

func (d *DistanceMatrix1D) AddAtomic(e *DistanceMatrix1D) {
	d.add(e, true)
}

func (d *DistanceMatrix1D) Clean() {
	for i := range (*d).Data {
		(*d).Data[i] = uint64(0)
	}
}

func (d *DistanceMatrix1D) Set(p1 uint64, p2 uint64, val uint64) {
	(*d).Data[fromMatrixToVector(p1, p2, d.Dimension)] += val
}

func (d *DistanceMatrix1D) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	return (*d).Data[fromMatrixToVector(p1, p2, (*d).Dimension)]
}
