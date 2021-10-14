// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azfile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/url"
)

// A ServiceClient represents a URL to the Azure Storage File service allowing you to manipulate file share.
type ServiceClient struct {
	client *serviceClient
	u      url.URL
	cred   azcore.Credential
}

// URL returns the URL endpoint used by the ServiceClient object.
func (s ServiceClient) URL() string {
	return s.client.con.u
}

// NewServiceClient creates a ServiceClient object using the specified URL, credential, and options.
// Example of serviceURL: https://<your_storage_account>.file.core.windows.net
func NewServiceClient(serviceURL string, cred azcore.Credential, options *ClientOptions) (ServiceClient, error) {
	u, err := url.Parse(serviceURL)
	if err != nil {
		return ServiceClient{}, err
	}

	return ServiceClient{client: &serviceClient{
		con: newConnection(serviceURL, cred, options.getConnectionOptions()),
	}, u: *u, cred: cred}, nil
}

// NewShareClient creates a new ShareURL object by concatenating shareName to the end of
// ServiceURL's URL. The new ShareURL uses the same request policy pipeline as the ServiceURL.
// To change the pipeline, create the ShareURL and then call its WithPipeline method passing in the
// desired pipeline object. Or, call this package's NewShareURL instead of calling this object's
// NewShareURL method.
func (s ServiceClient) NewShareClient(shareName string) ShareClient {
	shareURL := appendToURLPath(s.client.con.u, shareName)
	containerConnection := &connection{shareURL, s.client.con.p}
	return ShareClient{
		client: &shareClient{
			con: containerConnection,
		},
		cred: s.cred,
	}
}
