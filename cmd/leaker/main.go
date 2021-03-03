package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/vodafon/leaker"
	"github.com/vodafon/swork"
)

var (
	flagProcs = flag.Int("procs", 20, "concurrency")
	flagScore = flag.Int("score", 70, "score")
)

type Processor struct {
	validators []leaker.Validator
	w          io.Writer
}

func (obj Processor) Process(line string) {
	for _, validator := range obj.validators {
		if validator.IsValid(line) {
			fmt.Fprintf(obj.w, "%s\n", line)
		}
	}
}

func main() {
	flag.Parse()
	if *flagProcs < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	score := float64(*flagScore)
	if score < 1.0 {
		score = 70.0
	}

	processor := Processor{
		w: os.Stdout,
		validators: []leaker.Validator{
			leaker.NewZxcvbnValidator(score),
		},
	}
	w := swork.NewWorkerGroup(*flagProcs, processor)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		w.StringC <- sc.Text()
	}

	close(w.StringC)

	if err := sc.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to read input: %s\n", err)
	}

	w.Wait()
}
