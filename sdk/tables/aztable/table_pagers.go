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
	PageResponse() TableEntityQueryByteResponseResponse
}

type tableEntityQueryResponsePager struct {
	tableClient       *TableClient
	current           *TableEntityQueryByteResponseResponse
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

	c, err := castToByteResponse(&resp)
	if err != nil {
		p.err = nil
	}

	p.current = &c
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
func (p *tableEntityQueryResponsePager) PageResponse() TableEntityQueryByteResponseResponse {
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
				return fmt.Errorf("Invalid struct for entity field '%s' of type '%s'", k, tn)
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

func flattenEntity(entity reflect.Value, entityMap *map[string]interface{}) {
	for i := 0; i < entity.NumField(); i++ {
		if !entity.Field(i).IsZero() {
			fieldName := entity.Type().Field(i).Name
			if fieldName == "PartitionKey" {
				(*entityMap)["PartitionKey"] = entity.Field(i).String()
			} else if fieldName == "RowKey" {
				(*entityMap)["RowKey"] = entity.Field(i).String()
			}
		}
	}
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
