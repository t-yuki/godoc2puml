package ast

type Package struct {
	QualifiedName string
	Classes       []Class
}

type Class struct {
	Name      string
	Relations []Relation
	Fields    []Field
}

type Relation struct {
	Target       string
	RelType      RelationType
	Label        string
	Multiplicity string
	Composition  bool
}

type RelationType string

const (
	Association RelationType = "association"
	Extension   RelationType = "extension"
	Composition RelationType = "composition"
	Agregation  RelationType = "agregation"
)

type Field struct {
	Name         string
	Type         string
	Multiplicity string
}
