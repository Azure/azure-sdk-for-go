// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ********************** The following goes in the aztables package *****************
// https://docs.microsoft.com/en-us/rest/api/storageservices/payload-format-for-table-service-operations

type Entity struct {
	PartitionKey string
	RowKey       string
	Timestamp    EdmDateTime
}

type EdmEntity struct {
	Metadata string `json:"odata.metadata"`
	Id       string `json:"odata.id"`
	EditLink string `json:"odata.editlink"`
	Type     string `json:"odata.type"`
	Etag     string `json:"odata.etag"`
	Entity
	Properties map[string]interface{} // Type assert the value to 1 of these: bool, int32, float64, string, EdmDateTime, EdmBinary, EdmGuid, EdmInt64
}

func (e EdmEntity) MarshalJSON() ([]byte, error) {
	entity := map[string]interface{}{}
	entity["PartitionKey"], entity["RowKey"] = e.PartitionKey, e.RowKey

	for propName, propValue := range e.Properties {
		entity[propName] = propValue
		edmType := ""
		switch propValue.(type) {
		case EdmDateTime:
			edmType = "Edm.DateTime"
		case EdmBinary:
			edmType = "Edm.Binary"
		case EdmGuid:
			edmType = "Edm.Guid"
		case EdmInt64:
			edmType = "Edm.Int64"
		}
		if edmType != "" {
			entity[propName+"@odata.type"] = edmType
		}
	}
	return json.Marshal(entity)
}

func (e *EdmEntity) UnmarshalJSON(data []byte) (err error) {
	var entity map[string]json.RawMessage
	err = json.Unmarshal(data, &entity)
	if err != nil {
		return
	}
	e.Properties = map[string]interface{}{}
	for propName, propRawValue := range entity {
		if strings.Contains(propName, "@odata.type") {
			continue // Skip the @odata.type properties; we look them up explicitly later
		}
		switch propName {
		// Look for EdmEntity's specific fields first
		case "odata.metadata":
			err = json.Unmarshal(propRawValue, &e.Metadata)
		case "odata.id":
			err = json.Unmarshal(propRawValue, &e.Id)
		case "odata.editlink": // TODO: Verify the casing that the the service returns
			err = json.Unmarshal(propRawValue, &e.EditLink)
		case "odata.type":
			err = json.Unmarshal(propRawValue, &e.Type)
		case "odata.etag":
			err = json.Unmarshal(propRawValue, &e.Etag)
		case "PartitionKey":
			err = json.Unmarshal(propRawValue, &e.PartitionKey)
		case "RowKey":
			err = json.Unmarshal(propRawValue, &e.RowKey)
		case "Timestamp":
			err = json.Unmarshal(propRawValue, &e.Timestamp)
		default:
			// Try to find the EDM type for this property & get it's value
			var propertyEdmTypeValue string = ""
			if propertyEdmTypeRawValue, ok := entity[propName+"@odata.type"]; ok {
				if err = json.Unmarshal(propertyEdmTypeRawValue, &propertyEdmTypeValue); err != nil {
					return
				}
			}

			var propValue interface{} = nil
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
				propValue = &EdmDateTime{}
				err = json.Unmarshal(propRawValue, propValue)
			case "Edm.Binary":
				propValue = &EdmBinary{}
				err = json.Unmarshal(propRawValue, propValue)
			case "Edm.Guid":
				var v EdmGuid
				propValue = &v
				err = json.Unmarshal(propRawValue, propValue)
			case "Edm.Int64":
				var v EdmInt64
				propValue = &v
				err = json.Unmarshal(propRawValue, propValue)
				fmt.Println(propValue)
			}
			if err != nil {
				return
			}
			e.Properties[propName] = propValue
		}
	}
	return
}

type EdmBinary []byte

func (e EdmBinary) MarshalText() ([]byte, error) {
	return ([]byte)(base64.StdEncoding.EncodeToString(([]byte)(e))), nil
}

func (e *EdmBinary) UnmarshalText(data []byte) error {
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	*e = EdmBinary(decoded)
	return nil
}

type EdmInt64 int64

func (e EdmInt64) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(e), 10)), nil
}

func (e *EdmInt64) UnmarshalText(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*e = EdmInt64(i)
	return nil
}

type EdmGuid string

func (e EdmGuid) MarshalText() ([]byte, error) {
	return ([]byte)(e), nil
}

func (e *EdmGuid) UnmarshalText(data []byte) error {
	*e = EdmGuid(string(data))
	return nil
}

type EdmDateTime time.Time

const rfc3339 = "2006-01-02T15:04:05.9999999Z"

func (e EdmDateTime) MarshalText() ([]byte, error) {
	return ([]byte)(time.Time(e).Format(rfc3339)), nil
}

func (e *EdmDateTime) UnmarshalText(data []byte) error {
	t, err := time.Parse(rfc3339, string(data))
	if err != nil {
		return err
	}
	*e = EdmDateTime(t)
	return nil
}

// ********************** The following are examples of customer code *****************

func aztabletest() {
	jsonEntity := ([]byte)(`{
        "odata.type":"account.Customers",
        "odata.id":https://myaccount.table.core.windows.net/Customers(PartitionKey='Customer03',RowKey='Name'),
        "odata.etag":"W/\"0x5B168C7B6E589D2\"",
        "odata.editlink":"Customers(PartitionKey='Customer03',RowKey='Name')",
        "PartitionKey":"partitionkey",
        "RowKey":"rowkey",
        "Timestamp":"2013-08-09T18:55:48.3402073Z",
        "Bool":false,
        "Int32":1234,
        Int64@odata.type:"Edm.Int64",
        "Int64":"123456789012",
        "Double":1234.1234,
        "String":"test",
        Guid@odata.type:"Edm.Guid",
        "Guid":"4185404a-5818-48c3-b9be-f217df0dba6f",
        DateTime@odata.type:"Edm.DateTime",
        "DateTime":"2013-08-02T17:37:43.9004348Z",
        Binary@odata.type:"Edm.Binary",
        "Binary":"AQIDBA=="
    }`)

	// Customers can unmarshal the json entity in many ways:
	// 1. Unmarshal all properties to a map realizing all values
	var m1 map[string]interface{}
	err := json.Unmarshal(jsonEntity, &m1)
	_ = err
	fmt.Printf("%+v\n\n", m1)

	// Marshalling round-trips the exact data
	data, err := json.MarshalIndent(m1, "", "  ")
	fmt.Println(string(data) + "\n")

	// 2. Unmarshal all properties to a map without realizing any values
	var m2 map[string]json.RawMessage
	err = json.Unmarshal(([]byte)(jsonEntity), &m2)
	// Marshalling round-trips the exact data
	data, err = json.MarshalIndent(m2, "", "  ")
	fmt.Println(string(data))

	// 3. Unmarshal to an EdmEntity (adds some type-safe fields & maintains EDM type info)
	var edmEntity EdmEntity
	err = json.Unmarshal(([]byte)(jsonEntity), &edmEntity)
	data, err = json.MarshalIndent(edmEntity, "", "  ")
	fmt.Println(string(data))

	// Marshalling round-trips the data (except for response-only fields)
	// Customer can add/remove fields from client and set to nil to remove field from service
	edmEntity.Properties["New"] = EdmGuid("New guid")
	delete(edmEntity.Properties, "Bool")
	edmEntity.Properties["NullBool"] = nil
	data, err = json.MarshalIndent(edmEntity, "", "  ")
	fmt.Println(string(data))

	// 4. Unmarshal all properties to fields in a customer-defined struct
	type myEntity struct {
		Entity
		Bool         bool
		Int32        *int32 // Customer can use a pointer; if nil, null is sent
		Int32Null    *int32 `json:",omitempty"` // if nil, null not sent
		Int64        EdmInt64
		Int64EdmType string `json:"Int64@odata.type"` // Customer can do this but it's pretty annoying
		Double       float64
		String       string
		Guid         EdmGuid
		DateTime     *EdmDateTime
		Binary       EdmBinary
	}
	s := myEntity{
		Entity:       Entity{PartitionKey: "PK", RowKey: "RK", Timestamp: EdmDateTime(time.Now())},
		Bool:         true,
		Int32:        func(v int32) *int32 { return &v }(5), // nil sends null
		Int64:        654321,
		Int64EdmType: "Edm.Int64", // NOTE: The customer could do this but it's pretty annoying
		Double:       43.21,
		String:       "Some string",
		Guid:         EdmGuid("this is a guid"),
		DateTime:     func(v time.Time) *EdmDateTime { return (*EdmDateTime)(&v) }(time.Now()),
		Binary:       EdmBinary([]byte{4, 3, 2, 1}),
	}
	// Marshalling does NOT send the odata.type properties
	data, err = json.MarshalIndent(s, "", "  ")
	fmt.Println(string(data))

	// How to delete (pass a null) a field (String)
	data, err = json.MarshalIndent(struct {
		*myEntity
		String *string
	}{myEntity: &s}, "", "  ")
	fmt.Println(string(data))

	err = json.Unmarshal([]byte(jsonEntity), &s)
	data, err = json.MarshalIndent(s, "", "  ")
	fmt.Println(string(data))
}
