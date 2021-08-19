// Retrieves the version string from the version.go files of track 2 go packages.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func main() {
	fset := token.NewFileSet()
	src := os.Args[1]

	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, decl := range f.Decls {
		d, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range d.Specs {
			s, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, value := range s.Values {
				v, ok := value.(*ast.BasicLit)
				if !ok {
					continue
				}
				fmt.Println(v.Value)
			}
		}
	}
}
