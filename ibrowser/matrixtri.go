package ibrowser

import (
	"math"
)

//
// Calc
//

// i = n - 2 - floor(sqrt(-8*k + 4*n*(n-1)-7)/2.0 - 0.5)
// j = k + i + 1 - n*(n-1)/2 + (n-i)*((n-i)-1)/2

// https://en.wikipedia.org/wiki/Triangular_matrix
// Strictly triangular matrix - no diagonal
// Regular  triangular matrix - with diagonal

// https://stackoverflow.com/questions/27086195/linear-index-upper-triangular-matrix

// IJ stores i,j coordinates
type IJ struct {
	I uint64
	J uint64
}

var mapcacheP []IJ
var mapcacheIJ [][]uint64
var mapcacheInited bool

// StrictlyUpperTriangularMatrix holds a Strictly Upper Triagular Matrix
type StrictlyUpperTriangularMatrix struct {
	Dimension uint64
	dim       float64
	dv        float64
	dd        float64
	ds        float64
	dp        float64
	size      uint64
}

// NewStrictlyUpperTriangularMatrix instantiate a new StrictlyUpperTriangularMatrix
func NewStrictlyUpperTriangularMatrix(dimension uint64) (m *StrictlyUpperTriangularMatrix) {
	m = &StrictlyUpperTriangularMatrix{
		Dimension: dimension,
	}

	m.Update()

	return
}

// Update precalculates some variables
func (m *StrictlyUpperTriangularMatrix) Update() {
	m.dim = float64(m.Dimension)
	m.dv = (m.dim * (m.dim - 1) / 2.0)
	m.dd = m.dim - 2
	m.ds = 4*m.dim*(m.dim-1) - 7
	m.dp = m.dim * (m.dim - 1) / 2
	m.size = m.Dimension * (m.Dimension - 1) / 2

	if !mapcacheInited {
		mapcacheInited = true
		mapcacheP = make([]IJ, m.size, m.size)
		mapcacheIJ = make([][]uint64, m.Dimension, m.Dimension)

		var i, j, p uint64

		for i = 0; i < m.Dimension; i++ {
			// println("i ", i)
			mapcacheIJ[i] = make([]uint64, m.Dimension, m.Dimension)
			for j = i + 1; j < m.Dimension; j++ {
				// print(" j ", j)
				p = m.IJToPO(i, j)
				// println(" p ", p)
				mapcacheP[p] = IJ{I: i, J: j}
				mapcacheIJ[i][j] = p
			}
		}
	}
	// print("mapcache", mapcache)
}

// IJToPO converts i,j coordinates into serial coordinate
func (m *StrictlyUpperTriangularMatrix) IJToPO(i uint64, j uint64) uint64 {
	fi := float64(i)
	fj := float64(j)

	if fi > fj {
		fi, fj = fj, fi
	}

	// fk := (dim * (dim - 1) / 2) - (dim-fi)*((dim-fi)-1)/2 + fj - fi - 1
	// fp := m.dv - (m.dim-fi)*((m.dim-fi)-1)/2 + fj - fi - 1

	// https://stackoverflow.com/questions/7945722/efficient-way-to-represent-a-lower-upper-triangular-matrix
	// With the diagonal
	// i = (y*(2*n - y + 1))/2 + (x - y - 1)
	// Without the diagonal
	// (y*(2*n - y - 1))/2 + (x - y -1)
	fp := (fi*(2*m.dim-fi-1))/2 + (fj - fi - 1)

	return uint64(fp)
}

// IJToP converts i,j coordinates into serial coordinate
func (m *StrictlyUpperTriangularMatrix) IJToP(i uint64, j uint64) uint64 {
	return mapcacheIJ[i][j]
}

// PToIJO converts serial coordinate to i,j coordinates
func (m *StrictlyUpperTriangularMatrix) PToIJO(p uint64) (uint64, uint64) {
	idx := float64(p)

	// fi := dim - 2 - math.Floor(math.Sqrt(-8*idx+4*dim*(dim-1)-7)/2.0-0.5)
	// fj := idx + fi + 1 - dim*(dim-1)/2 + (dim-fi)*((dim-fi)-1)/2

	fi := m.dd - math.Floor((math.Sqrt(-8*idx+m.ds)/2.0)-0.5)
	fj := idx + fi + 1 - m.dp + (m.dim-fi)*((m.dim-fi)-1)/2

	return uint64(fi), uint64(fj)
}

// PToIJ converts serial coordinate to i,j coordinates
func (m *StrictlyUpperTriangularMatrix) PToIJ(p uint64) (uint64, uint64) {
	ij := mapcacheP[p]

	return ij.I, ij.J
}

// CalculateSize calculate matrix size
func (m *StrictlyUpperTriangularMatrix) CalculateSize() uint64 {
	// size = m.Dimension * (m.Dimension - 1) / 2
	return m.size
}

//
//
//
//
//
// https://codeforces.com/blog/entry/60423?locale=en

// C2 constant to conversion
func C2(n uint64) uint64 {
	return n * (n - 1) / 2
}

//
// Lower triangle
//

// LowerTriangle calculates a lower triangle
type LowerTriangle struct {
	Dimension uint64
	size      uint64
}

// NewLowerTriangle creates a new instance of LowerTriangle
func NewLowerTriangle(Dimension uint64) (m *LowerTriangle) {
	m = &LowerTriangle{
		Dimension: Dimension,
		size:      C2(Dimension + 1),
	}
	return m
}

// IJToP convertes i,j coorditates to 1D coordinates
func (m *LowerTriangle) IJToP(i uint64, j uint64) uint64 {
	if i < j {
		i, j = j, i
	}
	// assert( i >= j and i >= 1 ), this->at( C2( i ) + j - 1 ) = k;
	return C2(i) + j - 1
}

// PToIJ convertes i,j coorditates to 1D coordinates
func (m *LowerTriangle) PToIJ(p uint64) (i uint64, j uint64) {
	panic("not implemented")
	return 0, 0
}

// CalculateSize calculate matrix size
func (m *LowerTriangle) CalculateSize() uint64 {
	// size = m.Dimension * (m.Dimension - 1) / 2
	return m.size
}

//
//
//

// UpperTriangle calculates a upper triangle
type UpperTriangle struct {
	Dimension uint64
	size      uint64
	p         uint64
}

// NewUpperTriangle creates a new instance of UpperTriangle
func NewUpperTriangle(Dimension uint64) (m *UpperTriangle) {
	m = &UpperTriangle{
		Dimension: Dimension,
		size:      C2(Dimension + 1),
		p:         Dimension + 1,
	}
	return
}

// IJToP convertes i,j coorditates to 1D coordinates
func (m *UpperTriangle) IJToP(i uint64, j uint64) (p uint64) {
	// assert( j >= i and i >= 1 );  this->at(  ) = k;
	p = C2(m.p-i) + m.Dimension - j

	// https://gist.github.com/kylebgorman/8064310
	// TriangleDiagonal
	// row *  self.n      - (row - 1) * ((row - 1) + 1) / 2 + col - row
	// TriangleNoDiagonal
	// row * (self.n - 1) - (row - 1) * ((row - 1) + 1) / 2 + col - row - 1

	return
}

// PToIJ convertes i,j coorditates to 1D coordinates
func (m *UpperTriangle) PToIJ(p uint64) (i uint64, j uint64) {
	panic("not implemented")
	return 0, 0

	// https://gist.github.com/PhDP/2358809
	//
	// def getX(i):
	// 	return int(-0.5 + 0.5 * sqrt(1 + 8 * (i - 1))) + 2
	//
	// def getY(i):
	// 	return getX(i) * (3 - getX(i)) / 2 + i - 1
	//
	// def getXY(i):
	// 	return (getX(i + 1) - 1, getY(i + 1) - 1)

	// if p == 0 {
	// 	return 0, 0
	// }
	// i = uint64(math.Floor(-0.5+0.5*math.Sqrt(1+8*(float64(p)-1))) + 2)
	// j = i*(3-i)/2 + p - 1
	// return
}

// CalculateSize calculate matrix size
func (m *UpperTriangle) CalculateSize() uint64 {
	// size = m.Dimension * (m.Dimension - 1) / 2
	return m.size
}
