package ast

type Package struct {
	QualifiedName string
	Classes       []*Class
}

type Class struct {
	Name      string
	Methods   []*Method
	Fields    []*Field
	Relations []*Relation
}

type Method struct {
	Name      string
	Arguments []DeclPair
	Results   []DeclPair
	Public    bool
}

type DeclPair struct {
	Name string
	Type string
}

type Field struct {
	Name         string
	Type         string
	Multiplicity string
	Public       bool
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
