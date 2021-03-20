// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
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

func castAndRemoveAnnotations(entity *map[string]interface{}) {
	/*
			foreach (var propertyName in entity.Keys)
		            {
		                var spanPropertyName = propertyName.AsSpan();
		                var iSuffix = spanPropertyName.IndexOf(spanOdataSuffix);
		                if (iSuffix > 0)
		                {
		                    // This property is an Odata annotation. Save it in the typeAnnoations dictionary.
		                    typeAnnotationsWithKeys[spanPropertyName.Slice(0, iSuffix).ToString()] = ((entity[propertyName] as string)!, propertyName);
		                }
		            }

		            // Iterate through the types that are serialized as string by default and Parse them as the correct type, as indicated by the type annotations.
		            foreach (var annotation in typeAnnotationsWithKeys.Keys)
		            {
		                entity[annotation] = typeAnnotationsWithKeys[annotation].TypeAnnotation switch
		                {
		                    TableConstants.Odata.EdmBinary => Convert.FromBase64String(entity[annotation] as string),
		                    TableConstants.Odata.EdmDateTime => DateTimeOffset.Parse(entity[annotation] as string, CultureInfo.InvariantCulture, DateTimeStyles.RoundtripKind),
		                    TableConstants.Odata.EdmGuid => Guid.Parse(entity[annotation] as string),
		                    TableConstants.Odata.EdmInt64 => long.Parse(entity[annotation] as string, CultureInfo.InvariantCulture),
		                    _ => throw new NotSupportedException("Not supported type " + typeAnnotationsWithKeys[annotation])
		                };

		                // Remove the type annotation property from the dictionary.
		                entity.Remove(typeAnnotationsWithKeys[annotation].AnnotationKey);
		            }
	*/
	// for name, val := range *entity {

	// }
}
