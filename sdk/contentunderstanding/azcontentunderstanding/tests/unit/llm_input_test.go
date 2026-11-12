// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding_test

import (
	"encoding/json"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	sdk "github.com/Azure/azure-sdk-for-go/sdk/contentunderstanding/azcontentunderstanding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func llmDocument(markdown string, start, end int32) *sdk.DocumentContent {
	return &sdk.DocumentContent{
		MimeType: to.Ptr("application/pdf"), Markdown: to.Ptr(markdown),
		StartPageNumber: to.Ptr(start), EndPageNumber: to.Ptr(end),
	}
}

func llmResult(contents ...sdk.AnalysisContentClassification) *sdk.AnalysisResult {
	return &sdk.AnalysisResult{Contents: contents}
}

func llmSpan(offset, length int32) *sdk.ContentSpan {
	return &sdk.ContentSpan{Offset: to.Ptr(offset), Length: to.Ptr(length)}
}

func llmString(value string) sdk.ContentFieldClassification {
	return &sdk.StringField{ValueString: to.Ptr(value)}
}

func renderLLM(t *testing.T, result *sdk.AnalysisResult, metadata map[string]any, options *sdk.LLMInputOptions) string {
	t.Helper()
	output, err := sdk.ToLLMInput(result, metadata, options)
	require.NoError(t, err)
	return output
}

func llmFrontMatter(t *testing.T, output string) string {
	t.Helper()
	parts := strings.SplitN(output, "\n---", 2)
	require.Len(t, parts, 2)
	return strings.TrimPrefix(parts[0], "---\n")
}

func TestToLLMInputBasic(t *testing.T) {
	output, err := sdk.ToLLMInput(nil, nil, nil)
	require.ErrorContains(t, err, "result")
	assert.Empty(t, output)
	for _, result := range []*sdk.AnalysisResult{{}, {Contents: []sdk.AnalysisContentClassification{}}} {
		assert.Empty(t, renderLLM(t, result, nil, nil))
	}
	for _, mime := range []string{"application/pdf", "image/jpeg", "", "audio/wav"} {
		t.Run(mime, func(t *testing.T) {
			doc := llmDocument("CONTOSO LTD.\n\n# INVOICE", 1, 1)
			doc.MimeType = to.Ptr(mime)
			doc.Fields = map[string]sdk.ContentFieldClassification{
				"VendorName": llmString("CONTOSO LTD."), "InvoiceDate": llmString("2019-11-15"),
			}
			expectedMime := mime
			if expectedMime == "" {
				expectedMime = "unknown"
			}
			output := renderLLM(t, llmResult(doc), nil, nil)
			assert.Equal(t, "---\nmimeType: "+expectedMime+"\npages: 1\nfields:\n"+
				"  InvoiceDate: '2019-11-15'\n  VendorName: CONTOSO LTD.\n---\n"+
				"<!-- InputPageNumber: 1 -->\n\nCONTOSO LTD.\n\n# INVOICE", output)
			assert.NotContains(t, output, "contentType")
		})
	}
	doc := llmDocument("text", 0, 0)
	doc.MimeType = nil
	assert.Equal(t, "---\nmimeType: unknown\n---\n<!-- InputPageNumber: 1 -->\n\ntext",
		renderLLM(t, llmResult(doc), nil, nil))
	assert.Equal(t, "---\nmimeType: unknown\n---",
		renderLLM(t, llmResult(&sdk.AnalysisContent{}), nil, nil))
}

func TestToLLMInputIncludeFlags(t *testing.T) {
	for _, fields := range []*bool{nil, to.Ptr(true), to.Ptr(false)} {
		for _, markdown := range []*bool{nil, to.Ptr(true), to.Ptr(false)} {
			t.Run(boolName(fields)+"/"+boolName(markdown), func(t *testing.T) {
				doc := llmDocument("Body text", 1, 1)
				doc.Fields = map[string]sdk.ContentFieldClassification{"Name": llmString("Test")}
				doc.Metadata = map[string]*string{"author": to.Ptr("Jane")}
				doc.Category = to.Ptr("Invoice")
				result := llmResult(doc)
				result.Warnings = []*sdk.Error{{Code: to.Ptr("hate"), Message: to.Ptr("Safety warning.")}}
				output := renderLLM(t, result, map[string]any{"source": "file"}, &sdk.LLMInputOptions{
					IncludeFields: fields, IncludeMarkdown: markdown,
				})
				assert.Equal(t, fields == nil || *fields, strings.Contains(output, "Name: Test"))
				assert.Equal(t, markdown == nil || *markdown, strings.Contains(output, "Body text"))
				assert.Equal(t, markdown == nil || *markdown, strings.Contains(output, "<!-- InputPageNumber:"))
				for _, text := range []string{"mimeType: application/pdf", "metadata:", "customMetadata:", "pages: 1", "category: Invoice", "warnings:"} {
					assert.Contains(t, output, text)
				}
			})
		}
	}
	doc := llmDocument("", 1, 1)
	doc.Markdown = nil
	doc.Fields = map[string]sdk.ContentFieldClassification{"Name": llmString("Test")}
	assert.Equal(t, "---\nmimeType: application/pdf\npages: 1\nfields:\n  Name: Test\n---",
		renderLLM(t, llmResult(doc), nil, nil))
	assert.NotContains(t, renderLLM(t, llmResult(llmDocument("text", 1, 1)), nil, nil), "fields:")
}

func boolName(v *bool) string {
	if v == nil {
		return "default"
	}
	if *v {
		return "true"
	}
	return "false"
}

func TestToLLMInputMetadata(t *testing.T) {
	doc := llmDocument("text", 1, 1)
	doc.Metadata = map[string]*string{"title": to.Ptr("Report"), "author": to.Ptr("Jane"), "absent": nil}
	custom := map[string]any{
		"source": "invoice.pdf", "department": "finance", "mimeType": "caller",
		"metadata": "caller", "fields": "caller", "pages": "caller",
	}
	assert.Equal(t, "---\nmimeType: application/pdf\ncustomMetadata:\n"+
		"  department: finance\n  fields: caller\n  metadata: caller\n  mimeType: caller\n"+
		"  pages: caller\n  source: invoice.pdf\nmetadata:\n  author: Jane\n  title: Report\n"+
		"pages: 1\n---\n<!-- InputPageNumber: 1 -->\n\ntext",
		renderLLM(t, llmResult(doc), custom, nil))
	for _, metadata := range []map[string]*string{nil, {}, {"missing": nil}} {
		doc.Metadata = metadata
		output := renderLLM(t, llmResult(doc), nil, nil)
		assert.NotContains(t, output, "\nmetadata:")
		assert.NotContains(t, output, "customMetadata:")
	}
	assert.Contains(t, renderLLM(t, llmResult(doc), map[string]any{}, nil), "customMetadata: {}")
	assert.Contains(t, renderLLM(t, llmResult(doc), map[string]any{"missing": nil}, nil), "  missing: null")
	const opaque = `{"document":{"createdAt":"2026-07-16T19:00:00Z","tags":["finance","invoice"],"properties":{"pageCount":1}}}`
	doc.Metadata = map[string]*string{"xmp": to.Ptr(opaque)}
	output := renderLLM(t, llmResult(doc), nil, nil)
	assert.Contains(t, output, "  xmp: '"+opaque+"'")
	assert.NotContains(t, output, "\n  document:")
	assert.NotContains(t, output, "pageCount: 1")
}

func TestToLLMInputTypedMetadata(t *testing.T) {
	type label string
	type count uint64
	var nullString *string
	doc := llmDocument("", 0, 0)
	metadata := map[string]any{
		"values": []any{"first", nil, "third"}, "tags": [2]string{"finance", "invoice"},
		"properties": map[string]string{"department": "accounts-payable"},
		"emptyList":  []string{}, "emptyMap": map[string]string{},
		"missing": nullString, "nilMap": map[string]string(nil), "nilList": []string(nil),
		"": "empty-key", "emptyValue": "", "label": label("hello"),
		"large": count(math.MaxUint64), "signed": int64(9007199254740993),
		"raw": json.RawMessage(`{"missing":null,"list":[null,{},[],[1,2]],"integer":9007199254740993,"decimal":1.234567890123456789}`),
		"ptr": to.Ptr(42), "jsonNumber": json.Number("1e1000"),
		"records": []any{map[string]any{
			"list": []any{1, 2}, "nested": map[string]any{"x": true}, "tail": "end",
		}},
	}
	output := renderLLM(t, llmResult(doc), metadata, nil)
	for _, text := range []string{
		"  '': empty-key", "  emptyValue: ''", "  emptyList: []", "  emptyMap: {}",
		"  missing: null", "  nilMap: null", "  nilList: null", "  label: hello", "  ptr: 42",
		"  large: 18446744073709551615", "  signed: 9007199254740993", "  jsonNumber: 1e1000",
		"    integer: 9007199254740993", "    decimal: 1.234567890123456789", "    missing: null",
		"  values:\n  - first\n  - null\n  - third", "  tags:\n  - finance\n  - invoice",
		"  properties:\n    department: accounts-payable", "    - null\n    - {}\n    - []\n    -\n      - 1\n      - 2",
		"  records:\n  - list:\n      - 1\n      - 2\n    nested:\n      x: true\n    tail: end",
	} {
		assert.Contains(t, output, text)
	}
	expected, err := json.Marshal(map[string]any{"mimeType": "application/pdf", "customMetadata": metadata})
	require.NoError(t, err)
	assert.YAMLEq(t, string(expected), llmFrontMatter(t, output))
	// Repeated rendering neither mutates inputs nor depends on Go map iteration.
	for i := 0; i < 10; i++ {
		assert.Equal(t, output, renderLLM(t, llmResult(doc), metadata, nil))
	}
	assert.Equal(t, json.RawMessage(`{"missing":null,"list":[null,{},[],[1,2]],"integer":9007199254740993,"decimal":1.234567890123456789}`), metadata["raw"])
}

func TestToLLMInputYAMLScalars(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   any
		want string
	}{
		{"bool string", "true", "'true'"}, {"false string", "false", "'false'"},
		{"yes", "yes", "'yes'"}, {"no", "no", "'no'"}, {"on", "ON", "'ON'"},
		{"off", "off", "'off'"}, {"null string", "null", "'null'"},
		{"date", "2019-11-15", "'2019-11-15'"}, {"number string", "123", "'123'"},
		{"decimal string", "3.14", "'3.14'"}, {"exponent string", "1e6", "'1e6'"},
		{"plain", "hello", "hello"}, {"vendor", "CONTOSO LTD.", "CONTOSO LTD."},
		{"integer", int32(42), "42"}, {"whole", 165.0, "165"}, {"decimal", 3.14, "3.14"},
		{"scientific large", 1e20, "1E+20"}, {"scientific small", 1e-5, "1E-05"},
		{"fractional large", 1234567.89, "1234567.89"}, {"negative zero", math.Copysign(0, -1), "0"},
		{"int64 boundary", float64(0x1p63), "9.223372036854776E+18"},
		{"float32", float32(3.14), "3.14"}, {"bool", true, "true"}, {"false", false, "false"},
		{"infinity", math.Inf(1), "Infinity"}, {"negative infinity", math.Inf(-1), "-Infinity"},
		{"nan", math.NaN(), "NaN"}, {"dash", "- item", "'- item'"},
		{"comment", "# comment", "'# comment'"}, {"empty", "", "''"},
		{"apostrophe", "value: it's good", "'value: it''s good'"},
		{"space", " trailing ", "' trailing '"}, {"colon", "value:", "'value:'"},
		{"inline comment", "value #comment", "'value #comment'"},
		{"tilde", "~", "'~'"}, {"hex", "0xFF", "'0xFF'"},
		{"underscored number", "1_000", "'1_000'"},
		{"underscored decimal", "1.0_0", "'1.0_0'"},
		{"time string", "14:30:00", "14:30:00"}, {"yaml nan", ".nan", "'.nan'"},
		{"tab", "x\ty", `"x\ty"`}, {"carriage return", "x\ry", `"x\ry"`},
		{"control", "x\x00y", `"x\u0000y"`},
		{"NEL", "x\u0085---\ny", `"x\u0085---\ny"`},
		{"forbidden literal", "x\uffffy", `"x\uffffy"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			output := renderLLM(t, llmResult(llmDocument("", 0, 0)), map[string]any{"value": tt.in}, nil)
			assert.Contains(t, output, "  value: "+tt.want+"\n")
			expectedValue := tt.want
			if s, ok := tt.in.(string); ok {
				expectedValue = strconv.QuoteToASCII(s)
			}
			assert.YAMLEq(t, "mimeType: application/pdf\ncustomMetadata:\n  value: "+expectedValue, llmFrontMatter(t, output))
		})
	}
}

func TestToLLMInputYAMLKeysAndMultilineSafety(t *testing.T) {
	doc := llmDocument("Document body", 1, 3)
	doc.Metadata = map[string]*string{"description": to.Ptr("Q3 notes\n---\nreviewer: bob"), "author": to.Ptr("Jane")}
	custom := map[string]any{
		"comments":    "First section\n---\nSecond section",
		"with: colon": "value1", "with# hash": "value2", "- dash_start": "value3", "normal_key": "value4",
		"nested":              []any{map[string]any{"text": "One\n---\nTwo"}, []any{"One\n---\nTwo"}},
		"malicious\n---\nkey": "safe",
	}
	output := renderLLM(t, llmResult(doc), custom, nil)
	for _, text := range []string{
		"  description: 'Q3 notes\n    ---\n    reviewer: bob'",
		"  comments: 'First section\n    ---\n    Second section'",
		"'with: colon': value1", "'with# hash': value2", "'- dash_start': value3", "normal_key: value4",
		`"malicious\n---\nkey": safe`,
	} {
		assert.Contains(t, output, text)
	}
	assert.Equal(t, 2, strings.Count("\n"+output+"\n", "\n---\n"))
	assert.Contains(t, output, "Document body")
	assert.YAMLEq(t, `mimeType: application/pdf
customMetadata:
  comments: First section --- Second section
  "with: colon": value1
  "with# hash": value2
  "- dash_start": value3
  normal_key: value4
  nested:
    - text: One --- Two
    - [One --- Two]
  "malicious\n---\nkey": safe
metadata:
  description: "Q3 notes --- reviewer: bob"
  author: Jane
pages: 1-3`, llmFrontMatter(t, output))
}

func TestToLLMInputFields(t *testing.T) {
	date := time.Date(2019, 11, 15, 23, 30, 1, 0, time.FixedZone("local", -8*3600))
	doc := llmDocument("", 1, 1)
	doc.Fields = map[string]sdk.ContentFieldClassification{
		"TotalAmount": &sdk.ObjectField{ValueObject: map[string]sdk.ContentFieldClassification{
			"Amount": &sdk.NumberField{ValueNumber: to.Ptr(165.0)}, "CurrencyCode": llmString("USD"),
		}},
		"Topics": &sdk.ArrayField{ValueArray: []sdk.ContentFieldClassification{
			llmString("Topic A"), nil, &sdk.StringField{}, llmString("Topic B"),
		}},
		"InvoiceDate": &sdk.DateField{ValueDate: &date}, "MeetingTime": &sdk.TimeField{ValueTime: &date},
		"IsVerified": &sdk.BooleanField{ValueBoolean: to.Ptr(true)},
		"Quantity":   &sdk.IntegerField{ValueInteger: to.Ptr(int64(9007199254740993))},
		"Present":    llmString("exists"), "Empty": llmString(""),
		"Data":                &sdk.JSONField{ValueJSON: []byte(`{"key":"val","items":[1,2,null,{},[],[3,4]],"active":true,"nested":{"score":3.5,"missing":null},"missing":null,"large":9007199254740993}`)},
		"ResolvedEmptyObject": &sdk.ObjectField{ValueObject: map[string]sdk.ContentFieldClassification{"null": &sdk.StringField{}}},
	}
	for _, nullField := range []struct {
		name  string
		field sdk.ContentFieldClassification
	}{
		{"Missing", &sdk.StringField{}}, {"Nil", nil}, {"TypedNil", (*sdk.StringField)(nil)},
		{"EmptyObject", &sdk.ObjectField{}}, {"EmptyArray", &sdk.ArrayField{}},
		{"NullArray", &sdk.ArrayField{ValueArray: []sdk.ContentFieldClassification{&sdk.StringField{}}}},
		{"NullJSON", &sdk.JSONField{ValueJSON: []byte("null")}}, {"AbsentJSON", &sdk.JSONField{}},
		{"NullInteger", &sdk.IntegerField{}}, {"NullNumber", &sdk.NumberField{}}, {"NullDate", &sdk.DateField{}},
		{"NullTime", &sdk.TimeField{}}, {"NullBoolean", &sdk.BooleanField{}}, {"Unknown", &sdk.ContentField{}},
	} {
		doc.Fields[nullField.name] = nullField.field
	}
	output := renderLLM(t, llmResult(doc), nil, nil)
	for _, text := range []string{
		"  InvoiceDate: '2019-11-15'", "  MeetingTime: 23:30:01", "  IsVerified: true",
		"  Quantity: 9007199254740993", "  Present: exists", "  Empty: ''", "  ResolvedEmptyObject: {}",
		"  TotalAmount:\n    Amount: 165\n    CurrencyCode: USD",
		"  Topics:\n  - Topic A\n  - Topic B", "    active: true", "    key: val",
		"    large: 9007199254740993", "    nested:\n      score: 3.5",
		"    items:\n    - 1\n    - 2\n    - null\n    - {}\n    - []\n    -\n      - 3\n      - 4",
	} {
		assert.Contains(t, output, text)
	}
	for _, text := range []string{"Missing:", "Nil:", "\n  EmptyObject:", "EmptyArray:", "Null", "AbsentJSON:", "Unknown:", "missing:"} {
		assert.NotContains(t, output, text)
	}
	assert.YAMLEq(t, `{
  "mimeType": "application/pdf",
  "pages": 1,
  "fields": {
    "Data": {
      "key": "val", "items": [1, 2, null, {}, [], [3, 4]], "active": true,
      "nested": {"score": 3.5}, "large": 9007199254740993
    },
    "Empty": "", "InvoiceDate": "2019-11-15", "IsVerified": true,
    "MeetingTime": "23:30:01", "Present": "exists", "Quantity": 9007199254740993,
    "ResolvedEmptyObject": {}, "Topics": ["Topic A", "Topic B"],
    "TotalAmount": {"Amount": 165, "CurrencyCode": "USD"}
  }
}`, llmFrontMatter(t, output))
	assert.Contains(t, string(doc.Fields["Data"].(*sdk.JSONField).ValueJSON), `"missing":null`)
	empty := llmDocument("", 0, 0)
	empty.Fields = map[string]sdk.ContentFieldClassification{"Missing": &sdk.StringField{}}
	assert.NotContains(t, renderLLM(t, llmResult(empty), nil, nil), "fields:")
}

func TestToLLMInputGeneratedJSONModels(t *testing.T) {
	var result sdk.AnalysisResult
	require.NoError(t, json.Unmarshal([]byte(`{
		"stringEncoding":"codePoint",
		"contents":[{"kind":"document","mimeType":"application/pdf","markdown":"你好 🌍",
			"startPageNumber":2,"endPageNumber":2,"metadata":{"xmp":"{\"pageCount\":1}"},
			"fields":{"Amount":{"type":"integer","valueInteger":9007199254740993},
				"Data":{"type":"json","valueJson":{"large":18446744073709551615,"null":null}},
				"Date":{"type":"date","valueDate":"2019-11-15"},
				"Time":{"type":"time","valueTime":"14:30:00"}}}],
		"warnings":[{"code":"InvalidField","message":"Field is invalid.","target":"invoice.total"}]
	}`), &result))
	output := renderLLM(t, &result, nil, nil)
	for _, text := range []string{
		"Amount: 9007199254740993", "large: 18446744073709551615", "Date: '2019-11-15'",
		"Time: 14:30:00", `xmp: '{"pageCount":1}'`, "target: invoice.total", "<!-- InputPageNumber: 2 -->\n\n你好 🌍",
	} {
		assert.Contains(t, output, text)
	}
	assert.NotContains(t, output, "null:")
}

func TestToLLMInputWarnings(t *testing.T) {
	for _, tt := range []struct {
		name     string
		warnings []*sdk.Error
		want     string
	}{
		{"normal", []*sdk.Error{{Code: to.Ptr("hate"), Message: to.Ptr("Content flagged.")}},
			"warnings:\n- code: hate\n  message: Content flagged.\n"},
		{"target", []*sdk.Error{{Code: to.Ptr("InvalidField"), Message: to.Ptr("Field is invalid."), Target: to.Ptr("invoice.total")}},
			"warnings:\n- code: InvalidField\n  message: Field is invalid.\n  target: invoice.total\n"},
		{"telemetry", []*sdk.Error{{Code: to.Ptr("Telemetry"), Message: to.Ptr("LLMStats: completion calls: 2")}}, ""},
		{"whitespace", []*sdk.Error{{Message: to.Ptr(" \t\n\u2003LLMStats: counters")}}, ""},
		{"case sensitive", []*sdk.Error{{Message: to.Ptr("llmstats: keep as a real warning")}},
			"warnings:\n- message: 'llmstats: keep as a real warning'\n"},
		{"empty", []*sdk.Error{nil, {}, {Code: to.Ptr(""), Message: to.Ptr(""), Target: to.Ptr("")}}, ""},
		{"target only", []*sdk.Error{{Target: to.Ptr("field")}}, "warnings:\n- target: field\n"},
		{"mixed", []*sdk.Error{
			{Message: to.Ptr("LLMStats: remove")}, {Code: to.Ptr("ContentWarning"), Message: to.Ptr("Potentially sensitive content.")},
		}, "warnings:\n- code: ContentWarning\n  message: Potentially sensitive content.\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			doc := llmDocument("A log excerpt:\n- LLMStats: keep this body text", 1, 1)
			result := llmResult(doc)
			result.Warnings = tt.warnings
			output := renderLLM(t, result, nil, nil)
			if tt.want == "" {
				assert.NotContains(t, output, "warnings:")
			} else {
				assert.Contains(t, output, tt.want)
			}
			assert.NotContains(t, output, "LLMStats: remove")
			assert.Contains(t, output, "LLMStats: keep this body text")
		})
	}
}

func TestToLLMInputPages(t *testing.T) {
	for _, tt := range []struct {
		pages []int32
		want  string
	}{
		{[]int32{1}, "1"}, {[]int32{1, 2, 3}, "1-3"}, {[]int32{2, 3, 5}, "2-3, 5"},
		{[]int32{9, 5, 2, 1, 3, 8, 7}, "1-3, 5, 7-9"}, {[]int32{0, -1, 3}, "3"},
	} {
		t.Run(tt.want, func(t *testing.T) {
			doc := llmDocument("text", 2, 9)
			for _, number := range tt.pages {
				doc.Pages = append(doc.Pages, &sdk.DocumentPage{PageNumber: to.Ptr(number)})
			}
			assert.Contains(t, renderLLM(t, llmResult(doc), nil, nil), "\npages: "+tt.want+"\n")
		})
	}
	for _, tt := range []struct {
		start, end int32
		want       string
	}{
		{0, 0, ""}, {3, 3, "\npages: 3\n"}, {2, 5, "\npages: 2-5\n"},
	} {
		doc := llmDocument("text", tt.start, tt.end)
		doc.Pages = []*sdk.DocumentPage{nil, {PageNumber: to.Ptr(int32(0))}}
		output := renderLLM(t, llmResult(doc), nil, nil)
		if tt.want == "" {
			assert.NotContains(t, output, "pages:")
		} else {
			assert.Contains(t, output, tt.want)
		}
	}
	doc := llmDocument("Page 1 text\n\n<!-- PageBreak -->\n\nPage 2 text", 2, 3)
	output := renderLLM(t, llmResult(doc), nil, nil)
	assert.Equal(t, "---\nmimeType: application/pdf\npages: 2-3\n---\n"+
		"<!-- InputPageNumber: 2 -->\n\nPage 1 text\n\n<!-- InputPageNumber: 3 -->\n\nPage 2 text", output)
	assert.NotContains(t, output, "<!-- PageBreak -->")
	doc.Pages = []*sdk.DocumentPage{
		{PageNumber: to.Ptr(int32(2)), Spans: []*sdk.ContentSpan{llmSpan(0, 10)}},
		{PageNumber: to.Ptr(int32(3)), Spans: []*sdk.ContentSpan{llmSpan(30, 10)}},
	}
	output = renderLLM(t, llmResult(doc), nil, nil)
	assert.Equal(t, "---\nmimeType: application/pdf\npages: 2-3\n---\n"+
		"<!-- InputPageNumber: 2 -->\n\nPage 1 tex<!-- InputPageNumber: 3 -->\n\nt\n\nPage 2 text", output)
	assert.Equal(t, 1, strings.Count(output, "<!-- InputPageNumber: 2 -->"))
	assert.Equal(t, 1, strings.Count(output, "<!-- InputPageNumber: 3 -->"))
	assert.NotContains(t, output, "<!-- PageBreak -->")
	// Service-provided markers win over spans and break-based injection.
	body := "<!-- InputPageNumber: 2 -->\n\nFirst.\n\n<!-- InputPageNumber: 5 -->\n\nSecond."
	doc.Markdown = &body
	assert.True(t, strings.HasSuffix(renderLLM(t, llmResult(doc), nil, nil), "\n"+body))
}

func TestToLLMInputUnicodePageSpans(t *testing.T) {
	first, separator, second := "你好 🌍", "\n\n<!-- PageBreak -->\n\n", "café 🚀"
	md := first + separator + second + "third"
	for _, encoding := range []string{"", "codePoint", "utf8", "utf16"} {
		t.Run(encoding, func(t *testing.T) {
			doc := llmDocument(md, 2, 5)
			// Intentionally unsorted and non-contiguous original source pages.
			doc.Pages = []*sdk.DocumentPage{
				{PageNumber: to.Ptr(int32(5)), Spans: []*sdk.ContentSpan{llmSpan(spanUnits(first+separator+second, encoding), 5)}},
				{PageNumber: to.Ptr(int32(2)), Spans: []*sdk.ContentSpan{llmSpan(0, spanUnits(first, encoding))}},
				{PageNumber: to.Ptr(int32(3)), Spans: []*sdk.ContentSpan{llmSpan(spanUnits(first+separator, encoding), spanUnits(second, encoding))}},
			}
			result := llmResult(doc)
			result.StringEncoding = to.Ptr(encoding)
			assert.Equal(t, "---\nmimeType: application/pdf\npages: 2-3, 5\n---\n"+
				"<!-- InputPageNumber: 2 -->\n\n"+first+"\n\n<!-- InputPageNumber: 3 -->\n\n"+
				second+"<!-- InputPageNumber: 5 -->\n\nthird", renderLLM(t, result, nil, nil))
		})
	}
	// Oversized page offsets are clamped, matching reference page rendering.
	doc := llmDocument("text", 2, 5)
	for i, number := range []int32{2, 3, 5} {
		doc.Pages = append(doc.Pages, &sdk.DocumentPage{PageNumber: to.Ptr(number), Spans: []*sdk.ContentSpan{llmSpan(int32(i*5), 5)}})
	}
	assert.Contains(t, renderLLM(t, llmResult(doc), nil, nil), "pages: 2-3, 5")
}

func TestToLLMInputMedia(t *testing.T) {
	for _, mime := range []string{"audio/mpeg", "audio/wav", "video/mp4"} {
		av := &sdk.AudioVisualContent{MimeType: &mime, Markdown: to.Ptr("Speaker 1: Hello"),
			StartTimeMs: to.Ptr(int64(0)), EndTimeMs: to.Ptr(int64(23000)),
			Fields: map[string]sdk.ContentFieldClassification{"Summary": llmString("A short call.")}}
		output := renderLLM(t, llmResult(av), nil, nil)
		assert.Equal(t, "---\nmimeType: "+mime+"\nfields:\n  Summary: A short call.\n---\nSpeaker 1: Hello", output)
		assert.NotContains(t, output, "timeRange")
		assert.NotContains(t, output, "pages:")
	}
	av := func(start, end int64, md string) *sdk.AudioVisualContent {
		return &sdk.AudioVisualContent{
			MimeType: to.Ptr("video/mp4"), Markdown: &md, StartTimeMs: &start, EndTimeMs: &end,
		}
	}
	output := renderLLM(t, llmResult(av(0, 23099, "Opening"), av(24000, 43000, "Continuing")), nil, nil)
	assert.Equal(t, "---\nmimeType: video/mp4\ntimeRange: 00:00 – 00:23\n---\nOpening\n\n*****\n\n"+
		"---\nmimeType: video/mp4\ntimeRange: 00:24 – 00:43\n---\nContinuing", output)
	output = renderLLM(t, llmResult(av(300000, 630000, "Later"), av(0, 10000, "Earlier")), nil, nil)
	assert.Contains(t, output, "timeRange: 05:00 – 10:30")
	assert.Less(t, strings.Index(output, "Later"), strings.Index(output, "Earlier"))
	output = renderLLM(t, llmResult(av(3600000, 3660000, "hour"), llmDocument("doc", 1, 1)), nil, nil)
	assert.NotContains(t, output, "timeRange:") // Count AV blocks, not all contents.
	output = renderLLM(t, llmResult(av(3600000, 3660000, "hour"), av(0, 0, "zero")), nil, nil)
	assert.Contains(t, output, "timeRange: 60:00 – 61:00")
}

func TestToLLMInputFrontMatterOrderAndServiceOrder(t *testing.T) {
	doc := llmDocument("First", 2, 2)
	doc.Category = to.Ptr("Invoice")
	doc.Metadata = map[string]*string{"author": to.Ptr("Jane")}
	doc.Fields = map[string]sdk.ContentFieldClassification{"X": llmString("value")}
	result := llmResult(doc, llmDocument("Second", 1, 1))
	result.Warnings = []*sdk.Error{{Code: to.Ptr("warning")}}
	output := renderLLM(t, result, map[string]any{"source": "invoice.pdf"}, nil)
	prev := -1
	for _, key := range []string{"mimeType:", "customMetadata:", "\nmetadata:", "category:", "pages:", "fields:", "warnings:"} {
		pos := strings.Index(output, key)
		assert.Greater(t, pos, prev, key)
		prev = pos
	}
	assert.Less(t, strings.Index(output, "First"), strings.Index(output, "Second"))
	assert.Equal(t, 1, strings.Count(output, "\n\n*****\n\n"))
	assert.Equal(t, 2, strings.Count(output, "warnings:"))
}

func TestToLLMInputClassification(t *testing.T) {
	for _, encoding := range []string{"", "codePoint", "utf8", "utf16"} {
		for _, routing := range []string{"none", "partial", "all"} {
			t.Run(encoding+"/"+routing, func(t *testing.T) {
				first, second, third := "你好 🌍", "Bank text.", "Inv2."
				parent := llmDocument(first+"\n\n"+second+"\n\n"+third, 2, 5)
				parent.Path = to.Ptr("input1")
				parent.Metadata = map[string]*string{"parentMetadata": to.Ptr("not inherited")}
				parent.Fields = map[string]sdk.ContentFieldClassification{"ParentField": llmString("should be skipped")}
				parent.Segments = []*sdk.DocumentContentSegment{
					{SegmentID: to.Ptr("seg1"), Category: to.Ptr("Invoice"), StartPageNumber: to.Ptr(int32(2)),
						EndPageNumber: to.Ptr(int32(2)), Span: llmSpan(0, spanUnits(first, encoding))},
					{SegmentID: to.Ptr("seg2"), Category: to.Ptr("BankStatement"), StartPageNumber: to.Ptr(int32(3)),
						EndPageNumber: to.Ptr(int32(3)), Span: llmSpan(spanUnits(first+"\n\n", encoding), spanUnits(second, encoding))},
					{SegmentID: to.Ptr("seg3"), Category: to.Ptr("Invoice"), StartPageNumber: to.Ptr(int32(5)),
						EndPageNumber: to.Ptr(int32(5)), Span: llmSpan(spanUnits(first+"\n\n"+second+"\n\n", encoding), spanUnits(third, encoding))},
				}
				result := llmResult(parent)
				result.StringEncoding = to.Ptr(encoding)
				if routing != "none" {
					child := llmDocument("Own invoice markdown", 2, 2)
					child.Path, child.Category = to.Ptr("input1/seg1"), to.Ptr("Invoice")
					child.Fields = map[string]sdk.ContentFieldClassification{"Vendor": llmString("CONTOSO")}
					result.Contents = append(result.Contents, child)
				}
				if routing == "all" {
					for _, seg := range parent.Segments[1:] {
						child := llmDocument("Own "+*seg.Category, *seg.StartPageNumber, *seg.EndPageNumber)
						child.Path, child.Category = to.Ptr("input1/"+*seg.SegmentID), seg.Category
						result.Contents = append(result.Contents, child)
					}
				}
				before, err := json.Marshal(result)
				require.NoError(t, err)
				output := renderLLM(t, result, nil, nil)
				blocks := strings.Split(output, "\n\n*****\n\n")
				require.Len(t, blocks, 3)
				for i, category := range []string{"Invoice", "BankStatement", "Invoice"} {
					assert.Contains(t, blocks[i], "category: "+category)
				}
				for _, text := range []string{"ParentField", "should be skipped", "parentMetadata"} {
					assert.NotContains(t, output, text)
				}
				if routing == "none" {
					assert.Contains(t, blocks[0], first)
				} else {
					assert.Contains(t, blocks[0], "Own invoice markdown")
					assert.Contains(t, blocks[0], "Vendor: CONTOSO")
					assert.NotContains(t, blocks[0], first)
				}
				if routing != "all" {
					assert.Contains(t, blocks[1], second)
					assert.Contains(t, blocks[2], third)
					assert.NotContains(t, blocks[2], "fields:")
				}
				after, err := json.Marshal(result)
				require.NoError(t, err)
				assert.Equal(t, string(before), string(after))
			})
		}
	}
}

func TestToLLMInputClassificationStableSortingAndPaths(t *testing.T) {
	parent := llmDocument("Synthetic", 1, 1)
	parent.Path = to.Ptr("input1")
	parent.Segments = []*sdk.DocumentContentSegment{{
		SegmentID: to.Ptr("seg"), Category: to.Ptr("Invoice"), StartPageNumber: to.Ptr(int32(1)),
		EndPageNumber: to.Ptr(int32(1)), Span: llmSpan(0, 9),
	}}
	child := llmDocument("Routed other input", 1, 1)
	child.Path, child.Category = to.Ptr("input2/seg"), to.Ptr("Invoice")
	output := renderLLM(t, llmResult(parent, child, llmDocument("Last equal page", 1, 1)), nil, nil)
	assert.Equal(t, 2, strings.Count(output, "\n\n*****\n\n"))
	assert.Less(t, strings.Index(output, "Synthetic"), strings.Index(output, "Routed other input"))
	assert.Less(t, strings.Index(output, "Routed other input"), strings.Index(output, "Last equal page"))
	child.Path = nil // Matching categories alone do not deduplicate.
	assert.Contains(t, renderLLM(t, llmResult(parent, child), nil, nil), "Synthetic")
	parent.Category = to.Ptr("Already classified")
	assert.Contains(t, renderLLM(t, llmResult(parent), nil, nil), "category: Already classified")
	parent.Category = nil
	parent.Segments[0].SegmentID = nil
	assert.Contains(t, renderLLM(t, llmResult(parent), nil, nil), "Synthetic")
	parent.Segments[0].Span = nil
	assert.NotContains(t, renderLLM(t, llmResult(parent), nil, nil), "Synthetic")
}

func spanUnits(s, encoding string) int32 {
	switch encoding {
	case "utf8":
		return int32(len(s))
	case "utf16":
		return int32(len(utf16.Encode([]rune(s))))
	default:
		return int32(utf8.RuneCountInString(s))
	}
}

func TestSliceContentSpanUnicode(t *testing.T) {
	const text = "Aé中🌍e\u0301🚀Z"
	for _, encoding := range []string{"", "codePoint", "utf8", "utf16"} {
		t.Run(encoding, func(t *testing.T) {
			runes := []rune(text)
			for start := 0; start <= len(runes); start++ {
				for end := start; end <= len(runes); end++ {
					offset := spanUnits(string(runes[:start]), encoding)
					length := spanUnits(string(runes[start:end]), encoding)
					got, err := sdk.SliceContentSpan(text, *llmSpan(offset, length), encoding)
					require.NoError(t, err)
					assert.Equal(t, string(runes[start:end]), got)
				}
			}
		})
	}
	got, err := sdk.SliceContentSpan("", *llmSpan(0, 0), "")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestSliceContentSpanInvalid(t *testing.T) {
	for _, tt := range []struct {
		name, text, encoding string
		span                 sdk.ContentSpan
		message              string
	}{
		{"offset missing", "text", "", sdk.ContentSpan{Length: to.Ptr(int32(1))}, "required"},
		{"length missing", "text", "", sdk.ContentSpan{Offset: to.Ptr(int32(0))}, "required"},
		{"negative offset", "text", "", *llmSpan(-1, 1), "nonnegative"},
		{"negative length", "text", "", *llmSpan(0, -1), "nonnegative"},
		{"past end", "text", "", *llmSpan(5, 0), "exceeds"},
		{"long length", "text", "", *llmSpan(0, 5), "exceeds"},
		{"overflow", "text", "", *llmSpan(math.MaxInt32, math.MaxInt32), "exceeds"},
		{"encoding", "text", "utf32", *llmSpan(0, 1), "unsupported"},
		{"encoding case", "text", "UTF8", *llmSpan(0, 1), "unsupported"},
		{"utf8 start split", "é", "utf8", *llmSpan(1, 0), "splits"},
		{"utf8 end split", "🌍", "utf8", *llmSpan(0, 3), "splits"},
		{"utf16 start split", "🌍", "utf16", *llmSpan(1, 0), "splits"},
		{"utf16 end split", "🌍", "utf16", *llmSpan(0, 1), "splits"},
		{"invalid text", "\xff", "utf8", *llmSpan(0, 1), "not valid UTF-8"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			output, err := sdk.SliceContentSpan(tt.text, tt.span, tt.encoding)
			require.ErrorContains(t, err, tt.message)
			assert.Empty(t, output)
		})
	}
}

func TestToLLMInputErrors(t *testing.T) {
	for _, tt := range []struct {
		name     string
		metadata any
	}{
		{"channel", make(chan int)}, {"struct", struct{ Name string }{"not a scalar"}},
		{"map keys", map[int]string{1: "value"}}, {"invalid raw JSON", json.RawMessage(`{"bad":`)},
		{"multiple JSON values", json.RawMessage(`1 2`)}, {"invalid number", json.Number("true")},
		{"number whitespace", json.Number("1 ")},
	} {
		t.Run(tt.name, func(t *testing.T) {
			output, err := sdk.ToLLMInput(llmResult(llmDocument("body", 1, 1)), map[string]any{"bad": tt.metadata}, nil)
			require.Error(t, err)
			assert.Empty(t, output)
		})
	}
	cyclic := map[string]any{}
	cyclic["cycle"] = cyclic
	_, err := sdk.ToLLMInput(llmResult(llmDocument("", 0, 0)), cyclic, nil)
	require.ErrorContains(t, err, "nesting")
	doc := llmDocument("body", 1, 1)
	doc.Fields = map[string]sdk.ContentFieldClassification{"bad": &sdk.JSONField{ValueJSON: []byte("{")}}
	_, err = sdk.ToLLMInput(llmResult(doc), nil, nil)
	require.ErrorContains(t, err, "invalid JSON")
	assert.Contains(t, renderLLM(t, llmResult(doc), nil, &sdk.LLMInputOptions{IncludeFields: to.Ptr(false)}), "body")
	doc.Fields = map[string]sdk.ContentFieldClassification{}
	obj := &sdk.ObjectField{ValueObject: doc.Fields}
	doc.Fields["cycle"] = obj
	_, err = sdk.ToLLMInput(llmResult(doc), nil, nil)
	require.ErrorContains(t, err, "nesting")
	_, err = sdk.ToLLMInput(llmResult((*sdk.DocumentContent)(nil)), nil, nil)
	require.ErrorContains(t, err, "must not be nil")
	doc = llmDocument("🌍", 1, 1)
	doc.Segments = []*sdk.DocumentContentSegment{{Span: llmSpan(0, 1)}}
	result := llmResult(doc)
	result.StringEncoding = to.Ptr("utf16")
	_, err = sdk.ToLLMInput(result, nil, nil)
	require.ErrorContains(t, err, "splits")
	assert.Contains(t, renderLLM(t, result, nil, &sdk.LLMInputOptions{IncludeMarkdown: to.Ptr(false)}), "mimeType:")
	doc.Segments = nil
	doc.Pages = []*sdk.DocumentPage{{PageNumber: to.Ptr(int32(1)), Spans: []*sdk.ContentSpan{llmSpan(1, 0)}}}
	_, err = sdk.ToLLMInput(result, nil, nil)
	require.ErrorContains(t, err, "splits")
}

// This optional differential check consumes captures produced by executing the
// reference .NET tests. Expected is the actual .NET output after sorting input
// dictionaries, not YAML reserialization of the output. Default tests remain
// hermetic; set the environment variable to a locally generated capture file.
func TestToLLMInputDotNetParity(t *testing.T) {
	path := os.Getenv("AZURE_CONTENT_UNDERSTANDING_DOTNET_LLM_CASES")
	if path == "" {
		t.Skip("set AZURE_CONTENT_UNDERSTANDING_DOTNET_LLM_CASES to run the .NET differential")
	}
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var cases []struct {
		Name, Kind, ScalarType, Expected string
		Result, Metadata, Value          json.RawMessage
		Options                          *sdk.LLMInputOptions
	}
	require.NoError(t, json.Unmarshal(data, &cases))
	names := make(map[string]bool)
	intentional := 0
	for _, tt := range cases {
		names[tt.Name] = true
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Kind == "null" {
				_, err := sdk.ToLLMInput(nil, nil, nil)
				require.Error(t, err)
				return
			}
			if tt.Kind == "scalar" {
				var value any
				require.NoError(t, json.Unmarshal(tt.Value, &value))
				switch tt.ScalarType {
				case "Double", "Single":
					if s, ok := value.(string); ok {
						value, err = strconv.ParseFloat(s, 64)
						require.NoError(t, err)
					}
					if tt.ScalarType == "Single" {
						value = float32(value.(float64))
					}
				case "Int32", "Int64":
					var number json.Number
					require.NoError(t, json.Unmarshal(tt.Value, &number))
					value, err = number.Int64()
					require.NoError(t, err)
				case "String", "Boolean":
				default:
					t.Fatalf("unsupported reference scalar type %s", tt.ScalarType)
				}
				output := renderLLM(t, llmResult(&sdk.AnalysisContent{}), map[string]any{"value": value}, nil)
				expected := tt.Expected
				if tt.ScalarType == "Double" && expected == "9223372036854775807" {
					// .NET's float-to-long cast saturates at this boundary,
					// and its subsequent double comparison rounds the long
					// back up to 2^63. Keep the actual number instead.
					require.Equal(t, float64(0x1p63), value)
					expected = "9.223372036854776E+18"
					intentional++
				}
				assert.Equal(t, "---\nmimeType: unknown\ncustomMetadata:\n  value: "+expected+"\n---", output)
				return
			}
			require.Equal(t, "result", tt.Kind)
			var result sdk.AnalysisResult
			require.NoError(t, json.Unmarshal(tt.Result, &result))
			var metadata map[string]any
			if len(tt.Metadata) > 0 {
				decoder := json.NewDecoder(strings.NewReader(string(tt.Metadata)))
				decoder.UseNumber()
				require.NoError(t, decoder.Decode(&metadata))
			}
			output := renderLLM(t, &result, metadata, tt.Options)
			assert.Equal(t, tt.Expected, output)
		})
	}
	assert.Len(t, names, 84, "current reference has 80 methods / 84 parameter-expanded cases")
	t.Logf("compared %d captures covering %d .NET cases; %d declared numeric-safety differences", len(cases), len(names), intentional)
}
