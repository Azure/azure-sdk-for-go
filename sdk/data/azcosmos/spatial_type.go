// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// SpatialType defines supported values for spatial index types in Spatial Indexes
type SpatialType string

const (
	// Represents a point.
	SpatialTypePoint SpatialType = "Point"
	// Represents a polygon.
	SpatialTypePolygon SpatialType = "Polygon"
	// Represents a line string.
	SpatialTypeLineString SpatialType = "LineString"
	// Represents a multi polygon.
	SpatialTypeMultiPolygon SpatialType = "MultiPolygon"
)

// Returns a list of available data types
func SpatialTypeValues() []SpatialType {
	return []SpatialType{SpatialTypePoint, SpatialTypePolygon, SpatialTypeLineString, SpatialTypeMultiPolygon}
}

func (c SpatialType) ToPtr() *SpatialType {
	return &c
}
