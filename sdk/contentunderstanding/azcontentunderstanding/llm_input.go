// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// LLMInputOptions controls which parts of an analysis result are rendered.
// Nil options and nil individual flags default to true.
type LLMInputOptions struct {
	// IncludeFields includes extracted structured fields in the front matter.
	IncludeFields *bool
	// IncludeMarkdown includes the markdown body and document page markers.
	IncludeMarkdown *bool
}

// ToLLMInput formats analysis contents as YAML front matter followed by markdown.
// Blocks are separated by "\n\n*****\n\n". Front-matter keys are ordered as
// mimeType, customMetadata, metadata, timeRange, category, pages, fields, warnings.
// Unlike insertion-ordered dictionaries, Go maps have no defined order; all
// caller, service, field and JSON maps are therefore sorted lexicographically.
//
// A nil customMetadata omits that block; a non-nil empty map emits {}. Nested
// metadata supports string-keyed maps, slices, arrays, pointers, JSON raw messages,
// json.Number and scalar strings, booleans and numbers. Nulls and empty containers
// are preserved. Unsupported values, cyclic or excessively nested data, malformed
// JSON and invalid used spans return errors. Service metadata strings remain opaque.
// Inputs are not modified. Spans use result.StringEncoding (default codePoint).
// Numeric values are not saturated to int64: for example, a float64 equal to
// 2^63 is emitted in scientific notation rather than as math.MaxInt64.
// Unlike the reference's minimal emitter, ambiguous YAML scalars and multiline
// keys are escaped, nested sequences and empty sequence objects retain their
// structure, and raw JSON numbers retain their original precision.
func ToLLMInput(result *AnalysisResult, customMetadata map[string]any, options *LLMInputOptions) (string, error) {
	if result == nil {
		return "", fmt.Errorf("result must not be nil")
	}
	if len(result.Contents) == 0 {
		return "", nil
	}
	includeFields, includeMarkdown := true, true
	if options != nil {
		if options.IncludeFields != nil {
			includeFields = *options.IncludeFields
		}
		if options.IncludeMarkdown != nil {
			includeMarkdown = *options.IncludeMarkdown
		}
	}
	encoding := llmValue(result.StringEncoding)
	contents, err := llmRenderableContents(result.Contents, encoding, includeMarkdown)
	if err != nil {
		return "", err
	}
	var custom any
	if customMetadata != nil {
		custom, err = llmNormalizeMetadata(reflect.ValueOf(customMetadata), 0)
		if err != nil {
			return "", fmt.Errorf("customMetadata: %w", err)
		}
	}
	var warnings []any
	for _, w := range result.Warnings {
		if w == nil || strings.HasPrefix(strings.TrimLeftFunc(llmValue(w.Message), unicode.IsSpace), "LLMStats:") {
			continue
		}
		entry := make(map[string]any)
		for k, v := range map[string]*string{"code": w.Code, "message": w.Message, "target": w.Target} {
			if llmValue(v) != "" {
				entry[k] = *v
			}
		}
		if len(entry) > 0 {
			warnings = append(warnings, entry)
		}
	}
	avCount := 0
	for _, content := range contents {
		if _, ok := content.(*AudioVisualContent); ok {
			avCount++
		}
	}
	blocks := make([]string, 0, len(contents))
	for i, content := range contents {
		base := content.GetAnalysisContent()
		mimeType := llmValue(base.MimeType)
		if mimeType == "" {
			mimeType = "unknown"
		}
		front := []llmProperty{{"mimeType", mimeType}}
		if customMetadata != nil {
			front = append(front, llmProperty{"customMetadata", custom})
		}
		metadata := make(map[string]any)
		for k, v := range base.Metadata {
			if v != nil {
				metadata[k] = *v
			}
		}
		if len(metadata) > 0 {
			front = append(front, llmProperty{"metadata", metadata})
		}
		if av, ok := content.(*AudioVisualContent); ok && avCount > 1 {
			front = append(front, llmProperty{"timeRange", llmTime(llmValue(av.StartTimeMs)) + " – " + llmTime(llmValue(av.EndTimeMs))})
		}
		if llmValue(base.Category) != "" {
			front = append(front, llmProperty{"category", *base.Category})
		}
		if doc, ok := content.(*DocumentContent); ok {
			if pages := llmPages(doc); pages != nil {
				front = append(front, llmProperty{"pages", pages})
			}
		}
		if includeFields {
			fields, err := llmResolveFields(base.Fields, 0)
			if err != nil {
				return "", fmt.Errorf("contents[%d].fields: %w", i, err)
			}
			if len(fields) > 0 {
				front = append(front, llmProperty{"fields", fields})
			}
		}
		if len(warnings) > 0 {
			front = append(front, llmProperty{"warnings", warnings})
		}
		lines := []string{"---"}
		llmEmitMapping(&lines, front, 0)
		lines = append(lines, "---")
		block := strings.Join(lines, "\n")
		if includeMarkdown && llmValue(base.Markdown) != "" {
			md := *base.Markdown
			if doc, ok := content.(*DocumentContent); ok {
				md, err = llmPageMarkers(doc, md, encoding)
				if err != nil {
					return "", fmt.Errorf("contents[%d].pages: %w", i, err)
				}
			}
			block += "\n" + md
		}
		blocks = append(blocks, block)
	}
	return strings.Join(blocks, "\n\n*****\n\n"), nil
}

// SliceContentSpan returns the substring identified by a service content span.
// Encoding is "codePoint" (Unicode scalar values), "utf8" (bytes), or "utf16"
// (UTF-16 code units). Empty encoding defaults to "codePoint". The input must be
// valid UTF-8. Missing, negative, out-of-range or split-character boundaries and
// unrecognized encodings return errors; zero-length spans at valid boundaries
// (including the end of the text) are allowed.
func SliceContentSpan(text string, span ContentSpan, encoding string) (string, error) {
	if span.Offset == nil || span.Length == nil {
		return "", fmt.Errorf("span offset and length are required")
	}
	offset, length := int64(*span.Offset), int64(*span.Length)
	if offset < 0 || length < 0 {
		return "", fmt.Errorf("span offset and length must be nonnegative")
	}
	start, err := llmByteOffset(text, offset, encoding, false)
	if err != nil {
		return "", fmt.Errorf("span offset: %w", err)
	}
	end, err := llmByteOffset(text, offset+length, encoding, false)
	if err != nil {
		return "", fmt.Errorf("span end: %w", err)
	}
	return text[start:end], nil
}

func llmByteOffset(text string, offset int64, encoding string, clamp bool) (int, error) {
	if encoding == "" {
		encoding = "codePoint"
	}
	if encoding != "codePoint" && encoding != "utf8" && encoding != "utf16" {
		return 0, fmt.Errorf("unsupported string encoding %q", encoding)
	}
	if !utf8.ValidString(text) {
		return 0, fmt.Errorf("text is not valid UTF-8")
	}
	if offset < 0 {
		return 0, fmt.Errorf("offset must be nonnegative")
	}
	var units int64
	for i, r := range text {
		if units == offset {
			return i, nil
		}
		step := int64(1)
		if encoding == "utf8" {
			step = int64(utf8.RuneLen(r))
		} else if encoding == "utf16" && r > 0xffff {
			step = 2
		}
		units += step
		if units > offset {
			return 0, fmt.Errorf("offset %d splits a character in %s", offset, encoding)
		}
	}
	if units == offset || clamp {
		return len(text), nil
	}
	return 0, fmt.Errorf("offset %d exceeds text length %d in %s", offset, units, encoding)
}

func llmValue[T any](p *T) (v T) {
	if p != nil {
		return *p
	}
	return v
}

func llmIsNil(v any) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

func llmRenderableContents(contents []AnalysisContentClassification, encoding string, includeMarkdown bool) ([]AnalysisContentClassification, error) {
	routed := make(map[string]bool)
	for i, content := range contents {
		if llmIsNil(content) || content.GetAnalysisContent() == nil {
			return nil, fmt.Errorf("contents[%d] must not be nil", i)
		}
		if doc, ok := content.(*DocumentContent); ok && llmValue(doc.Category) != "" && llmValue(doc.Path) != "" {
			routed[*doc.Path] = true
		}
	}
	result := make([]AnalysisContentClassification, 0, len(contents))
	expanded := false
	for i, content := range contents {
		doc, ok := content.(*DocumentContent)
		if !ok || len(doc.Segments) == 0 || llmValue(doc.Category) != "" {
			result = append(result, content)
			continue
		}
		expanded = true
		for j, seg := range doc.Segments {
			if seg == nil {
				return nil, fmt.Errorf("contents[%d].segments[%d] must not be nil", i, j)
			}
			if llmValue(seg.SegmentID) != "" && routed[llmValue(doc.Path)+"/"+*seg.SegmentID] {
				continue
			}
			child := &DocumentContent{
				MimeType: doc.MimeType, Category: seg.Category,
				StartPageNumber: seg.StartPageNumber, EndPageNumber: seg.EndPageNumber,
			}
			if includeMarkdown && llmValue(doc.Markdown) != "" && seg.Span != nil {
				md, err := SliceContentSpan(*doc.Markdown, *seg.Span, encoding)
				if err != nil {
					return nil, fmt.Errorf("contents[%d].segments[%d]: %w", i, j, err)
				}
				child.Markdown = &md
			}
			result = append(result, child)
		}
	}
	if expanded {
		page := func(c AnalysisContentClassification) int32 {
			if doc, ok := c.(*DocumentContent); ok {
				return llmValue(doc.StartPageNumber)
			}
			return 0
		}
		sort.SliceStable(result, func(i, j int) bool { return page(result[i]) < page(result[j]) })
	}
	return result, nil
}

func llmResolveFields(fields map[string]ContentFieldClassification, depth int) (map[string]any, error) {
	result := make(map[string]any)
	for k, f := range fields {
		v, err := llmResolveField(f, depth+1)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", k, err)
		}
		if v != nil {
			result[k] = v
		}
	}
	return result, nil
}

func llmResolveField(field ContentFieldClassification, depth int) (any, error) {
	if depth > 100 {
		return nil, fmt.Errorf("fields exceed maximum nesting depth (possible cycle)")
	}
	if llmIsNil(field) {
		return nil, nil
	}
	switch f := field.(type) {
	case *ObjectField:
		if len(f.ValueObject) > 0 {
			return llmResolveFields(f.ValueObject, depth)
		}
	case *ArrayField:
		items := make([]any, 0, len(f.ValueArray))
		for _, item := range f.ValueArray {
			v, err := llmResolveField(item, depth+1)
			if err != nil {
				return nil, err
			}
			if v != nil {
				items = append(items, v)
			}
		}
		if len(items) > 0 {
			return items, nil
		}
	case *JSONField:
		if len(f.ValueJSON) > 0 {
			return llmDecodeJSON(f.ValueJSON, false)
		}
	case *StringField:
		if f.ValueString != nil {
			return *f.ValueString, nil
		}
	case *BooleanField:
		if f.ValueBoolean != nil {
			return *f.ValueBoolean, nil
		}
	case *IntegerField:
		if f.ValueInteger != nil {
			return *f.ValueInteger, nil
		}
	case *NumberField:
		if f.ValueNumber != nil {
			return *f.ValueNumber, nil
		}
	case *DateField:
		if f.ValueDate != nil {
			return f.ValueDate.Format("2006-01-02"), nil
		}
	case *TimeField:
		if f.ValueTime != nil {
			return f.ValueTime.Format("15:04:05"), nil
		}
	case *ContentField:
		// Unknown service field types have no typed value to render.
	default:
		return nil, fmt.Errorf("unsupported content field %T", field)
	}
	return nil, nil
}

func llmDecodeJSON(data []byte, preserveNull bool) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var v any
	if err := decoder.Decode(&v); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("invalid JSON: expected a single value")
	}
	if !preserveNull {
		llmRemoveNullProperties(v)
	}
	return v, nil
}

func llmRemoveNullProperties(v any) {
	switch v := v.(type) {
	case map[string]any:
		for k, item := range v {
			if item == nil {
				delete(v, k)
			} else {
				llmRemoveNullProperties(item)
			}
		}
	case []any:
		for _, item := range v {
			llmRemoveNullProperties(item)
		}
	}
}

func llmNormalizeMetadata(v reflect.Value, depth int) (any, error) {
	if depth > 100 {
		return nil, fmt.Errorf("metadata exceeds maximum nesting depth (possible cycle)")
	}
	if !v.IsValid() {
		return nil, nil
	}
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		return llmNormalizeMetadata(v.Elem(), depth+1)
	}
	if raw, ok := v.Interface().(json.RawMessage); ok {
		if raw == nil {
			return nil, nil
		}
		return llmDecodeJSON(raw, true)
	}
	if n, ok := v.Interface().(json.Number); ok {
		if len(n) == 0 || (n[0] != '-' && (n[0] < '0' || n[0] > '9')) ||
			strings.TrimSpace(string(n)) != string(n) || !json.Valid([]byte(n)) {
			return nil, fmt.Errorf("invalid JSON number %q", n)
		}
		return n, nil
	}
	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("metadata map keys must be strings, got %s", v.Type())
		}
		if v.IsNil() {
			return nil, nil
		}
		result := make(map[string]any, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			item, err := llmNormalizeMetadata(iter.Value(), depth+1)
			if err != nil {
				return nil, fmt.Errorf("%q: %w", iter.Key().String(), err)
			}
			result[iter.Key().String()] = item
		}
		return result, nil
	case reflect.Slice, reflect.Array:
		if v.Kind() == reflect.Slice && v.IsNil() {
			return nil, nil
		}
		result := make([]any, v.Len())
		for i := range result {
			item, err := llmNormalizeMetadata(v.Index(i), depth+1)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %w", i, err)
			}
			result[i] = item
		}
		return result, nil
	case reflect.String:
		return v.String(), nil
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint(), nil
	case reflect.Float32:
		return float32(v.Float()), nil
	case reflect.Float64:
		return v.Float(), nil
	default:
		return nil, fmt.Errorf("unsupported metadata type %s", v.Type())
	}
}

const llmPagePrefix = "<!-- InputPageNumber:"

var llmPageBreak = regexp.MustCompile(`\n*<!-- PageBreak -->\n*`)

func llmPageMarkers(doc *DocumentContent, md, encoding string) (string, error) {
	if strings.Contains(md, llmPagePrefix) {
		return md, nil
	}
	type marker struct {
		offset int
		page   int32
	}
	var markers []marker
	for _, page := range doc.Pages {
		if page == nil || len(page.Spans) == 0 || page.Spans[0] == nil {
			continue
		}
		span := page.Spans[0]
		if span.Offset == nil {
			return "", fmt.Errorf("page span offset is required")
		}
		// The service can place a page boundary past the end of its markdown;
		// like the reference formatter, place such markers at the end.
		offset, err := llmByteOffset(md, int64(*span.Offset), encoding, true)
		if err != nil {
			return "", err
		}
		markers = append(markers, marker{offset, llmValue(page.PageNumber)})
	}
	if len(markers) == 0 {
		start := llmValue(doc.StartPageNumber)
		if start <= 0 {
			start = 1
		}
		var parts []string
		for i, chunk := range llmPageBreak.Split(md, -1) {
			if text := strings.TrimSpace(chunk); text != "" {
				parts = append(parts, fmt.Sprintf("%s %d -->\n\n%s", llmPagePrefix, int64(start)+int64(i), text))
			}
		}
		return strings.Join(parts, "\n\n"), nil
	}
	sort.SliceStable(markers, func(i, j int) bool { return markers[i].offset < markers[j].offset })
	matches := llmPageBreak.FindAllStringIndex(md, -1)
	clean := llmPageBreak.ReplaceAllString(md, "\n\n")
	adjust := func(offset int) int {
		shift := 0
		for _, match := range matches {
			if offset > match[0] {
				shift += match[1] - match[0] - 2
			}
		}
		return offset - shift
	}
	var out strings.Builder
	prev := 0
	for _, m := range markers {
		offset := adjust(m.offset)
		if offset < 0 || (offset < len(clean) && !utf8.RuneStart(clean[offset])) {
			return "", fmt.Errorf("adjusted page offset %d is not a valid character boundary", offset)
		}
		if offset > prev {
			out.WriteString(clean[prev:offset])
		}
		fmt.Fprintf(&out, "%s %d -->\n\n", llmPagePrefix, m.page)
		prev = offset
	}
	out.WriteString(clean[prev:])
	return out.String(), nil
}

func llmPages(doc *DocumentContent) any {
	var nums []int32
	for _, page := range doc.Pages {
		if page != nil && llmValue(page.PageNumber) > 0 {
			nums = append(nums, *page.PageNumber)
		}
	}
	if len(nums) > 0 {
		sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
		if len(nums) == 1 {
			return int64(nums[0])
		}
		var ranges []string
		start, prev := nums[0], nums[0]
		add := func() {
			if start == prev {
				ranges = append(ranges, strconv.FormatInt(int64(start), 10))
			} else {
				ranges = append(ranges, fmt.Sprintf("%d-%d", start, prev))
			}
		}
		for _, n := range nums[1:] {
			if int64(n) == int64(prev)+1 {
				prev = n
			} else {
				add()
				start, prev = n, n
			}
		}
		add()
		return strings.Join(ranges, ", ")
	}
	start, end := llmValue(doc.StartPageNumber), llmValue(doc.EndPageNumber)
	if start == 0 && end == 0 {
		return nil
	}
	if start == end {
		return int64(start)
	}
	return fmt.Sprintf("%d-%d", start, end)
}

func llmTime(ms int64) string {
	seconds := ms / 1000
	return fmt.Sprintf("%02d:%02d", seconds/60, seconds%60)
}

type llmProperty struct {
	key   string
	value any
}

func llmSortedProperties(m map[string]any) []llmProperty {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	result := make([]llmProperty, 0, len(keys))
	for _, k := range keys {
		result = append(result, llmProperty{k, m[k]})
	}
	return result
}

func llmEmitMapping(lines *[]string, properties []llmProperty, indent int) {
	for _, p := range properties {
		llmEmitProperty(lines, strings.Repeat(" ", indent), p, indent, indent)
	}
}

func llmEmitProperty(lines *[]string, prefix string, p llmProperty, indent, sequenceIndent int) {
	key := llmYAMLString(p.key)
	// Escaped keys cannot introduce new YAML entries or front-matter delimiters.
	if strings.ContainsAny(p.key, "\n\r") {
		key = strconv.Quote(p.key)
	}
	head := prefix + key + ":"
	switch v := p.value.(type) {
	case map[string]any:
		if len(v) == 0 {
			*lines = append(*lines, head+" {}")
		} else {
			*lines = append(*lines, head)
			llmEmitMapping(lines, llmSortedProperties(v), indent+2)
		}
	case []any:
		if len(v) == 0 {
			*lines = append(*lines, head+" []")
		} else {
			*lines = append(*lines, head)
			llmEmitSequence(lines, v, sequenceIndent)
		}
	default:
		*lines = append(*lines, head+" "+llmScalar(v, indent+2))
	}
}

func llmEmitSequence(lines *[]string, items []any, indent int) {
	prefix := strings.Repeat(" ", indent)
	for _, item := range items {
		switch v := item.(type) {
		case map[string]any:
			if len(v) == 0 {
				*lines = append(*lines, prefix+"- {}")
				continue
			}
			for i, p := range llmSortedProperties(v) {
				tag := prefix + "  "
				if i == 0 {
					tag = prefix + "- "
				}
				llmEmitProperty(lines, tag, p, indent+2, indent+4)
			}
		case []any:
			if len(v) == 0 {
				*lines = append(*lines, prefix+"- []")
			} else {
				*lines = append(*lines, prefix+"-")
				llmEmitSequence(lines, v, indent+2)
			}
		default:
			*lines = append(*lines, prefix+"- "+llmScalar(v, indent+2))
		}
	}
}

func llmScalar(value any, indent int) string {
	var scalar string
	switch v := value.(type) {
	case nil:
		return "null"
	case bool:
		return strconv.FormatBool(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return llmFloat(float64(v), 32)
	case float64:
		return llmFloat(v, 64)
	case json.Number:
		return string(v)
	case string:
		scalar = llmYAMLString(v)
	}
	return strings.ReplaceAll(scalar, "\n", "\n"+strings.Repeat(" ", indent))
}

func llmFloat(v float64, bits int) string {
	if math.IsInf(v, 1) {
		return "Infinity"
	}
	if math.IsInf(v, -1) {
		return "-Infinity"
	}
	if math.IsNaN(v) {
		return "NaN"
	}
	if v >= -0x1p63 && v < 0x1p63 && v == math.Trunc(v) {
		return strconv.FormatInt(int64(v), 10)
	}
	// Match the reference's invariant general format after its integer check.
	// Go's shortest general format switches to scientific notation earlier.
	if abs := math.Abs(v); abs >= 1e-4 && abs < 1e17 {
		return strconv.FormatFloat(v, 'f', -1, bits)
	}
	return strings.ReplaceAll(strconv.FormatFloat(v, 'g', -1, bits), "e", "E")
}

var (
	llmYAMLNumber = regexp.MustCompile(`^[+-]?(\d[\d_]*\.?[\d_]*|\.[\d_]+)([eE][+-]?[\d_]+)?$`)
	llmYAMLDate   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)
	llmYAMLTyped  = regexp.MustCompile(`(?i)^(~|true|false|yes|no|on|off|null|[+-]?\.(inf|nan)|[+-]?0[xob][0-9a-f_]+)$`)
)

func llmYAMLString(s string) string {
	quote := s == "" || strings.TrimSpace(s) != s || llmYAMLTyped.MatchString(s) ||
		llmYAMLNumber.MatchString(s) || llmYAMLDate.MatchString(s) ||
		strings.HasSuffix(s, ":") || strings.Contains(s, ": ") || strings.Contains(s, "# ") || strings.Contains(s, " #") ||
		strings.ContainsAny(s, "\n\r")
	if s != "" && strings.ContainsRune("-?:,[]{}#&*!|>'\"%@`", rune(s[0])) {
		quote = true
	}
	for _, r := range s {
		if (r < 0x20 && r != '\n') || (r >= 0x7f && r <= 0x9f) || r == '\u2028' || r == '\u2029' || r == '\ufffe' || r == '\uffff' {
			// JSON double-quoted strings are also YAML double-quoted scalars.
			encoded, _ := json.Marshal(s)
			quoted := string(encoded)
			// JSON permits C1 characters literally, whereas YAML treats NEL as
			// a line break and rejects some of the other C1 controls.
			for _, control := range s {
				if (control >= 0x7f && control <= 0x9f) || control == '\ufffe' || control == '\uffff' {
					quoted = strings.ReplaceAll(quoted, string(control), fmt.Sprintf(`\u%04x`, control))
				}
			}
			return quoted
		}
	}
	if quote {
		return "'" + strings.ReplaceAll(s, "'", "''") + "'"
	}
	return s
}
