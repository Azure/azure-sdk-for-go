//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
)

type beginDeleteKeyPoller struct {
	resp *http.Response
	poll func(context.Context) (*http.Response, error)
}

func (b *beginDeleteKeyPoller) Done() bool {
	return b.resp != nil && b.resp.StatusCode == http.StatusOK
}

func (b *beginDeleteKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
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

func (b *beginDeleteKeyPoller) Result(ctx context.Context, out *DeleteKeyResponse) error {
	var resp generated.KeyVaultClientGetDeletedKeyResponse
	data, err := runtime.Payload(b.resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*out = DeleteKeyResponse(getDeletedKeyResponseFromGenerated(resp))
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type beginRecoverDeletedKeyPoller struct {
	resp *http.Response
	poll func(context.Context) (*http.Response, error)
}

func (b *beginRecoverDeletedKeyPoller) Done() bool {
	return b.resp != nil && b.resp.StatusCode == http.StatusOK
}

func (b *beginRecoverDeletedKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
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

func (b *beginRecoverDeletedKeyPoller) Result(ctx context.Context, out *RecoverDeletedKeyResponse) error {
	var resp generated.KeyVaultClientRecoverDeletedKeyResponse
	data, err := runtime.Payload(b.resp)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}
	*out = RecoverDeletedKeyResponse(recoverDeletedKeyResponseFromGenerated(resp))
	return nil
}
