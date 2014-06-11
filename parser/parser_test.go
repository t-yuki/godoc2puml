package parser_test

import (
	"testing"

	"github.com/t-yuki/godoc2puml/parser"
)

func TestParsePackage(t *testing.T) {
	pkg, err := parser.ParsePackage(".")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", pkg)
}

func TestParsePackageNetHttp(t *testing.T) {
	pkg, err := parser.ParsePackage("net/http")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", pkg)

}
