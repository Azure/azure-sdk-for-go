//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net"
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

func (policy *hmacAuthenticationPolicy) Do(req *policy.Request) (*http.Response, error) {
	if bd := req.Body(); bd != nil {
		if body, err := ioutil.ReadAll(bd); err == nil {
			h := hmac.New(sha256.New, policy.secret)
			if _, err := h.Write(body); err == nil {
				url := req.Raw().URL
				if host, _, err := net.SplitHostPort(url.Host); err == nil {
					contentHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
					utcNowString := time.Now().UTC().Format(http.TimeFormat)

					pathAndQuery := req.Raw().URL.Path
					if query := url.RawQuery; query != "" {
						pathAndQuery = pathAndQuery + "?" + query
					}

					stringToSign := req.Raw().Method + "\n" + pathAndQuery + "\n" + utcNowString + ";" + host + ";" + contentHash

					h = hmac.New(sha256.New, policy.secret)
					if _, err := h.Write([]byte(stringToSign)); err == nil {
						signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

						req.Raw().Header["x-ms-content-sha256"] = []string{contentHash}
						req.Raw().Header["Date"] = []string{utcNowString}

						req.Raw().Header["Authorizarion"] = []string{
							"HMAC-SHA256 Credential=" + policy.credential + "&SignedHeaders=date;host;x-ms-content-sha256&Signature=" + signature,
						}
					}
				}
			}
		}
	}

	return req.Next()
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
