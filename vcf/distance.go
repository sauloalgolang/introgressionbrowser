package vcf

import (
	"fmt"
	"os"
)

type GT struct {
	Position  uint64
	Gt        *VCFGTVal
	Lgt       int
	IsDiploid bool
}

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

var DistanceTableValues = DistanceTableValuesHetLow

func CalculateDistanceDiploid(a *VCFGTVal, b *VCFGTVal) uint64 {
	a0 := (*a)[0]
	a1 := (*a)[1]
	b0 := (*b)[0]
	b1 := (*b)[1]

	i := a0*8 + a1*4 + b0*2 + b1*1

	d := DistanceTableValues[i]

	// fmt.Println(a0, a1, a0*2+a1*1, b0, b1, b0*2+b1*1, i, d)

	return d
}

func GetValids(samples VCFSamplesGT) (valids []GT, numValids int) {
	numSamples := uint64(len(samples))
	numValids = 0
	valids = make([]GT, numSamples, numSamples)

	for samplePos := uint64(0); samplePos < numSamples; samplePos++ {
		sample := samples[samplePos]
		gt := &sample.GT
		lgt := len(*gt)

		if lgt == 0 { // wrong.
			fmt.Print(" samplePos ", samplePos, " GT ", gt, " ", "WRONG 0")
			os.Exit(1)
		} else if lgt == 1 { // maybe no call
			if (*gt)[0] == -1 { // is no call
				// fmt.Print(" 1 samplePos ", samplePos, " GT ", gt, " ", "NC")
				continue
			} else {
				fmt.Println(" samplePos ", samplePos, " GT ", gt, " ", "WRONG NOT -1")
				os.Exit(1)
			}
		} else if lgt == 2 { // alts
			// fmt.Println(" samplePos ", samplePos, " GT ", gt, " ", "DIPLOID")
			if (*gt)[0] == -1 {
				continue
			} else {
				valids[numValids] = GT{samplePos, gt, lgt, true}
				numValids++
			}
		} else { // weird
			// fmt.Println(" samplePos ", samplePos, " GT ", gt, " ", "POLYPLOYD")
			valids[numValids] = GT{samplePos, gt, lgt, false}
			numValids++
		}
	}

	return valids, numValids
}

func CalculateDistance(numSamples uint64, reg *VCFRegister) *DistanceMatrix {
	reg.TempDistance.Clean()

	valids, numValids := GetValids(reg.Samples)

	// fmt.Println("valids", numValids, valids, numSamples)

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
				// fmt.Print("    BOTH DIPLOYD ")
				dist := CalculateDistanceDiploid(gt1, gt2)
				reg.TempDistance.Set(samplePos1, samplePos2, dist)
				// fmt.Println(gt1, " ", gt2, " ", dist)
			}
		}
	}

	return reg.TempDistance
}
