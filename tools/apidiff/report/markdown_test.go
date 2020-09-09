package report

import "testing"

func TestMarkdownWriter_String(t *testing.T) {
	md := MarkdownWriter{}
	md.WriteHeader("Header1")
	md.WriteSubheader("Sub-header1")
	md.WriteLine("Foo")
	md.WriteSubheader("Sub-header2")
	md.WriteLine(getTable().String())
	result := md.String()
	expected := `## Header1

### Sub-header1

Foo

### Sub-header2

| packages | api-versions |
| :--- | :--- |
| compute | 2020-06-01 |
| network | 2020-06-01 |
| resources | 2020-09-01 |
`
	if result != expected {
		t.Fatalf("expected %s, but got %s", expected, result)
	}
}

func getTable() *MarkdownTable {
	table := NewMarkdownTable("ll", "packages", "api-versions")
	table.AddRow("compute", "2020-06-01")
	table.AddRow("network", "2020-06-01")
	table.AddRow("resources", "2020-09-01")
	return table
}

func TestMarkdownTable_String(t *testing.T) {
	table := getTable()
	result := table.String()
	expected := `| packages | api-versions |
| :--- | :--- |
| compute | 2020-06-01 |
| network | 2020-06-01 |
| resources | 2020-09-01 |`
	if result != expected {
		t.Fatalf("expected %s, but got %s", expected, result)
	}
}
