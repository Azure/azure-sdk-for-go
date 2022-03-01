//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type hmacAuthenticationPolicy struct {
	credential string
	secret     []byte
}

func newHmacAuthenticationPolicy(credential string, secret []byte) *hmacAuthenticationPolicy {
	return &hmacAuthenticationPolicy{
		credential: credential,
		secret:     secret,
	}
}

func (policy *hmacAuthenticationPolicy) Do(request *policy.Request) (*http.Response, error) {
	req := request.Raw()
	id := policy.credential
	key := policy.secret

	method := req.Method
	host := req.URL.Host
	pathAndQuery := req.URL.Path
	if req.URL.RawQuery != "" {
		pathAndQuery = pathAndQuery + "?" + req.URL.RawQuery
	}

	var content []byte
	if req.Body != nil {
		var err error
		if content, err = ioutil.ReadAll(req.Body); err != nil {
			return nil, err
		}
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(content))

	timestamp := time.Now().UTC().Format(http.TimeFormat)
	contentHash := getContentHashBase64(content)
	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s", strings.ToUpper(method), pathAndQuery, timestamp, host, contentHash)
	signature := getHmac(stringToSign, key)

	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("x-ms-date", timestamp)
	req.Header.Set("Authorization", "HMAC-SHA256 Credential="+id+", SignedHeaders=x-ms-date;host;x-ms-content-sha256, Signature="+signature)

	return request.Next()
}

func signRequest(id string, secret string, req *http.Request) error {
	method := req.Method
	host := req.URL.Host
	pathAndQuery := req.URL.Path
	if req.URL.RawQuery != "" {
		pathAndQuery = pathAndQuery + "?" + req.URL.RawQuery
	}

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(content))

	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(http.TimeFormat)
	contentHash := getContentHashBase64(content)
	stringToSign := fmt.Sprintf("%s\n%s\n%s;%s;%s", strings.ToUpper(method), pathAndQuery, timestamp, host, contentHash)
	signature := getHmac(stringToSign, key)

	req.Header.Set("x-ms-content-sha256", contentHash)
	req.Header.Set("x-ms-date", timestamp)
	req.Header.Set("Authorization", "HMAC-SHA256 Credential="+id+", SignedHeaders=x-ms-date;host;x-ms-content-sha256, Signature="+signature)

	return nil
}

func getContentHashBase64(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func getHmac(content string, key []byte) string {
	hmac := hmac.New(sha256.New, key)
	hmac.Write([]byte(content))
	return base64.StdEncoding.EncodeToString(hmac.Sum(nil))
}

func parseConnectionString(connectionString string) (endpoint string, credential string, secret []byte, err error) {
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
