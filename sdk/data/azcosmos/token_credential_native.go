// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

/*
#include <stdlib.h>
#include "azurecosmosdriver.h"

cosmos_token_provider_t cosmos_go_token_provider(void);
*/
import "C"

import (
	"context"
	"runtime/cgo"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// tokenProviderState is owned by Rust after cosmos_account_ref_with_credential succeeds. The
// user_data_free callback releases its cgo.Handle after the last credential reference is gone.
type tokenProviderState struct {
	credential azcore.TokenCredential
	ctx        context.Context
	cancel     context.CancelFunc
}

// buildTokenAccount adapts an azcore.TokenCredential to the driver's asynchronous callback ABI.
func (d *nativeDriver) buildTokenAccount(cfg driverConfig) error {
	endpoint := C.CString(cfg.endpoint)
	defer C.free(unsafe.Pointer(endpoint))

	ctx, cancel := context.WithCancel(context.Background())
	state := &tokenProviderState{
		credential: cfg.tokenCredential,
		ctx:        ctx,
		cancel:     cancel,
	}
	handle := cgo.NewHandle(state)
	var richErr *C.cosmos_error_t
	status := C.cosmos_account_ref_with_credential(endpoint, C.cosmos_go_token_provider(), C.intptr_t(handle), &d.account, &richErr) //nolint:gocritic // dupSubExpr is reported against cgo-generated code, not this call.
	if err := statusError(status, richErr, "building the account reference"); err != nil {
		// Rust takes ownership only on success.
		cancel()
		handle.Delete()
		return err
	}
	d.tokenProvider = state
	return nil
}

//export goCosmosTokenProviderGet
func goCosmosTokenProviderGet(
	userData C.intptr_t,
	request *C.cosmos_token_request_t,
) C.int32_t {
	if request == nil || (request.scope_len > 0 && request.scope == nil) {
		return 1
	}

	var scope string
	if request.scope_len > 0 {
		bytes := unsafe.Slice((*byte)(unsafe.Pointer(request.scope)), int(request.scope_len))
		scope = string(bytes) // The ABI borrows the scope only until this callback returns.
	}

	value, valid := cgoHandleValue(cgo.Handle(uintptr(userData)))
	if !valid {
		return 1
	}
	state, ok := value.(*tokenProviderState)
	if !ok {
		return 1
	}
	requestID := C.uint64_t(request.request_id)
	go completeTokenRequest(requestID, state, scope)
	return 0
}

//export goCosmosTokenProviderFree
func goCosmosTokenProviderFree(userData C.intptr_t) {
	handle := cgo.Handle(uintptr(userData))
	value, valid := cgoHandleValue(handle)
	if !valid {
		return
	}
	defer cgoHandleDelete(handle)
	state, ok := value.(*tokenProviderState)
	if !ok {
		return
	}
	state.cancel()
}

func completeTokenRequest(
	requestID C.uint64_t,
	state *tokenProviderState,
	scope string,
) {
	token, err := state.credential.GetToken(state.ctx, policy.TokenRequestOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		message := []byte(err.Error())
		C.cosmos_token_request_complete(
			requestID,
			1,
			nil,
			0,
			0,
			bytesPointer(message),
			C.uintptr_t(len(message)),
		)
		return
	}

	value := []byte(token.Token)
	C.cosmos_token_request_complete(
		requestID,
		0,
		bytesPointer(value),
		C.uintptr_t(len(value)),
		C.int64_t(token.ExpiresOn.Unix()),
		nil,
		0,
	)
}

// bytesPointer returns nil for an empty slice and otherwise a pointer valid for the duration of one
// C call. The completion ABI copies the bytes before returning.
func bytesPointer(value []byte) *C.uint8_t {
	if len(value) == 0 {
		return nil
	}
	return (*C.uint8_t)(unsafe.Pointer(&value[0]))
}
