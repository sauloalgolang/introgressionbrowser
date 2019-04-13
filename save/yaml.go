package save

import (
	"gopkg.in/yaml.v2"
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
