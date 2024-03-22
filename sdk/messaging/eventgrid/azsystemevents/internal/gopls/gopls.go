//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package gopls

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Symbol struct {
	Name string
	// Either:
	// - Constant
	// - Class (type ACSRouterJobStatus string)
	// - Struct
	// - Field
	// - Function
	// - Method
	Type     string
	Position string
}

func Rename(filename string, sym Symbol, newName string) error {
	if sym.Name == "" && sym.Position == "" && sym.Type == "" {
		return errors.New("attempted to rename with an empty symbol")
	}

	cmd := exec.Command("gopls", "rename", "-w", filename+":"+sym.Position, newName)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func Symbols(filename string) (map[string]Symbol, error) {
	if _, err := os.Stat(filename); err != nil {
		return nil, err
	}

	// ex: PossibleRecordingChannelTypeValues Function 1003:6-1003:40
	cmd := exec.Command("gopls", "symbols", filename)
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	allBytes := stdout.Bytes()

	scanner := bufio.NewScanner(bytes.NewReader(allBytes))

	m := map[string]Symbol{}
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
			Type:     fields[1],
			Position: fields[2],
		}

		if _, exists := m[sym.Name]; exists {
			return nil, fmt.Errorf("%s was in the map twice", sym.Name)
		}

		m[sym.Name] = sym

		if sym.Type != "Field" {
			prevParent = sym.Name
		}
	}

	return m, nil
}
