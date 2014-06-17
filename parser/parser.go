package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	. "github.com/t-yuki/godoc2puml/ast"
)

func ParsePackage(packagePath string) (*Package, error) {
	fset, pkg, err := importPackage(packagePath)
	if err != nil {
		return nil, err
	}
	p := NewPackage(packagePath)
	err = parseGoAST(p, fset, pkg)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func importPackage(path string) (*token.FileSet, *ast.Package, error) {
	buildPkg, err := build.Import(path, ".", 0)
	if err != nil {
		return nil, nil, err
	}
	dir := buildPkg.Dir

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		for _, ignored := range buildPkg.IgnoredGoFiles {
			if fi.Name() == ignored {
				return false
			}
		}
		for _, gofile := range buildPkg.GoFiles { // GoFiles doesn't contain tests
			if fi.Name() == gofile {
				return true
			}
		}
		// we can also parse buildPkg.CgoFiles
		// however, due to oracle can't parse CgoFiles, we don't parse them
		// so any cgo user should not contain type and interfaces in CgoFiles that is `import "C"` is declared
		return false
	}, 0)
	if err != nil {
		return nil, nil, err
	}
	if len(pkgs) != 1 {
		return nil, nil, fmt.Errorf("package %s contains %d packges, it must be 1", path, len(pkgs))
	}
	for _, pkg := range pkgs {
		return fset, pkg, nil
	}
	panic("unreachable code")
}

func parseGoAST(p *Package, fset *token.FileSet, pkg *ast.Package) error {
	name2class := make(map[string]*Class)
	tv := &typeVisitor{pkg: p, name2class: name2class, fileSet: fset}
	mv := &methodVisitor{pkg: p, name2class: name2class}
	ast.Walk(tv, pkg)
	ast.Walk(mv, pkg)
	for _, class := range p.Classes {
		sort.Sort(fieldSorter(class.Fields))
		sort.Sort(methodSorter(class.Methods))
	}
	for _, iface := range p.Interfaces {
		sort.Sort(methodSorter(iface.Methods))
	}
	return nil
}

func typeGoString(expr ast.Node) string {
	if expr == nil {
		return ""
	}
	// TODO: refactor using go/printer
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.String()
	case *ast.ArrayType:
		return "[]" + typeGoString(expr.Elt)
	case *ast.StarExpr:
		return "*" + typeGoString(expr.X)
	case *ast.SelectorExpr:
		return typeGoString(expr.X) + "." + expr.Sel.String()
	case *ast.FuncType:
		pre := "func(" + typeGoString(expr.Params) + ")"
		if expr.Results == nil || expr.Results.List == nil {
			return pre
		}
		post := typeGoString(expr.Results)
		if len(expr.Results.List) == 1 {
			return pre + " " + post
		}
		return pre + " (" + post + ")"
	case *ast.FieldList:
		if expr == nil {
			return ""
		}
		var buf bytes.Buffer
		for i, field := range expr.List {
			if i != 0 {
				buf.WriteString(", ")
			}
			if len(field.Names) == 0 {
				buf.WriteString(typeGoString(field.Type))
			}
			for i, name := range field.Names {
				if i != 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(name.String())
				buf.WriteString(" ")
				buf.WriteString(typeGoString(field.Type))
			}
		}
		return buf.String()
	case *ast.MapType:
		return "map[" + typeGoString(expr.Key) + "]" + typeGoString(expr.Value)
	case *ast.InterfaceType:
		return "interface{" + typeGoString(expr.Methods) + "}"
	case *ast.StructType:
		return "struct{" + typeGoString(expr.Fields) + "}"
	case *ast.ChanType:
		switch expr.Dir {
		case ast.SEND:
			return "chan<- " + typeGoString(expr.Value)
		case ast.RECV:
			return "<-chan " + typeGoString(expr.Value)
		default:
			return "chan " + typeGoString(expr.Value)
		}
	case *ast.Ellipsis:
		return "..." + typeGoString(expr.Elt)
	default:
		panic(fmt.Errorf("%#v", expr))
	}
}

func isPrimitive(name string) bool {
	switch name {
	case "bool", "int", "uint", "byte", "rune", "float",
		"uint8", "int8",
		"uint16", "int16",
		"uint32", "int32",
		"uint64", "int64",
		"float32", "float64",
		"complex64", "complex128",
		"uintptr",
		"string":
		return true
	default:
	}
	switch {
	case strings.ContainsAny(name, " [({"):
		return true
	default:
		return false
	}
}

type fieldSorter []*Field

func (s fieldSorter) Len() int {
	return len(s)
}

func (s fieldSorter) Less(i int, j int) bool {
	if s[i].Public != s[j].Public {
		return s[i].Public
	}
	return s[i].Name < s[j].Name
}

func (s fieldSorter) Swap(i int, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}

type methodSorter []*Method

func (s methodSorter) Len() int {
	return len(s)
}

func (s methodSorter) Less(i int, j int) bool {
	if s[i].Public != s[j].Public {
		return s[i].Public
	}
	return s[i].Name < s[j].Name
}

func (s methodSorter) Swap(i int, j int) {
	t := s[i]
	s[i] = s[j]
	s[j] = t
}
