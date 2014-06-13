package parser_test

import (
	"encoding/json"
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

func TestParsePackageIO(t *testing.T) {
	pkg, err := parser.ParsePackage("io")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(pkg, "", "\t")
	t.Logf("%s", b)

}

func TestParsePackageStdLibs(t *testing.T) {
	for _, name := range []string{"go/ast", "go/token", "reflect", "database/sql", "image", "image/color", "os", "net", "net/http"} {
		pkg, err := parser.ParsePackage(name)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.MarshalIndent(pkg, "", "\t")
		t.Logf("%s", b)
	}
}
