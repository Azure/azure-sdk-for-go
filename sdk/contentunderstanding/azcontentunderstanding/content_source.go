// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcontentunderstanding

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Point is a document grounding coordinate.
type Point struct{ X, Y float64 }

// BoundingBox is an axis-aligned grounding rectangle.
type BoundingBox struct{ X, Y, Width, Height float64 }

// ContentSource is one document or audio/visual grounding region.
// PageNumber/Polygon are document properties; Time is an audio/visual property.
type ContentSource struct {
	RawValue    string
	PageNumber  *int
	Polygon     []Point
	Time        *time.Duration
	BoundingBox *BoundingBox
}

// String returns the original service region expression.
func (s ContentSource) String() string { return s.RawValue }

// ParseContentSources parses D(page[,polygon...]) and AV(milliseconds[,x,y,w,h])
// regions separated by semicolons. Empty segments are ignored, matching .NET.
func ParseContentSources(source string) ([]ContentSource, error) {
	if source == "" {
		return nil, errors.New("source is required")
	}
	result := make([]ContentSource, 0)
	for _, segment := range strings.Split(source, ";") {
		segment = strings.TrimSpace(segment)
		if segment == "" {
			continue
		}
		parsed, err := parseContentSource(segment)
		if err != nil {
			return nil, err
		}
		result = append(result, parsed)
	}
	return result, nil
}

func parseContentSource(source string) (ContentSource, error) {
	result := ContentSource{RawValue: source}
	if !strings.HasSuffix(source, ")") {
		return result, errors.New("source must end in ')'")
	}
	var inner string
	switch {
	case strings.HasPrefix(source, "D("):
		inner = source[2 : len(source)-1]
		parts := strings.Split(inner, ",")
		page, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil || page < 1 {
			return result, errors.New("document page must be a positive integer")
		}
		result.PageNumber = &page
		if len(parts) == 1 {
			return result, nil
		}
		if len(parts) < 7 || len(parts)%2 != 1 {
			return result, errors.New("document polygon requires at least three coordinate pairs")
		}
		for i := 1; i < len(parts); i += 2 {
			x, err := parseCoordinate(parts[i])
			if err != nil {
				return result, err
			}
			y, err := parseCoordinate(parts[i+1])
			if err != nil {
				return result, err
			}
			result.Polygon = append(result.Polygon, Point{X: x, Y: y})
		}
		minX, minY, maxX, maxY := result.Polygon[0].X, result.Polygon[0].Y, result.Polygon[0].X, result.Polygon[0].Y
		for _, p := range result.Polygon {
			minX = math.Min(minX, p.X)
			minY = math.Min(minY, p.Y)
			maxX = math.Max(maxX, p.X)
			maxY = math.Max(maxY, p.Y)
		}
		result.BoundingBox = &BoundingBox{X: minX, Y: minY, Width: maxX - minX, Height: maxY - minY}
	case strings.HasPrefix(source, "AV("):
		inner = source[3 : len(source)-1]
		parts := strings.Split(inner, ",")
		if len(parts) != 1 && len(parts) != 5 {
			return result, errors.New("audio/visual source requires one or five integers")
		}
		values := make([]int64, len(parts))
		for i, part := range parts {
			value, err := strconv.ParseInt(strings.TrimSpace(part), 10, 32)
			if err != nil {
				return result, fmt.Errorf("invalid audio/visual integer: %w", err)
			}
			values[i] = value
		}
		timestamp := time.Duration(values[0]) * time.Millisecond
		result.Time = &timestamp
		if len(values) == 5 {
			result.BoundingBox = &BoundingBox{X: float64(values[1]), Y: float64(values[2]), Width: float64(values[3]), Height: float64(values[4])}
		}
	default:
		return result, errors.New("unrecognized content source prefix")
	}
	return result, nil
}

func parseCoordinate(value string) (float64, error) {
	result, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil || math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, errors.New("invalid document coordinate")
	}
	return result, nil
}

// FormatContentSources reconstructs the service source expression.
func FormatContentSources(sources []ContentSource) (string, error) {
	if sources == nil {
		return "", errors.New("sources are required")
	}
	parts := make([]string, len(sources))
	for i, source := range sources {
		parts[i] = source.RawValue
	}
	return strings.Join(parts, ";"), nil
}

// GetFieldSources returns parsed grounding, or nil when the field has no source.
func GetFieldSources(field ContentFieldClassification) ([]ContentSource, error) {
	if field == nil || field.GetContentField() == nil {
		return nil, errors.New("field is required")
	}
	source := field.GetContentField().Source
	if source == nil || *source == "" {
		return nil, nil
	}
	return ParseContentSources(*source)
}
