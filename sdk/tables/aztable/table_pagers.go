// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

// TableEntityQueryResponsePager is a Pager for Table entity query results.
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
//
// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
type TableEntityQueryResponsePager interface {
	azcore.Pager

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() TableEntityQueryResponseResponse
}

type tableEntityQueryResponsePager struct {
	tableClient       *TableClient
	current           *TableEntityQueryResponseResponse
	tableQueryOptions *TableQueryEntitiesOptions
	queryOptions      *QueryOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
func (p *tableEntityQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextPartitionKey == nil && p.current.XMSContinuationNextRowKey == nil) {
		return false
	}
	var resp TableEntityQueryResponseResponse
	resp, p.err = p.tableClient.client.QueryEntities(ctx, p.tableClient.Name, p.tableQueryOptions, p.queryOptions)
	if p.err == nil {
		castAndRemoveAnnotationsSlice(&resp.TableEntityQueryResponse.Value)
	}
	p.current = &resp
	p.tableQueryOptions.NextPartitionKey = resp.XMSContinuationNextPartitionKey
	p.tableQueryOptions.NextRowKey = resp.XMSContinuationNextRowKey
	return p.err == nil && resp.TableEntityQueryResponse != nil && len(resp.TableEntityQueryResponse.Value) > 0
}

// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
func (p *tableEntityQueryResponsePager) PageResponse() TableEntityQueryResponseResponse {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableEntityQueryResponsePager) Err() error {
	return p.err
}

// TableQueryResponsePager is a Pager for Table Queries
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
//
// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
type TableQueryResponsePager interface {
	azcore.Pager

	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() TableQueryResponseResponse
}

// AsModels converts each map[string]interface{} entity result into a strongly slice of strongly typed models
// The modelSlice parameter should be a pointer to a slice of struct types that match the entity model type in the table response.
func (r *TableEntityQueryResponse) AsModels(modelSlice interface{}) error {
	models := reflect.ValueOf(modelSlice).Elem()
	tt := getTypeArray(models.Interface())
	fmap := getTypeValueMap(tt)
	for i, e := range r.Value {
		err := fromMap(tt, fmap, &e, models.Index(i))
		if err != nil {
			return nil
		}
	}

	return nil
}

// EntityMapAsModel converts a table entity in the form of map[string]interface{} and converts it to a strongly typed model.
//
// Example:
// mapEntity, err := client.GetEntity("somePartition", "someRow")
// myEntityModel := MyModel{}
// err = EntityMapAsModel(mapEntity, &myEntityModel)
func EntityMapAsModel(entityMap map[string]interface{}, model interface{}) error {
	tt := getTypeArray(model)
	fmap := getTypeValueMap(tt)
	err := fromMap(reflect.TypeOf(model).Elem(), fmap, &entityMap, reflect.ValueOf(model).Elem())
	if err != nil {
		return nil
	}
	return err
}

type tableQueryResponsePager struct {
	client            *tableClient
	current           *TableQueryResponseResponse
	tableQueryOptions *TableQueryOptions
	queryOptions      *QueryOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaulated by calling PageResponse on this Pager.
func (p *tableQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextTableName == nil) {
		return false
	}
	var resp TableQueryResponseResponse
	resp, p.err = p.client.Query(ctx, p.tableQueryOptions, p.queryOptions)
	p.current = &resp
	p.tableQueryOptions.NextTableName = resp.XMSContinuationNextTableName
	return p.err == nil && resp.TableQueryResponse.Value != nil && len(resp.TableQueryResponse.Value) > 0
}

// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
func (p *tableQueryResponsePager) PageResponse() TableQueryResponseResponse {
	return *p.current
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableQueryResponsePager) Err() error {
	return p.err
}

func castAndRemoveAnnotationsSlice(entities *[]map[string]interface{}) {
	for _, e := range *entities {
		castAndRemoveAnnotations(&e)
	}
}

// TODO: The default behavior of json.Unmarshal is to deserialize all json numbers as Float64.
// This can be a problem for table entities which store float and int differently
func castAndRemoveAnnotations(entity *map[string]interface{}) error {
	//value := (*entity)["value"].([]interface{})[0].(map[string]interface{})
	value := *entity
	for k, v := range value {

		iSuffix := strings.Index(k, OdataType)
		if iSuffix > 0 {
			// Get the name of the property that this odataType key describes.
			valueKey := k[0:iSuffix]
			// get the string value of the value at the valueKey
			valAsString := value[valueKey].(string)

			switch v {
			case edmBinary:
				value[valueKey] = []byte(valAsString)
			case edmDateTime:
				t, err := time.Parse(ISO8601, valAsString)
				if err != nil {
					return err
				}
				value[valueKey] = t
			case edmGuid:
				value[valueKey] = uuid.Parse(valAsString)
			case edmInt64:
				i, err := strconv.ParseInt(valAsString, 10, 64)
				if err != nil {
					return err
				}
				value[valueKey] = i
			default:
				return errors.New(fmt.Sprintf("unsupported annotation found: %s", k))
			}
			// remove the annotation key
			delete(value, k)
		}
	}
	return nil
}

func toOdataAnnotatedDictionary(entity *map[string]interface{}) error {
	entMap := *entity
	for k, v := range entMap {
		t := reflect.TypeOf(v)
	Switch:
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			if getTypeArray(v) != reflect.TypeOf(byte(0)) {
				return errors.New("arrays and slices must be of type byte")
			}
			// check if this is a uuid
			uuidVal, ok := v.(uuid.UUID)
			if ok {
				entMap[k] = uuidVal.String()
				entMap[odataType(k)] = edmGuid
			} else {
				entMap[odataType(k)] = edmBinary
				b := v.([]byte)
				entMap[k] = base64.StdEncoding.EncodeToString(b)
			}
		case reflect.Struct:
			switch tn := reflect.TypeOf(v).String(); tn {
			case "time.Time":
				entMap[odataType(k)] = edmDateTime
				time := v.(time.Time)
				entMap[k] = time.UTC().Format(ISO8601)
				continue
			default:
				return errors.New(fmt.Sprintf("Invalid struct for entity field '%s' of type '%s'", k, tn))
			}
		case reflect.Float32, reflect.Float64:
			entMap[odataType(k)] = edmDouble
		case reflect.Int64:
			entMap[odataType(k)] = edmInt64
			i64 := v.(int64)
			entMap[k] = strconv.FormatInt(i64, 10)
		case reflect.Ptr:
			if v == nil {
				// if the pointer is nil, ignore it.
				continue
			}
			// follow the pointer to the type and re-run the switch
			t = reflect.ValueOf(v).Elem().Type()
			goto Switch
		}
	}
	return nil
}

func toMap(ent interface{}) (*map[string]interface{}, error) {
	// If we were provided a map already, send that back.
	if reflect.ValueOf(ent).Kind() == reflect.Map {
		entMap := ent.(map[string]interface{})
		err := toOdataAnnotatedDictionary(&entMap)
		if err != nil {
			return nil, err
		}
		return &entMap, nil
	}
	s := reflect.ValueOf(&ent).Elem().Elem()
	typeOfT := s.Type()
	nf := s.NumField()
	entMap := make(map[string]interface{}, nf)

	for i := 0; i < nf; i++ {
		v := s.Field(i)
	Switch:
		f := typeOfT.Field(i)
		name := f.Name
		if name == etag || name == timestamp {
			// we do not need to serialize ETag or TimeStamp
			continue
		}
		// add odata annotations for the types that require it.
		switch k := v.Type().Kind(); k {
		case reflect.Array, reflect.Slice:
			if getTypeArray(v.Interface()) != reflect.TypeOf(byte(0)) {
				return nil, errors.New("arrays and slices must be of type byte")
			}
			// check if this is a uuid field as decorated by a tag
			if _, ok := f.Tag.Lookup("uuid"); ok {
				entMap[odataType(name)] = edmGuid
				u := v.Interface().([16]byte)
				var uu uuid.UUID = u
				entMap[name] = uu.String()
				continue
			} else {
				entMap[odataType(name)] = edmBinary
				b := v.Interface().([]byte)
				entMap[name] = base64.StdEncoding.EncodeToString(b)
				continue
			}
		case reflect.Struct:
			switch tn := v.Type().String(); tn {
			case "time.Time":
				entMap[odataType(name)] = edmDateTime
				time := v.Interface().(time.Time)
				entMap[name] = time.UTC().Format(ISO8601)
				continue
			default:
				return nil, errors.New(fmt.Sprintf("Invalid struct for entity field '%s' of type '%s'", typeOfT.Field(i).Name, tn))
			}
		case reflect.Float32, reflect.Float64:
			entMap[odataType(name)] = edmDouble
		case reflect.Int64:
			entMap[odataType(name)] = edmInt64
			i64 := v.Interface().(int64)
			entMap[name] = strconv.FormatInt(i64, 10)
			continue
		case reflect.Ptr:
			if v.IsNil() {
				// if the pointer is nil, ignore it.
				continue
			}
			// follow the pointer to the type and re-run the switch
			v = v.Elem()
			goto Switch

			// typeOfT.Field(i).Name, f.Type(), f.Interface())
		}
		entMap[name] = v.Interface()
	}
	return &entMap, nil
}

// fromMap converts an entity map to a strongly typed model interface
// tt is the type of the model
// fmap is the result of getTypeValueMap for the model type
// src is the source map value
// srcVal is the the Value of the source map value
func fromMap(tt reflect.Type, fmap *map[string]int, src *map[string]interface{}, srcVal reflect.Value) error {
	for k, v := range *src {
		// skip if this is an OData type descriptor
		iSuffix := strings.Index(k, OdataType)
		if iSuffix > 0 {
			continue
		}
		// fetch the Field index by property name from the field map
		fIndex := (*fmap)[k]
		// Get the Value for the Field
		val := srcVal.Field(fIndex)
	Switch:
		switch val.Kind() {
		case reflect.String:
			val.SetString(v.(string))
		case reflect.Bool:
			val.SetBool(v.(bool))
		case reflect.Float64:
			val.SetFloat(v.(float64))
		case reflect.Int:
			val.SetInt(int64(v.(float64)))
		case reflect.Int64:
			i64, err := strconv.ParseInt(v.(string), 10, 64)
			if err != nil {
				return err
			}
			val.SetInt(i64)
		case reflect.Struct:
			switch tn := val.Type().String(); tn {
			case "time.Time":
				t, err := time.Parse(ISO8601, v.(string))
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(t))
			}
		case reflect.Ptr:
			if val.IsNil() {
				// populate the nil pointer with it's element type and re-run the type evaluation
				val.Set(reflect.New(val.Type().Elem()))
				val = val.Elem()
				goto Switch
			}
		case reflect.Array, reflect.Map, reflect.Slice:
			if getTypeArray(val.Interface()) != reflect.TypeOf(byte(0)) {
				return errors.New("arrays and slices must be of type byte")
			}
			// 	// check if this is a uuid field as decorated by a tag
			if _, ok := tt.Field(fIndex).Tag.Lookup("uuid"); ok {
				u := uuid.Parse(v.(string))
				val.Set(reflect.ValueOf(u))
			} else {
				b, err := base64.StdEncoding.DecodeString(v.(string))
				if err != nil {
					return err
				}
				val.SetBytes(b)
			}
		}
	}
	return nil
}

// getTypeValueMap - builds a map of Field names to their Field index for the given interface{}
func getTypeValueMap(tt reflect.Type) *map[string]int {
	nf := tt.NumField()
	fmap := make(map[string]int)
	// build a map of field types
	for i := 0; i < nf; i++ {
		f := tt.Field(i)
		fmap[f.Name] = i
		if f.Name == etag {
			fmap[etagOdata] = i
		}
	}
	return &fmap
}

func odataType(n string) string {
	var b strings.Builder
	b.Grow(len(n) + len(OdataType))
	b.WriteString(n)
	b.WriteString(OdataType)
	return b.String()
}

func getTypeArray(arr interface{}) reflect.Type {
	return reflect.TypeOf(arr).Elem()
}
