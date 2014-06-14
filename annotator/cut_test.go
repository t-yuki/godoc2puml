package annotator_test

import (
	"testing"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/ast"
)

func TestCut(t *testing.T) {
	scope := ast.NewScope()
	pkg := &ast.Package{
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
					{Target: "fmt.Stringer", RelType: ast.Implementation},
					{Target: "pkg1.XError", RelType: ast.Composition},
					{Target: "error", RelType: ast.Implementation},
				},
			},
			{
				Name: "XError",
				Relations: []*ast.Relation{
					{Target: "error", RelType: ast.Implementation},
				},
			},
		},
	}
	scope.Packages["pkg1"] = pkg
	scope.Packages["io"] = &ast.Package{
		Name:    "io",
		Classes: []*ast.Class{},
		Interfaces: []*ast.Interface{
			{Name: "Reader", Relations: []*ast.Relation{}},
			{Name: "Writer", Relations: []*ast.Relation{}},
			{Name: "Closer", Relations: []*ast.Relation{}},
			{Name: "ReadWriter", Relations: []*ast.Relation{
				{Target: "io.Reader", RelType: ast.Extension},
				{Target: "io.Writer", RelType: ast.Extension},
			}},
			{Name: "ReadCloser", Relations: []*ast.Relation{
				{Target: "io.Reader", RelType: ast.Extension},
				{Target: "io.Closer", RelType: ast.Extension},
			}},
			{Name: "WriteCloser", Relations: []*ast.Relation{
				{Target: "io.Writer", RelType: ast.Extension},
				{Target: "io.Closer", RelType: ast.Extension},
			}},
			{Name: "ReadWriteCloser", Relations: []*ast.Relation{
				{Target: "io.ReadWriter", RelType: ast.Extension},
				{Target: "io.WriteCloser", RelType: ast.Extension},
				{Target: "io.ReadCloser", RelType: ast.Extension},
			}},
		},
	}
	scope.Packages["fmt"] = &ast.Package{
		Name:    "fmt",
		Classes: []*ast.Class{},
		Interfaces: []*ast.Interface{
			{Name: "Stringer", Relations: []*ast.Relation{}},
		},
	}
	scope.Packages[""] = &ast.Package{
		Name:    "",
		Classes: []*ast.Class{},
		Interfaces: []*ast.Interface{
			{Name: "error", Relations: []*ast.Relation{}},
		},
	}
	err := annotator.Cut(scope)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.Classes[0].Relations) != 3 {
		for _, rel := range pkg.Classes[0].Relations {
			t.Logf("%+v", rel)
		}
		t.Fatalf("io.ReadWriteCloser, fmt.Stringer and XError should be preserved: %#v", pkg.Classes[0].Relations)
	}
	if pkg.Classes[0].Relations[0].Target != "io.ReadWriteCloser" {
		t.Fatalf("other io.* should be omitted: %#v", pkg.Classes[0].Relations[0])
	}
	if pkg.Classes[0].Relations[1].Target != "fmt.Stringer" {
		t.Fatalf("fmt.Stringer should be preserved: %#v", pkg.Classes[0].Relations[1])
	}
	if pkg.Classes[0].Relations[2].Target != "pkg1.XError" {
		t.Fatalf("XError should be preserved: %#v", pkg.Classes[0].Relations[2])
	}
	if len(pkg.Classes[1].Relations) != 1 {
		t.Fatalf(".error should be preserved: %#v", pkg.Classes[1].Relations)
	}
}
