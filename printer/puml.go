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

{{ range $p := .Packages }}
{{ range .Classes }}
	class {{ joinName $p.Name .Name }} {
{{ range .Fields}}
		{{ if .Public }}+{{ else }}~{{ end }}{{.Name}} {{.Type}}
{{ end }}
{{ range .Methods}}
		{{ if .Public }}+{{ else }}~{{ end }}{{.Name}}({{ methodArgs .Arguments}}) {{ methodResults .Results}}
{{ end }}
	}
{{ end }}
{{ range .Interfaces }}
	interface {{ joinName $p.Name .Name }} {
{{ range .Methods}}
		{{ if .Public }}+{{ else }}~{{ end }}{{.Name}}({{ methodArgs .Arguments}}) {{ methodResults .Results}}
{{ end }}
	}
{{ end }}

{{ range $cl := .Classes }} {{ range .Relations}}
	"{{ joinName $p.Name $cl.Name }}" {{relType .RelType}} {{if .Multiplicity}}"{{.Multiplicity}}" {{end}}"{{qualifiedName .Target}}" {{if .Label}}: {{.Label}}{{end}}
{{ end }} {{ end }}
{{ range $iface := .Interfaces }} {{ range .Relations}}
	"{{ joinName $p.Name $iface.Name }}" {{relType .RelType}} "{{qualifiedName .Target}}"
{{ end }} {{ end }}
{{ end }}

hide interface fields

@enduml
`))

var pumlFuncs = map[string]interface{}{
	"relType":       pumlRelType,
	"methodArgs":    pumlMethodArgs,
	"methodResults": pumlMethodResults,
	"qualifiedName": pumlQualifiedName,
	"joinName":      pumlJoinName,
}

func FprintPlantUML(w io.Writer, scope *ast.Scope) {
	err := pumlTemplate.Execute(w, scope)
	if err != nil {
		panic(err)
	}
}

func pumlRelType(relType ast.RelationType) string {
	switch relType {
	case ast.Association:
		return "-->"
	case ast.Extension:
		return "--|>"
	case ast.Composition:
		return "*--"
	case ast.Agregation:
		return "o--"
	case ast.Implementation:
		return "..|>" // lolipop style?: "-()"
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
