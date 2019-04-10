package save

import (
	"encoding/gob"
	"os"
)

// https://medium.com/@kpbird/golang-serialize-struct-using-gob-part-1-e927a6547c00

// type Marshaler func(val interface{}) ([]byte, error)
// type UnMarshaler func(data []byte, v interface{}) error

func globMarsheler(filePath string, object interface{}) ([]byte, error) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}

	return []byte{}, err
}

func globUnMarsheler(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	return err
}
