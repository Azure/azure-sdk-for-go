// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package rntbd

import (
	"fmt"
	"net/url"
	"strings"
)

// Protocol represents the communication protocol for Cosmos DB connections.
type Protocol int

const (
	// ProtocolHTTPS represents the HTTPS protocol.
	ProtocolHTTPS Protocol = iota
	// ProtocolTCP represents the TCP/RNTBD protocol.
	ProtocolTCP
)

// String returns the name of the protocol.
func (p Protocol) String() string {
	switch p {
	case ProtocolHTTPS:
		return "Https"
	case ProtocolTCP:
		return "Tcp"
	default:
		return "Unknown"
	}
}

// Scheme returns the URI scheme for the protocol.
func (p Protocol) Scheme() string {
	switch p {
	case ProtocolHTTPS:
		return "https"
	case ProtocolTCP:
		return "rntbd"
	default:
		return ""
	}
}

// SchemeToProtocol converts a URI scheme string to a Protocol.
func SchemeToProtocol(scheme string) (Protocol, error) {
	switch strings.ToLower(scheme) {
	case "https":
		return ProtocolHTTPS, nil
	case "rntbd":
		return ProtocolTCP, nil
	default:
		return ProtocolHTTPS, fmt.Errorf("unknown protocol scheme: %s", scheme)
	}
}

// Uri wraps a URI string with its parsed representation.
type Uri struct {
	uriAsString string
	uri         *url.URL
}

// NewUri creates a new Uri from a string.
func NewUri(uriString string) Uri {
	parsed, err := url.Parse(uriString)
	if err != nil {
		parsed = nil
	}
	return Uri{
		uriAsString: uriString,
		uri:         parsed,
	}
}

// Create creates a new Uri from a string (factory method).
func UriCreate(uriString string) Uri {
	return NewUri(uriString)
}

// GetURI returns the parsed URL.
func (u Uri) GetURI() *url.URL {
	return u.uri
}

// GetURIAsString returns the original URI string.
func (u Uri) GetURIAsString() string {
	return u.uriAsString
}

// String returns the URI as a string.
func (u Uri) String() string {
	return u.uriAsString
}

// Equal checks if two URIs are equal.
func (u Uri) Equal(other Uri) bool {
	return u.uriAsString == other.uriAsString
}

// AddressInformation encapsulates physical address information for Cosmos DB.
type AddressInformation struct {
	protocol    Protocol
	isPublic    bool
	isPrimary   bool
	physicalUri Uri
}

// NewAddressInformation creates a new AddressInformation with the given parameters.
func NewAddressInformation(isPublic, isPrimary bool, physicalUri string, protocol Protocol) AddressInformation {
	return AddressInformation{
		protocol:    protocol,
		isPublic:    isPublic,
		isPrimary:   isPrimary,
		physicalUri: NewUri(normalizePhysicalUri(physicalUri)),
	}
}

// NewAddressInformationFromScheme creates a new AddressInformation using a protocol scheme string.
func NewAddressInformationFromScheme(isPublic, isPrimary bool, physicalUri string, protocolScheme string) (AddressInformation, error) {
	protocol, err := SchemeToProtocol(protocolScheme)
	if err != nil {
		return AddressInformation{}, err
	}
	return NewAddressInformation(isPublic, isPrimary, physicalUri, protocol), nil
}

// normalizePhysicalUri normalizes the physical URI by trimming trailing slashes and adding a single one.
// Backend returns non-normalized URIs with "//" tail, e.g.,
// https://cdb-ms-prod-westus2-fd2.documents.azure.com:15248/apps/.../replicas/132077748219659199s//
// We should trim the tail double "//" and add a single "/"
func normalizePhysicalUri(physicalUri string) string {
	if physicalUri == "" {
		return physicalUri
	}

	// Trim all trailing slashes
	i := len(physicalUri) - 1
	for i >= 0 && physicalUri[i] == '/' {
		i--
	}

	return physicalUri[:i+1] + "/"
}

// IsPublic returns whether this is a public address.
func (a AddressInformation) IsPublic() bool {
	return a.isPublic
}

// IsPrimary returns whether this is the primary replica.
func (a AddressInformation) IsPrimary() bool {
	return a.isPrimary
}

// GetPhysicalUri returns the physical URI.
func (a AddressInformation) GetPhysicalUri() Uri {
	return a.physicalUri
}

// GetProtocol returns the protocol.
func (a AddressInformation) GetProtocol() Protocol {
	return a.protocol
}

// GetProtocolName returns the protocol name (e.g., "Https", "Tcp").
func (a AddressInformation) GetProtocolName() string {
	return a.protocol.String()
}

// GetProtocolScheme returns the protocol scheme (e.g., "https", "rntbd").
func (a AddressInformation) GetProtocolScheme() string {
	return a.protocol.Scheme()
}

// String returns a string representation of AddressInformation.
func (a AddressInformation) String() string {
	return fmt.Sprintf("AddressInformation{protocol='%s', isPublic=%v, isPrimary=%v, physicalUri='%s'}",
		a.protocol, a.isPublic, a.isPrimary, a.physicalUri)
}

// Equal checks if two AddressInformation are equal.
func (a AddressInformation) Equal(other AddressInformation) bool {
	return a.protocol == other.protocol &&
		a.isPublic == other.isPublic &&
		a.isPrimary == other.isPrimary &&
		a.physicalUri.Equal(other.physicalUri)
}
