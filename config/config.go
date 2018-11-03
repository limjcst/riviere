package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Config is used to store configuations
type Config struct {
	Prefix       string `yaml:"prefix"`
	DBDriver     string `yaml:"db_driver"`
	DBSourceName string `yaml:"db_source_name"`
	Spec         string `yaml:"spec"`
}

// ParseConfig parses configurations from yaml formated file
func (c *Config) ParseConfig(filename string) *Config {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}
