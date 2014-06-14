package annotator_test

import (
	"testing"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/ast"
)

func TestCut(t *testing.T) {
	pkg := &ast.Package{
		QualifiedName: "pkg1",
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
					{Target: "XError", RelType: ast.Composition},
					{Target: ".error", RelType: ast.Implementation},
				},
			},
			{
				Name: "XError",
				Relations: []*ast.Relation{
					{Target: ".error", RelType: ast.Implementation},
				},
			},
		},
		Interfaces: []*ast.Interface{
			{Name: "io.Reader", Relations: []*ast.Relation{}},
			{Name: "io.Writer", Relations: []*ast.Relation{}},
			{Name: "io.Closer", Relations: []*ast.Relation{}},
			{Name: "io.ReadWriter", Relations: []*ast.Relation{
				{Target: "io.Reader", RelType: ast.Extension},
				{Target: "io.Writer", RelType: ast.Extension},
			}},
			{Name: "io.ReadCloser", Relations: []*ast.Relation{
				{Target: "io.Reader", RelType: ast.Extension},
				{Target: "io.Closer", RelType: ast.Extension},
			}},
			{Name: "io.WriteCloser", Relations: []*ast.Relation{
				{Target: "io.Writer", RelType: ast.Extension},
				{Target: "io.Closer", RelType: ast.Extension},
			}},
			{Name: "io.ReadWriteCloser", Relations: []*ast.Relation{
				{Target: "io.ReadWriter", RelType: ast.Extension},
				{Target: "io.WriteCloser", RelType: ast.Extension},
				{Target: "io.ReadCloser", RelType: ast.Extension},
			}},
			{Name: "fmt.Stringer", Relations: []*ast.Relation{}},
			{Name: ".error", Relations: []*ast.Relation{}},
		},
	}
	err := annotator.Cut(pkg)
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.Classes[0].Relations) != 3 {
		t.Fatalf("io.ReadWriteCloser, fmt.Stringer and XError should be preserved: %#v", pkg.Classes[0].Relations)
	}
	if pkg.Classes[0].Relations[0].Target != "io.ReadWriteCloser" {
		t.Fatalf("other io.* should be omitted: %#v", pkg.Classes[0].Relations[0])
	}
	if pkg.Classes[0].Relations[1].Target != "fmt.Stringer" {
		t.Fatalf("fmt.Stringer should be preserved: %#v", pkg.Classes[0].Relations[1])
	}
	if pkg.Classes[0].Relations[2].Target != "XError" {
		t.Fatalf("XError should be preserved: %#v", pkg.Classes[0].Relations[2])
	}
	if len(pkg.Classes[1].Relations) != 1 {
		t.Fatalf(".error should be preserved: %#v", pkg.Classes[1].Relations)
	}
}
