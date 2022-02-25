// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName string, accountKey string) (*SharedKeyCredential, error) {
	c := SharedKeyCredential{accountName: accountName}
	if err := c.SetAccountKey(accountKey); err != nil {
		return nil, err
	}
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

// SetAccountKey replaces the existing account key with the specified account key.
func (c *SharedKeyCredential) SetAccountKey(accountKey string) error {
	bytes, err := base64.StdEncoding.DecodeString(accountKey)
	if err != nil {
		return fmt.Errorf("decode account key: %w", err)
	}
	c.accountKey.Store(bytes)
	return nil
}

// computeHMACSHA256 generates a hash signature for an HTTP request or for a SAS.
func (c *SharedKeyCredential) ComputeHMACSHA256(message string) (string, error) {
	h := hmac.New(sha256.New, c.accountKey.Load().([]byte))
	_, err := h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), err
}

func (c *SharedKeyCredential) buildStringToSign(req *http.Request) (string, error) {
	// https://docs.microsoft.com/en-us/rest/api/storageservices/authentication-for-the-azure-storage-services
	headers := req.Header

	canonicalizedResource, err := c.buildCanonicalizedResource(req.URL)
	if err != nil {
		return "", err
	}

	stringToSign := strings.Join([]string{
		headers.Get(headerXmsDate),
		canonicalizedResource,
	}, "\n")
	return stringToSign, nil
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

	if compVal, ok := params["comp"]; ok {
		//do something here
		cr.WriteString("?" + "comp=" + compVal[0])
	}
	return cr.String(), nil
}

type sharedKeyCredPolicy struct {
	cred *SharedKeyCredential
}

func newSharedKeyCredPolicy(cred *SharedKeyCredential) *sharedKeyCredPolicy {
	s := &sharedKeyCredPolicy{
		cred: cred,
	}

	return s
}

func (s *sharedKeyCredPolicy) Do(req *policy.Request) (*http.Response, error) {
	if d := req.Raw().Header.Get(headerXmsDate); d == "" {
		req.Raw().Header.Set(headerXmsDate, time.Now().UTC().Format(http.TimeFormat))
	}
	stringToSign, err := s.cred.buildStringToSign(req.Raw())
	if err != nil {
		return nil, err
	}
	signature, err := s.cred.ComputeHMACSHA256(stringToSign)
	if err != nil {
		return nil, err
	}
	authHeader := strings.Join([]string{"SharedKeyLite ", s.cred.AccountName(), ":", signature}, "")
	req.Raw().Header.Set(headerAuthorization, authHeader)

	response, err := req.Next()
	if err != nil && response != nil && response.StatusCode == http.StatusForbidden {
		// Service failed to authenticate request, log it
		log.Write(log.EventResponse, "===== HTTP Forbidden status, String-to-Sign:\n"+stringToSign+"\n===============================\n")
	}
	return response, err
}
