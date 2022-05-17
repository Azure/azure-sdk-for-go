//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcertificates

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azcertificates/internal/generated"
)

type beginCreateCertificateOperation struct {
	PollURL string
	Status  string
	poll    func(context.Context, string) (*http.Response, error)
	result  func(context.Context) (CreateCertificateResponse, error)
}

func (b *beginCreateCertificateOperation) Done() bool {
	return b.Status == "completed" || b.Status == "cancelled"
}

func (b *beginCreateCertificateOperation) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.poll(ctx, b.PollURL)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	payload, err := runtime.Payload(resp)
	if err != nil {
		return nil, err
	}

	var op Operation
	if err := json.Unmarshal(payload, &op); err != nil {
		return nil, err
	}

	if op.Status == nil {
		return nil, errors.New("missing status")
	}
	b.Status = *op.Status
	return resp, nil
}

func (b *beginCreateCertificateOperation) Result(ctx context.Context, out *CreateCertificateResponse) error {
	result, err := b.result(ctx)
	if err != nil {
		return err
	}
	*out = result
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type beginDeleteCertificateOperation struct {
	resp *http.Response
	poll func(context.Context) (*http.Response, error)
}

func (b *beginDeleteCertificateOperation) Done() bool {
	return b.resp != nil && b.resp.StatusCode == http.StatusOK
}

func (b *beginDeleteCertificateOperation) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.poll(ctx)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNotFound) {
		return nil, runtime.NewResponseError(resp)
	}
	b.resp = resp
	return b.resp, nil
}

func (b *beginDeleteCertificateOperation) Result(ctx context.Context, out *DeleteCertificateResponse) error {
	var resp generated.KeyVaultClientDeleteCertificateResponse
	data, err := runtime.Payload(b.resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*out = deleteCertificateResponseFromGenerated(resp)
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type beginRecoverDeletedCertificate struct {
	resp *http.Response
	poll func(context.Context) (*http.Response, error)
}

func (b *beginRecoverDeletedCertificate) Done() bool {
	return b.resp != nil && b.resp.StatusCode == http.StatusOK
}

func (b *beginRecoverDeletedCertificate) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.poll(ctx)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNotFound) {
		return nil, runtime.NewResponseError(resp)
	}
	b.resp = resp
	return b.resp, nil
}

func (b *beginRecoverDeletedCertificate) Result(ctx context.Context, out *RecoverDeletedCertificateResponse) error {
	var resp generated.KeyVaultClientRecoverDeletedCertificateResponse
	data, err := runtime.Payload(b.resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*out = recoverDeletedCertificateResponseFromGenerated(resp)
	return nil
}
