//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// Code generated by Microsoft (R) AutoRest Code Generator (autorest: 3.4.3, generator: @autorest/go@4.0.0-preview.27)
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package internal

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"reflect"
)

// KeyVaultClientGetCertificateIssuersPager provides operations for iterating over paged responses.
type KeyVaultClientGetCertificateIssuersPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetCertificateIssuersResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetCertificateIssuersResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetCertificateIssuersPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetCertificateIssuersPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.CertificateIssuerListResult.NextLink == nil || len(*p.current.CertificateIssuerListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getCertificateIssuersHandleError(resp)
		return false
	}
	result, err := p.client.getCertificateIssuersHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetCertificateIssuersResponse page.
func (p *KeyVaultClientGetCertificateIssuersPager) PageResponse() KeyVaultClientGetCertificateIssuersResponse {
	return p.current
}

// KeyVaultClientGetCertificateVersionsPager provides operations for iterating over paged responses.
type KeyVaultClientGetCertificateVersionsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetCertificateVersionsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetCertificateVersionsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetCertificateVersionsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetCertificateVersionsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.CertificateListResult.NextLink == nil || len(*p.current.CertificateListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getCertificateVersionsHandleError(resp)
		return false
	}
	result, err := p.client.getCertificateVersionsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetCertificateVersionsResponse page.
func (p *KeyVaultClientGetCertificateVersionsPager) PageResponse() KeyVaultClientGetCertificateVersionsResponse {
	return p.current
}

// KeyVaultClientGetCertificatesPager provides operations for iterating over paged responses.
type KeyVaultClientGetCertificatesPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetCertificatesResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetCertificatesResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetCertificatesPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetCertificatesPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.CertificateListResult.NextLink == nil || len(*p.current.CertificateListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getCertificatesHandleError(resp)
		return false
	}
	result, err := p.client.getCertificatesHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetCertificatesResponse page.
func (p *KeyVaultClientGetCertificatesPager) PageResponse() KeyVaultClientGetCertificatesResponse {
	return p.current
}

// KeyVaultClientGetDeletedCertificatesPager provides operations for iterating over paged responses.
type KeyVaultClientGetDeletedCertificatesPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetDeletedCertificatesResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetDeletedCertificatesResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetDeletedCertificatesPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetDeletedCertificatesPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.DeletedCertificateListResult.NextLink == nil || len(*p.current.DeletedCertificateListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getDeletedCertificatesHandleError(resp)
		return false
	}
	result, err := p.client.getDeletedCertificatesHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetDeletedCertificatesResponse page.
func (p *KeyVaultClientGetDeletedCertificatesPager) PageResponse() KeyVaultClientGetDeletedCertificatesResponse {
	return p.current
}

// KeyVaultClientGetDeletedKeysPager provides operations for iterating over paged responses.
type KeyVaultClientGetDeletedKeysPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetDeletedKeysResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetDeletedKeysResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetDeletedKeysPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetDeletedKeysPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.DeletedKeyListResult.NextLink == nil || len(*p.current.DeletedKeyListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getDeletedKeysHandleError(resp)
		return false
	}
	result, err := p.client.getDeletedKeysHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetDeletedKeysResponse page.
func (p *KeyVaultClientGetDeletedKeysPager) PageResponse() KeyVaultClientGetDeletedKeysResponse {
	return p.current
}

// KeyVaultClientGetDeletedSasDefinitionsPager provides operations for iterating over paged responses.
type KeyVaultClientGetDeletedSasDefinitionsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetDeletedSasDefinitionsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetDeletedSasDefinitionsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetDeletedSasDefinitionsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetDeletedSasDefinitionsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.DeletedSasDefinitionListResult.NextLink == nil || len(*p.current.DeletedSasDefinitionListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getDeletedSasDefinitionsHandleError(resp)
		return false
	}
	result, err := p.client.getDeletedSasDefinitionsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetDeletedSasDefinitionsResponse page.
func (p *KeyVaultClientGetDeletedSasDefinitionsPager) PageResponse() KeyVaultClientGetDeletedSasDefinitionsResponse {
	return p.current
}

// KeyVaultClientGetDeletedSecretsPager provides operations for iterating over paged responses.
type KeyVaultClientGetDeletedSecretsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetDeletedSecretsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetDeletedSecretsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetDeletedSecretsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetDeletedSecretsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.DeletedSecretListResult.NextLink == nil || len(*p.current.DeletedSecretListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getDeletedSecretsHandleError(resp)
		return false
	}
	result, err := p.client.getDeletedSecretsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetDeletedSecretsResponse page.
func (p *KeyVaultClientGetDeletedSecretsPager) PageResponse() KeyVaultClientGetDeletedSecretsResponse {
	return p.current
}

// KeyVaultClientGetDeletedStorageAccountsPager provides operations for iterating over paged responses.
type KeyVaultClientGetDeletedStorageAccountsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetDeletedStorageAccountsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetDeletedStorageAccountsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetDeletedStorageAccountsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetDeletedStorageAccountsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.DeletedStorageListResult.NextLink == nil || len(*p.current.DeletedStorageListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getDeletedStorageAccountsHandleError(resp)
		return false
	}
	result, err := p.client.getDeletedStorageAccountsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetDeletedStorageAccountsResponse page.
func (p *KeyVaultClientGetDeletedStorageAccountsPager) PageResponse() KeyVaultClientGetDeletedStorageAccountsResponse {
	return p.current
}

// KeyVaultClientGetKeyVersionsPager provides operations for iterating over paged responses.
type KeyVaultClientGetKeyVersionsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetKeyVersionsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetKeyVersionsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetKeyVersionsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetKeyVersionsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.KeyListResult.NextLink == nil || len(*p.current.KeyListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getKeyVersionsHandleError(resp)
		return false
	}
	result, err := p.client.getKeyVersionsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetKeyVersionsResponse page.
func (p *KeyVaultClientGetKeyVersionsPager) PageResponse() KeyVaultClientGetKeyVersionsResponse {
	return p.current
}

// KeyVaultClientGetKeysPager provides operations for iterating over paged responses.
type KeyVaultClientGetKeysPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetKeysResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetKeysResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetKeysPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetKeysPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.KeyListResult.NextLink == nil || len(*p.current.KeyListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getKeysHandleError(resp)
		return false
	}
	result, err := p.client.getKeysHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetKeysResponse page.
func (p *KeyVaultClientGetKeysPager) PageResponse() KeyVaultClientGetKeysResponse {
	return p.current
}

// KeyVaultClientGetSasDefinitionsPager provides operations for iterating over paged responses.
type KeyVaultClientGetSasDefinitionsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetSasDefinitionsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetSasDefinitionsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetSasDefinitionsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetSasDefinitionsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.SasDefinitionListResult.NextLink == nil || len(*p.current.SasDefinitionListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getSasDefinitionsHandleError(resp)
		return false
	}
	result, err := p.client.getSasDefinitionsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetSasDefinitionsResponse page.
func (p *KeyVaultClientGetSasDefinitionsPager) PageResponse() KeyVaultClientGetSasDefinitionsResponse {
	return p.current
}

// KeyVaultClientGetSecretVersionsPager provides operations for iterating over paged responses.
type KeyVaultClientGetSecretVersionsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetSecretVersionsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetSecretVersionsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetSecretVersionsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetSecretVersionsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.SecretListResult.NextLink == nil || len(*p.current.SecretListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getSecretVersionsHandleError(resp)
		return false
	}
	result, err := p.client.getSecretVersionsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetSecretVersionsResponse page.
func (p *KeyVaultClientGetSecretVersionsPager) PageResponse() KeyVaultClientGetSecretVersionsResponse {
	return p.current
}

// KeyVaultClientGetSecretsPager provides operations for iterating over paged responses.
type KeyVaultClientGetSecretsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetSecretsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetSecretsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetSecretsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetSecretsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.SecretListResult.NextLink == nil || len(*p.current.SecretListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getSecretsHandleError(resp)
		return false
	}
	result, err := p.client.getSecretsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetSecretsResponse page.
func (p *KeyVaultClientGetSecretsPager) PageResponse() KeyVaultClientGetSecretsResponse {
	return p.current
}

// KeyVaultClientGetStorageAccountsPager provides operations for iterating over paged responses.
type KeyVaultClientGetStorageAccountsPager struct {
	client *KeyVaultClient
	current KeyVaultClientGetStorageAccountsResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, KeyVaultClientGetStorageAccountsResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *KeyVaultClientGetStorageAccountsPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *KeyVaultClientGetStorageAccountsPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.StorageListResult.NextLink == nil || len(*p.current.StorageListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.Con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.getStorageAccountsHandleError(resp)
		return false
	}
	result, err := p.client.getStorageAccountsHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current KeyVaultClientGetStorageAccountsResponse page.
func (p *KeyVaultClientGetStorageAccountsPager) PageResponse() KeyVaultClientGetStorageAccountsResponse {
	return p.current
}

// RoleAssignmentsListForScopePager provides operations for iterating over paged responses.
type RoleAssignmentsListForScopePager struct {
	client *roleAssignmentsClient
	current RoleAssignmentsListForScopeResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, RoleAssignmentsListForScopeResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *RoleAssignmentsListForScopePager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *RoleAssignmentsListForScopePager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RoleAssignmentListResult.NextLink == nil || len(*p.current.RoleAssignmentListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listForScopeHandleError(resp)
		return false
	}
	result, err := p.client.listForScopeHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current RoleAssignmentsListForScopeResponse page.
func (p *RoleAssignmentsListForScopePager) PageResponse() RoleAssignmentsListForScopeResponse {
	return p.current
}

// RoleDefinitionsListPager provides operations for iterating over paged responses.
type RoleDefinitionsListPager struct {
	client *roleDefinitionsClient
	current RoleDefinitionsListResponse
	err error
	requester func(context.Context) (*policy.Request, error)
	advancer func(context.Context, RoleDefinitionsListResponse) (*policy.Request, error)
}

// Err returns the last error encountered while paging.
func (p *RoleDefinitionsListPager) Err() error {
	return p.err
}

// NextPage returns true if the pager advanced to the next page.
// Returns false if there are no more pages or an error occurred.
func (p *RoleDefinitionsListPager) NextPage(ctx context.Context) bool {
	var req *policy.Request
	var err error
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.RoleDefinitionListResult.NextLink == nil || len(*p.current.RoleDefinitionListResult.NextLink) == 0 {
			return false
		}
		req, err = p.advancer(ctx, p.current)
	} else {
		req, err = p.requester(ctx)
	}
	if err != nil {
		p.err = err
		return false
	}
	resp, err := p.	client.con.Pipeline().Do(req)
	if err != nil {
		p.err = err
		return false
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		p.err = p.client.listHandleError(resp)
		return false
	}
	result, err := p.client.listHandleResponse(resp)
	if err != nil {
		p.err = err
		return false
	}
	p.current = result
	return true
}

// PageResponse returns the current RoleDefinitionsListResponse page.
func (p *RoleDefinitionsListPager) PageResponse() RoleDefinitionsListResponse {
	return p.current
}

