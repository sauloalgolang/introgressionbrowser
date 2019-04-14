package save

import (
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

//
//
// Yaml
//
//
func yamlMarshaler(filePath string, object interface{}) ([]byte, error) {
	file, err := os.Create(filePath)
	defer file.Close()

	if err == nil {
		yamlMarshalerWriter(file, object)
	}

	return []byte{}, err
}

func yamlMarshalerWriter(file io.Writer, object interface{}) {
	encoder := yaml.NewEncoder(file)
	encoder.Encode(object)
	encoder.Close()
}

func yamlUnMarshaler(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	defer file.Close()

	if err == nil {
		err = yamlUnMarshalerReader(file, object)
	}

	return err
}

func yamlUnMarshalerReader(file io.Reader, object interface{}) (err error) {
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(object)
	return err
}
