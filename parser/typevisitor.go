package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	. "github.com/t-yuki/godoc2puml/ast"
)

type typeVisitor struct {
	fileSet    *token.FileSet
	pkg        *Package
	name2class map[string]*Class
}

func (v *typeVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
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
		v.parseFields(cl, typeNode.Fields)
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

func (v *typeVisitor) parseFields(cl *Class, fields *ast.FieldList) {
	for _, field := range fields.List {
		multiplicity := ""
		if _, ok := field.Type.(*ast.ArrayType); ok {
			multiplicity = "*"
		}
		elementType := elementType(field.Type)
		switch {
		case isPrimitive(elementType):
			f := &Field{Type: typeGoString(field.Type), Multiplicity: multiplicity}

			if len(field.Names) == 0 { // anonymous field
				cl.Fields = append(cl.Fields, f)
			}
			for _, name := range field.Names {
				f2 := *f
				f2.Name = name.String()
				f2.Public = isPublic(f2.Name)
				cl.Fields = append(cl.Fields, &f2)
			}
		default:
			if elementType == "error" {
				elementType = ".error"
			}
			rel := &Relation{Target: elementType, Multiplicity: multiplicity}

			if len(field.Names) == 0 { // anonymous field
				rel.RelType = Composition
				cl.Relations = append(cl.Relations, rel)
			}
			for _, name := range field.Names {
				rel2 := *rel
				rel2.Label = name.String()
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
			method.Public = isPublic(method.Name)
			parseFuncType(method, typeNode)
			iface.Methods = append(iface.Methods, method)
		case *ast.Ident:
			if len(field.Names) != 0 {
				panic(fmt.Errorf("unexpected named fields in interface type %#v", field))
			}
			elementType := typeNode.String()
			if elementType == "error" {
				elementType = ".error"
			}
			rel := &Relation{Target: elementType, RelType: Extension}
			iface.Relations = append(iface.Relations, rel)
		case *ast.SelectorExpr:
			if len(field.Names) != 0 {
				panic(fmt.Errorf("unexpected named fields in interface type %#v", field))
			}
			elementType := elementType(typeNode)
			compensateInterface(v.pkg, elementType)
			rel := &Relation{Target: elementType, RelType: Extension}
			iface.Relations = append(iface.Relations, rel)
		default:
			panic(fmt.Errorf("unexpected field type in interface type %#v on %+v", field, iface))
		}
	}
}

func toSourcePos(fset *token.FileSet, node ast.Node) SourcePos {
	file := fset.File(node.Pos())
	start := file.Offset(node.Pos())
	end := file.Offset(node.End())
	return SourcePos(fmt.Sprintf("%s:#%d,#%d", file.Name(), start, end))
}

func compensateInterface(pkg *Package, name string) {
	for _, iface := range pkg.Interfaces {
		if iface.Name == name {
			return
		}
	}
	iface := &Interface{
		Name:      name,
		Relations: make([]*Relation, 0),
		Methods:   make([]*Method, 0),
	}
	pkg.Interfaces = append(pkg.Interfaces, iface)
}
