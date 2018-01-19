package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	model "github.com/nicouch/CBLogForwarder/model"
)

type supervisor struct {
	Files []model.LogFile `json:"files"`
}

var config string

func main() {
	flag.StringVar(&config, "configuration", "CBLogForwarder.conf", "configuration file")
	flag.Parse()

	s, err := loadConfiguration(config)
	if err != nil {
		fmt.Println("invalid configuration")
		return
	}
}

func loadConfiguration(config string) (s supervisor, err error) {
	c, err := ioutil.ReadFile(config)
	if err != nil {
		return
	}
	err = json.Unmarshal(c, &s)

	return
}
