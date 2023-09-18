//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

// HMACPolicy is a pipeline policy that implements HMAC authentication.
// https://learn.microsoft.com/en-us/azure/azure-app-configuration/rest-api-authentication-hmac
type HMACPolicy struct {
	credential string
	secret     []byte
}

// NewHMACPolicy creates a new instance of [HMACPolicy].
func NewHMACPolicy(credential string, secret []byte) *HMACPolicy {
	return &HMACPolicy{
		credential: credential,
		secret:     secret,
	}
}

// Do implements the policy.Policy interface on the [HMACPolicy] type.
func (policy *HMACPolicy) Do(request *policy.Request) (*http.Response, error) {
	req := request.Raw()

	pathAndQuery := req.URL.EscapedPath()
	if req.URL.RawQuery != "" {
		pathAndQuery = pathAndQuery + "?" + req.URL.RawQuery
	}

	var content []byte
	if body := request.Body(); body != nil {
		var err error
		if content, err = io.ReadAll(body); err != nil {
			return nil, err
		}
		// restore the body after reading
		if err = request.RewindBody(); err != nil {
			return nil, err
		}
	}

	timestamp := time.Now().UTC().Format(http.TimeFormat)

	contentHash, err := getContentHashBase64(content)
	if err != nil {
		return nil, err
	}

	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s", strings.ToUpper(req.Method), pathAndQuery, timestamp, req.URL.Host, contentHash)

	signature, err := getHmac(stringToSign, policy.secret)
	if err != nil {
		return nil, err
	}

	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("Date", timestamp)
	req.Header.Set("Authorization", "HMAC-SHA256 Credential="+policy.credential+", SignedHeaders=date;host;x-ms-content-sha256, Signature="+signature)

	return request.Next()
}

func getContentHashBase64(content []byte) (string, error) {
	hasher := sha256.New()

	_, err := hasher.Write(content)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func getHmac(content string, key []byte) (string, error) {
	hmac := hmac.New(sha256.New, key)

	_, err := hmac.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hmac.Sum(nil)), nil
}

// ParseConnectionString parses the provided connection string.
func ParseConnectionString(connectionString string) (endpoint string, credential string, secret []byte, err error) {
	const connectionStringEndpointPrefix = "Endpoint="
	const connectionStringCredentialPrefix = "Id="
	const connectionStringSecretPrefix = "Secret="

	var er error = errors.New("error parsing connection string")
	var ept *string
	var cred *string
	var sec *[]byte
	for _, seg := range strings.Split(connectionString, ";") {
		if strings.HasPrefix(seg, connectionStringEndpointPrefix) {
			if ept != nil {
				return "", "", []byte{}, er
			}

			ep := strings.TrimPrefix(seg, connectionStringEndpointPrefix)
			ept = &ep
		} else if strings.HasPrefix(seg, connectionStringCredentialPrefix) {
			if cred != nil {
				return "", "", []byte{}, er
			}

			c := strings.TrimPrefix(seg, connectionStringCredentialPrefix)
			cred = &c
		} else if strings.HasPrefix(seg, connectionStringSecretPrefix) {
			if sec != nil {
				return "", "", []byte{}, er
			}

			s, e := base64.StdEncoding.DecodeString(strings.TrimPrefix(seg, connectionStringSecretPrefix))
			if e != nil {
				return "", "", []byte{}, e
			}

			sec = &s
		}
	}

	if ept == nil || cred == nil || sec == nil {
		return "", "", []byte{}, er
	}

	return *ept, *cred, *sec, nil
}
