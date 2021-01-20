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
package markdown

import "testing"

func TestMarkdownWriter_String(t *testing.T) {
	md := Writer{}
	md.WriteTitle("Title")
	md.WriteTopLevelHeader("Top-level header")
	md.WriteHeader("Header1")
	md.WriteSubheader("Sub-header1")
	md.WriteLine("Foo")
	md.WriteSubheader("Sub-header2")
	md.WriteLine(getTable().String())
	result := md.String()
	expected := `# Title

## Top-level header

### Header1

#### Sub-header1

Foo

#### Sub-header2

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

func getTable() *Table {
	table := NewTable("ll", "packages", "api-versions")
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
