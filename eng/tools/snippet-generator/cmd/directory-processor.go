//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
)

const (
	mdExt = ".md"
	goExt = ".go"
)

var (
	snippetStartRegex = regexp.MustCompile(`^\s*//\s*Snippet:\s*([A-Za-z0-9-_]+)\s*$`)
	snippetEndRegex   = regexp.MustCompile(`^\s*//\s*EndSnippet\s*$`)
	replaceStartRegex = regexp.MustCompile("^\\s*```\\s*go\\s+Snippet:\\s*([A-Za-z0-9-_]+)\\s*$")
	replaceEndRegex   = regexp.MustCompile("^\\s*```\\s*$")
	ignoreStartRegex = regexp.MustCompile(`^\s*//\s*SnippetIgnore\s*$`)
	ignoreEndRegex = regexp.MustCompile(`^\s*//\s*EndSnippetIgnore`)
)

type directoryProcessor struct {
	Directory  string
	snippets   map[string]*Snippet
	strictMode bool
}

func NewDirectoryProcessor(dir string, strictMode bool) *directoryProcessor {
	return &directoryProcessor{
		Directory:  dir,
		snippets:   make(map[string]*Snippet),
		strictMode: strictMode,
	}
}

func (p *directoryProcessor) Process() error {
	// first enumerate all the files under this directory (and all its sub-directories)
	mdFiles, goFiles, err := findAllFiles(p.Directory)
	if err != nil {
		return fmt.Errorf("cannot enumerate all files in directory %s: %+v", p.Directory, err)
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

	if p.strictMode {
		// check if there are unused snippets
		for name, snippet := range p.snippets {
			if !snippet.IsUsed {
				errResult = multierror.Append(errResult, fmt.Errorf("snippet '%s' is not used", name))
			}
		}
		if errResult != nil {
			return errResult
		}
	}

	return nil
}

func (p *directoryProcessor) processGoSourceFile(path string) error {
	snippetStack := NewStack() // need a stack to store the name of regions and ensure the regions are properly paired with each other
	ignoreStack := NewStack()
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file '%s': %+v", path, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ignoreCounter := 0
	var snippetLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if snippetStack.Len() == 0 {
			matches := snippetStartRegex.FindStringSubmatch(line)
			if len(matches) < 1 {
				continue
			}
			snippetStack.Push(matches[1])
		} else {
			// if the stack is not empty, we also need to ensure that the snippet definition cannot nest with each other
			if snippetStartRegex.MatchString(line) {
				return fmt.Errorf("snippet definition cannot be nested with each other. Filepath '%s'", path)
			}
			// we need to check if current line is an end of a region
			if snippetEndRegex.MatchString(line) {
				name, _ := snippetStack.Pop() // we have checked the len of the stack, there will never be an error here
				if ignoreStack.Len() != 0 {
					return fmt.Errorf("ignore region is not properly enclosed in the Snippet `%s` in file `%s`", name, path)
				}
				p.snippets[name] = NewSnippet(name, snippetLines, path)
				// clean our lines
				snippetLines = nil
				continue
			}
			// also we need to check if it is a region that we ignore in a snippet
			if ignoreStack.Len() == 0 {
				if ignoreStartRegex.MatchString(line) {
					ignoreStack.Push(strconv.Itoa(ignoreCounter))
					ignoreCounter++
					continue
				}

				// only add new lines to the snippet when we are in a region and not in an ignore region
				snippetLines = append(snippetLines, line)
			} else {
				if ignoreEndRegex.MatchString(line) {
					_, _ = ignoreStack.Pop()
					continue
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if snippetStack.Len() > 0 {
		name, _ := snippetStack.Pop()
		return fmt.Errorf("at least one snippet scope ('%s') does not have its corresponding closing", name)
	}

	return nil
}

func (p *directoryProcessor) processMdFile(path string) error {
	stack := NewStack() // need a stack to store the name of regions and ensure the regions are properly paired with each other
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file '%s': %+v", path, err)
	}
	scanner := bufio.NewScanner(file)
	var lines []string
	lineNumber := -1
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if stack.Len() == 0 {
			// this means this is a plain text for our markdown file
			lines = append(lines, line)
			matches := replaceStartRegex.FindStringSubmatch(line)
			if len(matches) < 1 {
				continue
			}
			stack.Push(matches[1]) // put the name into the stack, we will append all the content of the snippet when we pop it out
		} else {
			// if the stack is not empty, it means we are in the middle of replacing one snippet, therefore we should discard all the lines we have here
			// if the stack is not empty, we also need to ensure that the snippet definition cannot nest with each other
			if replaceStartRegex.MatchString(line) {
				return fmt.Errorf("snippet definition cannot be nested with each other. Filepath '%s'", path)
			}
			// if the stack is not empty, we need to check if current line is an end of a region
			if replaceEndRegex.MatchString(line) {
				name, _ := stack.Pop() // we have checked the len of the stack, there will never be an error here
				// find the corresponding snippet
				snippet := p.snippets[name]
				// append all the lines here
				snippet.IsUsed = true
				lines = append(lines, snippet.SourceText...)
				lines = append(lines, line)
				log.Printf("Successfully replaced snippet '%s' in file '%s'", name, path)
				continue
			}
		}
	}
	if err := file.Close(); err != nil {
		return err
	}

	// write things back to the file
	file, err = os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create file '%s': %+v", path, err)
	}
	defer file.Close()
	if _, err := file.WriteString(strings.Join(lines, "\n")); err != nil {
		return fmt.Errorf("cannot write file '%s': %+v", path, err)
	}

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
			if strings.HasSuffix(info.Name(), "_test"+goExt) {
				goFiles = append(goFiles, path)
			}
			break
		}
		return nil
	})
	return mdFiles, goFiles, err
}
