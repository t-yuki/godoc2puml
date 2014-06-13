package ast

type Package struct {
	Path          string
	QualifiedName string
	Classes       []*Class
	Interfaces    []*Interface
}

type Class struct {
	Name      string
	Methods   []*Method
	Fields    []*Field
	Relations []*Relation
	Pos       SourcePos
}

type Interface struct {
	Name      string
	Methods   []*Method
	Relations []*Relation
}

type Method struct {
	Name      string
	Arguments []DeclPair
	Results   []DeclPair
	Public    bool
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
}

type RelationType string

const (
	Association    RelationType = "association"
	Extension      RelationType = "extension"
	Composition    RelationType = "composition"
	Agregation     RelationType = "agregation"
	Implementation RelationType = "implementation"
)

type DeclPair struct {
	Name string
	Type string
}

// SourcePos is source file position as `file:#start,#end` form.
// For more details, see http://golang.org/s/oracle-user-manual
type SourcePos string
