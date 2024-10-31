// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
package transform

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Transformer struct {
	workingDir string
}

func New(workingDir string) (*Transformer, error) {
	if _, err := exec.LookPath("gopls"); err != nil {
		return nil, fmt.Errorf("gopls not found in PATH: %v", err)
	}

	return &Transformer{
		workingDir: workingDir,
	}, nil
}

func (t *Transformer) executeGopls(args ...string) (*bytes.Buffer, error) {
	cmd := exec.Command("gopls", args...)
	cmd.Dir = t.workingDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("gopls error: %v\nstderr: %s", err, stderr.String())
	}

	return &stdout, nil
}

func (t *Transformer) findSymbolPosition(filename, symbolName string, symbolKind string, start *int, end *int) (*position, error) {
	stdout, err := t.executeGopls("symbols", filepath.Clean(filename))
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) < 3 {
			continue
		}

		symbol := parts[0]
		kind := parts[1]
		position := parts[2]

		if symbol == symbolName && kind == symbolKind {
			absPath, err := filepath.Abs(filename)
			if err != nil {
				return nil, fmt.Errorf("failed to get absolute path: %v", err)
			}
			pos, err := newPosition(absPath, position)
			if err != nil {
				return nil, err
			}
			if start != nil && pos.Line < *start {
				continue
			}
			if end != nil && pos.Line >= *end {
				break
			}
			return pos, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading gopls output: %v", err)
	}

	return nil, fmt.Errorf("symbol %s not found in %s", symbolName, filename)

}

func (t *Transformer) RenameStruct(filename string, oldName string, newName string) error {
	pos, err := t.findSymbolPosition(filename, oldName, "Struct", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to find struct position: %v", err)
	}

	_, err = t.executeGopls("rename",
		"-w",
		pos.String(),
		newName)

	return err
}

func (t *Transformer) RenameMethod(filename string, methodName string, newName string) error {
	pos, err := t.findSymbolPosition(filename, methodName, "Method", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to find method position: %v", err)
	}

	_, err = t.executeGopls("rename",
		"-w",
		pos.String(),
		newName)

	return err
}

type position struct {
	FileName    string
	Line        int
	ColumnStart int
	ColumnEnd   int
}

func newPosition(fileName string, pos string) (*position, error) {
	start, end, found := strings.Cut(pos, "-")
	if !found {
		return nil, fmt.Errorf("unexpected position format: %s", pos)
	}
	line, columnStart, found := strings.Cut(start, ":")
	if !found {
		return nil, fmt.Errorf("unexpected position format: %s", pos)
	}
	_, columnEnd, found := strings.Cut(end, ":")
	if !found {
		return nil, fmt.Errorf("unexpected position format: %s", pos)
	}

	lineNumber, err := strconv.Atoi(line)
	if err != nil {
		return nil, fmt.Errorf("failed to parse line number: %v", err)
	}
	columnStartNumber, err := strconv.Atoi(columnStart)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start column: %v", err)
	}
	columnEndNumber, err := strconv.Atoi(columnEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end column: %v", err)
	}

	return &position{
		FileName:    fileName,
		Line:        lineNumber,
		ColumnStart: columnStartNumber,
		ColumnEnd:   columnEndNumber,
	}, nil
}

func (p *position) String() string {
	return fmt.Sprintf("%s:%d:%d-%d:%d", p.FileName, p.Line, p.ColumnStart, p.Line, p.ColumnEnd)
}

func (t *Transformer) findStructRange(filename string, structName string) (start int, end int, err error) {
	pos, err := t.findSymbolPosition(filename, structName, "Struct", nil, nil)
	if err != nil {
		return 0, 0, err
	}

	stdout, err := t.executeGopls("folding_ranges",
		filename)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to get folding ranges: %v", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, fmt.Sprint(pos.Line)) {

			start, end, _ := strings.Cut(line, "-")
			start, _, _ = strings.Cut(start, ":")
			end, _, _ = strings.Cut(end, ":")

			startInt, err := strconv.Atoi(start)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse start line: %v", err)
			}

			endInt, err := strconv.Atoi(end)
			if err != nil {
				return 0, 0, fmt.Errorf("failed to parse end line: %v", err)
			}

			return startInt, endInt, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, fmt.Errorf("error reading gopls output: %v", err)
	}

	return 0, 0, fmt.Errorf("struct %s not found in %s", structName, filename)
}

func (t *Transformer) ReadRange(filename string, start, end *int) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if start != nil && lineNumber < *start {
			continue
		}
		if end != nil && lineNumber > *end {
			break
		}
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return lines, nil
}

func (t *Transformer) CopyStruct(filename string, structName string, newStructName string) error {
	start, end, err := t.findStructRange(filename, structName)
	if err != nil {
		return err
	}

	lines, err := t.ReadRange(filename, &start, &end)
	if err != nil {
		return err
	}

	lines[0] = strings.Replace(lines[0], structName, newStructName, 1)

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	allLines := strings.Join(lines, "\n")

	if _, err := fmt.Fprintf(file, "\n%s\n", allLines); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func (t *Transformer) RemoveField(filename string, structName string, fieldName string) error {
	start, end, err := t.findStructRange(filename, structName)
	if err != nil {
		return err
	}

	pos, err := t.findSymbolPosition(filename, fieldName, "Field", &start, &end)
	if err != nil {
		return err
	}

	lines, err := t.ReadRange(filename, nil, nil)
	if err != nil {
		return err
	}

	// Remove the field line
	removeFrom := pos.Line - 1
	if pos.Line-2 >= 0 && strings.HasPrefix(strings.TrimSpace(lines[pos.Line-2]), "//") {
		removeFrom = pos.Line - 2
	}

	lines = append(lines[:removeFrom], lines[pos.Line:]...)

	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	allLines := strings.Join(lines, "\n")

	if _, err := fmt.Fprintf(file, "\n%s\n", allLines); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
