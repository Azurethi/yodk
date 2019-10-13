package parser_test

import (
	"testing"

	"github.com/dbaumgarten/yodk/parser"
	"github.com/dbaumgarten/yodk/testdata"
)

func TestGenerator(t *testing.T) {
	p := parser.NewParser()
	gen := parser.Printer{}
	parsed, err := p.Parse(testdata.TestProgram)
	if err != nil {
		t.Fatal(err)
	}
	generated, err := gen.Print(parsed)
	if err != nil {
		t.Fatal(err)
	}

	err = testdata.ExecuteTestProgram(generated)
	if err != nil {
		t.Fatal(err)
	}
}
