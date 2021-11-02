package cmd

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io/fs"
	"path/filepath"
	"strings"
)

const (
	snippetPrefix    = "Snippet:"
	endSnippetPrefix = "EndSnippet"

	mdExt = "md"
	goExt = "go"
)

type directoryProcessor struct {
	directory string
	snippets  map[string]*Snippet
}

func (p *directoryProcessor) Process() error {
	// first enumerate all the files under this directory (and all its sub-directories)
	mdFiles, goFiles, err := findAllFiles(p.directory)
	if err != nil {
		return fmt.Errorf("cannot enumerate all files in directory %s: %+v", p.directory, err)
	}

	// collect all the snippets from go source files
	var errResult error
	for _, goFile := range goFiles {
		if err := p.processGoSourceFile(goFile); err != nil {
			errResult = multierror.Append(errResult, err)
		}
	}

	if errResult != nil {
		return fmt.Errorf("cannot read snippets: %+v", errResult)
	}

	// apply those snippets to markdown files
	for _, mdFile := range mdFiles {
		if err := p.processMdFile(mdFile); err != nil {
			errResult = multierror.Append(errResult, err)
		}
	}

	if errResult != nil {
		return fmt.Errorf("cannot apply snippets: %+v", errResult)
	}

	return nil
}

func (p *directoryProcessor) processGoSourceFile(path string) error {
	// TODO
	return nil
}

func (p *directoryProcessor) processMdFile(path string) error {
	// TODO
	return nil
}

func findAllFiles(dir string) ([]string, []string, error) {
	var mdFiles []string
	var goFiles []string
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch filepath.Ext(info.Name()) {
		case mdExt:
			mdFiles = append(mdFiles, path)
			break
		case goExt:
			// the snippets can only present in the test files
			if strings.HasSuffix(info.Name(), "_test") {
				goFiles = append(goFiles, path)
			}
			break
		}
		return nil
	})
	return mdFiles, goFiles, err
}
