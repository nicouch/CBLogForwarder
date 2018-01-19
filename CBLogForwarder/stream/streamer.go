package streamer

import (
	"fmt"
	"log"

	"github.com/nicouch/gotail"
)

func stream() {
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
