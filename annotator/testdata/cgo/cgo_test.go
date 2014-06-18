package cgo

import "testing"
import _ "fmt"

func TestCGoStruct(t *testing.T) {
	x := &Struct1{}
	if x.CGoFunc() != 1 {
		t.Fatal(x.CGoFunc())
	}
	x.N1()
	x.N2()
}
