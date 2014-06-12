package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
	. "github.com/t-yuki/godoc2puml/ast"
)

func ParsePackage(packagePath string) (*Package, error) {
	p := &Package{}
	p.QualifiedName = strings.Replace(packagePath, "/", ".", -1)
	p.Classes = make([]*Class, 0, 10)

	buildPkg, err := build.Import(packagePath, ".", build.FindOnly)
	if err != nil {
		return nil, err
	}
	dir := buildPkg.Dir

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		return !fi.IsDir() && !strings.HasSuffix(fi.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		name2class := make(map[string]*Class)
		tv := &typeVisitor{pkg: p, name2class: name2class}
		ast.Walk(tv, pkg)
		mv := &methodVisitor{pkg: p, name2class: name2class}
		ast.Walk(mv, pkg)
	}
	return p, nil
}

func elementType(expr ast.Node) string {
	if expr == nil {
		return ""
	}
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.String()
	case *ast.ArrayType:
		return elementType(expr.Elt)
	case *ast.StarExpr:
		return elementType(expr.X)
	case *ast.SelectorExpr:
		return elementType(expr.X) + "." + expr.Sel.String()
	case *ast.FuncType:
		return "func " + elementType(expr.Params) + elementType(expr.Results)
	case *ast.FieldList:
		if expr == nil {
			return ""
		}
		var buf bytes.Buffer
		for _, field := range expr.List {
			buf.WriteString(elementType(field.Type))
		}
		return buf.String()
	case *ast.MapType:
		return "map[" + elementType(expr.Key) + "]" + elementType(expr.Value)
	case *ast.InterfaceType:
		return "interface {" + elementType(expr.Methods) + "}"
	case *ast.StructType:
		return "struct {" + elementType(expr.Fields) + "}"
	case *ast.ChanType:
		switch expr.Dir {
		case ast.SEND:
			return "chan<- " + elementType(expr.Value)
		case ast.RECV:
			return "<-chan " + elementType(expr.Value)
		default:
			return "chan " + elementType(expr.Value)
		}
	default:
		panic(fmt.Errorf("%#v", expr))
	}
}

func isPrimitive(name string) bool {
	switch name {
	case "bool", "int", "uint", "byte", "float",
		"uint8", "int8",
		"uint16", "int16",
		"uint32", "int32",
		"uint64", "int64",
		"string":
		return true
	default:
	}
	switch {
	case strings.ContainsAny(name, " ["):
		return true
	default:
		return false
	}
}

func isPublic(name string) bool {
	first, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(first)
}
