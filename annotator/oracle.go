package annotator

import (
	"fmt"
	"go/build"
	"log"

	"code.google.com/p/go.tools/go/loader"
	"code.google.com/p/go.tools/oracle"
	"code.google.com/p/go.tools/oracle/serial"
	"github.com/t-yuki/godoc2puml/ast"
)

// Oracle annotates `pkg` using go.tools/oracle interface implements detector.
// It uses `scopes` as analysis scope.
// If `scopes` is none of one of `scopes` is zero string, it uses unit tests as scope.
func Oracle(pkg *ast.Package, scopes ...string) error {
	settings := build.Default
	settings.BuildTags = []string{} // TODO
	conf := loader.Config{Build: &settings, SourceImports: true}

	withTests := false
	if len(scopes) == 0 {
		withTests = true
	}
	for _, scope := range scopes {
		if scope == "" {
			withTests = true
		} else {
			conf.Import(scope)
		}
	}
	if withTests {
		conf.ImportWithTests(pkg.Name)
	} else {
		conf.Import(pkg.Name)
	}

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
			return fmt.Errorf("oracle annotator: query error: %+v, %v", err, class.Pos)
		}
		impls := res.Serial().Implements
		for _, target := range impls.AssignableFromPtr {
			addImplements(class, target)
		}
		for _, target := range impls.AssignableFrom {
			addImplements(class, target)
		}
	}
	return nil
}

func addImplements(class *ast.Class, impl serial.ImplementsType) {
	name := impl.Name
	switch impl.Name {
	case "runtime.stringer":
		return // ignore runtime.stringer because fmt.Stringer is more generic
	default:
	}

	rel := &ast.Relation{Target: name, RelType: ast.Implementation}
	class.Relations = append(class.Relations, rel)
}
