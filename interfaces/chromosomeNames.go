package interfaces

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

import "github.com/sauloalgolang/introgressionbrowser/save"

const IndexExtension = ".ibindex"

type ChromosomeInfo struct {
	ChromosomeName string
	StartPosition  int64
	NumRegisters   int64
}

type ChromosomeNamesType struct {
	Infos []ChromosomeInfo
}

func NewChromosomeNames(size int, cap int) (cn *ChromosomeNamesType) {
	cn = &ChromosomeNamesType{
		Infos: make([]ChromosomeInfo, size, cap),
	}
	return cn
}

func (cn *ChromosomeNamesType) IndexFileName(outPrefix string) (indexFile string) {
	indexFile = save.GenFilename(outPrefix, IndexExtension)
	return indexFile
}

func (cn *ChromosomeNamesType) Save(outPrefix string) {
	save.SaveWithExtension(outPrefix, "yaml", IndexExtension, cn)
}

func (cn *ChromosomeNamesType) Load(outPrefix string) {
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

func (cn *ChromosomeNamesType) Add(chromosomeName string, startPosition int64) {
	if !(chromosomeName == "") { // valid chromosome name
		cn.Infos = append(cn.Infos, ChromosomeInfo{
			ChromosomeName: chromosomeName,
			StartPosition:  startPosition,
			NumRegisters:   -1,
		})

	} else {
		fmt.Println("got last chromosome", cn)

		for p := 0; p < len(cn.Infos)-1; p++ {
			infoC := &cn.Infos[p]
			infoN := &cn.Infos[p+1]
			infoC.NumRegisters = infoN.StartPosition - infoC.StartPosition
		}

		cn.Infos[len(cn.Infos)-1].NumRegisters = startPosition - cn.Infos[len(cn.Infos)-2].StartPosition

		fmt.Println("fixed chromosome sizes", cn)
	}
}
