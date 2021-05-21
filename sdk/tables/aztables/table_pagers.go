// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
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

type StructEntityQueryResponsePager interface {
	NextPage(context.Context) bool
	PageResponse() StructQueryResponseResponse
	Err() error
}

type StructQueryResponseResponse struct {
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string

	// Date contains the information returned from the Date header response.
	Date *time.Time

	// ETag contains the information returned from the ETag header response.
	ETag *string

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string

	// The properties for the table entity query response.
	StructQueryResponse *StructQueryResponse

	// Version contains the information returned from the x-ms-version header response.
	Version *string

	// XMSContinuationNextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	XMSContinuationNextPartitionKey *string

	// XMSContinuationNextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	XMSContinuationNextRowKey *string
}

type StructQueryResponse struct {
	// The metadata response of the table.
	OdataMetadata *string `json:"odata.metadata,omitempty"`

	// List of table entities.
	Value *[]interface{} `json:"value,omitempty"`
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
	castAndRemoveAnnotationsSlice(&resp.TableEntityQueryResponse.Value)
	p.current = &resp
	p.tableQueryOptions.NextPartitionKey = resp.XMSContinuationNextPartitionKey
	p.tableQueryOptions.NextRowKey = resp.XMSContinuationNextRowKey
	return p.err == nil && resp.TableEntityQueryResponse.Value != nil && len(resp.TableEntityQueryResponse.Value) > 0
}

func (p *tableEntityQueryResponsePager) PageResponse() TableEntityQueryResponseResponse {
	return *p.current
}

func (p *tableEntityQueryResponsePager) Err() error {
	return p.err
}

type structQueryResponsePager struct {
	mapper            FromMapper
	tableClient       *TableClient
	current           *StructQueryResponseResponse
	tableQueryOptions *TableQueryEntitiesOptions
	queryOptions      *QueryOptions
	err               error
}

func (p *structQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextPartitionKey == nil && p.current.XMSContinuationNextRowKey == nil) {
		return false
	}
	var resp TableEntityQueryResponseResponse
	resp, p.err = p.tableClient.client.QueryEntities(ctx, p.tableClient.name, p.tableQueryOptions, p.queryOptions)
	castAndRemoveAnnotationsSlice(&resp.TableEntityQueryResponse.Value)
	//p.current = &resp
	r := make([]interface{}, 0, len(resp.TableEntityQueryResponse.Value))
	for _, e := range resp.TableEntityQueryResponse.Value {
		r = append(r, p.mapper.FromMap(&e))
	}
	p.current = &StructQueryResponseResponse{StructQueryResponse: &StructQueryResponse{Value: &r}}
	p.tableQueryOptions.NextPartitionKey = resp.XMSContinuationNextPartitionKey
	p.tableQueryOptions.NextRowKey = resp.XMSContinuationNextRowKey
	return p.err == nil && resp.TableEntityQueryResponse.Value != nil && len(resp.TableEntityQueryResponse.Value) > 0
}

func (p *structQueryResponsePager) PageResponse() StructQueryResponseResponse {
	return *p.current
}

func (p *structQueryResponsePager) Err() error {
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

type FromMapper interface {
	FromMap(e *map[string]interface{}) interface{}
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
	return p.err == nil && resp.TableQueryResponse.Value != nil && len(resp.TableQueryResponse.Value) > 0
}

func (p *tableQueryResponsePager) PageResponse() TableQueryResponseResponse {
	return *p.current
}

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
			case EdmBinary:
				value[valueKey] = []byte(valAsString)
			case EdmDateTime:
				t, err := time.Parse(ISO8601, valAsString)
				if err != nil {
					return err
				}
				value[valueKey] = t
			case EdmGuid:
				value[valueKey] = uuid.Parse(valAsString)
			case EdmInt64:
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
			if GetTypeArray(v) != reflect.TypeOf(byte(0)) {
				return errors.New("arrays and slices must be of type byte")
			}
			// check if this is a uuid
			uuidVal, ok := v.(uuid.UUID)
			if ok {
				entMap[k] = uuidVal.String()
				entMap[odataType(k)] = EdmGuid
			} else {
				entMap[odataType(k)] = EdmBinary
				b := v.([]byte)
				entMap[k] = base64.StdEncoding.EncodeToString(b)
			}
		case reflect.Struct:
			switch tn := reflect.TypeOf(v).String(); tn {
			case "time.Time":
				entMap[odataType(k)] = EdmDateTime
				time := v.(time.Time)
				entMap[k] = time.UTC().Format(ISO8601)
				continue
			default:
				return errors.New(fmt.Sprintf("Invalid struct for entity field '%s' of type '%s'", k, tn))
			}
		case reflect.Float32, reflect.Float64:
			entMap[odataType(k)] = EdmDouble
		case reflect.Int64:
			entMap[odataType(k)] = EdmInt64
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
	var s reflect.Value
	if reflect.ValueOf(ent).Kind() == reflect.Ptr {
		s = reflect.ValueOf(ent).Elem()
	} else {
		s = reflect.ValueOf(&ent).Elem().Elem()
	}
	typeOfT := s.Type()
	nf := s.NumField()
	entMap := make(map[string]interface{}, nf)

	for i := 0; i < nf; i++ {
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
		case reflect.Array, reflect.Slice:
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

func fromMap(src interface{}, fmap *map[string]int, m *map[string]interface{}) (interface{}, error) {
	tt := reflect.TypeOf(src)
	srcVal := reflect.New(tt).Elem()

	for k, v := range *m {
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
		case reflect.Float64:
			val.SetFloat(v.(float64))
		case reflect.Int:
			val.SetInt(int64(v.(float64)))
		case reflect.Int64:
			i64, err := strconv.ParseInt(v.(string), 10, 64)
			if err != nil {
				return nil, err
			}
			val.SetInt(i64)
		case reflect.Struct:
			switch tn := val.Type().String(); tn {
			case "time.Time":
				t, err := time.Parse(ISO8601, v.(string))
				if err != nil {
					return nil, err
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
			if GetTypeArray(val.Interface()) != reflect.TypeOf(byte(0)) {
				return nil, errors.New("arrays and slices must be of type byte")
			}
			// 	// check if this is a uuid field as decorated by a tag
			if _, ok := tt.Field(fIndex).Tag.Lookup("uuid"); ok {
				u := uuid.Parse(v.(string))
				val.Set(reflect.ValueOf(u))
			} else {
				b, err := base64.StdEncoding.DecodeString(v.(string))
				if err != nil {
					return nil, err
				}
				val.SetBytes(b)
			}
		}
	}
	return srcVal.Interface(), nil
}

// getTypeValueMap - builds a map of Field names to their Field index for the given interface{}
func getTypeValueMap(i interface{}) *map[string]int {
	tt := reflect.TypeOf(complexEntity{})
	nf := tt.NumField()
	fmap := make(map[string]int)
	// build a map of field types
	for i := 0; i < nf; i++ {
		f := tt.Field(i)
		fmap[f.Name] = i
		if f.Name == ETag {
			fmap[EtagOdata] = i
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

func GetTypeArray(arr interface{}) reflect.Type {
	return reflect.TypeOf(arr).Elem()
}
