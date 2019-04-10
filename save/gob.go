package save

import (
	"encoding/gob"
	"gopkg.in/yaml.v2"
	"os"
)

// https://medium.com/@kpbird/golang-serialize-struct-using-gob-part-1-e927a6547c00

// type Marshaler func(val interface{}) ([]byte, error)
// type UnMarshaler func(data []byte, v interface{}) error

func gobMarshaler(filePath string, object interface{}) ([]byte, error) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}

	return []byte{}, err
}

func gobUnMarshaler(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}

	return err
}

//
//
// Yaml
//
//
func yamlMarshaler(filePath string, object interface{}) ([]byte, error) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		encoder := yaml.NewEncoder(file)
		encoder.Encode(object)
		encoder.Close()
	}

	return []byte{}, err
}

func yamlUnMarshaler(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(object)
	}

	return err
}
