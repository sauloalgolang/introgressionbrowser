package interfaces

import (
	"sync/atomic"
)

type DistanceRow []uint64
type DistanceTable []uint64
type DistanceMatrix2D [][]uint64
type DistanceMatrix1D []uint64

type DistanceMatrix = DistanceMatrix2D

//
//
// Matrix 2D
//
//

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

func fromMatrixToVector(i int, j int, N int) int {
	if i <= j {
		return i*N - (i-1)*i/2 + j - i
	} else {
		return j*N - (j-1)*j/2 + i - j
	}
}
