package save

import (
	"fmt"
	"os"
)

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

var CompressorNames = []string{"none", "snappy", "gzip"}
var DefaultCompressor = "none"

//
//
// Compress Types
//
//

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

func GetCompressInformation(compressor string) *CompressFormat {
	sf, ok := Compressors[compressor]

	if !ok {
		fmt.Println("Unknown compressor: ", compressor, ". valid compressor are:")
		for k, _ := range Compressors {
			fmt.Println(" ", k)
		}
		os.Exit(1)
	}

	return &sf
}

func GetCompressExtension(compressor string) string {
	sc := GetCompressInformation(compressor)
	return sc.Extension
}

func GetCompressInterface(compressor string) CompressorInterface {
	sc := GetCompressInformation(compressor)
	return sc.Interface
}

func GetCompressInterfaceReader(compressor string) GenericNewReader {
	sc := GetCompressInterface(compressor)
	return sc.NewReader
}

func GetCompressInterfaceWriter(compressor string) GenericNewWriter {
	sc := GetCompressInterface(compressor)
	return sc.NewWriter
}

func GetCompressIsCompressed(compressor string) bool {
	sf := GetCompressInformation(compressor)
	return sf.Compressor != "none"
}
