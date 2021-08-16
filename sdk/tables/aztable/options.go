// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// Options for TableClient.Create and TableServiceClient.CreateTable method
type CreateTableOptions struct {
}

func (c *CreateTableOptions) toGenerated() *generated.TableCreateOptions {
	return &generated.TableCreateOptions{}
}

// Options for TableClient.Delete and TableServiceClient.DeleteTable methods
type DeleteTableOptions struct {
}

func (c *DeleteTableOptions) toGenerated() *generated.TableDeleteOptions {
	return &generated.TableDeleteOptions{}
}

// ListEntitiesOptions contains a group of parameters for the Table.Query method.
type ListEntitiesOptions struct {
	// OData filter expression.
	Filter *string
	// Specifies the media type for the response.
	Format *generated.ODataMetadataFormat
	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string
	// Maximum number of records to return.
	Top *int32
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

// ListEntitiesOptions contains a group of parameters for the TableServiceClient.QueryTables method.
type ListTablesOptions struct {
	// OData filter expression.
	Filter *string
	// Specifies the media type for the response.
	Format *generated.ODataMetadataFormat
	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string
	// Maximum number of records to return.
	Top *int32
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

type ODataMetadataFormat string

const (
	FullODataMetadata    ODataMetadataFormat = "application/json;odata=fullmetadata"
	MinimalODataMetadata ODataMetadataFormat = "application/json;odata=minimalmetadata"
	NoOdataMetadata      ODataMetadataFormat = "application/json;odata=nometadata"
)

// PossibleODataMetadataFormatValues returns the possible values for the ODataMetadataFormat const type.
func PossibleODataMetadataFormatValues() []ODataMetadataFormat {
	return []ODataMetadataFormat{
		FullODataMetadata,
		MinimalODataMetadata,
		NoOdataMetadata,
	}
}

// ToPtr returns a *ODataMetadataFormat pointing to the current value.
func (c ODataMetadataFormat) ToPtr() *ODataMetadataFormat {
	return &c
}

func (c ODataMetadataFormat) toInternal() generated.ODataMetadataFormat {
	switch c {
	case FullODataMetadata:
		return generated.ODataMetadataFormatApplicationJSONODataFullmetadata
	case MinimalODataMetadata:
		return generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata
	case NoOdataMetadata:
		return generated.ODataMetadataFormatApplicationJSONODataNometadata
	}
	return generated.ODataMetadataFormatApplicationJSONODataFullmetadata
}

// Options for TableClient.GetEntity method
type GetEntityOptions struct {
	Format ODataMetadataFormat
}

func (g *GetEntityOptions) toGenerated() (*generated.TableQueryEntityWithPartitionAndRowKeyOptions, *generated.QueryOptions) {
	if g.Format == FullODataMetadata {
		return &generated.TableQueryEntityWithPartitionAndRowKeyOptions{}, &generated.QueryOptions{
			Format: generated.ODataMetadataFormatApplicationJSONODataFullmetadata.ToPtr(),
		}
	}

	if g.Format == MinimalODataMetadata {
		return &generated.TableQueryEntityWithPartitionAndRowKeyOptions{}, &generated.QueryOptions{
			Format: generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata.ToPtr(),
		}
	}

	if g.Format == MinimalODataMetadata {
		return &generated.TableQueryEntityWithPartitionAndRowKeyOptions{}, &generated.QueryOptions{
			Format: generated.ODataMetadataFormatApplicationJSONODataNometadata.ToPtr(),
		}
	}
	return &generated.TableQueryEntityWithPartitionAndRowKeyOptions{}, &generated.QueryOptions{}
}

// Options for the TableClient.AddEntity operation
type AddEntityOptions struct {
	// Specifies whether the response should include the inserted entity in the payload. Possible values are return-no-content and return-content.
	ResponsePreference *generated.ResponseFormat
	// The properties for the table entity.
	TableEntityProperties map[string]interface{}
}

func (a *AddEntityOptions) toGenerated() *generated.TableInsertEntityOptions {
	return &generated.TableInsertEntityOptions{
		ResponsePreference:    a.ResponsePreference,
		TableEntityProperties: a.TableEntityProperties,
	}
}

type DeleteEntityOptions struct{}

func (d *DeleteEntityOptions) toGenerated() *generated.TableDeleteEntityOptions {
	return &generated.TableDeleteEntityOptions{}
}

type UpdateEntityOptions struct {
	IfMatch *string
}

func (u *UpdateEntityOptions) toGeneratedMergeEntity(m map[string]interface{}) *generated.TableMergeEntityOptions {
	if u == nil {
		return &generated.TableMergeEntityOptions{}
	}
	return &generated.TableMergeEntityOptions{
		IfMatch:               u.IfMatch,
		TableEntityProperties: m,
	}
}

func (u *UpdateEntityOptions) toGeneratedUpdateEntity(m map[string]interface{}) *generated.TableUpdateEntityOptions {
	if u == nil {
		return &generated.TableUpdateEntityOptions{}
	}
	return &generated.TableUpdateEntityOptions{
		IfMatch:               u.IfMatch,
		TableEntityProperties: m,
	}
}

type InsertEntityOptions struct {
	IfMatch *string
}

func (i *InsertEntityOptions) toGeneratedMergeEntity(m map[string]interface{}) *generated.TableMergeEntityOptions {
	if i == nil {
		return &generated.TableMergeEntityOptions{}
	}
	return &generated.TableMergeEntityOptions{
		IfMatch:               i.IfMatch,
		TableEntityProperties: m,
	}
}

func (i *InsertEntityOptions) toGeneratedUpdateEntity(m map[string]interface{}) *generated.TableUpdateEntityOptions {
	if i == nil {
		return &generated.TableUpdateEntityOptions{}
	}
	return &generated.TableUpdateEntityOptions{
		TableEntityProperties: m,
	}
}

type GetAccessPolicyOptions struct {
}

func (g *GetAccessPolicyOptions) toGenerated() *generated.TableGetAccessPolicyOptions {
	return &generated.TableGetAccessPolicyOptions{}
}

type SetAccessPolicyOptions struct {
	TableACL []*generated.SignedIdentifier
}

func (s *SetAccessPolicyOptions) toGenerated() *generated.TableSetAccessPolicyOptions {
	return &generated.TableSetAccessPolicyOptions{
		TableACL: s.TableACL,
	}
}
