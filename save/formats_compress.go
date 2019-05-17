package save

import (
	log "github.com/sirupsen/logrus"	
	"os"
)

// Compressors holds the available compressors
var Compressors = map[string]CompressFormat{
	"none": CompressFormat{
		Compressor: "none",
		Extension:  "",
		Interface:  emptyInterface,
	},
	"snappy": CompressFormat{
		Compressor: "snappy",
		Extension:  "snappy",
		Interface:  snappyInterface,
	},
	"gzip": CompressFormat{
		Compressor: "gzip",
		Extension:  "gz",
		Interface:  gzipInterface,
	},
}

// CompressorNames holds the compressor names
var CompressorNames = []string{"none", "snappy", "gzip"}

// DefaultCompressor holds the name of the default compressor
var DefaultCompressor = "none"

//
//
// Compress Types
//
//

// CompressFormat struct holding the information about a available compressor
type CompressFormat struct {
	Compressor string
	Extension  string
	Interface  CompressorInterface
}

//
//
// Compress functions
//
//

// GetCompressInformation returns the information regarding a given compressor
func GetCompressInformation(compressor string) *CompressFormat {
	sf, ok := Compressors[compressor]

	if !ok {
		log.Println("Unknown compressor: ", compressor, ". valid compressor are:")
		for k := range Compressors {
			log.Println(" ", k)
		}
		os.Exit(1)
	}

	return &sf
}

// GetCompressExtension returns the file extension for a given compressor
func GetCompressExtension(compressor string) string {
	sc := GetCompressInformation(compressor)
	return sc.Extension
}

// GetCompressInterface returns the CompressorInterface for a given compressor
func GetCompressInterface(compressor string) CompressorInterface {
	sc := GetCompressInformation(compressor)
	return sc.Interface
}

// GetCompressInterfaceReader returns the reader for a given compressor
func GetCompressInterfaceReader(compressor string) GenericNewReader {
	sc := GetCompressInterface(compressor)
	return sc.NewReader
}

// GetCompressInterfaceWriter returns the writer for a given compressor
func GetCompressInterfaceWriter(compressor string) GenericNewWriter {
	sc := GetCompressInterface(compressor)
	return sc.NewWriter
}

// GetCompressIsCompressed checks if a compressor compresses or not
func GetCompressIsCompressed(compressor string) bool {
	sf := GetCompressInformation(compressor)
	return sf.Compressor != "none"
}
