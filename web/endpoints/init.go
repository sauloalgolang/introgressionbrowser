package endpoints

import (
	// "fmt"
	"log"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/ibrowser"
	"github.com/sauloalgolang/introgressionbrowser/save"
)

var DATABASE_DIR = "res/"

type Parameters = ibrowser.Parameters
type IBrowser = ibrowser.IBrowser

type DbDb struct {
	Db map[string]*IBrowser
}

var GuessPrefixFormat = save.GuessPrefixFormat
var GuessFormat = save.GuessFormat
var NewIBrowser = ibrowser.NewIBrowser

var databases DbDb

func init() {
	databases.Db = make(map[string]*IBrowser, 0)
}

func (d *DbDb) Register(fileName string, path string) (*IBrowser, bool) {
	log.Print("Registering db :: filename: '", fileName, "' path: '", path, "'")

	if ib, ok := d.Db[fileName]; ok {
		log.Println(" - Exists")
		return ib, true
	} else {
		log.Println(" - Loading")
		ib = NewIBrowser(Parameters{})
		ib.EasyLoadFile(path)
		d.Db[fileName] = ib
		return ib, true
	}
}

func (d *DbDb) Get(fileName string) (*IBrowser, bool) {
	if ib, ok := d.Db[fileName]; ok {
		return ib, ok
	} else {
		return ib, ok
	}
}
