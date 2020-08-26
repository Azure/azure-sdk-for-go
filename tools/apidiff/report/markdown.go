// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package report

import (
	"strings"
)

// MarkdownWriter is a writer to write contents in markdown format
type MarkdownWriter struct {
	sb strings.Builder
	nl bool
}

func (md *MarkdownWriter) checkNL() {
	if md.nl {
		md.sb.WriteString("\n")
		md.nl = false
	}
}

// WriteHeader writes a header to the markdown document
func (md *MarkdownWriter) WriteHeader(h string) {
	md.checkNL()
	md.sb.WriteString("## ")
	md.sb.WriteString(h)
	md.sb.WriteString("\n\n")
}

// WriteSubheader writes a sub-header to the markdown document
func (md *MarkdownWriter) WriteSubheader(sh string) {
	md.checkNL()
	md.sb.WriteString("### ")
	md.sb.WriteString(sh)
	md.sb.WriteString("\n\n")
}

// WriteLine writes a line to the markdown document
func (md *MarkdownWriter) WriteLine(s string) {
	md.nl = true
	md.sb.WriteString(s)
	md.sb.WriteString("\n")
}

// WriteTable writes a table to the markdown document
func (md *MarkdownWriter) WriteTable(table MarkdownTable) {
	md.WriteLine(table.String())
}

// EmptyLine inserts an empty line to the markdown document
func (md *MarkdownWriter) EmptyLine() {
	md.WriteLine("")
}

// String outputs the markdown document as a string
func (md *MarkdownWriter) String() string {
	return md.sb.String()
}

// MarkdownTable describes a table in a markdown document
type MarkdownTable struct {
	sb        *strings.Builder
	headers   []string
	alignment string
	rows      []markdownTableRow
}

const (
	leftAlignment   = ":---"
	centerAlignment = ":---:"
	rightAlignment  = "---:"
)

// NewMarkdownTable creates a new table with given alignments and headers
func NewMarkdownTable(alignment string, headers ...string) *MarkdownTable {
	alignment, headers = checkAlignmentAndHeaders(alignment, headers)
	t := MarkdownTable{
		sb:        &strings.Builder{},
		headers:   headers,
		alignment: alignment,
	}
	return &t
}

// Columns returns the number of columns in this table
func (t *MarkdownTable) Columns() int {
	return len(t.headers)
}

// AddRow adds a new row to the table
func (t *MarkdownTable) AddRow(items ...string) {
	t.rows = append(t.rows, markdownTableRow{
		items: items,
	})
}

func checkAlignmentAndHeaders(alignment string, headers []string) (string, []string) {
	if len(alignment) == len(headers) {
		return alignment, headers
	}
	if len(alignment) > len(headers) {
		return alignment[:len(headers)], headers
	}
	// default alignment is l
	return alignment + strings.Repeat("l", len(headers)-len(alignment)), headers
}

type markdownTableRow struct {
	items []string
}

func (r *markdownTableRow) ensureItems(count int) []string {
	if len(r.items) < count {
		var items []string
		for i := 0; i < count-len(r.items); i++ {
			items = append(r.items, "")
		}
		return items
	}
	if len(r.items) > count {
		return r.items[:count]
	}
	return r.items
}

func (t *MarkdownTable) writeHeader() {
	if len(t.headers) == 0 {
		return
	}
	t.sb.WriteString("| ")
	t.sb.WriteString(strings.Join(t.headers, " | "))
	t.sb.WriteString(" |\n")
}

func (t *MarkdownTable) writeAlignment() {
	if len(t.alignment) == 0 {
		return
	}
	var alignments []string
	for _, a := range t.alignment {
		switch a {
		case 'c':
			alignments = append(alignments, centerAlignment)
		case 'r':
			alignments = append(alignments, rightAlignment)
		default:
			alignments = append(alignments, leftAlignment)
		}
	}
	t.sb.WriteString("| ")
	t.sb.WriteString(strings.Join(alignments, " | "))
	t.sb.WriteString(" |\n")
}

func (t *MarkdownTable) writeRows() {
	for _, r := range t.rows {
		t.writeRow(r)
	}
}

func (t *MarkdownTable) writeRow(row markdownTableRow) {
	items := row.ensureItems(t.Columns())
	t.sb.WriteString("| ")
	t.sb.WriteString(strings.Join(items, " | "))
	t.sb.WriteString(" |\n")
}

// String outputs the markdown table to a string
func (t *MarkdownTable) String() string {
	t.writeHeader()
	t.writeAlignment()
	t.writeRows()
	return t.sb.String()
}
