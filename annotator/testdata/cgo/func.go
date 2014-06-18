package cgo

type Struct1 struct {
	string
	cgotype
}

func (*Struct1) String() string {
	return "struct1"
}

func (*Struct1) CGoFunc() int {
	return CGoFunc()
}

func (s *Struct1) N1() int {
	return s.cgotype.N1()
}
