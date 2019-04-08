package interfaces

import (
	"fmt"
	"math"
	"sync/atomic"
)

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

func NewDistanceMatrix1D(dimension uint64) *DistanceMatrix1D {
	size := dimension * (dimension - 1) / 2

	fmt.Println("NewDistanceMatrix1D :: dimension:", dimension, "size:", size)

	r := DistanceMatrix1D{
		Data:      make(DistanceRow, size, size),
		Size:      size,
		Dimension: dimension,
	}

	r.Clean()

	return &r
}

func (d *DistanceMatrix1D) add(e *DistanceMatrix1D, isAtomic bool) {
	if isAtomic {
		for i := range (*d).Data {
			atomic.AddUint64(&(*d).Data[i], atomic.LoadUint64(&(*e).Data[i]))
		}
	} else {
		for i := range (*d).Data {
			(*d).Data[i] += (*e).Data[i]
		}
	}
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

// # https://stackoverflow.com/questions/27086195/linear-index-upper-triangular-matrix

func (d *DistanceMatrix1D) ijToK(i uint64, j uint64) uint64 {
	dim := float64(d.Dimension)
	fi := float64(i)
	fj := float64(j)

	fk := (dim * (dim - 1) / 2) - (dim-fi)*((dim-fi)-1)/2 + fj - fi - 1

	return uint64(fk)
}

func (d *DistanceMatrix1D) kToIJ(k uint64) (uint64, uint64) {
	dim := float64(d.Dimension)
	idx := float64(k)

	fi := dim - 2 - math.Floor(math.Sqrt(-8*idx+4*dim*(dim-1)-7)/2.0-0.5)
	fj := idx + fi + 1 - dim*(dim-1)/2 + (dim-fi)*((dim-fi)-1)/2

	return uint64(fi), uint64(fj)
}

func (d *DistanceMatrix1D) Set(p1 uint64, p2 uint64, val uint64) {
	(*d).Data[d.ijToK(p1, p2)] += val
}

func (d *DistanceMatrix1D) Get(p1 uint64, p2 uint64, dim uint64) uint64 {
	return (*d).Data[d.ijToK(p1, p2)]
}
