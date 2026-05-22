// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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

	signature, err := getHMAC(stringToSign, policy.secret)
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

func getHMAC(content string, key []byte) (string, error) {
	hmac := hmac.New(sha256.New, key)

	_, err := hmac.Write([]byte(content))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hmac.Sum(nil)), nil
}

// ParseConnectionString parses the provided connection string.
// Returns endpoint, cred, secret or an error.
func ParseConnectionString(connectionString string) (string, string, []byte, error) {
	const (
		endpointPrefix   = "Endpoint="
		credentialPrefix = "Id="
		secretPrefix     = "Secret="
	)

	var (
		ept  string
		cred string
		sec  []byte
	)

	const duplicateSection = "duplicate %s section"

	for _, seg := range strings.Split(connectionString, ";") {
		if strings.HasPrefix(seg, endpointPrefix) {
			if ept != "" {
				return "", "", nil, fmt.Errorf(duplicateSection, endpointPrefix)
			}

			ep := strings.TrimPrefix(seg, endpointPrefix)
			ept = ep
		} else if strings.HasPrefix(seg, credentialPrefix) {
			if cred != "" {
				return "", "", nil, fmt.Errorf(duplicateSection, credentialPrefix)
			}

			c := strings.TrimPrefix(seg, credentialPrefix)
			cred = c
		} else if strings.HasPrefix(seg, secretPrefix) {
			if sec != nil {
				return "", "", nil, fmt.Errorf(duplicateSection, secretPrefix)
			}

			s, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(seg, secretPrefix))
			if err != nil {
				return "", "", nil, err
			}

			sec = s
		}
	}

	const missingSection = "missing %s section"

	if ept == "" {
		return "", "", nil, fmt.Errorf(missingSection, endpointPrefix)
	}

	if cred == "" {
		return "", "", nil, fmt.Errorf(missingSection, credentialPrefix)
	}

	if sec == nil {
		return "", "", nil, fmt.Errorf(missingSection, secretPrefix)
	}

	return ept, cred, sec, nil
}
