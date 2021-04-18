package validators

import (
	"go/parser"
	"go/token"
)

func CheckGoSyntax(src string) error {
	_, err := parser.ParseFile(token.NewFileSet(), "", src, parser.ParseComments)

	return err
}
