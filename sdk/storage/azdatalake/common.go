// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azdatalake

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/lease"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/shared"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/sas"
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

// ===================================== LEASE CONSTANTS ============================================================

// StatusType defines values for StatusType
type StatusType = lease.StatusType

const (
	StatusTypeLocked   StatusType = lease.StatusTypeLocked
	StatusTypeUnlocked StatusType = lease.StatusTypeUnlocked
)

// PossibleStatusTypeValues returns the possible values for the StatusType const type.
func PossibleStatusTypeValues() []StatusType {
	return lease.PossibleStatusTypeValues()
}

// DurationType defines values for DurationType
type DurationType = lease.DurationType

const (
	DurationTypeInfinite DurationType = lease.DurationTypeInfinite
	DurationTypeFixed    DurationType = lease.DurationTypeFixed
)

// PossibleDurationTypeValues returns the possible values for the DurationType const type.
func PossibleDurationTypeValues() []DurationType {
	return lease.PossibleDurationTypeValues()
}

// StateType defines values for StateType
type StateType = lease.StateType

const (
	StateTypeAvailable StateType = lease.StateTypeAvailable
	StateTypeLeased    StateType = lease.StateTypeLeased
	StateTypeExpired   StateType = lease.StateTypeExpired
	StateTypeBreaking  StateType = lease.StateTypeBreaking
	StateTypeBroken    StateType = lease.StateTypeBroken
)

// PossibleStateTypeValues returns the possible values for the StateType const type.
func PossibleStateTypeValues() []StateType {
	return lease.PossibleStateTypeValues()
}

// SessionMode specifies whether eligible requests are authenticated with a session. Sessions apply
// only to clients authenticated with an azcore.TokenCredential, and only to the file reads they send
// to the blob endpoint; requests to the DFS endpoint always use the bearer token.
type SessionMode = azblob.SessionMode

const (
	// SessionModeAuto is the zero value, and therefore the default when no value is specified.
	// The client library decides whether sessions are used, and that decision may change in a
	// future release. Currently, SessionModeAuto resolves to SessionModeDisabled.
	SessionModeAuto = azblob.SessionModeAuto

	// SessionModeDisabled authenticates every request with a bearer token.
	SessionModeDisabled = azblob.SessionModeDisabled

	// SessionModeEnabled authenticates file reads with a session, created and cached per file
	// system. It requires the storage account name: when it can be determined from neither
	// SessionOptions.AccountName nor the client's URL, client construction fails.
	SessionModeEnabled = azblob.SessionModeEnabled
)

// PossibleSessionModeValues returns the possible values for the SessionMode const type.
func PossibleSessionModeValues() []SessionMode {
	return azblob.PossibleSessionModeValues()
}

// SessionOptions configures session authentication; set it as ClientOptions.Session. Clients using
// a shared key, a SAS or no credential ignore it.
type SessionOptions = azblob.SessionOptions

// SessionProvider provides and caches the sessions used to authenticate eligible requests. Share
// one across clients through SessionOptions.Provider. Create one with NewContainerSessionProvider.
type SessionProvider = azblob.SessionProvider

// ContainerSessionProvider is a SessionProvider that creates sessions with an
// azcore.TokenCredential and caches one per file system.
type ContainerSessionProvider = azblob.ContainerSessionProvider

// NewContainerSessionProvider creates a ContainerSessionProvider. Pass it as
// SessionOptions.Provider to clients that should share its cached sessions, including clients
// created independently of one another or after others have been discarded.
//   - serviceURL - the URL of the storage account, e.g. https://<account>.dfs.core.windows.net/. A
//     DFS URL is converted to the account's blob endpoint, where sessions are created, and a file
//     system or path URL is reduced to its service URL.
//   - cred - an Azure AD credential, typically obtained via the azidentity module
//   - options - options for the pipeline that creates sessions; pass nil to accept the default values
func NewContainerSessionProvider(serviceURL string, cred azcore.TokenCredential, options *azcore.ClientOptions) (*ContainerSessionProvider, error) {
	blobURL, _ := shared.GetURLs(serviceURL)
	var blobOpts azblob.ClientOptions
	if options != nil {
		blobOpts.ClientOptions = *options
	}
	return azblob.NewContainerSessionProvider(blobURL, cred, &blobOpts)
}
