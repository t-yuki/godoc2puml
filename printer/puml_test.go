package printer_test

import (
	"bytes"
	"testing"

	"github.com/t-yuki/godoc2puml/ast"
	"github.com/t-yuki/godoc2puml/parser"
	"github.com/t-yuki/godoc2puml/printer"
)

func TestFprintPlantUML(t *testing.T) {
	buf := &bytes.Buffer{}
	scope := ast.NewScope()
	scope.Packages["pkg1"] = &ast.Package{
		Name:    "pkg1",
		Classes: []*ast.Class{},
	}
	printer.FprintPlantUML(buf, scope, []string{})
}

func TestFprintPlantUMLStdIO(t *testing.T) {
	testFprintPlantUML(t, "io")
}
func TestFprintPlantUMLStdNet(t *testing.T) {
	testFprintPlantUML(t, "net")
}

func TestFprintPlantUMLStdNet(t *testing.T) {
	testFprintPlantUML(t, "net/http")
}

func testFprintPlantUML(t *testing.T, name string) {
	scope := ast.NewScope()
	pkg, err := parser.ParsePackage(name)
	if err != nil {
		t.Fatal(err)
	}
	scope.Packages[pkg.Name] = pkg
	buf := &bytes.Buffer{}
	printer.FprintPlantUML(buf, scope, []string{})
	t.Log(buf)
}
