package network

import "encoding/json"

type Package struct {
	Option int
	Data   string
}

func (pack *Package) SerializePack() (string, error) {
	jsonData, err := json.MarshalIndent(*pack, " ", "\t")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func DeserializePack(data string) (*Package, error) {
	var pack Package
	err := json.Unmarshal([]byte(data), &pack)
	if err != nil {
		return nil, err
	}
	return &pack, nil
}
