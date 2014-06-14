package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/parser"
	"github.com/t-yuki/godoc2puml/printer"
)

type Config struct {
	Help     bool
	Format   string
	Filter   string
	HTTPAddr string
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "show this help")
	flag.StringVar(&config.Format, "t", "text", `output format
        puml:  write PlantUML format`)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if config.Help {
		flag.Usage()
		return
	}
	packages := flag.Args()
	if len(packages) == 0 {
		panic("godoc2uml without explicit package path is not implemented yet")
	}
	if len(packages) != 1 {
		panic("godoc2uml with multiple packages is not implemented yet")
	}
	pkg, err := parser.ParsePackage(packages[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "package parse error:%#v", err)
		return
	}
	err = annotator.Oracle(pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "annotate error:%#v", err)
		return
	}
	err = annotator.Cut(pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "annotate error:%#v", err)
		return
	}
	printer.FprintPlantUML(os.Stdout, pkg)
}
