// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import "context"

// AnonymousCredential is for use with HTTP(S) requests that read public resource
// or for use with Shared Access Signatures (SAS).
func AnonymousCredential() Credential {
	return credentialFunc(func(AuthenticationPolicyOptions) Policy {
		return PolicyFunc(func(ctx context.Context, req *Request) (*Response, error) {
			return req.Next(ctx)
		})
	})
}
