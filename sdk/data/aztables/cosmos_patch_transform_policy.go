// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// cosmosPatchTransformPolicy transforms PATCH requests into POST requests with the "X-HTTP-Method":"MERGE" header set.
type cosmosPatchTransformPolicy struct{}

func (p cosmosPatchTransformPolicy) Do(req *policy.Request) (*http.Response, error) {
	transformPatchToCosmosPost(req)
	return req.Next()
}

func transformPatchToCosmosPost(req *policy.Request) {
	if req.Raw().Method == http.MethodPatch {
		req.Raw().Method = http.MethodPost
		req.Raw().Header.Set("X-HTTP-Method", "MERGE")
	}
}
