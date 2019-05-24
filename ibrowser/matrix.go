package ibrowser

// DistanceType16 1d array of uint16
type DistanceType16 = uint16

// DistanceType32 1d array of uint32
type DistanceType32 = uint32

// DistanceType64 1d array of uint64
type DistanceType64 = uint64

// DistanceRow16 1d array of uint16
type DistanceRow16 = []DistanceType16

// DistanceRow32 1d array of uint32
type DistanceRow32 = []DistanceType32

// DistanceRow64 1d array of uint64
type DistanceRow64 = []DistanceType64

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
//
//
//
//

// TriangularMatrix Alias to triangular matrix
type TriangularMatrix = StrictlyUpperTriangularMatrix

// type TriangularMatrix = LowerTriangle
// type TriangularMatrix = UpperTriangle

// NewTriangularMatrix creates a new instance of the triangular matrix
var NewTriangularMatrix = NewStrictlyUpperTriangularMatrix

// var NewTriangularMatrix = NewLowerTriangle
// var NewTriangularMatrix = NewUpperTriangle
