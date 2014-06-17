package cgoimport

import "testing"

func TestCGoStruct(t *testing.T) {
	x := &Struct1{}
	if x.CGoFunc() != 1 {
		t.Fatal(x.CGoFunc())
	}
}
