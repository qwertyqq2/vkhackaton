package configs

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type NodeConfig struct {
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
	Port string `json:"port"`
}

type Config struct {
	Me        NodeConfig   `json:"me"`
	Neigboors []NodeConfig `json:"neigboors"`
}

func read() (*Config, error) {
	jsonFile, err := os.Open("addresses.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var conf *Config
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
