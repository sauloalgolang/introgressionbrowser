package interfaces

import (
	"sync/atomic"
)

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
