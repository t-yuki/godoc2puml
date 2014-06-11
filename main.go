package main

import (
	"flag"
	"fmt"
	"os"
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
}
