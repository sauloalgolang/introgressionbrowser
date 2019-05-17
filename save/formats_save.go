package save

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

import (
	"gopkg.in/yaml.v2"
)

//
//
// Available formats
//
//

// Formats holds the available file formats
var Formats = map[string]Format{
	"yaml": Format{
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
	"gob": Format{
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
}

// FormatNames holds the name of the available file formats
var FormatNames = []string{"yaml"}

// DefaultFormat holds the name of the default file format
var DefaultFormat = "yaml"

//
//
// Format Types
//
//

// Marshaler marshaler function
type Marshaler func(interface{}) ([]byte, error)

// UnMarshaler unmarshaler function
type UnMarshaler func([]byte, interface{}) error

// MarshalerStreamer marshaler function which saves to stream
type MarshalerStreamer func(string, interface{}) ([]byte, error)

// UnMarshalerStreamer unmarshaler function which reads from stream
type UnMarshalerStreamer func(string, interface{}) error

// MarshalerStreamerWriter marshaler function which saves to stream accepting a writer
type MarshalerStreamerWriter func(io.Writer, interface{})

// UnMarshalerStreamerReader unmarshaler function which loads from stream accepting a reader
type UnMarshalerStreamerReader func(io.Reader, interface{}) error

// Format struct containing a extension and the converters to save
type Format struct {
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

// GetFormatInformation returns the information for a given format
func GetFormatInformation(format string) Format {
	sf, ok := Formats[format]

	if !ok {
		log.Println("Unknown format: ", format, ". valid formats are:")
		for k := range Formats {
			log.Println(" ", k)
		}
		os.Exit(1)
	}

	return sf
}

// GetFormatHasMarshal returns whether a format has a marshaller
func GetFormatHasMarshal(format string) bool {
	sf := GetFormatInformation(format)
	return sf.HasMarshal
}

// GetFormatHasStreamer returns whether a format has a streamer marshaller
func GetFormatHasStreamer(format string) bool {
	sf := GetFormatInformation(format)
	return sf.HasStreamer
}

// GetFormatExtension returns the extension for a given format
func GetFormatExtension(format string) string {
	sf := GetFormatInformation(format)
	return sf.Extension
}

// GetFormatMarshaler returns the marshaller for a given format
func GetFormatMarshaler(format string) Marshaler {
	sf := GetFormatInformation(format)
	return sf.Marshaler
}

// GetFormatMarshalerStreamer returns the streamer marshaller for a given format
func GetFormatMarshalerStreamer(format string) MarshalerStreamer {
	sf := GetFormatInformation(format)
	return sf.MarshalerStreamer
}

// GetFormatMarshalerStreamerWriter returns the streamer marshaller writer for a given format
func GetFormatMarshalerStreamerWriter(format string) MarshalerStreamerWriter {
	sf := GetFormatInformation(format)
	return sf.MarshalerStreamerWriter
}

// GetFormatUnMarshaler returns the unmarshaller for a given format
func GetFormatUnMarshaler(format string) UnMarshaler {
	sf := GetFormatInformation(format)
	return sf.UnMarshaler
}

// GetFormatUnMarshalerStreamer returns the streamer unmarshaller for a given format
func GetFormatUnMarshalerStreamer(format string) UnMarshalerStreamer {
	sf := GetFormatInformation(format)
	return sf.UnMarshalerStreamer
}

// GetFormatUnMarshalerStreamerReader returns the streamer marshaller reader for a given format
func GetFormatUnMarshalerStreamerReader(format string) UnMarshalerStreamerReader {
	sf := GetFormatInformation(format)
	return sf.UnMarshalerStreamerReader
}
