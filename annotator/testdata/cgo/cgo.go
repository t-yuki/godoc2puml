package cgo

// CGoFiles is not parsed by oracle

// #include "cgo.h"
import "C"

func CGoFunc() int {
	return int(C.CFunc())
}
