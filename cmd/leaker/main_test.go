package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vodafon/leaker"
)

func TestKeys(t *testing.T) {
	testFiles(t, "1")
}

func testFiles(t *testing.T, suf string) {
	buf := &bytes.Buffer{}

	processor := Processor{
		w:  buf,
		bl: blackList("./testdata/bl.txt"),
		validators: []leaker.Validator{
			leaker.NewZxcvbnValidator(80.0),
		},
	}

	inl, err := lines("testdata/" + suf + ".in")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range inl {
		processor.Process(v)
	}

	outl, err := lines("testdata/" + suf + ".out")
	if err != nil {
		t.Fatal(err)
	}
	exp := strings.Join(outl, "\n") + "\n"
	res := buf.String()

	if exp != res {
		t.Errorf("Incorrect result. Expected\n%s\n, got\n%s\n", exp, res)
	}
}
