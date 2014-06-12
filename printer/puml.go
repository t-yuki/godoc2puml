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

namespace {{.QualifiedName}} {
{{ range .Classes }}
	class {{.Name}} {
{{ range .Fields}}
		{{ if .Public }}+{{ else }}~{{ end }}{{.Name}} {{.Type}}
{{ end }}
{{ range .Methods}}
		{{ if .Public }}+{{ else }}~{{ end }}{{.Name}}({{ methodArgs .Arguments}}) {{ methodResults .Results}}
{{ end }}
	}
{{ end }}

{{ range $cl := .Classes }} {{ range .Relations}}
	{{$cl.Name}} {{relType .RelType}} {{if .Multiplicity}}"{{.Multiplicity}}" {{end}}{{.Target}} {{if .Label}}: {{.Label}}{{end}}
{{ end }} {{ end }}
}

@enduml
`))

var pumlFuncs = map[string]interface{}{
	"relType":       pumlRelType,
	"methodArgs":    pumlMethodArgs,
	"methodResults": pumlMethodResults,
}

func FprintPlantUML(w io.Writer, pkg *ast.Package) {
	err := pumlTemplate.Execute(w, pkg)
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
