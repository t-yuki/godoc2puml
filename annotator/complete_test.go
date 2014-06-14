package annotator_test

import (
	"testing"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/ast"
)

func TestComplete(t *testing.T) {
	scope := ast.NewScope()
	scope.Packages["pkg1"] = &ast.Package{
		Name: "pkg1",
		Classes: []*ast.Class{
			{
				Name: "A",
				Relations: []*ast.Relation{
					{Target: "io.Reader", RelType: ast.Implementation},
					{Target: "io.Writer", RelType: ast.Implementation},
					{Target: "io.Closer", RelType: ast.Implementation},
					{Target: "io.ReadWriter", RelType: ast.Implementation},
					{Target: "io.ReadCloser", RelType: ast.Implementation},
					{Target: "io.WriteCloser", RelType: ast.Implementation},
					{Target: "io.ReadWriteCloser", RelType: ast.Implementation},
				},
			},
			{
				Name: "XError",
				Relations: []*ast.Relation{
					{Target: "error", RelType: ast.Implementation},
					{Target: "fmt.Stringer", RelType: ast.Implementation},
					{Target: "os.PathError", RelType: ast.Association},
				},
				Methods: []*ast.Method{
					{Name: "Error"},
				},
			},
		},
	}
	err := annotator.Complete(scope)
	if err != nil {
		t.Fatal(err)
	}
	if len(scope.Packages) != 5 {
		for _, pkg := range scope.Packages {
			t.Logf("%+v", pkg)
		}
		t.Fatalf("pkg1, io, fmt, os, \"\" should be completed: %#v", scope.Packages)
	}
	if pkg, ok := scope.Packages["io"]; !ok || len(pkg.Interfaces) != 7 {
		t.Fatalf("io interfaces should be completed", pkg)
	}
	if pkg, ok := scope.Packages["fmt"]; !ok || len(pkg.Interfaces) != 1 || pkg.Interfaces[0].Name != "Stringer" {
		t.Fatalf("fmt.Stringer interface should be completed", pkg)
	}
	if pkg, ok := scope.Packages[""]; !ok || len(pkg.Interfaces) != 1 || pkg.Interfaces[0].Name != "error" {
		t.Fatalf(".error interface should be completed", pkg)
	}
	if pkg, ok := scope.Packages["os"]; !ok || len(pkg.Classes) != 1 || pkg.Classes[0].Name != "PathError" {
		t.Fatalf("os.PathError class should be completed", pkg)
	}
	pkg, ok := scope.Packages["pkg1"]
	if !ok || len(pkg.Classes) != 2 || pkg.Classes[1].Name != "XError" {
		t.Fatalf("pkg1.A and pkg1.XError class should be completed", pkg)
	}
	if cl := pkg.Classes[1]; len(cl.Methods) != 1 || cl.Methods[0].Name != "Error" {
		t.Fatalf("pkg1.XError methods should be completed", cl)
	}
}
