package interfaces

import (
	log "github.com/sirupsen/logrus"
)

// IndexExtension is the default extension from the ibrowser save file
const IndexExtension = "ibindex"

// ChromosomeInfo holds the chromosome information
type ChromosomeInfo struct {
	ChromosomeName string
	StartPosition  int64
	NumRegisters   int64
}

// ChromosomeNamesType holds the chromosome position information
type ChromosomeNamesType struct {
	Infos          []ChromosomeInfo
	NumChromosomes int64
	StartPosition  int64
	EndPosition    int64
	NumRegisters   int64
}

// NewChromosomeNames creates a new ChromosomeNamesType instance
func NewChromosomeNames(size int, cap int) (cn *ChromosomeNamesType) {
	cn = &ChromosomeNamesType{
		Infos: make([]ChromosomeInfo, size, cap),
	}
	return cn
}

// Save saves in YAML format with the IndexExtension
func (cn *ChromosomeNamesType) Save(outPrefix string) {
	saver := NewSaver(outPrefix, "yaml")
	saver.SetExtension(IndexExtension)
	saver.Save(cn)
}

// Load loads the data dump
func (cn *ChromosomeNamesType) Load(outPrefix string) {
	saver := NewSaver(outPrefix, "yaml")
	saver.SetExtension(IndexExtension)
	saver.Load(cn)
}

// Exists checks whether dump file exists
func (cn *ChromosomeNamesType) Exists(outPrefix string) (bool, error) {
	saver := NewSaver(outPrefix, "yaml")
	saver.SetExtension(IndexExtension)
	return saver.Exists()
}

// Add adds chromosome
func (cn *ChromosomeNamesType) Add(chromosomeName string, startPosition int64) {
	if !(chromosomeName == "") { // valid chromosome name
		cn.Infos = append(cn.Infos, ChromosomeInfo{
			ChromosomeName: chromosomeName,
			StartPosition:  startPosition,
			NumRegisters:   -1,
		})

	} else {
		log.Println("got last chromosome", cn)

		cn.NumChromosomes = int64(len(cn.Infos))
		cn.NumRegisters = 0

		for p := int64(0); p < cn.NumChromosomes-1; p++ {
			infoC := &cn.Infos[p]
			infoN := &cn.Infos[p+1]
			infoC.NumRegisters = infoN.StartPosition - infoC.StartPosition
			cn.NumRegisters += infoC.NumRegisters
		}

		cn.Infos[cn.NumChromosomes-1].NumRegisters = startPosition - cn.Infos[cn.NumChromosomes-2].StartPosition

		cn.StartPosition = cn.Infos[0].StartPosition
		cn.EndPosition = cn.Infos[cn.NumChromosomes-1].StartPosition

		log.Println("fixed chromosome sizes", cn)
	}
}
