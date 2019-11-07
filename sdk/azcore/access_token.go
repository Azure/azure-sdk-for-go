// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"context"
	"encoding/json"
	"time"
)

// TokenCredential interface serves as an anonymous field for all other credentials to use the GetToken method
type TokenCredential interface {
	GetToken(ctx context.Context, scopes []string) (*AccessToken, error)
}

// AccessToken is used to set and maintain tokens for authentication
type AccessToken struct {
	Token     string      `json:"access_token"`
	ExpiresIn json.Number `json:"expires_in"`
	ExpiresOn time.Time
}
