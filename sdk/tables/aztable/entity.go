// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"encoding/base64"
	"encoding/json"
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
	EditLink string `json:"odata.editLink"`
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
		case "odata.editLink": // TODO: Verify the casing that the the service returns
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
				var v EdmDateTime
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			case "Edm.Binary":
				var v EdmBinary
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			case "Edm.Guid":
				var v EdmGuid
				err = json.Unmarshal(propRawValue, &v)
				propValue = v
			case "Edm.Int64":
				var v EdmInt64
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
