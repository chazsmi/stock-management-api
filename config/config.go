package config

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Rabbit struct {
		Host string
	}
}

var File string

// Read and retrun config
func ReadReturn(path string) (*Config, error) {
	f := OpenConfig(&path)
	c, err := ReadConfig(f)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Open the config file for reading
func OpenConfig(configFile *string) *os.File {
	f, err := os.Open(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// ReadConfig config from io.Reader
func ReadConfig(r io.Reader) (*Config, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var c Config
	yaml.Unmarshal(data, &c)
	return &c, nil
}
