// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

const lenBearerTokenPrefix = len("Bearer ")

type cosmosBearerTokenPolicy struct {
}

func (b *cosmosBearerTokenPolicy) Do(req *policy.Request) (*http.Response, error) {
	currentAuthorization := req.Raw().Header.Get(headerAuthorization)
	if currentAuthorization == "" {
		return nil, errors.New("authorization header is missing")
	}

	token := currentAuthorization[lenBearerTokenPrefix:]
	req.Raw().Header.Set(headerAuthorization, fmt.Sprintf("type=aad&ver=1.0&sig=%v", token))
	return req.Next()
}
