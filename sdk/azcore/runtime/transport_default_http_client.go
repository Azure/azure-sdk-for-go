//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

var defaultHTTPClient *http.Client

func init() {
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.TLSClientConfig.MinVersion = tls.VersionTLS12
	defaultHTTPClient = &http.Client{
		Transport: defaultTransport,
	}
}

// used to adapt a TransportPolicy to a Policy
type transportPolicy struct {
	trans policy.Transporter
}

func (tp transportPolicy) Do(req *policy.Request) (*http.Response, error) {
	resp, err := tp.trans.Do(req.Raw())
	if err != nil {
		return nil, err
	} else if resp == nil {
		// there was no response and no error (rare but can happen)
		// this ensures the retry policy will retry the request
		return nil, errors.New("received nil response")
	}
	return resp, nil
}
