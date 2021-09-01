// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// Options for Client.Create and ServiceClient.CreateTable method
type CreateTableOptions struct {
}

func (c *CreateTableOptions) toGenerated() *generated.TableCreateOptions {
	return &generated.TableCreateOptions{}
}

// Options for Client.Delete and ServiceClient.DeleteTable methods
type DeleteTableOptions struct {
}

func (c *DeleteTableOptions) toGenerated() *generated.TableDeleteOptions {
	return &generated.TableDeleteOptions{}
}

// ListEntitiesOptions contains a group of parameters for the Table.Query method.
type ListEntitiesOptions struct {
	// OData filter expression.
	Filter *string
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
		Format: generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata.ToPtr(),
		Select: l.Select,
		Top:    l.Top,
	}
}

// ListEntitiesOptions contains a group of parameters for the ServiceClient.QueryTables method.
type ListTablesOptions struct {
	// OData filter expression.
	Filter *string
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
		Format: generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata.ToPtr(),
		Select: l.Select,
		Top:    l.Top,
	}
}

type ResponseFormat string

const (
	ResponseFormatReturnContent   ResponseFormat = "return-content"
	ResponseFormatReturnNoContent ResponseFormat = "return-no-content"
)

// PossibleResponseFormatValues returns the possible values for the ResponseFormat const type.
func PossibleResponseFormatValues() []ResponseFormat {
	return []ResponseFormat{
		ResponseFormatReturnContent,
		ResponseFormatReturnNoContent,
	}
}

// ToPtr returns a *ResponseFormat pointing to the current value.
func (c ResponseFormat) ToPtr() *ResponseFormat {
	return &c
}

// Options for Client.GetEntity method
type GetEntityOptions struct {
}

func (g *GetEntityOptions) toGenerated() (*generated.TableQueryEntityWithPartitionAndRowKeyOptions, *generated.QueryOptions) {
	return &generated.TableQueryEntityWithPartitionAndRowKeyOptions{}, &generated.QueryOptions{Format: generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata.ToPtr()}
}

// Options for the Client.AddEntity operation
type AddEntityOptions struct {
	// Specifies whether the response should include the inserted entity in the payload. Possible values are return-no-content and return-content.
	ResponsePreference *ResponseFormat
}

type DeleteEntityOptions struct {
	IfMatch *azcore.ETag
}

func (d *DeleteEntityOptions) toGenerated() *generated.TableDeleteEntityOptions {
	return &generated.TableDeleteEntityOptions{}
}

type UpdateEntityOptions struct {
	IfMatch    *azcore.ETag
	UpdateMode EntityUpdateMode
}

func (u *UpdateEntityOptions) toGeneratedMergeEntity(m map[string]interface{}) *generated.TableMergeEntityOptions {
	if u == nil {
		return &generated.TableMergeEntityOptions{}
	}
	return &generated.TableMergeEntityOptions{
		IfMatch:               to.StringPtr(string(*u.IfMatch)),
		TableEntityProperties: m,
	}
}

func (u *UpdateEntityOptions) toGeneratedUpdateEntity(m map[string]interface{}) *generated.TableUpdateEntityOptions {
	if u == nil {
		return &generated.TableUpdateEntityOptions{}
	}
	return &generated.TableUpdateEntityOptions{
		IfMatch:               to.StringPtr(string(*u.IfMatch)),
		TableEntityProperties: m,
	}
}

type InsertEntityOptions struct {
	ETag       azcore.ETag
	UpdateMode EntityUpdateMode
}

type GetAccessPolicyOptions struct {
}

func (g *GetAccessPolicyOptions) toGenerated() *generated.TableGetAccessPolicyOptions {
	return &generated.TableGetAccessPolicyOptions{}
}

type SetAccessPolicyOptions struct {
	TableACL []*SignedIdentifier
}

func (s *SetAccessPolicyOptions) toGenerated() *generated.TableSetAccessPolicyOptions {
	var sis []*generated.SignedIdentifier
	for _, t := range s.TableACL {
		sis = append(sis, toGeneratedSignedIdentifier(t))
	}
	return &generated.TableSetAccessPolicyOptions{
		TableACL: sis,
	}
}

type GetStatisticsOptions struct {
}

func (g *GetStatisticsOptions) toGenerated() *generated.ServiceGetStatisticsOptions {
	return &generated.ServiceGetStatisticsOptions{}
}

type GetPropertiesOptions struct {
}

func (g *GetPropertiesOptions) toGenerated() *generated.ServiceGetPropertiesOptions {
	return &generated.ServiceGetPropertiesOptions{}
}

type SetPropertiesOptions struct{}

func (s *SetPropertiesOptions) toGenerated() *generated.ServiceSetPropertiesOptions {
	return &generated.ServiceSetPropertiesOptions{}
}

// ListOptions contains a group of parameters for the ServiceClient.Query method.
type ListOptions struct {
	// OData filter expression.
	Filter *string
	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string
	// Maximum number of records to return.
	Top *int32
}
