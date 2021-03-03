package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/vodafon/leaker"
)

func TestAWSKeys(t *testing.T) {
	testFiles(t, "1")
}

func testFiles(t *testing.T, suf string) {
	buf := &bytes.Buffer{}

	processor := Processor{
		w: buf,
		validators: []leaker.Validator{
			leaker.NewZxcvbnValidator(70.0),
		},
	}

	for _, v := range lines(t, "testdata/"+suf+".in") {
		processor.Process(v)
	}

	exp := strings.Join(lines(t, "testdata/"+suf+".out"), "\n") + "\n"
	res := buf.String()

	if exp != res {
		t.Errorf("Incorrect result. Expected\n%s\n, got\n%s\n", exp, res)
	}
}

func lines(t *testing.T, fp string) []string {
	_, err := os.Stat(fp)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(fp)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res := []string{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		res = append(res, sc.Text())
	}
	return res
}
