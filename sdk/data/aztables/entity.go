// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// https://learn.microsoft.com/rest/api/storageservices/payload-format-for-table-service-operations

// Entity is the bare minimum properties for a valid Entity. These should be embedded in a custom struct.
type Entity struct {
	PartitionKey string
	RowKey       string
	Timestamp    EDMDateTime
}

// EDMEntity is an entity that embeds the azcore.Entity type and has a Properties map for user defined entity properties
type EDMEntity struct {
	Entity
	Metadata string `json:"odata.metadata"`
	ID       string `json:"odata.id"`
	EditLink string `json:"odata.editLink"`
	Type     string `json:"odata.type"`
	ETag     string `json:"odata.etag"`

	// Properties contains user-defined entity properties and values.
	// The value can be one of the following types:
	//
	//   - bool
	//   - int32 for 32-bit numeric values without a decimal point
	//   - float64 for numeric values with a decimal point
	//   - string
	//   - EDMDateTime
	//   - EDMBinary
	//   - EDMGUID
	//   - EDMInt64 for 64-bit numeric values without a decimal point
	//
	// See https://learn.microsoft.com/rest/api/storageservices/payload-format-for-table-service-operations#property-types-in-a-json-feed
	Properties map[string]any
}

// MarshalJSON implements the json.Marshal method
func (e EDMEntity) MarshalJSON() ([]byte, error) {
	entity := map[string]any{}
	entity["PartitionKey"], entity["RowKey"] = prepareKey(e.PartitionKey), prepareKey(e.RowKey)

	for propName, propValue := range e.Properties {
		entity[propName] = propValue
		edmType := ""
		switch propValue.(type) {
		case EDMDateTime:
			edmType = "Edm.DateTime"
		case EDMBinary:
			edmType = "Edm.Binary"
		case EDMGUID:
			edmType = "Edm.Guid"
		case EDMInt64:
			edmType = "Edm.Int64"
		}
		if edmType != "" {
			entity[propName+"@odata.type"] = edmType
		}
	}
	return json.Marshal(entity)
}

// UnmarshalJSON implements the json.Unmarshal method
func (e *EDMEntity) UnmarshalJSON(data []byte) (err error) {
	var entity map[string]json.RawMessage
	err = json.Unmarshal(data, &entity)
	if err != nil {
		return
	}
	e.Properties = map[string]any{}
	for propName, propRawValue := range entity {
		if strings.Contains(propName, "@odata.type") {
			continue // Skip the @odata.type properties; we look them up explicitly later
		}
		switch propName {
		// Look for EDMEntity's specific fields first
		case "odata.metadata":
			err = json.Unmarshal(propRawValue, &e.Metadata)
		case "odata.id":
			err = json.Unmarshal(propRawValue, &e.ID)
		case "odata.editLink":
			err = json.Unmarshal(propRawValue, &e.EditLink)
		case "odata.type":
			err = json.Unmarshal(propRawValue, &e.Type)
		case "odata.etag":
			err = json.Unmarshal(propRawValue, &e.ETag)
		case "PartitionKey":
			err = json.Unmarshal(propRawValue, &e.PartitionKey)
		case "RowKey":
			err = json.Unmarshal(propRawValue, &e.RowKey)
		case "Timestamp":
			err = json.Unmarshal(propRawValue, &e.Timestamp)
		default:
			// Try to find the EDM type for this property & get it's value
			var propertyEdmTypeValue string
			if propertyEdmTypeRawValue, ok := entity[propName+"@odata.type"]; ok {
				if err = json.Unmarshal(propertyEdmTypeRawValue, &propertyEdmTypeValue); err != nil {
					return
				}
			}

			var propValue any = nil
			switch propertyEdmTypeValue {
			case "": // "<property>@odata.type" doesn't exist, infer the EDM type from the JSON type
				// Try to unmarshal this property value as an int32 first
				var i32 int32
				if err = json.Unmarshal(propRawValue, &i32); err == nil {
					propValue = i32
				} else { // Failed to parse number as an int32; unmarshal as usual
					err = json.Unmarshal(propRawValue, &propValue)
				}
			case "Edm.DateTime":
				var v EDMDateTime
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			case "Edm.Binary":
				var v EDMBinary
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			case "Edm.Guid":
				var v EDMGUID
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			case "Edm.Int64":
				var v EDMInt64
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			}
			if err != nil {
				return
			}
			e.Properties[propName] = propValue
		}
	}
	return
}

// EDMBinary represents an Entity Property that is a byte slice. A byte slice wrapped in
// EDMBinary will also receive the correct odata annotation for round-trip accuracy.
type EDMBinary []byte

// MarshalText implements the encoding.TextMarshaler interface
func (e EDMBinary) MarshalText() ([]byte, error) {
	return ([]byte)(base64.StdEncoding.EncodeToString(([]byte)(e))), nil
}

// UnmarshalText implements the encoding.TextMarshaler interface
func (e *EDMBinary) UnmarshalText(data []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	*e = EDMBinary(decoded)
	return nil
}

// EDMInt64 represents an entity property that is a 64-bit integer. Using EDMInt64 guarantees
// proper odata type annotations.
type EDMInt64 int64

// MarshalText implements the encoding.TextMarshaler interface
func (e EDMInt64) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(e), 10)), nil
}

// UnmarshalText implements the encoding.TextMarshaler interface
func (e *EDMInt64) UnmarshalText(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*e = EDMInt64(i)
	return nil
}

// EDMGUID represents an entity property that is a GUID wrapped in a string. Using EDMGUID guarantees
// proper odata type annotations.
type EDMGUID string

// MarshalText implements the encoding.TextMarshaler interface
func (e EDMGUID) MarshalText() ([]byte, error) {
	return ([]byte)(e), nil
}

// UnmarshalText implements the encoding.TextMarshaler interface
func (e *EDMGUID) UnmarshalText(data []byte) error {
	*e = EDMGUID(string(data))
	return nil
}

// EDMDateTime represents an entity property that is a time.Time object. Using EDMDateTime guarantees
// proper odata type annotations.
type EDMDateTime time.Time

// MarshalText implements the encoding.TextMarshaler interface
func (e EDMDateTime) MarshalText() ([]byte, error) {
	return ([]byte)(time.Time(e).Format(rfc3339)), nil
}

// UnmarshalText implements the encoding.TextMarshaler interface
func (e *EDMDateTime) UnmarshalText(data []byte) error {
	t, err := time.Parse(rfc3339, string(data))
	if err != nil {
		return err
	}
	*e = EDMDateTime(t)
	return nil
}

func prepareKey(key string) string {
	// escape any single-quotes
	return strings.ReplaceAll(key, "'", "''")
}
