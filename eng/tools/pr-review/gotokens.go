// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/scanner"
	"go/token"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var copyrightHeader = regexp.MustCompile(`^// Copyright \(c\) (?:[0-9]{4}(?:-[0-9]{4})? )?Microsoft Corporation\.(?: All rights reserved\.)?$`)

type sourceTokens struct {
	generated bool
	tokens    []string
	err       error
}

func scanSource(source string) sourceTokens {
	result := sourceTokens{}
	files := token.NewFileSet()
	var lexer scanner.Scanner
	lexer.Init(files.AddFile("source.go", -1, len(source)), []byte(source), func(pos token.Position, message string) {
		if result.err == nil {
			result.err = fmt.Errorf("%s: %s", pos, message)
		}
	}, scanner.ScanComments)
	header := true
	for {
		pos, kind, literal := lexer.Scan()
		if kind == token.EOF {
			break
		}
		if kind == token.COMMENT {
			if header && strings.Contains(literal, "DO NOT EDIT") {
				result.generated = true
			}
			if header && copyrightHeader.MatchString(literal) {
				continue
			}
			if strings.HasPrefix(literal, "//go:") || strings.HasPrefix(literal, "// +build") ||
				strings.HasPrefix(literal, "//line ") || strings.HasPrefix(literal, "/*line ") ||
				strings.HasPrefix(literal, "//export ") {
				offset := files.Position(pos).Offset
				end := offset + len(literal)
				startGap, endGap := offset, end
				for startGap > 0 && strings.ContainsRune(" \t\r\n", rune(source[startGap-1])) {
					startGap--
				}
				for endGap < len(source) && strings.ContainsRune(" \t\r\n", rune(source[endGap])) {
					endGap++
				}
				// Directive placement can be meaningful even when the tokens are unchanged.
				literal += fmt.Sprintf("@%t:%d:%d", files.Position(pos).Column == 1,
					min(strings.Count(source[startGap:offset], "\n"), 2),
					min(strings.Count(source[end:endGap], "\n"), 2))
			}
		} else {
			header = false
		}
		if kind == token.SEMICOLON {
			literal = ";"
		}
		result.tokens = append(result.tokens, fmt.Sprintf("%d:%s", kind, literal))
	}
	return result
}

func compareSources(before, after *string) (bool, bool, error) {
	var sources [2]sourceTokens
	for i, source := range []*string{before, after} {
		if source != nil {
			sources[i] = scanSource(*source)
		}
	}
	generated := sources[0].generated || sources[1].generated
	if !generated {
		return false, false, nil
	}
	for i, source := range []*string{before, after} {
		if source == nil {
			continue
		}
		if sources[i].err != nil {
			return true, false, sources[i].err
		}
		file, err := parser.ParseFile(token.NewFileSet(), "source.go", *source, parser.AllErrors|parser.SkipObjectResolution)
		if err != nil {
			return true, false, err
		}
		for _, imported := range file.Imports {
			path, err := strconv.Unquote(imported.Path.Value)
			if err != nil {
				return true, false, err
			}
			if path == "C" {
				return true, false, fmt.Errorf("cgo preamble attachment requires manual review")
			}
		}
	}
	return true, !slices.Equal(sources[0].tokens, sources[1].tokens), nil
}

func run() error {
	var input struct {
		Before *string `json:"before"`
		After  *string `json:"after"`
	}
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		return err
	}
	generated, changed, err := compareSources(input.Before, input.After)
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(struct {
		Generated bool `json:"generated"`
		Changed   bool `json:"changed"`
	}{generated, changed})
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
