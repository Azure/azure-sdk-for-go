// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ContentRange is a page or millisecond range expression. A direct conversion
// preserves a raw service expression; constructors validate their numeric inputs.
type ContentRange string

func (r ContentRange) String() string { return string(r) }

// Ptr supplies a generated AnalysisInput.ContentRange or binary Range option.
func (r ContentRange) Ptr() *string { value := string(r); return &value }

// Page selects one 1-based page.
func Page(page int) (ContentRange, error) {
	if page < 1 {
		return "", errors.New("page must be at least 1")
	}
	return ContentRange(fmt.Sprint(page)), nil
}

// Pages selects an inclusive, 1-based page range.
func Pages(start, end int) (ContentRange, error) {
	if start < 1 || end < start {
		return "", errors.New("pages require start >= 1 and end >= start")
	}
	return ContentRange(fmt.Sprintf("%d-%d", start, end)), nil
}

// PagesFrom selects every page from start onward.
func PagesFrom(start int) (ContentRange, error) {
	if start < 1 {
		return "", errors.New("start page must be at least 1")
	}
	return ContentRange(fmt.Sprintf("%d-", start)), nil
}

// TimeRange selects an inclusive media interval, serialized in integer milliseconds.
func TimeRange(start, end time.Duration) (ContentRange, error) {
	if start < 0 || end < start {
		return "", errors.New("time range requires 0 <= start <= end")
	}
	return ContentRange(fmt.Sprintf("%d-%d", start.Milliseconds(), end.Milliseconds())), nil
}

// TimeRangeFrom selects media from start onward.
func TimeRangeFrom(start time.Duration) (ContentRange, error) {
	if start < 0 {
		return "", errors.New("start time must be nonnegative")
	}
	return ContentRange(fmt.Sprintf("%d-", start.Milliseconds())), nil
}

// CombineRanges joins page/time/raw expressions without changing their order.
func CombineRanges(ranges ...ContentRange) (ContentRange, error) {
	if len(ranges) == 0 {
		return "", errors.New("at least one range is required")
	}
	parts := make([]string, len(ranges))
	for i, value := range ranges {
		parts[i] = string(value)
	}
	return ContentRange(strings.Join(parts, ",")), nil
}
