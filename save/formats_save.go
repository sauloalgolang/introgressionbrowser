package save

import (
	"fmt"
	"io"
	"os"
)

import (
	// "gopkg.in/mgo.v2/bson"
	"gopkg.in/yaml.v2"
)

// import "github.com/kelindar/binary"

//
//
// Available formats
//
//

var Formats = map[string]SaveFormat{
	"yaml": SaveFormat{
		Extension:                 "yaml",
		HasMarshal:                true,
		HasStreamer:               true,
		Marshaler:                 yaml.Marshal,
		UnMarshaler:               yaml.Unmarshal,
		MarshalerStreamer:         yamlMarshaler,
		UnMarshalerStreamer:       yamlUnMarshaler,
		MarshalerStreamerWriter:   yamlMarshalerWriter,
		UnMarshalerStreamerReader: yamlUnMarshalerReader,
	},
	"gob": SaveFormat{
		Extension:                 "gob",
		HasMarshal:                false,
		HasStreamer:               true,
		Marshaler:                 emptyMarshaler,
		UnMarshaler:               emptyUnMarshaler,
		MarshalerStreamer:         gobMarshaler,
		UnMarshalerStreamer:       gobUnMarshaler,
		MarshalerStreamerWriter:   gobMarshalerWriter,
		UnMarshalerStreamerReader: gobUnMarshalerReader,
	},
	// "bson": SaveFormat{
	// 	Extension:           "bson",
	// 	HasMarshal:          true,
	// 	HasStreamer:         false,
	// 	Marshaler:           bson.Marshal,
	// 	UnMarshaler:         bson.Unmarshal,
	// 	MarshalerStreamer:   emptyMarshalerStreamer,
	// 	UnMarshalerStreamer: emptyUnMarshalerStreamer,
	// },
	// "json": SaveFormat{".json", true, false, json.Marshal, json.Unmarshal, emptyMarshalerStreamer, emptyUnMarshalerStreamer}, // no numerical key
	// "binary": SaveFormat{"bin", true, false, binary.Marshal, binary.Unmarshal, emptyMarshalerStreamer, emptyUnMarshalerStreamer}, // fail to export reader
}

var FormatNames = []string{"yaml", "gob"} // "bson",
var DefaultFormat = "yaml"

//
//
// Format Types
//
//

type Marshaler func(interface{}) ([]byte, error)
type UnMarshaler func([]byte, interface{}) error
type MarshalerStreamer func(string, interface{}) ([]byte, error)
type UnMarshalerStreamer func(string, interface{}) error
type MarshalerStreamerWriter func(io.Writer, interface{})
type UnMarshalerStreamerReader func(io.Reader, interface{}) error

type SaveFormat struct {
	Extension                 string
	HasMarshal                bool // returns bytes
	HasStreamer               bool // write directly to stream
	Marshaler                 Marshaler
	UnMarshaler               UnMarshaler
	MarshalerStreamer         MarshalerStreamer
	UnMarshalerStreamer       UnMarshalerStreamer
	MarshalerStreamerWriter   MarshalerStreamerWriter
	UnMarshalerStreamerReader UnMarshalerStreamerReader
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
// Format functions
//
//

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

func GetFormatHasMarshal(format string) bool {
	sf := GetFormatInformation(format)
	return sf.HasMarshal
}

func GetFormatHasStreamer(format string) bool {
	sf := GetFormatInformation(format)
	return sf.HasStreamer
}

func GetFormatExtension(format string) string {
	sf := GetFormatInformation(format)
	return sf.Extension
}

func GetFormatMarshaler(format string) Marshaler {
	sf := GetFormatInformation(format)
	return sf.Marshaler
}

func GetFormatMarshalerStreamer(format string) MarshalerStreamer {
	sf := GetFormatInformation(format)
	return sf.MarshalerStreamer
}

func GetFormatMarshalerStreamerWriter(format string) MarshalerStreamerWriter {
	sf := GetFormatInformation(format)
	return sf.MarshalerStreamerWriter
}

func GetFormatUnMarshaler(format string) UnMarshaler {
	sf := GetFormatInformation(format)
	return sf.UnMarshaler
}

func GetFormatUnMarshalerStreamer(format string) UnMarshalerStreamer {
	sf := GetFormatInformation(format)
	return sf.UnMarshalerStreamer
}

func GetFormatUnMarshalerStreamerReader(format string) UnMarshalerStreamerReader {
	sf := GetFormatInformation(format)
	return sf.UnMarshalerStreamerReader
}
