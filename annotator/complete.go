package annotator

import (
	"log"

	"github.com/t-yuki/godoc2puml/ast"
	"github.com/t-yuki/godoc2puml/parser"
)

func Complete(scope *ast.Scope) error {
	packageRefs := map[string]map[string]bool{}
	for _, pkg := range scope.Packages {
		for _, class := range pkg.Classes {
			for _, rel := range class.Relations {
				addReference(packageRefs, rel)
			}
		}
		for _, iface := range pkg.Interfaces {
			for _, rel := range iface.Relations {
				addReference(packageRefs, rel)
			}
		}
	}
	for _, pkg := range scope.Packages {
		delete(packageRefs, pkg.Name)
	}

	for path, refs := range packageRefs {
		if path == "" {
			pkg := ast.NewPackage("")
			pkg.Interfaces = append(pkg.Interfaces,
				&ast.Interface{Name: "error", Relations: []*ast.Relation{}})
			scope.Packages[path] = pkg
			continue
		}
		pkg, err := parser.ParsePackage(path)
		if err != nil {
			log.Printf("Complete ParsePackage err: %#v", err)
			continue
		}

		newclasses := make([]*ast.Class, 0, len(pkg.Classes))
		for _, class := range pkg.Classes {
			if refs[path+"."+class.Name] {
				// TODO: hide details, really?
				class.Fields = []*ast.Field{}
				class.Methods = []*ast.Method{}
				class.Relations = []*ast.Relation{}
				newclasses = append(newclasses, class)
			}
		}
		pkg.Classes = newclasses

		newifaces := make([]*ast.Interface, 0, len(pkg.Interfaces))
		for _, iface := range pkg.Interfaces {
			if refs[path+"."+iface.Name] {
				iface.Methods = []*ast.Method{}
				iface.Relations = []*ast.Relation{}
				newifaces = append(newifaces, iface)
			}
		}
		pkg.Interfaces = newifaces
		scope.Packages[path] = pkg
	}
	return nil
}

func addReference(refs map[string]map[string]bool, rel *ast.Relation) {
	name := rel.Target
	pkg := packageName(name)
	if refs[pkg] == nil {
		refs[pkg] = make(map[string]bool)
	}
	refs[pkg][name] = true
}
