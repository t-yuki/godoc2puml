package annotator

import (
	"fmt"
	"go/build"
	"log"
	"strings"

	"code.google.com/p/go.tools/go/loader"
	"code.google.com/p/go.tools/oracle"
	"code.google.com/p/go.tools/oracle/serial"
	"github.com/t-yuki/godoc2puml/ast"
)

func Oracle(pkg *ast.Package) error {
	settings := build.Default
	settings.BuildTags = []string{} // TODO
	conf := loader.Config{Build: &settings, SourceImports: true}
	conf.ImportWithTests(pkg.Path)
	iprog, err := conf.Load()
	if err != nil {
		return fmt.Errorf("oracle annotator: conf load error: %+v", err)
	}
	o, err := oracle.New(iprog, nil, false)
	if err != nil {
		return fmt.Errorf("oracle annotator: create error: %+v", err)
	}
	for _, class := range pkg.Classes {
		qpos, err := oracle.ParseQueryPos(iprog, string(class.Pos), false)
		if err != nil {
			log.Printf("oracle annotator: parse query pos error: %+v, %+v", err, class.Pos)
			continue
		}

		res, err := o.Query("implements", qpos)
		if err != nil {
			return fmt.Errorf("oracle annotator: query error: %+v", err)
		}
		impls := res.Serial().Implements
		for _, target := range impls.AssignableFromPtr {
			addImplements(pkg, class, target)
		}
		for _, target := range impls.AssignableFrom {
			addImplements(pkg, class, target)
		}
	}
	return nil
}

func addImplements(pkg *ast.Package, class *ast.Class, impl serial.ImplementsType) {
	name := qualifiedName(pkg.Path, impl.Name)
	switch impl.Name {
	case "runtime.stringer":
		return // ignore runtime.stringer because fmt.Stringer is more generic
	case "error":
		name = ".error" // error is buildin - global scope
	default:
		compensateInterface(pkg, name)
	}

	rel := &ast.Relation{Target: name, RelType: ast.Implementation}
	class.Relations = append(class.Relations, rel)
}

func qualifiedName(basepath string, name string) string {
	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			if name[:i] == basepath {
				name = name[i+1:]
			}
			break
		}
	}
	return strings.Replace(name, "/", ".", -1)
}

func compensateInterface(pkg *ast.Package, name string) {
	for _, iface := range pkg.Interfaces {
		if iface.Name == name {
			return
		}
	}
	iface := &ast.Interface{
		Name:      name,
		Relations: make([]*ast.Relation, 0),
		Methods:   make([]*ast.Method, 0),
	}
	// TODO: append iface to proper package, not current `pkg`
	pkg.Interfaces = append(pkg.Interfaces, iface)
}
