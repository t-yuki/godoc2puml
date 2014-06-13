package parser

import (
	"go/ast"
	. "github.com/t-yuki/godoc2puml/ast"
)

type methodVisitor struct {
	pkg        *Package
	name2class map[string]*Class
}

func (v *methodVisitor) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.FuncDecl:
		v.visitFuncDecl(node)
	default:
		return v
	}
	return nil
}

func (v *methodVisitor) visitFuncDecl(node *ast.FuncDecl) {
	if node.Recv == nil {
		return
	}
	recv := node.Recv.List[0]
	typeName := elementType(recv.Type)
	class, ok := v.name2class[typeName]
	if !ok {
		// unknown method receiver
		return
	}
	method := &Method{
		Name:      node.Name.String(),
		Arguments: make([]DeclPair, 0, 10),
		Results:   make([]DeclPair, 0, 10),
	}
	method.Public = isPublic(method.Name)
	parseFuncType(method, node.Type)
	class.Methods = append(class.Methods, method)
}

func parseFuncType(method *Method, node *ast.FuncType) {
	for _, field := range node.Params.List {
		argType := typeGoString(field.Type)
		if len(field.Names) == 0 {
			method.Arguments = append(method.Arguments, DeclPair{"", argType})
		}
		for _, name := range field.Names {
			method.Arguments = append(method.Arguments, DeclPair{name.String(), argType})
		}
	}
	if node.Results != nil {
		for _, field := range node.Results.List {
			argType := typeGoString(field.Type)
			if len(field.Names) == 0 {
				method.Results = append(method.Results, DeclPair{"", argType})
			}
			for _, name := range field.Names {
				method.Results = append(method.Results, DeclPair{name.String(), argType})
			}
		}
	}
}
