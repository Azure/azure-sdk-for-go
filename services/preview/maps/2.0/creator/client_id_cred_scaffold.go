// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package creator

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type ClientIdCredScaffold struct {
	azcore.Credential
	XMsClientId *string
}

func (policy ClientIdCredScaffold) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	base_policy := policy.Credential.AuthenticationPolicy(options)
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		if len(req.Header.Get("x-ms-client-id")) == 0 && policy.XMsClientId != nil {
			req.Header.Set("x-ms-client-id", *policy.XMsClientId)
		}
		return base_policy.Do(req)
	})
}
