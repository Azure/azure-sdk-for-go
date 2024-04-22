// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// AddEntityOptions contains optional parameters for Client.AddEntity
type AddEntityOptions struct {
	// Format specifies the amount of metadata returned.
	// The default is MetadataFormatMinimal.
	Format *MetadataFormat
}

// CreateTableOptions contains optional parameters for Client.Create and ServiceClient.CreateTable
type CreateTableOptions struct {
	// placeholder for future optional parameters
}

func (c *CreateTableOptions) toGenerated() *generated.TableClientCreateOptions {
	return &generated.TableClientCreateOptions{}
}

// DeleteEntityOptions contains optional parameters for Client.DeleteEntity
type DeleteEntityOptions struct {
	IfMatch *azcore.ETag
}

func (d *DeleteEntityOptions) toGenerated() *generated.TableClientDeleteEntityOptions {
	return &generated.TableClientDeleteEntityOptions{}
}

// DeleteTableOptions contains optional parameters for Client.Delete and ServiceClient.DeleteTable
type DeleteTableOptions struct {
	// placeholder for future optional parameters
}

func (c *DeleteTableOptions) toGenerated() *generated.TableClientDeleteOptions {
	return &generated.TableClientDeleteOptions{}
}

// GetAccessPolicyOptions contains optional parameters for Client.GetAccessPolicy
type GetAccessPolicyOptions struct {
	// placeholder for future optional parameters
}

func (g *GetAccessPolicyOptions) toGenerated() *generated.TableClientGetAccessPolicyOptions {
	return &generated.TableClientGetAccessPolicyOptions{}
}

// GetEntityOptions contains optional parameters for Client.GetEntity
type GetEntityOptions struct {
	// Format specifies the amount of metadata returned.
	// The default is MetadataFormatMinimal.
	Format *MetadataFormat
}

// GetPropertiesOptions contains optional parameters for Client.GetProperties
type GetPropertiesOptions struct {
	// placeholder for future optional parameters
}

func (g *GetPropertiesOptions) toGenerated() *generated.ServiceClientGetPropertiesOptions {
	return &generated.ServiceClientGetPropertiesOptions{}
}

// GetStatisticsOptions contains optional parameters for ServiceClient.GetStatistics
type GetStatisticsOptions struct {
	// placeholder for future optional parameters
}

func (g *GetStatisticsOptions) toGenerated() *generated.ServiceClientGetStatisticsOptions {
	return &generated.ServiceClientGetStatisticsOptions{}
}

// ListEntitiesOptions contains optional parameters for Table.Query
type ListEntitiesOptions struct {
	// OData filter expression.
	Filter *string

	// Select expression using OData notation. Limits the columns on each record
	// to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string

	// Maximum number of records to return.
	Top *int32

	// The NextPartitionKey to start paging from
	NextPartitionKey *string

	// The NextRowKey to start paging from
	NextRowKey *string

	// Format specifies the amount of metadata returned.
	// The default is MetadataFormatMinimal.
	Format *MetadataFormat
}

func (l *ListEntitiesOptions) toQueryOptions() *generated.QueryOptions {
	if l == nil {
		return &generated.QueryOptions{}
	}

	return &generated.QueryOptions{
		Filter: l.Filter,
		Format: l.Format,
		Select: l.Select,
		Top:    l.Top,
	}
}

// ListTablesOptions contains optional parameters for ServiceClient.QueryTables
type ListTablesOptions struct {
	// OData filter expression.
	Filter *string

	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string

	// Maximum number of records to return.
	Top *int32

	// NextTableName is the continuation token for the next table to page from
	NextTableName *string

	// Format specifies the amount of metadata returned.
	// The default is MetadataFormatMinimal.
	Format *MetadataFormat
}

func (l *ListTablesOptions) toQueryOptions() *generated.QueryOptions {
	if l == nil {
		return &generated.QueryOptions{}
	}

	return &generated.QueryOptions{
		Filter: l.Filter,
		Format: l.Format,
		Select: l.Select,
		Top:    l.Top,
	}
}

// SetAccessPolicyOptions contains optional parameters for Client.SetAccessPolicy
type SetAccessPolicyOptions struct {
	TableACL []*SignedIdentifier
}

func (s *SetAccessPolicyOptions) toGenerated() *generated.TableClientSetAccessPolicyOptions {
	if len(s.TableACL) == 0 {
		return &generated.TableClientSetAccessPolicyOptions{}
	}
	sis := make([]*generated.SignedIdentifier, len(s.TableACL))
	for i := range s.TableACL {
		sis[i] = toGeneratedSignedIdentifier(s.TableACL[i])
	}
	return &generated.TableClientSetAccessPolicyOptions{
		TableACL: sis,
	}
}

// SetPropertiesOptions contains optional parameters for Client.SetProperties
type SetPropertiesOptions struct {
	// placeholder for future optional parameters
}

func (s *SetPropertiesOptions) toGenerated() *generated.ServiceClientSetPropertiesOptions {
	return &generated.ServiceClientSetPropertiesOptions{}
}

// SubmitTransactionOptions contains optional parameters for Client.SubmitTransaction
type SubmitTransactionOptions struct {
	// placeholder for future optional parameters
}

// UpdateEntityOptions contains optional parameters for Client.UpdateEntity
type UpdateEntityOptions struct {
	IfMatch    *azcore.ETag
	UpdateMode UpdateMode
}

func (u *UpdateEntityOptions) toGeneratedMergeEntity(m map[string]any) *generated.TableClientMergeEntityOptions {
	if u == nil {
		return &generated.TableClientMergeEntityOptions{}
	}
	return &generated.TableClientMergeEntityOptions{
		IfMatch:               (*string)(u.IfMatch),
		TableEntityProperties: m,
	}
}

func (u *UpdateEntityOptions) toGeneratedUpdateEntity(m map[string]any) *generated.TableClientUpdateEntityOptions {
	if u == nil {
		return &generated.TableClientUpdateEntityOptions{}
	}
	return &generated.TableClientUpdateEntityOptions{
		IfMatch:               (*string)(u.IfMatch),
		TableEntityProperties: m,
	}
}

// UpsertEntityOptions contains optional parameters for Client.InsertEntity
type UpsertEntityOptions struct {
	// ETag is the optional etag for the Table
	ETag azcore.ETag

	// UpdateMode is the desired mode for the Update. Use UpdateModeReplace to replace fields on
	// the entity, use UpdateModeMerge to merge fields of the entity.
	UpdateMode UpdateMode
}
