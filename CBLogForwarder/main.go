package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/nicouch/gotail"
)

type supervisor struct {
	Files []logFile `json:"files"`
}

type logFile struct {
	FileName      string `json:"file"`
	SplitOn       string `json:"splitOn"`
	OutputIndices []int  `json:"outputIndices"`
}

var (
	config       string
	selectedMode *int
)

const (
	stream = iota
	batch

	usage = `
	Export all log files with transformation line by line
	Transformations are based on a config file, example:
			{	
				"files": [
						{
							"file": "fakeFile1.log",
							"splitOn": " ",
							"outputIndices": [3, 4, 0, 1, 2]
						}
					]
			}
		Each file JSON object must contains:
			- a file attribute which reprensents the file path
			- a splitOn attribute to split the line
			- an outputIndices to reorder the split. You don't need to specify all indices

	Arguments
			- configuration: specifies the configuration file path (no default)
			- mode: accepts 0 or 1. O is stream mode, 1 is batch (no default)
	`
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.StringVar(&config, "configuration", "", "configuration file")
	selectedMode = flag.Int("mode", 2, "mode for reading file between stream (continuous) and batch (one shot)")
	flag.Parse()

	if config == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

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
		fmt.Println(usage)
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

func streamer(s supervisor) {
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

func batcher(s supervisor) {
	var wg sync.WaitGroup
	wg.Add(len(s.Files))

	for _, f := range s.Files {
		go readFile(f, &wg)
	}

	wg.Wait()

	os.Exit(0)
}

func readFile(logFile logFile, wg *sync.WaitGroup) {
	defer wg.Done()
	f, err := os.Open(logFile.FileName)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	line, err := r.ReadString('\n')
	for err == nil {
		fmt.Println(transform(line, logFile))
		line, err = r.ReadString('\n')
	}
}

func transform(line string, file logFile) string {
	l := strings.Split(line, file.SplitOn)
	format := make([]string, 0)
	for i := 0; i < len(file.OutputIndices); i++ {
		format = append(format, l[file.OutputIndices[i]])
	}
	s := strings.Join(format, " ")

	return s
}
