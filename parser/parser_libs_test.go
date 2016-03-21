package parser_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/t-yuki/godoc2puml/parser"
)

func TestParsePackage_io(t *testing.T) {
	testParsePackage(t, "io")
}

func TestParsePackage_image(t *testing.T) {
	testParsePackage(t, "image")
}

func testParsePackage(t *testing.T, name string) {
	pkg, err := parser.ParsePackage(name)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(pkg, "", "\t")
	t.Logf("%s", strings.SplitAfterN(string(b), `",`, 2)[0])

}

func TestParsePackage_StdLibs(t *testing.T) {
	for _, name := range []string{"go/ast", "go/token", "reflect", "database/sql", "image", "image/color", "os", "net", "net/http"} {
		pkg, err := parser.ParsePackage(name)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := json.MarshalIndent(pkg, "", "\t")
		t.Logf("%s", strings.SplitAfterN(string(b), `",`, 2)[0])
	}
}
