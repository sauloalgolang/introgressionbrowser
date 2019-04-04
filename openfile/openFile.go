package openfile

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

func OpenFile(sourceFile string, isTar bool, isGz bool, continueOnError bool, callBack interfaces.VCFMaskedReaderType) {
	f, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	if !isTar && !isGz {
		callBack(io.Reader(f), continueOnError)
	} else {
		gzReader, err := gzip.NewReader(f)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer gzReader.Close()

		if !isTar {
			callBack(gzReader, continueOnError)
		} else {
			tarReader := tar.NewReader(gzReader)

			i := 0
			for {
				header, err := tarReader.Next()

				if err == io.EOF {
					break
				}

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				name := header.Name

				switch header.Typeflag {
				case tar.TypeDir:
					continue
				case tar.TypeReg:
					fmt.Println("(", i, ")", "Name: ", name)
					callBack(tarReader, continueOnError)
				default:
					fmt.Printf("%s : %c %s %s\n",
						"Yikes! Unable to figure out type",
						header.Typeflag,
						"in file",
						name,
					)
				}

				i++
			}
		}
	}
}
