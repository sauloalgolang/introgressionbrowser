package save

import (
	"encoding/gob"
	"io"
	"os"
)

// https://medium.com/@kpbird/golang-serialize-struct-using-gob-part-1-e927a6547c00

func gobMarshaler(filePath string, object interface{}) ([]byte, error) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		gobMarshalerWriter(file, object)
	}

	return []byte{}, err
}

func gobMarshalerWriter(file io.Writer, object interface{}) {
	encoder := gob.NewEncoder(file)
	encoder.Encode(object)
}

func gobUnMarshaler(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		err = gobUnMarshalerReader(file, object)
	}

	return err
}

func gobUnMarshalerReader(file io.Reader, object interface{}) (err error) {
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(object)
	return err
}
