package interfaces

type DistanceTable []uint64

// type DistanceMatrix = DistanceMatrix1D64

// var NewDistanceMatrix = NewDistanceMatrix1D64

type DistanceMatrix = DistanceMatrix1Dg

var NewDistanceMatrix = NewDistanceMatrix1Dg

type DistanceMatrix1D_T interface {
	// Exported Methods
	Add(*DistanceMatrix1D_T)
	AddAtomic(*DistanceMatrix1D_T)
	Clean()
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
