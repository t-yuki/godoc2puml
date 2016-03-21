package annotator

import (
	"regexp"

	"github.com/t-yuki/godoc2puml/ast"
)

// Filter removes nodes matched to `regexString`.
func Filter(scope *ast.Scope, regexString string, whitelistString string) error {
	regex, err := regexp.Compile(regexString)
	if err != nil {
		return err
	}
	whitelist, err := regexp.Compile(whitelistString)
	if err != nil {
		return err
	}

	for _, pkg := range scope.Packages {
		newclasses := make([]*ast.Class, 0, len(pkg.Classes))
		for _, class := range pkg.Classes {
			if regex.MatchString(pkg.Name+"."+class.Name) &&
				!whitelist.MatchString(pkg.Name+"."+class.Name) {
				continue
			}
			newrels := make([]*ast.Relation, 0, len(class.Relations))
			for _, rel := range class.Relations {
				if regex.MatchString(rel.Target) &&
					!whitelist.MatchString(rel.Target) {
					continue
				}
				newrels = append(newrels, rel)
			}
			class.Relations = newrels
			newclasses = append(newclasses, class)
		}
		pkg.Classes = newclasses

		newifaces := make([]*ast.Interface, 0, len(pkg.Interfaces))
		for _, iface := range pkg.Interfaces {
			if regex.MatchString(pkg.Name+"."+iface.Name) &&
				!whitelist.MatchString(pkg.Name+"."+iface.Name) {
				continue
			}
			newrels := make([]*ast.Relation, 0, len(iface.Relations))
			for _, rel := range iface.Relations {
				if regex.MatchString(rel.Target) &&
					!whitelist.MatchString(rel.Target) {
					continue
				}
				newrels = append(newrels, rel)
			}
			iface.Relations = newrels
			newifaces = append(newifaces, iface)
		}
		pkg.Interfaces = newifaces
	}
	return nil
}
