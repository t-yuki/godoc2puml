package annotator

import "github.com/t-yuki/godoc2puml/ast"

// Cut removes (probably) unnecessary relations preserving longest path.
func Cut(scope *ast.Scope) error {
	backproj := buildBackProjections(scope)
	for _, pkg := range scope.Packages {
		for _, class := range pkg.Classes {
			newrels := make([]*ast.Relation, 0, len(class.Relations))
			for _, rel := range class.Relations {
				if rel.RelType == ast.Implementation &&
					!isLongestPath(backproj, class.Relations, rel) {
					continue
				}
				newrels = append(newrels, rel)
			}
			class.Relations = newrels
		}
		for _, iface := range pkg.Interfaces {
			newrels := make([]*ast.Relation, 0, len(iface.Relations))
			for _, rel := range iface.Relations {
				if rel.RelType == ast.Extension &&
					!isLongestPath(backproj, iface.Relations, rel) {
					continue
				}
				newrels = append(newrels, rel)
			}
			iface.Relations = newrels
		}

	}
	return nil
}

func buildBackProjections(scope *ast.Scope) (backproj map[string][]string) {
	backproj = map[string][]string{}
	for _, pkg := range scope.Packages {
		for _, iface := range pkg.Interfaces {
			name := iface.Name
			if pkg.Name != "" {
				name = pkg.Name + "." + name
			}
			for _, rel := range iface.Relations {
				if rel.RelType != ast.Extension {
					continue
				}
				addPath(backproj, rel, name)
			}
		}
		for _, class := range pkg.Classes {
			name := class.Name
			if pkg.Name != "" {
				name = pkg.Name + "." + name
			}
			for _, rel := range class.Relations {
				if rel.RelType != ast.Composition && rel.RelType != ast.Implementation {
					continue
				}
				addPath(backproj, rel, name)
			}
		}
	}
	return
}

func addPath(backproj map[string][]string, to *ast.Relation, from string) {
	target := to.Target
	if backproj[target] == nil {
		backproj[target] = []string{}
	}
	backproj[target] = append(backproj[target], from)
}

func isLongestPath(backproj map[string][]string, rootRels []*ast.Relation, goal *ast.Relation) bool {
	roots := map[string]*ast.Relation{}
	for _, rel := range rootRels {
		if rel == goal || (rel.RelType != ast.Implementation && rel.RelType != ast.Composition) {
			continue
		}
		target := rel.Target
		roots[target] = rel
	}
	return !findRouteToGoalRecursive(backproj, roots, goal.Target)
}

func findRouteToGoalRecursive(backproj map[string][]string, roots map[string]*ast.Relation, goal string) (reachable bool) {
	names := backproj[goal]
	for _, name := range names {
		if roots[name] != nil {
			return true
		}
		if findRouteToGoalRecursive(backproj, roots, name) {
			return true
		}
	}
	return false
}
