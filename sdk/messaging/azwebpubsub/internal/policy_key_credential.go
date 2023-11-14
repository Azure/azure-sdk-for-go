// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/errorinfo"
	"github.com/golang-jwt/jwt/v5"
)

type KeyCredentialPolicy struct {
	key string
}

func NewWebPubSubKeyCredentialPolicy(key string) *KeyCredentialPolicy {
	return &KeyCredentialPolicy{
		key: key,
	}
}

// Do implementes the Do method on the [policy.Polilcy] interface.
func (k *KeyCredentialPolicy) Do(req *policy.Request) (*http.Response, error) {
	val := k.key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": req.Raw().URL.String(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	secretKey := []byte(val) // Replace with your secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, errorinfo.NonRetriableError(errors.New("error signing the token"))
	}
	req.Raw().Header.Add("Authorization", "Bearer "+tokenString)
	return req.Next()
}
