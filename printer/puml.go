package printer

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/t-yuki/godoc2puml/ast"
)

var pumlTemplate = template.Must(template.New("plantuml").Funcs(pumlFuncs).Parse(`
@startuml

set namespaceSeparator /

{{- $lolipop := .Lolipop }}

{{- range $p := .Packages -}}
  {{- range .Classes }}

class {{ joinName $p.Name .Name }} {
    {{- range .Fields }}
	{{ if .Public }}+{{ else }}~{{ end }}{{ .Name }} {{ .Type -}}
    {{- end }}
    {{- range .Methods }}
	{{ if .Public }}+{{ else }}~{{ end }}{{ .Name }}({{ methodArgs .Arguments }})
	{{- if .Results }} {{ methodResults .Results }}{{ end }}
    {{- end }}
}
  {{- end }}

  {{- range .Interfaces }}

interface {{ joinName $p.Name .Name }} {
    {{- range .Methods }}
	{{ if .Public }}+{{ else }}~{{ end }}{{ .Name }}({{ methodArgs .Arguments }})
	{{- if .Results }} {{ methodResults .Results }}{{ end }}
    {{- end }}
}
  {{- end }}

  {{- range $cl := .Classes }}
    {{- range .Relations }}
"{{ joinName $p.Name $cl.Name }}" {{ relType .RelType (isLolipop $lolipop .Target) }}
{{- if .Multiplicity }} "{{- .Multiplicity }}" {{ end }} "{{ qualifiedName .Target }}"
{{- if .Label }}: {{ .Label }}{{ end -}}
    {{- end }}
  {{- end }}
  {{- range $iface := .Interfaces }}
    {{- range .Relations }}
"{{ joinName $p.Name $iface.Name }}" {{ relType .RelType false }} "{{ qualifiedName .Target -}}"
    {{- end }}
  {{- end }}
{{- end }}

hide interface fields

@enduml
`))

var pumlFuncs = map[string]interface{}{
	"relType":       pumlRelType,
	"methodArgs":    pumlMethodArgs,
	"methodResults": pumlMethodResults,
	"qualifiedName": pumlQualifiedName,
	"joinName":      pumlJoinName,
	"isLolipop":     pumlIsLolipop,
}

func FprintPlantUML(w io.Writer, scope *ast.Scope, lolipopPackages []string) {
	packages := make(ast.PackageSlice, 0, len(scope.Packages))
	for _, p := range scope.Packages {
		packages = append(packages, p)
	}
	packages.Sort()

	err := pumlTemplate.Execute(w, map[string]interface{}{"Packages": packages, "Lolipop": lolipopPackages})
	if err != nil {
		panic(err)
	}
}

func pumlRelType(relType ast.RelationType, lolipop bool) string {
	switch relType {
	case ast.Association:
		return "->"
	case ast.Extension:
		return "-|>"
	case ast.Composition:
		return "*-"
	case ast.Agregation:
		return "o-"
	case ast.Implementation:
		if lolipop {
			return "-()"
		}
		return ".|>"
	}
	panic(relType)
}

func pumlMethodArgs(decls []ast.DeclPair) string {
	b := &bytes.Buffer{}
	for i, v := range decls {
		if i != 0 {
			b.WriteString(", ")
		}
		if v.Name == "" {
			fmt.Fprintf(b, "%s", v.Type)
		} else {
			fmt.Fprintf(b, "%s %s", v.Name, v.Type)
		}
	}
	return b.String()
}

func pumlMethodResults(decls []ast.DeclPair) string {
	if len(decls) >= 2 {
		return "(" + pumlMethodArgs(decls) + ")"
	}
	return pumlMethodArgs(decls)
}

func pumlQualifiedName(name string) string {
	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			return name[:i] + "/" + name[i+1:]
		}
	}
	return name
}

func pumlJoinName(name1 string, name2 string) string {
	if name1 == "" {
		return name2
	}
	if name2 == "" {
		return name1
	}
	return name1 + "/" + name2
}

func pumlIsLolipop(packages []string, name string) bool {
	actual := ""
	for i := len(name) - 1; i >= 0 && name[i] != '/'; i-- {
		if name[i] == '.' {
			actual = name[:i]
		}
	}
	for _, pkg := range packages {
		if pkg == actual {
			return true
		}
	}
	return false
}
