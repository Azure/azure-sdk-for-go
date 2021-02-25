// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func NewAnonymousCredential() AnonymousCredential {
	return AnonymousCredential{}
}

type AnonymousCredential struct{}

func (AnonymousCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return anonymousCredentialPolicy{}
}

type anonymousCredentialPolicy struct {
}

func (p anonymousCredentialPolicy) Do(request *azcore.Request) (*azcore.Response, error) {
	return request.Next()
}
