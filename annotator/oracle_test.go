package annotator_test

import (
	"encoding/json"
	"testing"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/parser"
)

func TestOracleCGo(t *testing.T) {
	pkg, err := parser.ParsePackage("github.com/t-yuki/godoc2puml/annotator/testdata/cgo")
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

func TestOracleCGoImport(t *testing.T) {
	pkg, err := parser.ParsePackage("github.com/t-yuki/godoc2puml/annotator/testdata/cgoimport")
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

func TestOracleStdIO(t *testing.T) {
	testOracle(t, "io")
}

func TestOracleStdImage(t *testing.T) {
	testOracle(t, "image")
}

func TestOracleStdNet(t *testing.T) {
	testOracle(t, "net")
}

func TestOracleStdNetHTTP(t *testing.T) {
	testOracle(t, "net/http")
}

func testOracle(t *testing.T, name string) {
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
