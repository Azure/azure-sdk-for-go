// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"
	"fmt"
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
				t, err := time.Parse(time.RFC3339Nano, valAsString)
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
