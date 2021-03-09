package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/vodafon/leaker"
	"github.com/vodafon/swork"
)

var (
	flagProcs     = flag.Int("procs", 20, "concurrency")
	flagScore     = flag.Int("score", 80, "score")
	flagBlackList = flag.String("bl", "", "blacklist")
)

type Processor struct {
	validators []leaker.Validator
	w          io.Writer
	bl         []string
}

func (obj Processor) Process(line string) {
	for _, bl := range obj.bl {
		if strings.Contains(line, bl) {
			return
		}
	}
	for _, validator := range obj.validators {
		if validator.IsValid(line) {
			// fmt.Printf("%q - %v\n", line, zxcvbn.PasswordStrength(line, nil).Entropy)
			fmt.Fprintf(obj.w, "%s\n", line)
		}
	}
}

func blackList(fp string) []string {
	if fp == "" {
		return []string{}
	}
	res, err := lines(fp)
	if err != nil {
		log.Fatalf("read %q error: %v", fp, err)
	}
	return res
}

func main() {
	flag.Parse()
	if *flagProcs < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	score := float64(*flagScore)
	if score < 1.0 {
		score = 80.0
	}

	processor := Processor{
		w:  os.Stdout,
		bl: blackList(*flagBlackList),
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

func lines(fp string) ([]string, error) {
	_, err := os.Stat(fp)
	if err != nil {
		return []string{}, err
	}

	f, err := os.Open(fp)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	res := []string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		res = append(res, sc.Text())
	}
	return res, nil
}
