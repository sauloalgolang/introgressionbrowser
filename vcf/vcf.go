package vcf

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	// "sync"
)

import (
	"github.com/brentp/vcfgo"
	"github.com/remeh/sizedwaitgroup"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
	"github.com/sauloalgolang/introgressionbrowser/openfile"
)

const DEBUG = false
const BREAKAT = 500000

// https://github.com/brentp/vcfgo/blob/master/examples/main.go

//
//
// Chromosome Gather
//
//

type chromosomeNamesType struct {
	Names []string
}

func (cn *chromosomeNamesType) IndexFileName(outPrefix string) (indexFile string) {
	indexFile = outPrefix + ".ibindex"
	return indexFile
}

func (cn *chromosomeNamesType) Save(outPrefix string) {
	d, err := yaml.Marshal(cn)

	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	// fmt.Printf("--- dump:\n%s\n\n", d)
	outfile := cn.IndexFileName(outPrefix)
	fmt.Println(" saving index to ", outfile)
	err = ioutil.WriteFile(outfile, d, 0644)
	fmt.Println(" done")
}

func (cn *chromosomeNamesType) Load(outPrefix string) {
	outfile := cn.IndexFileName(outPrefix)

	data, err := ioutil.ReadFile(outfile)

	if err != nil {
		fmt.Printf("yamlFile. Get err   #%v ", err)
	}

	err = yaml.Unmarshal(data, &cn)

	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
	}
}

func (cn *chromosomeNamesType) Add(chromosomeName string) {
	cn.Names = append(cn.Names, chromosomeName)
}

func GatherChromosomeNames(sourceFile string, isTar bool, isGz bool, continueOnError bool) (chromosomeNames chromosomeNamesType) {
	indexfile := chromosomeNames.IndexFileName(sourceFile)

	if _, err := os.Stat(indexfile); err == nil {
		// path/to/whatever exists
		fmt.Println(" exists")
		chromosomeNames.Load(sourceFile)
		// fmt.Println(chromosomeNames)

	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println(" creating")

		addToNames := func(SampleNames *interfaces.VCFSamples, register *interfaces.VCFRegister) {
			fmt.Println("adding chromosome ", register.Chromosome, chromosomeNames)
			chromosomeNames.Add(register.Chromosome)
		}

		getNames := func(r io.Reader, continueOnError bool) {
			ProcessVcfRaw(r, addToNames, continueOnError, "")
		}

		openfile.OpenFile(sourceFile, isTar, isGz, continueOnError, getNames)

		chromosomeNames.Save(sourceFile)
	}

	return chromosomeNames
}

//
//
// Chrosmosome Callback
//
//

type ChromosomeCallbackRegister struct {
	callBack       interfaces.VCFMaskedReaderChromosomeType
	chromosomeName string
	wg             *sizedwaitgroup.SizedWaitGroup
	// wg             *sync.WaitGroup
}

func (cc *ChromosomeCallbackRegister) ChromosomeCallback(r io.Reader, continueOnError bool) {
	defer cc.wg.Done()
	cc.callBack(r, continueOnError, cc.chromosomeName)
}

//
//
// Open VCF
//
//

func OpenVcfFile(sourceFile string, continueOnError bool, numThreads int, callBack interfaces.VCFMaskedReaderChromosomeType) {
	fmt.Println("OpenVcfFile :: ",
		"sourceFile", sourceFile,
		"continueOnError", continueOnError,
		"numThreads", numThreads)

	isTar := false
	isGz := false

	if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf.tar.gz") {
		fmt.Println(" .tar.gz format")
		isTar = true
		isGz = true
	} else if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf.gz") {
		fmt.Println(" .gz format")
		isTar = false
		isGz = true
	} else if strings.HasSuffix(strings.ToLower(sourceFile), ".vcf") {
		fmt.Println(" .vcf format")
		isTar = false
		isGz = false
	} else {
		fmt.Println("unknown file suffix!")
		os.Exit(1)
	}

	chromosomeNames := GatherChromosomeNames(sourceFile, isTar, isGz, continueOnError)
	for _, chromosomeName := range chromosomeNames.Names {
		fmt.Println("Gathered Chromosome Names :: ", chromosomeName)
	}

	// wg := sync.WaitGroup
	wg := sizedwaitgroup.New(numThreads)
	for _, chromosomeName := range chromosomeNames.Names {
		ccr := ChromosomeCallbackRegister{
			callBack:       callBack,
			chromosomeName: chromosomeName,
			wg:             &wg,
		}

		// wg.Add(1)
		wg.Add()

		go openfile.OpenFile(sourceFile, isTar, isGz, continueOnError, ccr.ChromosomeCallback)
	}

	fmt.Println("Waiting for all chromosomes to complete")
	wg.Wait()
	fmt.Println("All chromosomes completed")
}

//
//
// Process VCF
//
//

func SliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func ProcessVcf(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool, chromosomeName string) {
	bufreader := bufio.NewReader(r)
	ProcessVcfRaw(bufreader, callback, continueOnError, chromosomeName)
	// ProcessVcfVcfGo(bufreader, callback, continueOnError, chromosomeName)
}

func ProcessVcfRaw(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool, chromosomeName string) {
	fmt.Println("Opening file to read chromosome:", chromosomeName)

	contents := bufio.NewScanner(r)
	cbuffer := make([]byte, 0, bufio.MaxScanTokenSize)

	contents.Buffer(cbuffer, bufio.MaxScanTokenSize*50) // Otherwise long lines crash the scanner.

	var SampleNames []string

	gtIndex := -1
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
					// fmt.Println("SampleNames", SampleNames)
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

		if chromosomeName == "" { // return only chromosome names
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
			if chrom != chromosomeName {
				if foundChromosome { // already found, therefore finished
					fmt.Println("Finished reading chromosome", chromosomeName)
					return
				} else { // not found yet, therefore continue
					continue
				}
			} else {
				if !foundChromosome { // first time found. let system know
					fmt.Println("Found chromosome", chromosomeName)
					foundChromosome = true
				}
			}
		}

		registerNumber++

		if BREAKAT > 0 && registerNumber >= BREAKAT {
			return
		}

		pos, pos_err := strconv.ParseUint(cols[1], 10, 64)
		alt := cols[4]
		altCols := strings.Split(alt, ",")
		info := cols[8]
		infoCols := strings.Split(info, ";")

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
		numSamples := len(samples)
		samplesGT := make([]interfaces.VCFGT, numSamples, numSamples)

		if len(samples) != len(SampleNames) {
			if continueOnError {
				continue
			} else {
				fmt.Println("wrong number of columns: expected ", len(SampleNames), " got ", len(samples))
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

		register := interfaces.VCFRegisterRaw{
			LineNumber: lineNumber,
			Chromosome: chrom,
			Position:   pos,
			Alt:        altCols,
			Samples:    samplesGT,
		}

		callback(&SampleNames, &register)
	}
}

func ProcessVcfVcfGo(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool, chromosomeName string) {
	vr, err := vcfgo.NewReader(r, false)

	if err != nil {
		panic(err)
	}

	// fmt.Printf("VR %v\n", vr)

	header := vr.Header
	SampleNames := header.SampleNames // []string
	numSamples := len(SampleNames)

	if DEBUG {
		Infos := header.Infos // map[string]*Info
		// // Id          string
		// // Description string
		// // Number      string // A G R . ''
		// // Type        string // STRING INTEGER FLOAT FLAG CHARACTER UNKONWN
		SampleFormats := header.SampleFormats // map[string]*SampleFormat
		Filters := header.Filters             // map[string]string
		Extras := header.Extras               // map[string]string
		FileFormat := header.FileFormat       // string
		Contigs := header.Contigs             // map[string]map[string]string

		fmt.Println("FileFormat:", FileFormat)

		fmt.Println("SAMPLES")
		for samplePos, sampleName := range SampleNames {
			fmt.Println(samplePos, sampleName)
		}

		fmt.Println("INFO")
		for infoID, info := range Infos {
			fmt.Println(infoID, "Description:", info.Description, "Number:", info.Number, "Type:", info.Type)
		}

		fmt.Println("SAMPLE FORMATS")
		for sampleID, sampleFmt := range SampleFormats {
			fmt.Println(sampleID, "Description:", sampleFmt.Description, "Number:", sampleFmt.Number, "Type:", sampleFmt.Type)
		}

		fmt.Println("FILTERS")
		for filterId, filterName := range Filters {
			fmt.Println(filterId, filterName)
		}

		fmt.Println("CONTIGS")
		for contigsId, contigsName := range Contigs {
			fmt.Println(contigsId, contigsName)
		}

		fmt.Println("EXTRAS")
		for extrasId, extrasName := range Extras {
			fmt.Println(extrasId, extrasName)
		}
	}

	var rowNum int64

	for {
		variant := vr.Read()

		// if vr.LineNumber >= 30000 {
		// 	break
		// }

		if vr.LineNumber%100000 == 0 && vr.LineNumber != 0 {
			fmt.Println(vr.LineNumber,
				variant.Chromosome,
				variant.Pos,
			)

			if BREAKAT > 0 && vr.LineNumber >= BREAKAT {
				return
			}
		}

		// continue

		if variant == nil {
			if e := vr.Error(); e != io.EOF && e != nil {
				vr.Clear()
			}
			break
		}

		if vr.Error() != nil {
			// fmt.Println(" -- vr error", vr.Error())
			vr.Clear()
			if continueOnError {
				continue
			} else {
				break
			}
		}

		if len(variant.Samples) == 0 {
			if DEBUG {
				fmt.Println("NO VARIANTS")
			}
		} else {

			// reg := new(interfaces.VCFRegister)

			// type VCFRegister struct {
			//  Samples      *[]string
			// 	IsHomozygous bool
			// 	IsIndel      bool
			// 	IsMNP        bool
			//  Row          uint64
			// 	Chromosome   string
			// 	Position     uint64
			// 	Quality      float32
			// 	Info         map[string]interface{}
			// 	Filter       string
			// 	NumAlt       uint64
			// 	Phased       bool
			// 	GT           [][]int
			// 	Fields       map[string]string
			// }

			if DEBUG {
				fmt.Printf("%d\t%s\t%d\t%s\t%s\t%v\n",
					rowNum,
					variant.Chromosome,
					variant.Pos,
					variant.Id(),
					variant.Ref(),
					variant.Alt())
				fmt.Printf(" Qual: %v\n", variant.Quality)
				fmt.Printf(" Filter: %v\n", variant.Filter)
				fmt.Printf(" Info: %v\n", variant.Info())
				fmt.Printf(" Format: %v\n", variant.Format)
				fmt.Printf(" Samples: %v\n", variant.Samples)

				// type Variant struct {
				// 	Chromosome      string
				// 	Pos        		uint64
				// 	Id         		string
				// 	Ref        		string
				// 	Alt        		[]string
				// 	Quality    		float32
				// 	Filter     		string
				// 	Info       		InfoMap
				// 	Format     		[]string
				// 	Samples    		[]*SampleGenotype
				// 	Header     		*Header
				// 	LineNumber 		int64
				// }

				vinfo := variant.Info()

				fmt.Println(" INFO:")
				for _, infoKey := range vinfo.Keys() {
					nfo, _ := vinfo.Get(infoKey)
					fmt.Println("  ", infoKey, ":", nfo)
				}

				for samplePos, sampleName := range SampleNames {
					sample := variant.Samples[samplePos]

					if sample != nil {
						fmt.Println("", "sample: #", samplePos,
							"name:", sampleName,
							"Phased:", sample.Phased,
							"GT:", sample.GT,
							"DP:", sample.DP,
							"GL:", sample.GL,
							"GQ:", sample.GQ,
							"MQ:", sample.MQ,
							"Fields:", sample.Fields)

						// &{false [] 0 [] 0 0 map[]}
						// type SampleGenotype struct {
						// 	Phased bool
						// 	GT     []int
						// 	DP     int
						// 	GL     []float32
						// 	GQ     int
						// 	MQ     int
						// 	Fields map[string]string
						// }

						// var pl interface{}

						// if hasPL {
						// 	pl, err = variant.GetGenotypeField(sample, "PL", plFmt)

						// 	if err != nil && sample != nil {
						// 		fmt.Println("", "ERR PL:", err)
						// 		log.Fatal(err)
						// 	}
						// } else {
						// 	pl = nil
						// }

						// fmt.Println("", "PL:", pl, "GQ:", sample.GQ, "DP:", sample.DP)

						fmt.Print(" FIELDS:")
						for fieldId, fieldVal := range sample.Fields {
							fmt.Print(" ", fieldId, ":", fieldVal)
						}
						fmt.Println("")
					} // if sample
					vr.Clear()
				} // for sample
			} // if debug

			samplesGT := make([]interfaces.VCFGT, numSamples, numSamples)

			for samplePos, _ := range variant.Samples {
				sample := variant.Samples[samplePos]

				if sample != nil {
					// fmt.Println("", "sample: #", samplePos,
					// 	"name:", sampleName,
					// 	"Phased:", sample.Phased,
					// 	"GT:", sample.GT,
					// 	"DP:", sample.DP,
					// 	"GL:", sample.GL,
					// 	"GQ:", sample.GQ,
					// 	"MQ:", sample.MQ,
					// 	"Fields:", sample.Fields)
					samplesGT[samplePos].GT = sample.GT
				} else {
					samplesGT[samplePos].GT = interfaces.VCFGTVal{-1, -1}
				}
			}

			register := interfaces.VCFRegisterRaw{
				LineNumber: variant.LineNumber,
				Chromosome: variant.Chromosome,
				Position:   variant.Pos,
				Alt:        variant.Alt(),
				Samples:    samplesGT,
			}
			callback(&SampleNames, &register)
		} // if has variant
	} //for variant

	fmt.Println("Finished")
}
