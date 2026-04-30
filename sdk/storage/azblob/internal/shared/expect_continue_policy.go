// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const (
	// expectContinueThreshold is the minimum Content-Length (in bytes) for a PUT
	// request to receive the "Expect: 100-continue" header. Requests with a body
	// larger than 8 MiB get the header by default.
	expectContinueThreshold int64 = 8 * 1024 * 1024

	// EnvExpectContinueDisabled is the environment variable that, when set to "1"
	// or "true", disables the Expect: 100-continue policy entirely.
	EnvExpectContinueDisabled = "AZURE_STORAGE_DISABLE_EXPECT_CONTINUE"
)

// ExpectContinuePolicy is a per-retry pipeline policy that adds the
// "Expect: 100-continue" header to PUT requests with Content-Length > 8 MiB.
//
// This allows the server to reject the request (e.g. with 429 or 403) before
// the client sends the full body, saving bandwidth under throttling or error
// conditions.
//
// The policy is disabled entirely if the environment variable
// AZURE_STORAGE_DISABLE_EXPECT_CONTINUE is set to "1" or "true".
type ExpectContinuePolicy struct {
	disabled bool
}

// NewExpectContinuePolicy creates a new ExpectContinuePolicy. It reads the
// environment variable at construction time so it does not check the env on
// every request.
func NewExpectContinuePolicy() *ExpectContinuePolicy {
	v := os.Getenv(EnvExpectContinueDisabled)
	return &ExpectContinuePolicy{
		disabled: v == "1" || v == "true",
	}
}

// Do implements the policy.Policy interface.
func (p *ExpectContinuePolicy) Do(req *policy.Request) (*http.Response, error) {
	if !p.disabled && req.Raw().Method == http.MethodPut && req.Raw().ContentLength >= expectContinueThreshold {
		req.Raw().Header.Set("Expect", "100-continue")
	}
	return req.Next()
}
