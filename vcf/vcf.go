package vcf

import (
	"fmt"
	"io"
)

import (
	"github.com/brentp/vcfgo"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/interfaces"
)

const DEBUG = false

// https://github.com/brentp/vcfgo/blob/master/examples/main.go

func ProcessVcf(r io.Reader, callback interfaces.VCFCallBack, continueOnError bool) {
	vr, err := vcfgo.NewReader(r, false)

	if err != nil {
		panic(err)
	}

	fmt.Printf("VR %v\n", vr)

	header := vr.Header
	SampleNames := header.SampleNames // []string
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

		vr.Clear()

		rowNum = vr.LineNumber

		// if rowNum >= 30000 {
		// 	break
		// }

		if rowNum%100000 == 0 && rowNum != 0 {
			fmt.Println(rowNum)
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
			callback(&SampleNames, variant)
		} // if has variant
	} //for variant

	fmt.Println("Finished")
}
