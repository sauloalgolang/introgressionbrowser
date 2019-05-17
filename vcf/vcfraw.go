package vcf

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"strings"
)

// ProcessVcfRaw process vcf manually
func ProcessVcfRaw(r io.Reader, callBackParameters CallBackParameters, callback RegisterCallBack, chromosomeNames []string) {
	log.Println("Opening file to read chromosome:", chromosomeNames)

	contents := bufio.NewScanner(r)
	cbuffer := make([]byte, 0, bufio.MaxScanTokenSize)

	contents.Buffer(cbuffer, bufio.MaxScanTokenSize*50) // Otherwise long lines crash the scanner.

	SampleNames := make([]string, 0, 100)
	numSampleNames := uint64(0)

	register := RegisterRaw{
		LineNumber:       0,
		Chromosome:       "",
		ChromosomeNumber: 0,
		Position:         0,
	}

	sendOnlyChromosomeNames := len(chromosomeNames) == 1 && chromosomeNames[0] == ""

	gtIndex := -1
	lastChrom := ""
	chromosomeNumber := -1
	chromIndex := -1
	// chromIndexOk := false
	lastChromosomeName := ""
	lineNumber := int64(0)
	registerNumberThread := int64(0)
	registerNumberChrom := int64(0)
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
					// log.Println("columnNames", columnNames)

					SampleNames = columnNames[9:]
					numSampleNames = uint64(len(SampleNames))
					register.TempDistance = NewDistanceMatrix(numSampleNames)
					// log.Println("SampleNames", SampleNames, "chromosomeNames", chromosomeNames)
				}
			}
			continue
		}

		cols := strings.Split(row, "\t")

		if len(cols) < 9 {
			log.Println("less than 9 columns. can't continue")
			os.Exit(1)
		}

		chrom := cols[0]

		if chrom != lastChrom {
			chromosomeNumber++
			registerNumberChrom = 0
			chromIndex, _ = SliceIndex(len(chromosomeNames), func(i int) bool { return chromosomeNames[i] == chrom })
			log.Println("  new chromosome ", chrom, " index ", chromIndex, " in ", chromosomeNames)
		}

		lastChrom = chrom

		if sendOnlyChromosomeNames { // return only chromosome names
			if chrom != lastChromosomeName { // first time to see it
				lastChromosomeName = chrom

				register := RegisterRaw{
					LineNumber:       lineNumber,
					Chromosome:       chrom,
					ChromosomeNumber: chromosomeNumber,
					Position:         0,
					Alt:              nil,
					Samples:          nil,
				}

				callback(&SampleNames, &register)
			}
			continue
		}

		if chromIndex == -1 {
			if foundChromosome { // already found, therefore finished
				log.Println("Finished reading chromosome", chromosomeNames, " now at ", chrom, registerNumberThread, " registers ")
				return
			}
			
			// not found yet, therefore continue
			continue
		} else {
			if !foundChromosome {
				log.Println("Found first chromosome from list:", chromosomeNames, " now at ", chrom, registerNumberThread, " registers ")
				foundChromosome = true
			}
		}

		registerNumberChrom++

		if BreakAtChrom > 0 && registerNumberChrom >= BreakAtChrom {
			// log.Println(" BREAKING ", chrom, " at register ", registerNumberChrom)
			continue
		}

		registerNumberThread++

		if BreakAtThread > 0 && registerNumberThread >= BreakAtThread {
			log.Println(" BREAKING ", chromosomeNames, " at register ", registerNumberThread)
			return
		}

		pos, posErr := strconv.ParseUint(cols[1], 10, 64)
		alt := cols[4]
		altCols := strings.Split(alt, ",")
		info := cols[8]
		infoCols := strings.Split(info, ";")

		if len(altCols) > 1 { // no polymorphic SNPs
			continue
		}

		if gtIndex == -1 || infoCols[gtIndex] != "GT" {
			gtIndex, _ = SliceIndex(len(infoCols), func(i int) bool { return infoCols[i] == "GT" })
		}

		if posErr != nil {
			if callBackParameters.ContinueOnError {
				continue
			} else {
				log.Println(posErr)
				os.Exit(1)
			}
		}

		if gtIndex == -1 {
			if callBackParameters.ContinueOnError {
				continue
			} else {
				log.Println("no genotype info field")
				os.Exit(1)
			}
		}

		samples := cols[9:]
		numSamples := uint64(len(samples))
		samplesGT := make([]RegisterGenotype, numSamples, numSamples)

		if numSamples != numSampleNames {
			if callBackParameters.ContinueOnError {
				continue
			} else {
				log.Println("wrong number of columns: expected ", numSampleNames, " got ", numSamples)
				os.Exit(1)
			}
		}

		for samplePos, sample := range samples {
			sampleCols := strings.Split(sample, ";")
			sampleGT := sampleCols[gtIndex]
			sampleGTVal := make(GenotypeVal, 2, 2)

			if sampleGT[0] == '.' {
				sampleGTVal[0] = -1
				sampleGTVal[1] = -1

			} else {
				if len(sampleGT) == 3 {
					sampleGT0, sampleGT0Err := strconv.Atoi(string(sampleGT[0]))
					sampleGT1, sampleGT1Err := strconv.Atoi(string(sampleGT[2]))

					if sampleGT0Err != nil {
						if callBackParameters.ContinueOnError {
							continue
						} else {
							log.Println(sampleGT0Err)
							os.Exit(1)
						}
					}

					if sampleGT1Err != nil {
						if callBackParameters.ContinueOnError {
							continue
						} else {
							log.Println(sampleGT1Err)
							os.Exit(1)
						}
					}

					sampleGTVal[0] = sampleGT0
					sampleGTVal[1] = sampleGT1
				}
			}
			samplesGT[samplePos].Genotype = sampleGTVal
		}

		//  0          1        2 3 4 5      6
		// [SL2.50ch01 73633505 . G A 120.91 .
		//
		// 7
		// AC1=2;AC=60;AF1=0.5086;AN=70;DP4=46,40,297,318;DP=1223;MQ=117;SF=5,7,19,27,39,45,54,55,56,67,78,123,130,134,156,161,164,186,216,223,252,266,269,271,272,276,278,287,288,298,299,307,336,338,358
		// 8  9
		// GT . . . . . 0|1

		if lineNumber%100000 == 0 && lineNumber != 0 {
			log.Println(lineNumber,
				registerNumberThread,
				chrom,
				pos,
				// row,
				// cols,
				// info,
				// samples,
				// samplesGT,
			)
		}

		register.LineNumber = lineNumber
		register.Chromosome = chrom
		register.ChromosomeNumber = chromosomeNumber
		register.Position = pos
		register.Alt = altCols
		register.Samples = samplesGT
		register.Distance = CalculateDistance(numSampleNames, &register)

		callback(&SampleNames, &register)
	}

	if sendOnlyChromosomeNames { // return only chromosome names
		// return final count

		register := RegisterRaw{
			LineNumber: lineNumber,
			Chromosome: "",
			Position:   0,
			Alt:        nil,
			Samples:    nil,
		}

		callback(&SampleNames, &register)
	}

}
