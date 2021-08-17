// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// cosmosPatchTransformPolicy transforms PATCH requests into POST requests with the "X-HTTP-Method":"MERGE" header set.
type cosmosPatchTransformPolicy struct{}

func (p cosmosPatchTransformPolicy) Do(req *azcore.Request) (*azcore.Response, error) {
	transformPatchToCosmosPost(req)
	return req.Next()
}

func transformPatchToCosmosPost(req *azcore.Request) {
	if req.Method == http.MethodPatch {
		req.Method = http.MethodPost
		req.Header.Set("X-HTTP-Method", "MERGE")
	}
}
