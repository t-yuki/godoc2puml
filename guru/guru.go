// Package guru is a wrapper.
package guru

import (
	"go/build"
	"io"

	"golang.org/x/tools/cmd/guru/serial"
)

// A Query specifies a single guru query.
type Query struct {
	Pos   string         // query position
	Build *build.Context // package loading configuration

	// pointer analysis options
	Scope      []string  // main packages in (*loader.Config).FromArgs syntax
	PTALog     io.Writer // (optional) pointer-analysis log file
	Reflection bool      // model reflection soundly (currently slow).
}

type QueryResult struct {
	Implements serial.Implements
}

// TODO: implement

func (q *Query) Serial() Serial {
	return &QueryResult{}
}

func Run(mode string, q *Query) error {
	return nil
}
