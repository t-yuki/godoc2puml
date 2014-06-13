package annotator

import (
	"fmt"
	"go/build"
	"log"
	"strings"

	"code.google.com/p/go.tools/go/loader"
	"code.google.com/p/go.tools/oracle"
	"github.com/t-yuki/godoc2puml/ast"
)

func Annotate(pkg *ast.Package) error {
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
			name := strings.Replace(target.Name, "/", ".", -1)
			rel := &ast.Relation{Target: name, RelType: ast.Implementation}
			class.Relations = append(class.Relations, rel)
			compensateInterface(pkg, name)
		}
		for _, target := range impls.AssignableFrom {
			name := strings.Replace(target.Name, "/", ".", -1)
			rel := &ast.Relation{Target: name, RelType: ast.Implementation}
			class.Relations = append(class.Relations, rel)
			compensateInterface(pkg, name)
		}
	}
	return nil
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
	pkg.Interfaces = append(pkg.Interfaces, iface)
}
