// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// KeyCredentialPolicy authorizes requests with a [azcore.KeyCredential].
type KeyCredentialPolicy struct {
	cred   *exported.KeyCredential
	header string
	format func(string) string
}

// KeyCredentialPolicyOptions contains the optional values configuring [KeyCredentialPolicy].
type KeyCredentialPolicyOptions struct {
	// Format is used if the key needs special formatting (e.g. a prefix) before it's inserted into the HTTP request.
	// The value passed to the callback is the raw key and the return value is the augmented key.
	Format func(string) string
}

// NewKeyCredentialPolicy creates a new instance of [KeyCredentialPolicy].
//   - cred is the [azcore.KeyCredential] used to authenticate with the service
//   - header is the name of the HTTP request header in which the key is placed
//   - options contains optional configuration, pass nil to accept the default values
func NewKeyCredentialPolicy(cred *exported.KeyCredential, header string, options *KeyCredentialPolicyOptions) *KeyCredentialPolicy {
	if options == nil {
		options = &KeyCredentialPolicyOptions{}
	}
	return &KeyCredentialPolicy{
		cred:   cred,
		header: header,
		format: options.Format,
	}
}

// Do implementes the Do method on the [policy.Polilcy] interface.
func (k *KeyCredentialPolicy) Do(req *policy.Request) (*http.Response, error) {
	val := exported.KeyCredentialGet(k.cred)
	if k.format != nil {
		val = k.format(val)
	}
	req.Raw().Header.Add(k.header, val)
	return req.Next()
}
