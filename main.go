// Copyright (C) 2015, 2016 Yukinari Toyota. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"go/build"

	"github.com/t-yuki/godoc2puml/annotator"
	"github.com/t-yuki/godoc2puml/ast"
	"github.com/t-yuki/godoc2puml/parser"
	"github.com/t-yuki/godoc2puml/printer"
)

type Config struct {
	Help       bool
	Scope      string
	Lolipop    string
	Field      string
	Format     string
	Ignore     string
	DontIgnore string
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "show this help")
	flag.StringVar(&config.Scope, "scope", "", `set analysis scope (main) package. if it is omitted, scope is tests in the same directory`)
	flag.StringVar(&config.Lolipop, "lolipop", "", `set package names in comma-separated strings that use lolipop-style interface relationship instead of implementation`)
	flag.StringVar(&config.Field, "field", "", `set package names in comma-separated strings that use field relationship instead of association`)
	flag.StringVar(&config.Format, "t", "text", `output format
        puml:  write PlantUML format`)
	flag.StringVar(&config.Ignore, "ignore", "(fmt\\.Stringer|\\w+\\.[a-z][\\w]*)$", `name filter to ignore. default value removes fmt.String and private declarations except specified packages`)
	flag.StringVar(&config.DontIgnore, "dont-ignore", "", `white-list for ignore. default/empty value means packages of arg`)
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

	if config.DontIgnore == "" {
		v := strings.Replace(strings.Join(packages, "|"), ".", "\\.", -1)
		config.DontIgnore = "^(" + v + ")\\.(\\w+)$"
	}

	scope := ast.NewScope()

	for _, path := range packages {
		packageList := matchPackages(path)
		for _, path := range packageList {
			pkg, err := parser.ParsePackage(path, strings.Split(config.Field, ",")...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "package %s parse error:%#v\n", path, err)
				continue
			}

			err = annotator.Oracle(pkg, config.Scope)
			if err != nil {
				fmt.Fprintf(os.Stderr, "annotate error %s: %#v\n", path, err)
				continue
			}
			scope.Packages[pkg.Name] = pkg
		}
	}

	err := annotator.Complete(scope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "annotate error:%#v\n", err)
		return
	}

	err = annotator.Cut(scope)
	if err != nil {
		fmt.Fprintf(os.Stderr, "annotate error:%#v\n", err)
		return
	}

	err = annotator.Filter(scope, config.Ignore, config.DontIgnore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "filter error:%#v\n", err)
		return
	}

	printer.FprintPlantUML(os.Stdout, scope, strings.Split(config.Lolipop, ","))
}

// from go/src/cmd/go/main.go
// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

var (
	goroot       = filepath.Clean(runtime.GOROOT())
	gorootSrcPkg = filepath.Join(goroot, "src/pkg")
)

// hasPathPrefix reports whether the path s begins with the
// elements in prefix.
func hasPathPrefix(s, prefix string) bool {
	switch {
	default:
		return false
	case len(s) == len(prefix):
		return s == prefix
	case len(s) > len(prefix):
		if prefix != "" && prefix[len(prefix)-1] == '/' {
			return strings.HasPrefix(s, prefix)
		}
		return s[len(prefix)] == '/' && s[:len(prefix)] == prefix
	}
}

// treeCanMatchPattern(pattern)(name) reports whether
// name or children of name can possibly match pattern.
// Pattern is the same limited glob accepted by matchPattern.
func treeCanMatchPattern(pattern string) func(name string) bool {
	wildCard := false
	if i := strings.Index(pattern, "..."); i >= 0 {
		wildCard = true
		pattern = pattern[:i]
	}
	return func(name string) bool {
		return len(name) <= len(pattern) && hasPathPrefix(pattern, name) ||
			wildCard && strings.HasPrefix(name, pattern)
	}
}

func matchPackages(pattern string) []string {
	match := func(string) bool { return true }
	treeCanMatch := func(string) bool { return true }
	if pattern != "all" && pattern != "std" {
		match = matchPattern(pattern)
		treeCanMatch = treeCanMatchPattern(pattern)
	}

	have := map[string]bool{
		"builtin": true, // ignore pseudo-package that exists only for documentation
	}
	if !build.Default.CgoEnabled {
		have["runtime/cgo"] = true // ignore during walk
	}
	var pkgs []string

	for _, src := range build.Default.SrcDirs() {
		if pattern == "std" && src != gorootSrcPkg {
			continue
		}
		src = filepath.Clean(src) + string(filepath.Separator)
		filepath.Walk(src, func(path string, fi os.FileInfo, err error) error {
			if err != nil || !fi.IsDir() || path == src {
				return nil
			}

			// Avoid .foo, _foo, and testdata directory trees.
			_, elem := filepath.Split(path)
			if strings.HasPrefix(elem, ".") || strings.HasPrefix(elem, "_") || elem == "testdata" {
				return filepath.SkipDir
			}

			name := filepath.ToSlash(path[len(src):])
			if pattern == "std" && strings.Contains(name, ".") {
				return filepath.SkipDir
			}
			if !treeCanMatch(name) {
				return filepath.SkipDir
			}
			if have[name] {
				return nil
			}
			have[name] = true
			if !match(name) {
				return nil
			}
			_, err = build.Default.ImportDir(path, 0)
			if err != nil {
				if _, noGo := err.(*build.NoGoError); noGo {
					return nil
				}
			}
			pkgs = append(pkgs, name)
			return nil
		})
	}
	return pkgs
}

// matchPattern(pattern)(name) reports whether
// name matches pattern.  Pattern is a limited glob
// pattern in which '...' means 'any string' and there
// is no other special syntax.
func matchPattern(pattern string) func(name string) bool {
	re := regexp.QuoteMeta(pattern)
	re = strings.Replace(re, `\.\.\.`, `.*`, -1)
	// Special case: foo/... matches foo too.
	if strings.HasSuffix(re, `/.*`) {
		re = re[:len(re)-len(`/.*`)] + `(/.*)?`
	}
	reg := regexp.MustCompile(`^` + re + `$`)
	return func(name string) bool {
		return reg.MatchString(name)
	}
}
