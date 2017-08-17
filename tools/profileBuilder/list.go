package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"

	"github.com/marstr/collection"
)

// ListStrategy allows a mechanism for a list of packages that should be included in a profile.
type ListStrategy struct {
	io.Reader
}

// Enumerate reads a new line delimited list of packages names relative to $GOPATH
func (list ListStrategy) Enumerate(cancel <-chan struct{}) collection.Enumerator {
	results := make(chan interface{})

	go func() {
		defer close(results)

		var currentLine string
		for _, err := fmt.Fscanln(list, &currentLine); err == nil; {
			var pkg map[string]*ast.Package
			files := token.NewFileSet()
			pkg, err = parser.ParseDir(files, currentLine, nil, parser.ParseComments)

			for _, entry := range pkg {
				select {
				case results <- entry:
					// Intentionally Left Blank
				case <-cancel:
					return
				}
			}
		}
	}()

	return results
}
