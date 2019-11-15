// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azstorage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	bytes, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return nil, fmt.Errorf("decode account key: %w", err)
	}
	c := SharedKeyCredential{accountName: accountName}
	c.accountKey.Store(bytes)
	return &c, nil
}

// SharedKeyCredential contains an account's name and its primary or secondary key.
// It is immutable making it shareable and goroutine-safe.
type SharedKeyCredential struct {
	// Only the NewSharedKeyCredential method should set these; all other methods should treat them as read-only
	accountName string
	accountKey  atomic.Value // []byte
}

// AccountName returns the Storage account's name.
func (c *SharedKeyCredential) AccountName() string {
	return c.accountName
}

// ComputeHMACSHA256 generates a hash signature for an HTTP request or for a SAS.
func (c *SharedKeyCredential) ComputeHMACSHA256(message string) (base64String string) {
	h := hmac.New(sha256.New, c.accountKey.Load().([]byte))
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func (c *SharedKeyCredential) buildStringToSign(req *http.Request) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	headers := req.Header
	contentLength := headers.Get(azcore.HeaderContentLength)
	if contentLength == "0" {
		contentLength = ""
	}

	canonicalizedResource, err := c.buildCanonicalizedResource(req.URL)
	if err != nil {
		return "", err
	}

	stringToSign := strings.Join([]string{
		req.Method,
		headers.Get(azcore.HeaderContentEncoding),
		headers.Get(azcore.HeaderContentLanguage),
		contentLength,
		headers.Get(azcore.HeaderContentMD5),
		headers.Get(azcore.HeaderContentType),
		"", // Empty date because x-ms-date is expected (as per web page above)
		headers.Get(azcore.HeaderIfModifiedSince),
		headers.Get(azcore.HeaderIfMatch),
		headers.Get(azcore.HeaderIfNoneMatch),
		headers.Get(azcore.HeaderIfUnmodifiedSince),
		headers.Get(azcore.HeaderRange),
		c.buildCanonicalizedHeader(headers),
		canonicalizedResource,
	}, "\n")
	return stringToSign, nil
}

func (c *SharedKeyCredential) buildCanonicalizedHeader(headers http.Header) string {
	cm := map[string][]string{}
	for k, v := range headers {
		headerName := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(headerName, "x-ms-") {
			cm[headerName] = v // NOTE: the value must not have any whitespace around it.
		}
	}
	if len(cm) == 0 {
		return ""
	}

	keys := make([]string, 0, len(cm))
	for key := range cm {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ch := bytes.NewBufferString("")
	for i, key := range keys {
		if i > 0 {
			ch.WriteRune('\n')
		}
		ch.WriteString(key)
		ch.WriteRune(':')
		ch.WriteString(strings.Join(cm[key], ","))
	}
	return string(ch.Bytes())
}

func (c *SharedKeyCredential) buildCanonicalizedResource(u *url.URL) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	cr := bytes.NewBufferString("/")
	cr.WriteString(c.accountName)

	if len(u.Path) > 0 {
		// Any portion of the CanonicalizedResource string that is derived from
		// the resource's URI should be encoded exactly as it is in the URI.
		// -- https://msdn.microsoft.com/en-gb/library/azure/dd179428.aspx
		cr.WriteString(u.EscapedPath())
	} else {
		// a slash is required to indicate the root path
		cr.WriteString("/")
	}

	// params is a map[string][]string; param name is key; params values is []string
	params, err := url.ParseQuery(u.RawQuery) // Returns URL decoded values
	if err != nil {
		return "", fmt.Errorf("failed to parse query params: %w", err)
	}

	if len(params) > 0 { // There is at least 1 query parameter
		paramNames := []string{} // We use this to sort the parameter key names
		for paramName := range params {
			paramNames = append(paramNames, paramName) // paramNames must be lowercase
		}
		sort.Strings(paramNames)

		for _, paramName := range paramNames {
			paramValues := params[paramName]
			sort.Strings(paramValues)

			// Join the sorted key values separated by ','
			// Then prepend "keyName:"; then add this string to the buffer
			cr.WriteString("\n" + paramName + ":" + strings.Join(paramValues, ","))
		}
	}
	return string(cr.Bytes()), nil
}

// AuthenticationPolicy implements the Credential interface on SharedKeyCredential.
func (c *SharedKeyCredential) AuthenticationPolicy(azcore.AuthenticationPolicyOptions) azcore.Policy {
	return azcore.PolicyFunc(func(ctx context.Context, req *azcore.Request) (*azcore.Response, error) {
		// Add a x-ms-date header if it doesn't already exist
		if d := req.Request.Header.Get(azcore.HeaderXmsDate); d == "" {
			req.Request.Header[azcore.HeaderXmsDate] = []string{time.Now().UTC().Format(http.TimeFormat)}
		}
		stringToSign, err := c.buildStringToSign(req.Request)
		if err != nil {
			return nil, err
		}
		signature := c.ComputeHMACSHA256(stringToSign)
		authHeader := strings.Join([]string{"SharedKey ", c.AccountName(), ":", signature}, "")
		req.Request.Header[azcore.HeaderAuthorization] = []string{authHeader}

		response, err := req.Do(ctx)
		if err != nil && response != nil && response.StatusCode == http.StatusForbidden {
			// Service failed to authenticate request, log it
			azcore.Log().Write(azcore.LogError, "===== HTTP Forbidden status, String-to-Sign:\n"+stringToSign+"\n===============================\n")
		}
		return response, err
	})
}
