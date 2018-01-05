package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nicouch/gotail"
)

type logFile struct {
	FileName      string `json:"file"`
	SplitOn       string `json:"splitOn"`
	OutputIndices []int  `json:"outputIndices"`
}

type supervisor struct {
	Files []logFile `json:"files"`
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

	tails := make([]*gotail.Tail, 0)
	for _, f := range s.Files {
		tail, err := gotail.NewTail(f.FileName, gotail.Config{Timeout: 10, SplitOn: f.SplitOn, Output: f.OutputIndices})
		if err != nil {
			log.Fatalln(err)
		}
		tails = append(tails, tail)
	}

	// lines on the tail.Lines channel for new lines.
	agg := make(chan string)
	for _, t := range tails {
		go func(c chan string) {
			for l := range c {
				agg <- l
			}
		}(t.Lines)
	}

	for {
		fmt.Println(<-agg)
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
