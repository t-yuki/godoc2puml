package annotator

import "github.com/t-yuki/godoc2puml/ast"

func Cut(pkg *ast.Package) error {
	backproj := buildBackProjections(pkg)
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
	return nil
}

func buildBackProjections(scope *ast.Package) (backproj map[string][]string) {
	backproj = map[string][]string{}
	for _, iface := range scope.Interfaces {
		for _, rel := range iface.Relations {
			if rel.RelType != ast.Extension {
				continue
			}
			if backproj[rel.Target] == nil {
				backproj[rel.Target] = []string{}
			}
			backproj[rel.Target] = append(backproj[rel.Target], iface.Name)
		}
	}
	for _, class := range scope.Classes {
		for _, rel := range class.Relations {
			if rel.RelType != ast.Composition && rel.RelType != ast.Implementation {
				continue
			}
			if backproj[rel.Target] == nil {
				backproj[rel.Target] = []string{}
			}
			backproj[rel.Target] = append(backproj[rel.Target], class.Name)
		}
	}
	return
}

func isLongestPath(backproj map[string][]string, rootRels []*ast.Relation, goal *ast.Relation) bool {
	roots := map[string]*ast.Relation{}
	for _, rel := range rootRels {
		if rel == goal || (rel.RelType != ast.Implementation && rel.RelType != ast.Composition) {
			continue
		}
		roots[rel.Target] = rel
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
