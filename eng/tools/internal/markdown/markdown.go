// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package markdown

import (
	"fmt"
	"strings"
)

// RenderLink returns a rendered markdown link
func RenderLink(name, link string) string {
	return fmt.Sprintf("[%s](%s)", name, link)
}

// Writer is a writer to write contents in markdown format
type Writer struct {
	sb strings.Builder
	nl bool
}

func (md *Writer) checkNL() {
	if md.nl {
		md.sb.WriteString("\n")
		md.nl = false
	}
}

// WriteTitle writes a title to the markdown document
func (md *Writer) WriteTitle(h string) {
	md.checkNL()
	md.sb.WriteString("# ")
	md.sb.WriteString(h)
	md.sb.WriteString("\n\n")
}

// WriteTopLevelHeader writes a header to the markdown document
func (md *Writer) WriteTopLevelHeader(h string) {
	md.checkNL()
	md.sb.WriteString("## ")
	md.sb.WriteString(h)
	md.sb.WriteString("\n\n")
}

// WriteHeader writes a header to the markdown document
func (md *Writer) WriteHeader(h string) {
	md.checkNL()
	md.sb.WriteString("### ")
	md.sb.WriteString(h)
	md.sb.WriteString("\n\n")
}

// WriteSubheader writes a sub-header to the markdown document
func (md *Writer) WriteSubheader(sh string) {
	md.checkNL()
	md.sb.WriteString("#### ")
	md.sb.WriteString(sh)
	md.sb.WriteString("\n\n")
}

// WriteLine writes a line to the markdown document
func (md *Writer) WriteLine(s string) {
	md.nl = true
	md.sb.WriteString(s)
	md.sb.WriteString("\n")
}

// WriteListItem writes a line in a list
func (md *Writer) WriteListItem(item string) {
	md.WriteLine(fmt.Sprintf("- %s", item))
}

// WriteTable writes a table to the markdown document
func (md *Writer) WriteTable(table Table) {
	md.WriteLine(table.String())
}

// EmptyLine inserts an empty line to the markdown document
func (md *Writer) EmptyLine() {
	md.WriteLine("")
}

// String outputs the markdown document as a string
func (md *Writer) String() string {
	return md.sb.String()
}

// Table describes a table in a markdown document
type Table struct {
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

// NewTable creates a new table with given alignments and headers
func NewTable(alignment string, headers ...string) *Table {
	alignment, headers = checkAlignmentAndHeaders(alignment, headers)
	t := Table{
		sb:        &strings.Builder{},
		headers:   headers,
		alignment: alignment,
	}
	return &t
}

// Columns returns the number of columns in this table
func (t *Table) Columns() int {
	return len(t.headers)
}

// Rows returns the number of rows in this table
func (t *Table) Rows() int {
	return len(t.rows)
}

// AddRow adds a new row to the table
func (t *Table) AddRow(items ...string) {
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

func (t *Table) writeHeader() {
	if len(t.headers) == 0 {
		return
	}
	t.sb.WriteString("| ")
	t.sb.WriteString(strings.Join(t.headers, " | "))
	t.sb.WriteString(" |")
}

func (t *Table) writeAlignment() {
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
	t.sb.WriteString("\n| ")
	t.sb.WriteString(strings.Join(alignments, " | "))
	t.sb.WriteString(" |")
}

func (t *Table) writeRows() {
	for _, r := range t.rows {
		t.writeRow(r)
	}
}

func (t *Table) writeRow(row markdownTableRow) {
	items := row.ensureItems(t.Columns())
	t.sb.WriteString("\n| ")
	t.sb.WriteString(strings.Join(items, " | "))
	t.sb.WriteString(" |")
}

// String outputs the markdown table to a string
func (t *Table) String() string {
	t.writeHeader()
	t.writeAlignment()
	t.writeRows()
	return t.sb.String()
}
