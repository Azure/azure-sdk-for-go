package azcore

import (
	"context"
	"errors"
	"time"
)

type bearerTokenAuthPolicy struct {
	credential TokenCredential
	refreshOn  time.Time
	scopes     []string
}

// NewBearerTokenAuthPolicy creates a BearerTokenAuthPolicy object configured using the specified scopes and Token Credential.
func NewBearerTokenAuthPolicy(scopes []string, cred TokenCredential) (Policy, error) {
	if cred == nil {
		return nil, errors.New("cred cannot be nil")
	}

	if len(scopes) == 0 {
		return nil, errors.New("scopes cannot be empty")
	}

	return &bearerTokenAuthPolicy{scopes: scopes, credential: cred}, nil
}

func (c *bearerTokenAuthPolicy) Do(ctx context.Context, req *Request) (*Response, error) {
	if req.Request.URL.Scheme != "https" {
		// HTTPS must be used, otherwise the tokens are at the risk of being exposed
		return nil, errors.New("token credentials require a URL using the https protocol scheme")
	}

	token, err := c.credential.GetToken(ctx, c.scopes)
	if err != nil {
		return nil, err
	}

	req.Request.Header.Set(headerAuthorization, "Bearer "+token.Token)
	return req.Do(ctx)
}
