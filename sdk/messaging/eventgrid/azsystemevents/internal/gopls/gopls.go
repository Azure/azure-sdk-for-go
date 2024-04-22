//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package gopls

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// SymbolType maps to the second field in the `gopls symbols` output.
type SymbolType string

const (
	SymbolTypeConstant SymbolType = "Constant"
	SymbolTypeClass    SymbolType = "Class" // (ex: ACSRouterJobStatus, from 'type ACSRouterJobStatus string')
	SymbolTypeStruct   SymbolType = "Struct"
	SymbolTypeField    SymbolType = "Field"
	SymbolTypeFunction SymbolType = "Function"
	SymbolTypeMethod   SymbolType = "Method" // ex: '(*WebBackupOperationCompletedEventData).UnmarshalJSON'
)

type Symbol struct {
	Name string
	Type SymbolType

	// Position is a start and end range for a symbol
	// ex: 36:2-36:7
	// (line 36, columns 2 - 7)
	Position string

	File string
}

type SymbolMap struct {
	m map[string]*Symbol
}

func NewSymbolMap(m map[string]*Symbol) *SymbolMap {
	newM := map[string]*Symbol{}

	for k, v := range m {
		newM[strings.ToLower(k)] = v
	}

	return &SymbolMap{m: newM}
}

// Get retrieves a symbol, ignoring casing.
func (s SymbolMap) Get(name string) *Symbol {
	v, exists := s.m[strings.ToLower(name)]

	if !exists {
		log.Fatalf("No key exists for %q", name)
	}

	return v
}

// All gets all the symbols within our map, in no specific order.
func (s SymbolMap) All() []*Symbol {
	var all []*Symbol

	for _, v := range s.m {
		all = append(all, v)
	}

	return all
}

// ex: 36:2-36:7
var posRE = regexp.MustCompile(`(\d+):(\d+)-(\d+):(\d+)`)

// StartLine returns the starting line for this symbol.
func (s Symbol) StartLine() (int64, error) {
	var matches = posRE.FindStringSubmatch(s.Position)

	if matches == nil {
		return 0, fmt.Errorf("couldn't parse position %q", s.Position)
	}

	return strconv.ParseInt(matches[1], 10, 64)
}

// Rename renames a given symbol to a new name. This naming will propagate
// throughout, including references to the type, the type name itself and
// any comments that contain the type name in godoc's reference format
// with '[typename]'.
func Rename(sym *Symbol, newName string) error {
	cmd := exec.Command("gopls", "rename", "-w", sym.File+":"+sym.Position, newName)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Symbols parses the output of `gopls symbols`.
func Symbols(filename string) (*SymbolMap, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	// ex: PossibleRecordingChannelTypeValues Function 1003:6-1003:40
	cmd := exec.Command("gopls", "symbols", filename)
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed getting symbols for %s: %w", filename, err)
	}

	allBytes := stdout.Bytes()

	scanner := bufio.NewScanner(bytes.NewReader(allBytes))

	m := map[string]*Symbol{}
	prevParent := ""

	for scanner.Scan() {
		// ex:
		//
		// Top level types are listed like this:
		// RecordingFormatTypeMp3 Constant 1030:2-1030:24
		//
		// Fields are listed like this - skipping those for now since I'm mostly concerned with
		// renaming types.
		//         Address Field 5379:2-5379:9

		var fields []string

		if len(scanner.Text()) > 0 && scanner.Text()[0] == '\t' {
			// this is a field - we'll add it and "parent" it to the last non-field symbol before it
			// (and clip out the first char, which is a tab)
			fields = strings.Split(scanner.Text()[1:], " ")

			fields[0] = prevParent + "." + fields[0]
		} else {
			fields = strings.Split(scanner.Text(), " ")
		}

		if len(fields) != 3 {
			return nil, fmt.Errorf("failed to parse %q into three fields, got %d", scanner.Text(), len(fields))
		}

		sym := Symbol{
			Name:     fields[0],
			Type:     SymbolType(fields[1]),
			Position: fields[2],
			File:     filename,
		}

		if _, exists := m[sym.Name]; exists {
			return nil, fmt.Errorf("%s was in the map twice", sym.Name)
		}

		m[sym.Name] = &sym

		if sym.Type != SymbolTypeField {
			prevParent = sym.Name
		}
	}

	return NewSymbolMap(m), nil
}
