package main

import (
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type config struct {
	Prefix       string `yaml:"prefix"`
	DBDriver     string `yaml:"db_driver"`
	DBSourceName string `yaml:"db_source_name"`
}

func (c *config) parseConfig(filename string) *config {
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
