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
	pkg := &ast.Package{
		QualifiedName: "pkg1",
		Classes:       []*ast.Class{},
	}
	printer.FprintPlantUML(buf, pkg)
}

func TestFprintPlantUMLStdLibs(t *testing.T) {
	for _, name := range []string{"io", "net", "net/http"} {
		pkg, err := parser.ParsePackage(name)
		if err != nil {
			t.Fatal(err)
		}
		buf := &bytes.Buffer{}
		printer.FprintPlantUML(buf, pkg)
		t.Log(buf)
	}

}
