package tools

import (
	"fmt"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

// type DistanceRow = interfaces.DistanceRow
// type DistanceMatrix = interfaces.DistanceMatrix
// type DistanceTable = interfaces.DistanceTable

type GT struct {
	Position  uint64
	Gt        *interfaces.VCFGTVal
	Lgt       int
	IsDiploid bool
}

// var TempDistanceMatrix DistanceMatrix

var DistanceTableValues = interfaces.DistanceTable{
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

func Min64(a uint64, b uint64) uint64 {
	if a < b {
		return a
	} else if a == b {
		return a
	} else {
		return b
	}
}

func Max64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	} else if a == b {
		return a
	} else {
		return b
	}
}

func CalculateDistanceDiploid(a *interfaces.VCFGTVal, b *interfaces.VCFGTVal) uint64 {
	// fmt.Println("DistanceTableValues", DistanceTableValues)

	a0 := (*a)[0]
	a1 := (*a)[1]
	b0 := (*b)[0]
	b1 := (*b)[1]

	i := a0*8 + a1*4 + b0*2 + b1*0

	// fmt.Print(a0, a1, b0, b1, x, y)

	d := DistanceTableValues[i]

	// fmt.Println(d, DistanceTableValues)

	return d

	// if a0 == a1 { // a homo - AA
	// 	if b0 == b1 { // b homo - AA
	// 		if a0 == b0 {
	// 			return 3 // equal homo - AA x AA
	// 		} else {
	// 			return 0 // diff homo - AA x BB
	// 		}
	// 	} else { // b hete - AB
	// 		return 2 // homo hete - AA x AB
	// 	}
	// } else { // a hete - AB
	// 	if b0 == b1 { // b homo - BB
	// 		return 2 // hete homo - AB x AA
	// 	} else { // b hete - AB
	// 		return 1 // both hete - AB x AB
	// 	}
	// }

	// return uint64(0)
}

func GetValids(samples interfaces.VCFSamplesGT) (valids []GT, numValids int) {
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

func CalculateDistance(numSamples uint64, reg *interfaces.VCFRegister) *interfaces.DistanceMatrix {
	// if uint64(len(TempDistanceMatrix)) != numSamples {
	// 	fmt.Println("CalculateDistance NewDistanceMatrix")
	// 	TempDistanceMatrix = *NewDistanceMatrix(numSamples)
	// } else {
	// 	TempDistanceMatrix.Clean()
	// }
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
				// TempDistanceMatrix[samplePos2][samplePos1] += dist
				// fmt.Println(gt1, " ", gt2, " ", dist)
			}
		}
	}

	return reg.TempDistance
}
