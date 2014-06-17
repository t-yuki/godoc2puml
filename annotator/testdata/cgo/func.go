package cgo

type Struct1 struct {
	string
}

func (*Struct1) String() string {
	return "struct1"
}

func (*Struct1) CGoFunc() int {
	return CGoFunc()
}
