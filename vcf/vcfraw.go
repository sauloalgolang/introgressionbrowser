package vcf

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
)

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func ProcessVcfRaw(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool, chromosomeNames []string) {
	fmt.Println("Opening file to read chromosome:", chromosomeNames)

	contents := bufio.NewScanner(r)
	cbuffer := make([]byte, 0, bufio.MaxScanTokenSize)

	contents.Buffer(cbuffer, bufio.MaxScanTokenSize*50) // Otherwise long lines crash the scanner.

	SampleNames := make([]string, 0, 100)
	numSampleNames := uint64(0)

	register := interfaces.VCFRegisterRaw{
		LineNumber: 0,
		Chromosome: "",
		Position:   0,
		// TempDistance: interfaces.NewDistanceMatrix(0),
	}

	sendOnlyChromosomeNames := len(chromosomeNames) == 1 && chromosomeNames[0] == ""

	gtIndex := -1
	lastChrom := ""
	chromIndex := -1
	lastChromosomeName := ""
	lineNumber := int64(0)
	registerNumber := int64(0)
	foundChromosome := false

	for contents.Scan() {
		lineNumber++

		row := contents.Text()
		rowLen := len(row)

		if rowLen == 0 {
			continue
		}

		if row[0] == '#' {
			if rowLen > 1 {
				if row[1] == '#' {

				} else {
					columnNames := strings.Split(row, "\t")
					// fmt.Println("columnNames", columnNames)

					SampleNames = columnNames[9:]
					numSampleNames = uint64(len(SampleNames))
					register.TempDistance = interfaces.NewDistanceMatrix("_tmp_"+strings.Join(chromosomeNames, "_"), 0, 0, 0, numSampleNames)
					// fmt.Println("SampleNames", SampleNames, "chromosomeNames", chromosomeNames)
				}
			}
			continue
		}

		cols := strings.Split(row, "\t")

		if len(cols) < 9 {
			fmt.Println("less than 9 columns. can't continue")
			os.Exit(1)
		}

		chrom := cols[0]

		if chrom != lastChrom {
			chromIndex = SliceIndex(len(chromosomeNames), func(i int) bool { return chromosomeNames[i] == chrom })
			fmt.Println("new chromosome ", chrom, " index ", chromIndex, " in ", chromosomeNames)
		}

		lastChrom = chrom

		if sendOnlyChromosomeNames { // return only chromosome names
			if chrom != lastChromosomeName { // first time to see it
				lastChromosomeName = chrom

				register := interfaces.VCFRegisterRaw{
					LineNumber: lineNumber,
					Chromosome: chrom,
					Position:   0,
					Alt:        nil,
					Samples:    nil,
				}

				callback(&SampleNames, &register)
			}
			continue
		} else {
			if chromIndex == -1 {
				if foundChromosome { // already found, therefore finished
					fmt.Println("Finished reading chromosome", chromosomeNames, " now at ", chrom, registerNumber, " registers ")
					return
				} else { // not found yet, therefore continue
					continue
				}
				// } else {
				// 	if !foundChromosome { // first time found. let system know
				// 		foundChromosome = true
				// 	}
				// 	fmt.Println("Found chromosome", chrom, " register ", registerNumber)
			}
		}

		registerNumber++

		if BREAKAT > 0 && registerNumber >= BREAKAT {
			fmt.Println(" BREAKING ", chromosomeNames, " at register ", registerNumber)
			return
		}

		pos, pos_err := strconv.ParseUint(cols[1], 10, 64)
		alt := cols[4]
		altCols := strings.Split(alt, ",")
		info := cols[8]
		infoCols := strings.Split(info, ";")

		if len(altCols) > 1 { // no polymorphic SNPs
			continue
		}

		if gtIndex == -1 || infoCols[gtIndex] != "GT" {
			gtIndex = SliceIndex(len(infoCols), func(i int) bool { return infoCols[i] == "GT" })
		}

		if pos_err != nil {
			if continueOnError {
				continue
			} else {
				fmt.Println(pos_err)
				os.Exit(1)
			}
		}

		if gtIndex == -1 {
			if continueOnError {
				continue
			} else {
				fmt.Println("no genotype info field")
				os.Exit(1)
			}
		}

		samples := cols[9:]
		numSamples := uint64(len(samples))
		samplesGT := make([]interfaces.VCFGT, numSamples, numSamples)

		if numSamples != numSampleNames {
			if continueOnError {
				continue
			} else {
				fmt.Println("wrong number of columns: expected ", numSampleNames, " got ", numSamples)
				os.Exit(1)
			}
		}

		for samplePos, sample := range samples {
			sampleCols := strings.Split(sample, ";")
			sampleGT := sampleCols[gtIndex]
			sampleGTVal := make(interfaces.VCFGTVal, 2, 2)

			if sampleGT[0] == '.' {
				sampleGTVal[0] = -1
				sampleGTVal[1] = -1

			} else {
				if len(sampleGT) == 3 {
					sampleGT0, sampleGT0_err := strconv.Atoi(string(sampleGT[0]))
					sampleGT1, sampleGT1_err := strconv.Atoi(string(sampleGT[2]))

					if sampleGT0_err != nil {
						if continueOnError {
							continue
						} else {
							fmt.Println(sampleGT0_err)
							os.Exit(1)
						}
					}

					if sampleGT1_err != nil {
						if continueOnError {
							continue
						} else {
							fmt.Println(sampleGT1_err)
							os.Exit(1)
						}
					}

					sampleGTVal[0] = sampleGT0
					sampleGTVal[1] = sampleGT1
				}
			}
			samplesGT[samplePos].GT = sampleGTVal
		}

		//  0          1        2 3 4 5      6
		// [SL2.50ch01 73633505 . G A 120.91 .
		//
		// 7
		// AC1=2;AC=60;AF1=0.5086;AN=70;DP4=46,40,297,318;DP=1223;MQ=117;SF=5,7,19,27,39,45,54,55,56,67,78,123,130,134,156,161,164,186,216,223,252,266,269,271,272,276,278,287,288,298,299,307,336,338,358
		// 8  9
		// GT . . . . . 0|1

		if lineNumber%100000 == 0 && lineNumber != 0 {
			fmt.Println(lineNumber,
				registerNumber,
				// row,
				// cols,
				chrom,
				pos,
				// info,
				// samples,
				// samplesGT,
			)
		}

		register.LineNumber = lineNumber
		register.Chromosome = chrom
		register.Position = pos
		register.Alt = altCols
		register.Samples = samplesGT
		register.Distance = CalculateDistance(numSampleNames, &register)

		callback(&SampleNames, &register)
	}

	if sendOnlyChromosomeNames { // return only chromosome names
		// return final count

		register := interfaces.VCFRegisterRaw{
			LineNumber: lineNumber,
			Chromosome: "",
			Position:   0,
			Alt:        nil,
			Samples:    nil,
		}

		callback(&SampleNames, &register)
	}

}
