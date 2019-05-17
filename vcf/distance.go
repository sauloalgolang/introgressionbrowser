package vcf

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// Genotype holds genotipical information
type Genotype struct {
	Position  uint64
	Gt        *GenotypeVal
	Lgt       int
	IsDiploid bool
}

// DistanceTableValuesHetLow alias to a distance table where heterozygous
// have lower values than homozygous
var DistanceTableValuesHetLow = DistanceTable{
	3, 1, 1, 0, //  0  1  2  3
	1, 2, 2, 1, //  4  5  6  7
	1, 2, 2, 1, //  8  9 10 11
	0, 1, 1, 3, // 12 13 14 15
	//      | AA AB BA BB
	//      |  0  1  2  3
	// -----|------------
	// AA 0 |  3  1  1  0
	// AB 1 |  1  2  2  1
	// BA 2 |  1  2  2  1
	// BB 3 |  0  1  1  3
	//-------------------
}

// DistanceTableValuesHetEqual alias to a distance table where heterozygous
// have the same value as homozygous
var DistanceTableValuesHetEqual = DistanceTable{
	2, 1, 1, 0, //  0  1  2  3
	1, 2, 2, 1, //  4  5  6  7
	1, 2, 2, 1, //  8  9 10 11
	0, 1, 1, 2, // 12 13 14 15
	//      | AA AB BA BB
	//      |  0  1  2  3
	// -----|------------
	// AA 0 |  3  1  1  0
	// AB 1 |  1  2  2  1
	// BA 2 |  1  2  2  1
	// BB 3 |  0  1  1  2
	//-------------------
}

// DistanceTableValues alias to the default distance table
var DistanceTableValues = DistanceTableValuesHetLow

// DistanceTable alias to the default distance table
type DistanceTable = []uint64

// DistanceRow alias to the default ditance row in a distance matrix
type DistanceRow = []uint64

// DistanceMatrix distance matrix type
type DistanceMatrix []DistanceRow

// NewDistanceMatrix creates a new distance matrix
func NewDistanceMatrix(numSampleNames uint64) *DistanceMatrix {
	res := make(DistanceMatrix, numSampleNames, numSampleNames)

	for i := uint64(0); i < numSampleNames; i++ {
		res[i] = make(DistanceRow, numSampleNames, numSampleNames)
	}

	return &res
}

// Clean clears the distance matrix, zeroing it
func (d *DistanceMatrix) Clean() {
	j := uint64(0)
	le := uint64(len(*d))
	for i := uint64(0); i < le; i++ {
		for j = i; j < le; j++ {
			(*d)[i][j] = 0
		}
	}
}

// Set sets a specific cell value
func (d *DistanceMatrix) Set(p1 uint64, p2 uint64, val uint64) {
	if p1 > p2 {
		p1, p2 = p2, p1
	}
	(*d)[p1][p2] = val
}

// CalculateDistanceDiploid calculates the distance between two diploid snp calls
func CalculateDistanceDiploid(a *GenotypeVal, b *GenotypeVal) uint64 {
	a0 := (*a)[0]
	a1 := (*a)[1]
	b0 := (*b)[0]
	b1 := (*b)[1]

	i := a0*8 + a1*4 + b0*2 + b1*1

	d := DistanceTableValues[i]

	// log.Println(a0, a1, a0*2+a1*1, b0, b1, b0*2+b1*1, i, d)

	return d
}

// GetValids returns a list of all valid snp calls
func GetValids(samples SamplesGenotype) (valids []Genotype, numValids int) {
	numSamples := uint64(len(samples))
	numValids = 0
	valids = make([]Genotype, numSamples, numSamples)

	for samplePos := uint64(0); samplePos < numSamples; samplePos++ {
		sample := samples[samplePos]
		gt := &sample.Genotype
		lgt := len(*gt)

		if lgt == 0 { // wrong.
			log.Print(" samplePos ", samplePos, " GT ", gt, " ", "WRONG 0")
			os.Exit(1)
		} else if lgt == 1 { // maybe no call
			if (*gt)[0] == -1 { // is no call
				// log.Print(" 1 samplePos ", samplePos, " GT ", gt, " ", "NC")
				continue
			} else {
				log.Println(" samplePos ", samplePos, " GT ", gt, " ", "WRONG NOT -1")
				os.Exit(1)
			}
		} else if lgt == 2 { // alts
			// log.Println(" samplePos ", samplePos, " GT ", gt, " ", "DIPLOID")
			if (*gt)[0] == -1 {
				continue
			} else {
				valids[numValids] = Genotype{samplePos, gt, lgt, true}
				numValids++
			}
		} else { // weird
			// log.Println(" samplePos ", samplePos, " GT ", gt, " ", "POLYPLOYD")
			valids[numValids] = Genotype{samplePos, gt, lgt, false}
			numValids++
		}
	}

	return valids, numValids
}

// CalculateDistance calculates the distance between two snp calls
func CalculateDistance(numSamples uint64, reg *Register) *DistanceMatrix {
	reg.TempDistance.Clean()

	valids, numValids := GetValids(reg.Samples)

	// log.Println("valids", numValids, valids, numSamples)

	for validPos1 := 0; validPos1 < numValids; validPos1++ {
		valid1 := valids[validPos1]
		samplePos1 := valid1.Position
		gt1 := valid1.Gt
		isDiploid1 := valid1.IsDiploid
		// lgt1 := valid1.Lgt

		for validPos2 := validPos1 + 1; validPos2 < numValids; validPos2++ {
			valid2 := valids[validPos2]
			samplePos2 := valid2.Position
			gt2 := valid2.Gt
			isDiploid2 := valid2.IsDiploid
			// lgt2 := valid2.Lgt

			if isDiploid1 && isDiploid2 {
				// log.Print("    BOTH DIPLOYD ")
				dist := CalculateDistanceDiploid(gt1, gt2)
				reg.TempDistance.Set(samplePos1, samplePos2, dist)
				// log.Println(gt1, " ", gt2, " ", dist)
			}
		}
	}

	return reg.TempDistance
}
