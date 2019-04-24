package annotator

import (
	"fmt"
	"go/build"
	"strings"

	"github.com/t-yuki/godoc2puml/ast"
	"github.com/t-yuki/godoc2puml/guru"
	"golang.org/x/tools/cmd/guru/serial"
	"golang.org/x/tools/go/loader"
)

// Oracle annotates `pkg` using go.tools/guru interface implements detector.
// It uses `scopes` as analysis scope.
// If `scopes` is none of one of `scopes` is zero string, it uses unit tests as scope.
func Oracle(pkg *ast.Package, scopes ...string) error {
	settings := build.Default
	settings.BuildTags = []string{} // TODO
	conf := loader.Config{Build: &settings}

	if len(scopes) == 0 || scopes[0] == "" {
		conf.ImportWithTests(pkg.Name)
		scopes = []string{pkg.Name, pkg.Name + "_test"}
	} else {
		for _, scope := range scopes {
			conf.Import(scope)
		}
		conf.Import(pkg.Name)
	}

	_, err := conf.Load()
	if err != nil {
		return fmt.Errorf("guru annotator: conf load error: %+v", err)
	}

	for _, class := range pkg.Classes {
		query := guru.Query{
			Mode:  "implements",
			Pos:   string(class.Pos),
			Build: conf.Build,
			Scope: scopes,
		}
		err = guru.Run(&query)
		if err != nil {
			return fmt.Errorf("guru annotator: query error: %+v, %v", err, class.Pos)
		}
		impls := query.Serial().Implements
		for _, target := range impls.AssignableFromPtr {
			if strings.Contains(target.Pos, "_test.go:") { // ignore relations to _test
				continue
			}
			addImplements(class, target)
		}
		for _, target := range impls.AssignableFrom {
			if strings.Contains(target.Pos, "_test.go:") { // ignore relations to _test
				continue
			}
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
