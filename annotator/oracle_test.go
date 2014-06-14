package annotator_test

import (
	"encoding/json"
	"testing"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/parser"
)

func TestOracleStdLibs(t *testing.T) {
	for _, name := range []string{"io", "net", "net/http"} {
		pkg, err := parser.ParsePackage(name)
		if err != nil {
			t.Fatal(err)
		}
		err = annotator.Oracle(pkg)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.MarshalIndent(pkg, "", "\t")
		t.Logf("%s", b)
	}
}

func TestOracleGoAST(t *testing.T) {
	name := "go/ast"
	pkg, err := parser.ParsePackage(name)
	if err != nil {
		t.Fatal(err)
	}
	err = annotator.Oracle(pkg)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(pkg, "", "\t")
	t.Logf("%s", b)
}
