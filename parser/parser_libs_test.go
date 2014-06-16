package parser_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/t-yuki/godoc2puml/parser"
)

func TestParsePackage_IO(t *testing.T) {
	pkg, err := parser.ParsePackage("io")
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
