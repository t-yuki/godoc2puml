package cgo

// CGoFiles is not parsed by oracle

// #include "cgo.h"
import "C"

type cgotype struct {
	n C.int
}

func CGoFunc() int {
	return int(C.CFunc())
}

func (c *cgotype) N1() int {
	return int(c.n)
}

func (s *Struct1) N2() int {
	return int(s.cgotype.n)
}
