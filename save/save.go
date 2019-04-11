package save

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
	// "encoding/json"
	"io/ioutil"
	"os"
)

// import "github.com/kelindar/binary"

//
//
// Available formats
//
//

var Formats = map[string]SaveFormat{
	"yaml": SaveFormat{"yaml", true, true, yaml.Marshal, yaml.Unmarshal, yamlMarshaler, yamlUnMarshaler},
	"bson": SaveFormat{"bson", true, false, bson.Marshal, bson.Unmarshal, emptyMarshalerStreamer, emptyUnMarshalerStreamer},
	// "json": SaveFormat{".json", true, false, json.Marshal, json.Unmarshal, emptyMarshalerStreamer, emptyUnMarshalerStreamer}, // no numerical key
	// "binary": SaveFormat{"bin", true, false, binary.Marshal, binary.Unmarshal, emptyMarshalerStreamer, emptyUnMarshalerStreamer}, // fail to export reader
	"gob": SaveFormat{"gob", false, true, emptyMarshaler, emptyUnMarshaler, gobMarshaler, gobUnMarshaler},
}

var FormatNames = []string{"yaml", "bson", "gob"}

//
//
// Types
//
//

type Marshaler func(interface{}) ([]byte, error)
type UnMarshaler func([]byte, interface{}) error
type MarshalerStreamer func(string, interface{}) ([]byte, error)
type UnMarshalerStreamer func(string, interface{}) error

type SaveFormat struct {
	Extension           string
	HasMarshal          bool // returns bytes
	HasStreamer         bool // write directly to stream
	Marshaler           Marshaler
	UnMarshaler         UnMarshaler
	MarshalerStreamer   MarshalerStreamer
	UnMarshalerStreamer UnMarshalerStreamer
}

func emptyMarshaler(val interface{}) ([]byte, error) {
	return []byte{}, *new(error)
}

func emptyUnMarshaler(data []byte, val interface{}) error {
	return *new(error)
}

func emptyMarshalerStreamer(filename string, val interface{}) ([]byte, error) {
	return []byte{}, *new(error)
}

func emptyUnMarshalerStreamer(filename string, val interface{}) error {
	return *new(error)
}

//
//
// General Functions
//
//

func GenFilename(outPrefix string, extension string) string {
	return outPrefix + "." + extension
}

func GetFormatInformation(format string) SaveFormat {
	sf, ok := Formats[format]

	if !ok {
		fmt.Println("Unknown format: ", format, ". valid formats are:")
		for k, _ := range Formats {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}

	return sf
}

func GetExtension(format string) string {
	sf := GetFormatInformation(format)
	return sf.Extension
}

func GetMarshaler(format string) Marshaler {
	sf := GetFormatInformation(format)
	return sf.Marshaler
}

func GetMarshalerStreamer(format string) MarshalerStreamer {
	sf := GetFormatInformation(format)
	return sf.MarshalerStreamer
}

func GetUnMarshaler(format string) UnMarshaler {
	sf := GetFormatInformation(format)
	return sf.UnMarshaler
}

func GetUnMarshalerStreamer(format string) UnMarshalerStreamer {
	sf := GetFormatInformation(format)
	return sf.UnMarshalerStreamer
}

func GetHasMarshal(format string) bool {
	sf := GetFormatInformation(format)
	return sf.HasMarshal
}

func GetHasStreamer(format string) bool {
	sf := GetFormatInformation(format)
	return sf.HasStreamer
}

//
//
// Save
//
//

//
// Save
//

func Save(outPrefix string, format string, val interface{}) {
	extension := GetExtension(format)
	SaveWithExtension(outPrefix, format, extension, val)
}

func SaveWithExtension(outPrefix string, format string, extension string, val interface{}) {
	hasStreamer := GetHasStreamer(format)
	hasMarshal := GetHasMarshal(format)

	if hasStreamer {
		marshaler := GetMarshalerStreamer(format)
		saveDataStream(outPrefix, format, extension, marshaler, val)

	} else if hasMarshal {
		marshaler := GetMarshaler(format)
		saveData(outPrefix, format, extension, marshaler, val)
	}
}

func saveData(outPrefix string, format string, extension string, marshaler Marshaler, val interface{}) {
	d, err := marshaler(val)

	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	outfile := GenFilename(outPrefix, extension)
	fmt.Println("saving data to ", outfile)

	err = ioutil.WriteFile(outfile, d, 0644)
	fmt.Println("  done")
}

func saveDataStream(outPrefix string, format string, extension string, marshaler MarshalerStreamer, val interface{}) {
	outfile := GenFilename(outPrefix, extension)
	fmt.Println("saving stream to ", outfile)
	marshaler(outfile, val)
}

//
//
// Load
//
//

//
// Load Marshal
//

func Load(outPrefix string, format string, val interface{}) {
	extension := GetExtension(format)
	LoadWithExtension(outPrefix, format, extension, val)
}

func LoadWithExtension(outPrefix string, format string, extension string, val interface{}) {
	hasStreamer := GetHasStreamer(format)
	hasMarshal := GetHasMarshal(format)

	if hasStreamer {
		unmarshaler := GetUnMarshalerStreamer(format)
		loadDataStream(outPrefix, format, extension, unmarshaler, val)

	} else if hasMarshal {
		unmarshaler := GetUnMarshaler(format)
		loadData(outPrefix, format, extension, unmarshaler, val)
	}
}

func loadData(outPrefix string, format string, extension string, unmarshaler UnMarshaler, val interface{}) {
	outfile := GenFilename(outPrefix, extension)

	data, err := ioutil.ReadFile(outfile)

	if err != nil {
		fmt.Printf("dump file. Get err   #%v ", err)
	}

	err = unmarshaler(data, val)

	if err != nil {
		fmt.Printf("cannot unmarshal data: %v", err)
	}
}

func loadDataStream(outPrefix string, format string, extension string, unmarshaler UnMarshalerStreamer, val interface{}) {
	outfile := GenFilename(outPrefix, extension)
	fmt.Println("loading from ", outfile)
	unmarshaler(outfile, val)
}
