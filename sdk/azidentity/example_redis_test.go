//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity_test

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/redis/go-redis/v9"
)

// This example demonstrates how to use azidentity to authenticate a [go-redis] client
// connecting to Azure Cache for Redis. See the [Azure Cache for Redis documentation]
// for information on configuring a cache to use Entra ID authentication.
//
// [Azure Cache for Redis documentation]: https://learn.microsoft.com/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication
// [go-redis]: https://pkg.go.dev/github.com/redis/go-redis/v9
func Example_redis() {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		// TODO: handle error
	}
	client := redis.NewClient(&redis.Options{
		Addr:                       fmt.Sprintf("%s:6380", "TODO: redis host"),
		CredentialsProviderContext: redisCredentialProvider(credential),
		TLSConfig:                  &tls.Config{MinVersion: tls.VersionTLS12},
	})
	// TODO: use the client
	_ = client
}

// redisCredentialProvider returns a function that provides a username and password to a Redis
// client. The password is an Entra ID access token acquired from the given credential. The
// username is the object ID of the principal to whom Entra issued that token.
func redisCredentialProvider(credential azcore.TokenCredential) func(context.Context) (string, string, error) {
	return func(ctx context.Context) (string, string, error) {
		// get an access token for Azure Cache for Redis
		tk, err := credential.GetToken(ctx, policy.TokenRequestOptions{
			// Azure Cache for Redis uses the same scope in all clouds
			Scopes: []string{"https://redis.azure.com/.default"},
		})
		if err != nil {
			return "", "", err
		}
		// the token is a JWT; get the principal's object ID from its payload
		parts := strings.Split(tk.Token, ".")
		if len(parts) != 3 {
			return "", "", errors.New("token must have 3 parts")
		}
		payload, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err != nil {
			return "", "", fmt.Errorf("couldn't decode payload: %s", err)
		}
		claims := struct {
			OID string `json:"oid"`
		}{}
		err = json.Unmarshal(payload, &claims)
		if err != nil {
			return "", "", fmt.Errorf("couldn't unmarshal payload: %s", err)
		}
		if claims.OID == "" {
			return "", "", errors.New("missing object ID claim")
		}
		return claims.OID, tk.Token, nil
	}
}
