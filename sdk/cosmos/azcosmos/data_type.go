// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// DataType defines supported values for data types in Spatial Indexes
type DataType string

const (
	// Represents a line.
	DataTypeString DataType = "String"
	// Represents a number.
	DataTypeNumber DataType = "Number"
	// Represents a point.
	DataTypePoint DataType = "Point"
	// Represents a polygon.
	DataTypePolygon DataType = "Polygon"
	// Represents a line string.
	DataTypeLineString DataType = "LineString"
	// Represents a multi polygon.
	DataTypeMultiPolygon DataType = "MultiPolygon"
)

// Returns a list of available data types
func DataTypeValues() []DataType {
	return []DataType{DataTypeString, DataTypeNumber, DataTypePoint, DataTypePolygon, DataTypeLineString, DataTypeMultiPolygon}
}

func (c DataType) ToPtr() *DataType {
	return &c
}
