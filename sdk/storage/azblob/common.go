// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

// SharedKeyCredential contains an account's name and its primary or secondary key.
type SharedKeyCredential = exported.SharedKeyCredential

// NewSharedKeyCredential creates an immutable SharedKeyCredential containing the
// storage account's name and either its primary or secondary key.
func NewSharedKeyCredential(accountName, accountKey string) (*SharedKeyCredential, error) {
	return exported.NewSharedKeyCredential(accountName, accountKey)
}

// URLParts object represents the components that make up an Azure Storage Container/Blob URL.
// NOTE: Changing any SAS-related field requires computing a new SAS signature.
type URLParts = sas.URLParts

// ParseURL parses a URL initializing URLParts' fields including any SAS-related & snapshot query parameters. Any other
// query parameters remain in the UnparsedParams field. This method overwrites all fields in the URLParts object.
func ParseURL(u string) (URLParts, error) {
	return sas.ParseURL(u)
}

// HTTPRange defines a range of bytes within an HTTP resource, starting at offset and
// ending at offset+count. A zero-value HTTPRange indicates the entire resource. An HTTPRange
// which has an offset and zero value count indicates from the offset to the resource's end.
type HTTPRange = exported.HTTPRange

// ExpectContinueMode is the mode for applying the HTTP "Expect: 100-continue" header to
// operations that include a request body.
type ExpectContinueMode = exported.ExpectContinueMode

const (
	// ExpectContinueModeApplyOnThrottle indicates that Expect-Continue will not be applied
	// until specific errors are encountered from the service, at which point it will be
	// applied for a fixed window of time after the last triggering error. This is the default.
	ExpectContinueModeApplyOnThrottle = exported.ExpectContinueModeApplyOnThrottle

	// ExpectContinueModeOn indicates Expect-Continue will be applied regardless of recent
	// error status. The ContentLengthThreshold option still applies.
	ExpectContinueModeOn = exported.ExpectContinueModeOn

	// ExpectContinueModeOff indicates Expect-Continue will never be applied.
	ExpectContinueModeOff = exported.ExpectContinueModeOff
)

// ExpectContinueOptions configures the behavior for applying the HTTP "Expect: 100-continue"
// header to operations that include a request body.
type ExpectContinueOptions = exported.ExpectContinueOptions

// SessionMode specifies whether eligible requests are authenticated with a session. Sessions
// apply only to clients authenticated with an azcore.TokenCredential.
type SessionMode = exported.SessionMode

const (
	// SessionModeAuto is the zero value, and therefore the default when no value is specified.
	// The client library decides whether sessions are used, and that decision may change in a
	// future release. Currently, SessionModeAuto resolves to SessionModeDisabled.
	SessionModeAuto = exported.SessionModeAuto

	// SessionModeDisabled authenticates every request with a bearer token.
	SessionModeDisabled = exported.SessionModeDisabled

	// SessionModeEnabled authenticates eligible requests with a session, created and cached per
	// container. Currently, Get Blob (blob.Client.DownloadStream and the downloads built on it) is
	// the only eligible operation. It requires the storage account name: when it can be
	// determined from neither SessionOptions.AccountName nor the client's URL, client
	// construction fails.
	SessionModeEnabled = exported.SessionModeEnabled
)

// PossibleSessionModeValues returns the possible values for the SessionMode const type.
func PossibleSessionModeValues() []SessionMode {
	return exported.PossibleSessionModeValues()
}

// SessionOptions configures session authentication. Sessions apply only to clients
// authenticated with an azcore.TokenCredential; clients using a shared key, a SAS or no
// credential ignore these options.
type SessionOptions = exported.SessionOptions

// SessionProvider provides and caches the sessions used to authenticate eligible requests. Share
// one across clients through SessionOptions.Provider. It can't be implemented outside this
// module; create one with NewContainerSessionProvider.
type SessionProvider = exported.SessionProvider

// ContainerSessionProvider is a SessionProvider that creates sessions with an
// azcore.TokenCredential and caches one per container.
type ContainerSessionProvider = exported.ContainerSessionProvider

// NewContainerSessionProvider creates a ContainerSessionProvider. Pass it as
// SessionOptions.Provider to clients that should share its cached sessions, including clients
// created independently of one another or after others have been discarded.
//   - serviceURL - the URL of the blob service, e.g. https://<account>.blob.core.windows.net/. A
//     container or blob URL is reduced to its service URL.
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - client options for the pipeline that creates sessions; pass nil to accept the
//     default values
//
// The provider retains one cached session per container it is used with.
func NewContainerSessionProvider(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*ContainerSessionProvider, error) {
	return base.NewContainerSessionProvider(serviceURL, cred, (*base.ClientOptions)(options))
}
