package ast

type Package struct {
	QualifiedName string
	Classes       []Class
}

type Class struct {
	Name      string
	Relations []Relation
}

type Relation struct {
	Target       string
	Label        string
	Multiplicity string
}
