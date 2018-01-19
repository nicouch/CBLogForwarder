package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	config       string
	selectedMode *int
)

const (
	stream = iota
	batch
)

func main() {
	flag.StringVar(&config, "configuration", "CBLogForwarder.conf", "configuration file")
	selectedMode = flag.Int("mode", 0, "mode for reading file between stream (continuous) and batch (one shot)")
	flag.Parse()

	s, err := loadConfiguration(config)
	if err != nil {
		fmt.Println("invalid configuration")
		return
	}

	switch *selectedMode {
	case stream:
		streamer(s)
	case batch:
		batcher(s)
	default:
		os.Exit(0)
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
