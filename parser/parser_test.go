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
		{"xi8", "int8"},
		{"yembed.int16", "*int16"},
		{"yembed.xi32", "[]int32"},
		{"yembed.yi32", "[]int32"},
		{"yembed.znested.i64", "int64"},
		{"yyembed.int16", "*int16"},
		{"yyembed.xi32", "[]int32"},
		{"yyembed.yi32", "[]int32"},
		{"yyembed.znested.i64", "int64"},
		{"zembed.bool", "*bool"},
		{"zzembed", "[]struct{int8, u8 uint8, uchar uint8}"}, // TODO: parse more?
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
	path := "./testdata/relations"
	pkg, err := parser.ParsePackage(path)
	if err != nil {
		t.Fatal(err)
	}
	relations := pkg.Classes[0].Relations

	table := []struct {
		Target       string
		Label        string
		RelType      ast.RelationType
		Multiplicity string
	}{
		{"struct1", "", ast.Composition, ""},
		{"iface1", "", ast.Composition, ""}, // even if iface1 is interface, it can't be detected on parse
		{"struct1", "xassociation", ast.Association, ""},
		{"struct1", "yembed.struct1", ast.Association, ""},
		{"iface1", "yembed.iface1", ast.Association, ""},
		{"iface1", "yembed.nested.if1", ast.Association, "*"},
		{"struct1", "zembed.struct1", ast.Association, ""},
		{"iface1", "zembed.iface1", ast.Association, ""},
		{"iface1", "zembed.nested.if1", ast.Association, "*"},
	}

	for i, v := range table {
		rel := relations[i]
		if rel.Label != v.Label {
			t.Fatalf("%d Label want:%v but:%v", i, v.Label, rel.Label)
		}
		if rel.Target != path+"."+v.Target {
			t.Fatalf("%d Target want:%v but:%v", i, path+"."+v.Target, rel.Target)
		}
		if rel.RelType != v.RelType {
			t.Fatalf("%d RelType want:%v but:%v", i, v.RelType, rel.RelType)
		}
		if rel.Multiplicity != v.Multiplicity {
			t.Fatalf("%d Multiplicity want:%v but:%v", i, v.Multiplicity, rel.Multiplicity)
		}
	}

}
