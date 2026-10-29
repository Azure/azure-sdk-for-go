// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package container

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
)

// SessionProvider is an azblob.SessionProvider that creates sessions with an
// azcore.TokenCredential and caches one per container.
type SessionProvider = exported.ContainerSessionProvider

// NewSessionProvider creates a SessionProvider. Pass it as azblob.SessionOptions.Provider to
// clients that should share its cached sessions, including clients created independently of one
// another or after others have been discarded.
//   - serviceURL - the URL of the blob service, e.g. https://<account>.blob.core.windows.net/. A
//     container or blob URL is reduced to its service URL.
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options for the pipeline that creates sessions; pass nil to accept the
//     default values
//
// The provider retains one cached session per container it is used with.
func NewSessionProvider(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*SessionProvider, error) {
	return base.NewContainerSessionProvider(serviceURL, cred, (*base.ClientOptions)(options))
}
