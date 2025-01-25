package main

import (
	"flag"
	"fmt"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var isMGMT = false

func filter(f fs.FileInfo) bool {
	return !strings.HasSuffix(f.Name(), "_test.go")
}

func findAllSubDirectories(root string) []string {
	var ret []string

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if strings.Contains(path, "resourcemanager") {
			isMGMT = true
		}
		if strings.Contains(path, "eng/tools") {
			return filepath.SkipDir
		}
		if d.IsDir() && strings.HasSuffix(path, "internal") {
			return filepath.SkipDir
		} else if d.IsDir() {
			ret = append(ret, path)
		}
		return nil
	})

	return ret
}

// validateDirectory counts the number of missing doc comments
func validateDirectory(directory string) int {
	missingDocCount := 0

	fset := token.NewFileSet() // positions are relative to fset
	d, err := parser.ParseDir(fset, directory, filter, parser.ParseComments)
	if err != nil {
		panic(fmt.Errorf("could not parse directory: %w", err))
	}

	for k, f := range d {
		if strings.Contains("_test", k) {
			continue
		}
		fmt.Println("package", k)
		p := doc.New(f, "./", 2)

		for _, t := range p.Types {
			if !strings.HasPrefix(t.Doc, t.Name) {
				fmt.Println("missing or invalid doc comment. all docs should start with the type they are documenting")
				fmt.Printf("type: '%s' docs: '%s'\n", t.Name, t.Doc)
				missingDocCount += 1
			}

			for _, m := range t.Methods {
				if !strings.HasPrefix(m.Doc, m.Name) {
					fmt.Println("missing or invalid doc comment. all docs should start with the function they are documenting")
					fmt.Printf("type: '%s' method: '%s' docs: '%s'\n", t.Name, m.Name, m.Doc)
					missingDocCount += 1
				}
			}
		}

		for _, f := range p.Funcs {
			if strings.HasPrefix(f.Name, "Example") {
				continue
			}
			if !strings.HasPrefix(f.Doc, f.Name) {
				fmt.Println("missing or invalid doc comment. all docs should start with the function they are documenting")
				if f.Recv != "" {
					fmt.Printf("type: %s receiver: %s docs: %s\n", f.Name, f.Recv, f.Doc)
				} else {
					fmt.Printf("type: %s docs: %s\n", f.Name, f.Doc)
				}
				missingDocCount += 1
			}
		}
	}
	return missingDocCount
}

func main() {
	var root string
	flag.StringVar(&root, "directory", ".", "directory to check docs for")
	flag.Parse()

	fmt.Printf("checking documentation in %s\n", root)
	var totalMissing = 0
	for _, dir := range findAllSubDirectories(root) {
		totalMissing += validateDirectory(dir)
	}

	if totalMissing > 0 {
		fmt.Printf("Found %d missing doc comments\n", totalMissing)
		if !isMGMT {
			os.Exit(1)
		}

	} else {
		fmt.Println("There are no public methods/functions/types with missing documentation")
	}
}
