// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type AzureSasCredential struct {
	signature string
}

func NewAzureSasCredential(signature string) (*AzureSasCredential, error) {
	return &AzureSasCredential{
		signature: signature,
	}, nil
}

func (a *AzureSasCredential) Update(signature string) {
	a.signature = signature
}

// AuthenticationPolicy implements the Credential interface on SharedKeyCredential.
func (a *AzureSasCredential) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(req *azcore.Request) (*azcore.Response, error) {
		currentUrl := req.URL.String()
		query := req.URL.Query()

		signature := strings.TrimPrefix(a.signature, "?")

		if query.Encode() != "" {
			if !strings.Contains(currentUrl, signature) {
				currentUrl = currentUrl + "?" + signature
			}
		} else {
			if strings.HasSuffix(currentUrl, "?") {
				currentUrl = currentUrl + signature
			} else {
				currentUrl = currentUrl + "?" + signature
			}
		}

		newUrl, err := url.Parse(currentUrl)
		if err != nil {
			return nil, err
		}
		req.URL = newUrl

		return req.Next()
	})
}
