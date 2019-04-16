package openfile

// https://www.avitzurel.com/blog/2015/09/16/read-gzip-file-content-with-golang/
// https://gist.github.com/indraniel/1a91458984179ab4cf80

import (
	"archive/tar"
	"fmt"
	// "compress/gzip"
	gzip "github.com/klauspost/pgzip"
	"io"
	"os"
	"runtime"
)

import "github.com/sauloalgolang/introgressionbrowser/interfaces"

func OpenFile(sourceFile string, isTar bool, isGz bool, callBackParameters interfaces.CallBackParameters, callBack interfaces.VCFMaskedReaderType) {
	f, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	if !isTar && !isGz {
		callBack(io.Reader(f), callBackParameters)
	} else {
		runtime.GOMAXPROCS(runtime.NumCPU())

		gzReader, err := gzip.NewReaderN(f, 2500000, 32)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer gzReader.Close()

		if !isTar {
			callBack(gzReader, callBackParameters)
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
					callBack(tarReader, callBackParameters)
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
