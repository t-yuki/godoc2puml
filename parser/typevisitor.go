package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	. "github.com/t-yuki/godoc2puml/ast"
)

type typeVisitor struct {
	astPackage    *ast.Package
	astFile       *ast.File
	fileSet       *token.FileSet
	pkg           *Package
	name2class    map[string]*Class
	fieldPackages []string
}

func (v *typeVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.Package:
		v.astPackage = node
		return v
	case *ast.File:
		v.astFile = node
		return v
	case *ast.TypeSpec:
		v.visitTypeSpec(node)
	default:
		return v
	}
	return nil
}

func (v *typeVisitor) visitTypeSpec(node *ast.TypeSpec) {
	switch typeNode := node.Type.(type) {
	case *ast.StructType:
		cl := &Class{
			Name:      node.Name.Name,
			Relations: make([]*Relation, 0, 10),
			Methods:   make([]*Method, 0, 10),
			Pos:       toSourcePos(v.fileSet, node),
		}
		v.parseFields(cl, "", typeNode.Fields)
		v.pkg.Classes = append(v.pkg.Classes, cl)
		v.name2class[cl.Name] = cl
	case *ast.InterfaceType:
		iface := &Interface{
			Name:      node.Name.Name,
			Relations: make([]*Relation, 0, 10),
			Methods:   make([]*Method, 0, 10),
		}
		v.parseMethods(iface, typeNode.Methods)
		v.pkg.Interfaces = append(v.pkg.Interfaces, iface)
	default:
		return
	}
}

func (v *typeVisitor) parseFields(cl *Class, prefix string, fields *ast.FieldList) {
	for _, field := range fields.List {
		multiplicity := ""
		switch fType := field.Type.(type) {
		case *ast.ArrayType:
			multiplicity = "*"
			// TODO: we should nested array struct in runtime.Memstats like:
			// type A struct { X []struct{ bool }; Y []*struct{ *A }; Z *[]struct{ []*A } }
		case *ast.StarExpr:
			if fType, ok := fType.X.(*ast.StructType); ok { // case: type A struct { X *struct{ bool } }
				if len(field.Names) == 0 { // anonymous field
					panic("anonymous struct is prohibited in go spec")
				}
				for _, name := range field.Names {
					v.parseFields(cl, prefix+name.String()+".", fType.Fields)
				}
				continue
			}
		case *ast.StructType:
			if len(field.Names) == 0 { // anonymous field
				panic("anonymous struct is prohibited in go spec")
			}
			for _, name := range field.Names {
				v.parseFields(cl, prefix+name.String()+".", fType.Fields)
			}
			continue
		}
		elementType := v.elementType(field.Type)
		switch {
		case isPrimitive(elementType) || isField(v.fieldPackages, elementType):
			f := &Field{Type: typeGoString(field.Type), Multiplicity: multiplicity}

			if len(field.Names) == 0 { // anonymous field
				f.Name = prefix + elementType
				cl.Fields = append(cl.Fields, f)
			}
			for _, name := range field.Names {
				f2 := *f
				f2.Name = prefix + name.String()
				f2.Public = ast.IsExported(f2.Name)
				cl.Fields = append(cl.Fields, &f2)
			}
		default:
			rel := &Relation{Target: elementType, Multiplicity: multiplicity}

			if len(field.Names) == 0 { // anonymous field
				if prefix == "" {
					rel.RelType = Composition
				} else {
					rel.Label = prefix + path.Ext(elementType)[1:]
					rel.RelType = Association
				}
				cl.Relations = append(cl.Relations, rel)
			}
			for _, name := range field.Names {
				rel2 := *rel
				rel2.Label = prefix + name.String()
				rel2.RelType = Association
				cl.Relations = append(cl.Relations, &rel2)
			}
		}
	}
}

func (v *typeVisitor) parseMethods(iface *Interface, fields *ast.FieldList) {
	for _, field := range fields.List {
		switch typeNode := field.Type.(type) {
		case *ast.FuncType:
			if len(field.Names) != 1 {
				panic(fmt.Errorf("unexpected named fields in interface type %#v", field))
			}
			method := &Method{
				Name:      field.Names[0].String(),
				Arguments: make([]DeclPair, 0, 10),
				Results:   make([]DeclPair, 0, 10),
			}
			method.Public = ast.IsExported(method.Name)
			parseFuncType(method, typeNode)
			iface.Methods = append(iface.Methods, method)
		case *ast.Ident:
			if len(field.Names) != 0 {
				panic(fmt.Errorf("unexpected named fields in interface type %#v", field))
			}
			elementType := v.elementType(typeNode)
			rel := &Relation{Target: elementType, RelType: Extension}
			iface.Relations = append(iface.Relations, rel)
		case *ast.SelectorExpr:
			if len(field.Names) != 0 {
				panic(fmt.Errorf("unexpected named fields in interface type %#v", field))
			}
			elementType := v.elementType(typeNode)
			rel := &Relation{Target: elementType, RelType: Extension}
			iface.Relations = append(iface.Relations, rel)
		default:
			panic(fmt.Errorf("unexpected field type in interface type %#v on %+v", field, iface))
		}
	}
}

func (v *typeVisitor) elementType(expr ast.Node) string {
	if expr == nil {
		return ""
	}
	p := v.astPackage
	f := v.astFile
	switch expr := expr.(type) {
	case *ast.Ident:
		name := expr.String()
		if name == "error" || isPrimitive(name) {
			return name
		}
		for _, f := range p.Files {
			if f.Scope.Lookup(name) != nil {
				return v.pkg.Name + "." + name
			}
		}
		// TODO: refactor
		for _, imp := range f.Imports {
			local := imp.Name
			if local == nil || local.String() != `.` {
				continue
			}
			importPath, _ := strconv.Unquote(imp.Path.Value)
			if local == nil && path.Base(importPath) != name {
				continue
			}

			buildPkg, err := build.Import(importPath, ".", build.FindOnly)
			if err != nil {
				log.Printf("%#v", err)
				continue
			}
			dir := buildPkg.Dir

			fset := token.NewFileSet()
			pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
				return !fi.IsDir() && !strings.HasSuffix(fi.Name(), "_test.go")
			}, 0)
			if err != nil {
				log.Printf("%#v", err)
				continue
			}
			for _, pkg := range pkgs {
				for _, f := range pkg.Files {
					if f.Scope.Lookup(name) == nil {
						continue
					}
					break
				}
				break // break?
			}
			return importPath + "." + name
		}
		log.Printf("can't resolve name", name)
		return name
	case *ast.ArrayType:
		return v.elementType(expr.Elt)
	case *ast.StarExpr:
		return v.elementType(expr.X)
	case *ast.SelectorExpr:
		name := expr.X.(*ast.Ident).String()
		for _, imp := range f.Imports {
			local := imp.Name
			if local != nil && local.String() != name {
				continue
			}
			importPath, _ := strconv.Unquote(imp.Path.Value)
			if local == nil && path.Base(importPath) != name {
				continue
			}
			name = importPath
			break
		}
		return name + "." + expr.Sel.String()
	case *ast.FuncType:
		return strings.TrimSpace("func(" + v.elementType(expr.Params) + ") " + v.elementType(expr.Results))
	case *ast.FieldList:
		if expr == nil {
			return ""
		}
		var buf bytes.Buffer
		for _, field := range expr.List {
			buf.WriteString(v.elementType(field.Type))
		}
		return buf.String()
	case *ast.MapType:
		return "map[" + v.elementType(expr.Key) + "]" + v.elementType(expr.Value)
	case *ast.InterfaceType:
		return "interface {" + v.elementType(expr.Methods) + "}"
	case *ast.StructType:
		return "struct {" + v.elementType(expr.Fields) + "}"
	case *ast.ChanType:
		switch expr.Dir {
		case ast.SEND:
			return "chan<- " + v.elementType(expr.Value)
		case ast.RECV:
			return "<-chan " + v.elementType(expr.Value)
		default:
			return "chan " + v.elementType(expr.Value)
		}
	case *ast.Ellipsis:
		return "..." + v.elementType(expr.Elt)
	default:
		panic(fmt.Errorf("%#v", expr))
	}
}

func toSourcePos(fset *token.FileSet, node ast.Node) SourcePos {
	file := fset.File(node.Pos())
	start := file.Offset(node.Pos())
	end := file.Offset(node.End())
	return SourcePos(fmt.Sprintf("%s:#%d,#%d", file.Name(), start, end))
}

func isField(packages []string, name string) bool {
	actual := ""
	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			actual = name[:i]
		}
	}
	for _, pkg := range packages {
		if pkg == actual {
			return true
		}
	}
	return false
}
