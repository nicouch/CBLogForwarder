package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

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
