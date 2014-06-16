package parser_test

import (
	"testing"

	"github.com/t-yuki/godoc2puml/ast"
	"github.com/t-yuki/godoc2puml/parser"
)

func TestParsePackage(t *testing.T) {
	pkg, err := parser.ParsePackage(".")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", pkg)
}

func TestParsePackage_NoFiles(t *testing.T) {
	pkg, err := parser.ParsePackage("./testdata/nofiles")
	if err == nil {
		t.Fatalf("it should be error %+v", pkg)
	}
}

func TestParsePackage_BuildConstrains(t *testing.T) {
	pkg, err := parser.ParsePackage("./testdata/buildconstrains")
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.Classes) != 1 || pkg.Classes[0].Name != "ok" {
		t.Fatalf("it should contain ok but:%+v", pkg.Classes)
	}
}

func TestParsePackage_BuildTags(t *testing.T) {
	pkg, err := parser.ParsePackage("./testdata/buildtags")
	if err != nil {
		t.Fatal(err)
	}
	if len(pkg.Classes) != 1 || pkg.Classes[0].Name != "ok" {
		t.Fatalf("it should contain ok but:%+v", pkg.Classes)
	}
}

func TestParsePackage_Fields(t *testing.T) {
	t.Skip("not implemented yet")
	pkg, err := parser.ParsePackage("./testdata/fields")
	if err != nil {
		t.Fatal(err)
	}
	fields := pkg.Classes[0].Fields

	table := []struct {
		Name string
		Type string
	}{
		{"uint8", "uint8"},
		{"i8", "int8"},
		{"embed.int16", "*int16"},
		{"embed.i32", "[]int32"},
		{"embed.embed.nested.i64", "int64"},
	}

	for i, v := range table {
		field := fields[i]
		if field.Name != v.Name {
			t.Fatalf("%d Name want:%v but:%v", i, v.Name, field.Name)
		}
		if field.Type != v.Type {
			t.Fatalf("%d Type want:%v but:%v", i, v.Type, field.Type)
		}
	}
}

func TestParsePackage_Relations(t *testing.T) {
	t.Skip("not implemented yet")
	path := "./testdata/relations"
	pkg, err := parser.ParsePackage(path)
	if err != nil {
		t.Fatal(err)
	}
	relations := pkg.Classes[0].Relations

	table := []struct {
		target       string
		reltype      ast.RelationType
		multiplicity string
	}{
		{"struct1", ast.Composition, ""},
		{"iface1", ast.Implementation, ""},
		{"association", ast.Association, ""},
		{"embed.struct1", ast.Association, ""},
		{"embed.iface1", ast.Association, ""},
		{"embed.nested.if1", ast.Association, "*"},
	}

	for i, v := range table {
		rel := relations[i]
		if rel.Target != path+"."+v.target {
			t.Fatalf("%d target want:%v but:%v", i, path+"."+v.target, rel.Target)
		}
		if rel.RelType != v.reltype {
			t.Fatalf("%d reltype want:%v but:%v", i, v.reltype, rel.RelType)
		}
	}

}
