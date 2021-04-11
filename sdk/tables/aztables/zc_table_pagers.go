// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

// Pager for Table entity queries
type TableEntityQueryResponsePager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// Page returns the current TableQueryResponseResponse.
	PageResponse() TableEntityQueryResponseResponse

	// Err returns the last error encountered while paging.
	Err() error
}

type tableEntityQueryResponsePager struct {
	tableClient       *TableClient
	current           *TableEntityQueryResponseResponse
	tableQueryOptions *TableQueryEntitiesOptions
	queryOptions      *QueryOptions
	err               error
}

func (p *tableEntityQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextPartitionKey == nil && p.current.XMSContinuationNextRowKey == nil) {
		return false
	}
	var resp TableEntityQueryResponseResponse
	resp, p.err = p.tableClient.client.QueryEntities(ctx, p.tableClient.name, p.tableQueryOptions, p.queryOptions)
	p.current = &resp
	p.tableQueryOptions.NextPartitionKey = resp.XMSContinuationNextPartitionKey
	p.tableQueryOptions.NextRowKey = resp.XMSContinuationNextRowKey
	return p.err == nil && resp.TableEntityQueryResponse.Value != nil && len(*resp.TableEntityQueryResponse.Value) > 0
}

func (p *tableEntityQueryResponsePager) PageResponse() TableEntityQueryResponseResponse {
	return *p.current
}

func (p *tableEntityQueryResponsePager) Err() error {
	return p.err
}

// Pager for Table Queries
type TableQueryResponsePager interface {
	// NextPage returns true if the pager advanced to the next page.
	// Returns false if there are no more pages or an error occurred.
	NextPage(context.Context) bool

	// Page returns the current TableQueryResponseResponse.
	PageResponse() TableQueryResponseResponse

	// Err returns the last error encountered while paging.
	Err() error
}

type tableQueryResponsePager struct {
	client            *tableClient
	current           *TableQueryResponseResponse
	tableQueryOptions *TableQueryOptions
	queryOptions      *QueryOptions
	err               error
}

func (p *tableQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextTableName == nil) {
		return false
	}
	var resp TableQueryResponseResponse
	resp, p.err = p.client.Query(ctx, p.tableQueryOptions, p.queryOptions)
	p.current = &resp
	p.tableQueryOptions.NextTableName = resp.XMSContinuationNextTableName
	return p.err == nil && resp.TableQueryResponse.Value != nil && len(*resp.TableQueryResponse.Value) > 0
}

func (p *tableQueryResponsePager) PageResponse() TableQueryResponseResponse {
	return *p.current
}

func (p *tableQueryResponsePager) Err() error {
	return p.err
}

func castAndRemoveAnnotationsSlice(entities *[]map[string]interface{}) {

}

// TODO: The default behavior of json.Unmarshal is to deserialize all json numbers as Float64.
// This can be a problem for table entities which store float and int differently
func castAndRemoveAnnotations(entity *map[string]interface{}) (*map[string]interface{}, error) {
	value := (*entity)["value"].([]interface{})[0].(map[string]interface{})
	for k, v := range value {

		iSuffix := strings.Index(k, OdataType)
		if iSuffix > 0 {
			// Get the name of the property that this odataType key describes.
			valueKey := k[0:iSuffix]
			// get the string value of the value at the valueKey
			valAsString := value[valueKey].(string)

			switch v {
			case EdmBinary:
				value[valueKey] = []byte(valAsString)
			case EdmDateTime:
				t, err := time.Parse(ISO8601, valAsString)
				if err != nil {
					return nil, err
				}
				value[valueKey] = t
			case EdmGuid:
				value[valueKey] = uuid.Parse(valAsString)
			case EdmInt64:
				i, err := strconv.ParseInt(valAsString, 10, 64)
				if err != nil {
					return nil, err
				}
				value[valueKey] = i
			default:
				return nil, errors.New(fmt.Sprintf("unsupported annotation found: %s", k))
			}
			// remove the annotation key
			delete(value, k)
		}
	}
	return &value, nil
}

func toMap(ent interface{}) (*map[string]interface{}, error) {
	entMap := make(map[string]interface{})
	var s reflect.Value
	if reflect.ValueOf(ent).Kind() == reflect.Ptr {
		s = reflect.ValueOf(ent).Elem()
	} else {
		s = reflect.ValueOf(&ent).Elem().Elem()
	}
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		v := s.Field(i)
	Switch:
		f := typeOfT.Field(i)
		name := f.Name
		if name == ETag || name == Timestamp {
			// we do not need to serialize ETag or TimeStamp
			continue
		}
		// add odata annotations for the types that require it.
		switch k := v.Type().Kind(); k {
		case reflect.Array, reflect.Map, reflect.Slice:
			if GetTypeArray(v.Interface()) != reflect.TypeOf(byte(0)) {
				return nil, errors.New("arrays and slices must be of type byte")
			}
			// check if this is a uuid field as decorated by a tag
			if _, ok := f.Tag.Lookup("uuid"); ok {
				entMap[odataType(name)] = EdmGuid
				u := v.Interface().([16]byte)
				var uu uuid.UUID = u
				entMap[name] = uu.String()
				continue
			} else {
				entMap[odataType(name)] = EdmBinary
				b := v.Interface().([]byte)
				entMap[name] = base64.StdEncoding.EncodeToString(b)
				continue
			}
		case reflect.Struct:
			switch tn := v.Type().String(); tn {
			case "time.Time":
				entMap[odataType(name)] = EdmDateTime
				time := v.Interface().(time.Time)
				entMap[name] = time.UTC().Format(ISO8601)
				continue
			default:
				return nil, errors.New(fmt.Sprintf("Invalid struct for entity field '%s' of type '%s'", typeOfT.Field(i).Name, tn))
			}
		case reflect.Float32, reflect.Float64:
			entMap[odataType(name)] = EdmDouble
		case reflect.Int64:
			entMap[odataType(name)] = EdmInt64
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

func odataType(n string) string {
	return fmt.Sprintf("%s%s", n, OdataType)
}

func GetTypeArray(arr interface{}) reflect.Type {
	return reflect.TypeOf(arr).Elem()
}
