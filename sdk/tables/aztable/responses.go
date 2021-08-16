// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"net/http"

	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

type CreateTableResponse struct {
	RawResponse *http.Response
}

func createTableResponseFromGen(g *generated.TableCreateResponse) *CreateTableResponse {
	if g == nil {
		return &CreateTableResponse{}
	}
	return &CreateTableResponse{
		RawResponse: g.RawResponse,
	}
}

type DeleteTableResponse struct {
	RawResponse *http.Response
}

func deleteTableResponseFromGen(g *generated.TableDeleteResponse) *DeleteTableResponse {
	if g == nil {
		return &DeleteTableResponse{}
	}
	return &DeleteTableResponse{
		RawResponse: g.RawResponse,
	}
}

type AddEntityResponse struct {
	RawResponse       *http.Response
	ETag              *string
	PreferenceApplied *string
	Value             map[string]interface{}
}

func addEntityResponseFromGenerated(g *generated.TableInsertEntityResponse) *AddEntityResponse {
	if g == nil {
		return &AddEntityResponse{}
	}

	return &AddEntityResponse{
		RawResponse:       g.RawResponse,
		ETag:              g.ETag,
		PreferenceApplied: g.PreferenceApplied,
		Value:             g.Value,
	}
}

type DeleteEntityResponse struct {
	RawResponse *http.Response
}

func deleteEntityResponseFromGenerated(g *generated.TableDeleteEntityResponse) *DeleteEntityResponse {
	if g == nil {
		return &DeleteEntityResponse{}
	}
	return &DeleteEntityResponse{
		RawResponse: g.RawResponse,
	}
}

type UpdateEntityResponse struct {
	RawResponse *http.Response
	ETag        *string
}

func updateEntityResponseFromMergeGenerated(g *generated.TableMergeEntityResponse) *UpdateEntityResponse {
	if g == nil {
		return &UpdateEntityResponse{}
	}

	return &UpdateEntityResponse{
		RawResponse: g.RawResponse,
		ETag:        g.ETag,
	}
}

func updateEntityResponseFromUpdateGenerated(g *generated.TableUpdateEntityResponse) *UpdateEntityResponse {
	if g == nil {
		return &UpdateEntityResponse{}
	}

	return &UpdateEntityResponse{
		RawResponse: g.RawResponse,
		ETag:        g.ETag,
	}
}
