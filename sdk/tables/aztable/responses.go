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
	RawResponse *http.Response
	ETag        *string
}

func addEntityResponseFromGenerated(g *generated.TableInsertEntityResponse) *AddEntityResponse {
	if g == nil {
		return &AddEntityResponse{}
	}

	return &AddEntityResponse{
		RawResponse: g.RawResponse,
		ETag:        g.ETag,
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

type InsertEntityResponse struct {
	RawResponse *http.Response
	ETag        *string
}

func insertEntityFromGeneratedMerge(g *generated.TableMergeEntityResponse) *InsertEntityResponse {
	if g == nil {
		return &InsertEntityResponse{}
	}

	return &InsertEntityResponse{
		RawResponse: g.RawResponse,
		ETag:        g.ETag,
	}
}

func insertEntityFromGeneratedUpdate(g *generated.TableUpdateEntityResponse) *InsertEntityResponse {
	if g == nil {
		return &InsertEntityResponse{}
	}

	return &InsertEntityResponse{
		RawResponse: g.RawResponse,
		ETag:        g.ETag,
	}
}

type GetAccessPolicyResponse struct {
	RawResponse       *http.Response
	SignedIdentifiers []*SignedIdentifier
}

func getAccessPolicyResponseFromGenerated(g *generated.TableGetAccessPolicyResponse) *GetAccessPolicyResponse {
	if g == nil {
		return &GetAccessPolicyResponse{}
	}

	var sis []*SignedIdentifier
	for _, s := range g.SignedIdentifiers {
		sis = append(sis, fromGeneratedSignedIdentifier(s))
	}
	return &GetAccessPolicyResponse{
		RawResponse:       g.RawResponse,
		SignedIdentifiers: sis,
	}
}

type SetAccessPolicyResponse struct {
	RawResponse *http.Response
}

func setAccessPolicyResponseFromGenerated(g *generated.TableSetAccessPolicyResponse) *SetAccessPolicyResponse {
	if g == nil {
		return &SetAccessPolicyResponse{}
	}
	return &SetAccessPolicyResponse{
		RawResponse: g.RawResponse,
	}
}

type GetStatisticsResponse struct {
	RawResponse *http.Response
}

func getStatisticsResponseFromGenerated(g *generated.ServiceGetStatisticsResponse) *GetStatisticsResponse {
	return &GetStatisticsResponse{
		RawResponse: g.RawResponse,
	}
}

type GetPropertiesResponse struct {
	RawResponse *http.Response
	// The set of CORS rules.
	Cors []*CorsRule `xml:"Cors>CorsRule"`

	// A summary of request statistics grouped by API in hourly aggregates for tables.
	HourMetrics *Metrics `xml:"HourMetrics"`

	// Azure Analytics Logging settings.
	Logging *Logging `xml:"Logging"`

	// A summary of request statistics grouped by API in minute aggregates for tables.
	MinuteMetrics *Metrics `xml:"MinuteMetrics"`
}

func getPropertiesResponseFromGenerated(g *generated.ServiceGetPropertiesResponse) *GetPropertiesResponse {
	var cors []*CorsRule
	for _, c := range g.Cors {
		cors = append(cors, fromGeneratedCors(c))
	}
	return &GetPropertiesResponse{
		RawResponse:   g.RawResponse,
		Cors:          cors,
		HourMetrics:   fromGeneratedMetrics(g.HourMetrics),
		Logging:       fromGeneratedLogging(g.Logging),
		MinuteMetrics: fromGeneratedMetrics(g.MinuteMetrics),
	}
}

type SetPropertiesResponse struct {
	RawResponse *http.Response
}

func setPropertiesResponseFromGenerated(g *generated.ServiceSetPropertiesResponse) *SetPropertiesResponse {
	return &SetPropertiesResponse{
		RawResponse: g.RawResponse,
	}
}
