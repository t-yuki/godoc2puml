package parser

import (
	"go/ast"
	. "github.com/t-yuki/godoc2puml/ast"
)

type typeVisitor struct {
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
	st, ok2 := node.Type.(*ast.StructType)
	if !ok2 {
		return
	}

	cl := &Class{
		Name:      node.Name.Name,
		Relations: make([]*Relation, 0, 10),
		Methods:   make([]*Method, 0, 10),
	}
	parseFields(cl, st.Fields)
	v.pkg.Classes = append(v.pkg.Classes, cl)
	v.name2class[cl.Name] = cl
}

func parseFields(cl *Class, fields *ast.FieldList) {
	for _, field := range fields.List {
		multiplicity := ""
		if _, ok := field.Type.(*ast.ArrayType); ok {
			multiplicity = "0..*"
		}
		elementType := elementType(field.Type)
		switch {
		case isPrimitive(elementType):
			f := &Field{Type: elementType, Multiplicity: multiplicity}

			if len(field.Names) == 0 { // anonymous field
				cl.Fields = append(cl.Fields, f)
			}
			for _, name := range field.Names {
				f2 := *f
				f2.Name = name.String()
				f.Public = isPublic(f2.Name)
				cl.Fields = append(cl.Fields, &f2)
			}
		default:
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
