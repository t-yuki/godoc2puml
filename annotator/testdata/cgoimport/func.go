package cgoimport

import "github.com/t-yuki/godoc2puml/annotator/testdata/cgo"

type Struct1 struct {
	string
}

func (*Struct1) String() string {
	return "struct1"
}

func (*Struct1) CGoFunc() int {
	return cgo.CGoFunc()
}
