// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"encoding/xml"
	"strings"
)

// ETag is an entity tag.
// TODO: does this belong in azcore?
type ETag string

const (
	// ETagNone represents an empty entity tag.
	ETagNone ETag = ""

	// ETagAny matches any entity tag.
	ETagAny ETag = "*"
)

// Metadata contains metadata key/value pairs.
// TODO: does this belong in azcore?
type Metadata map[string]string

const mdPrefix = "x-ms-meta-"

const mdPrefixLen = len(mdPrefix)

// UnmarshalXML implements the xml.Unmarshaler interface for Metadata.
func (md *Metadata) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tokName := ""
	for t, err := d.Token(); err == nil; t, err = d.Token() {
		switch tt := t.(type) {
		case xml.StartElement:
			tokName = strings.ToLower(tt.Name.Local)
			break
		case xml.CharData:
			if *md == nil {
				*md = Metadata{}
			}
			(*md)[tokName] = string(tt)
			break
		}
	}
	return nil
}

// LeaseDurationType enumerates the values for lease duration type.
type LeaseDurationType string

const (
	// LeaseDurationFixed ...
	LeaseDurationFixed LeaseDurationType = "fixed"
	// LeaseDurationInfinite ...
	LeaseDurationInfinite LeaseDurationType = "infinite"
	// LeaseDurationNone represents an empty LeaseDurationType.
	LeaseDurationNone LeaseDurationType = ""
)

// PossibleLeaseDurationTypeValues returns an array of possible values for the LeaseDurationType const type.
func PossibleLeaseDurationTypeValues() []LeaseDurationType {
	return []LeaseDurationType{LeaseDurationFixed, LeaseDurationInfinite, LeaseDurationNone}
}

// LeaseStateType enumerates the values for lease state type.
type LeaseStateType string

const (
	// LeaseStateAvailable ...
	LeaseStateAvailable LeaseStateType = "available"
	// LeaseStateBreaking ...
	LeaseStateBreaking LeaseStateType = "breaking"
	// LeaseStateBroken ...
	LeaseStateBroken LeaseStateType = "broken"
	// LeaseStateExpired ...
	LeaseStateExpired LeaseStateType = "expired"
	// LeaseStateLeased ...
	LeaseStateLeased LeaseStateType = "leased"
	// LeaseStateNone represents an empty LeaseStateType.
	LeaseStateNone LeaseStateType = ""
)

// PossibleLeaseStateTypeValues returns an array of possible values for the LeaseStateType const type.
func PossibleLeaseStateTypeValues() []LeaseStateType {
	return []LeaseStateType{LeaseStateAvailable, LeaseStateBreaking, LeaseStateBroken, LeaseStateExpired, LeaseStateLeased, LeaseStateNone}
}

// LeaseStatusType enumerates the values for lease status type.
type LeaseStatusType string

const (
	// LeaseStatusLocked ...
	LeaseStatusLocked LeaseStatusType = "locked"
	// LeaseStatusNone represents an empty LeaseStatusType.
	LeaseStatusNone LeaseStatusType = ""
	// LeaseStatusUnlocked ...
	LeaseStatusUnlocked LeaseStatusType = "unlocked"
)

// ListContainersIncludeType enumerates the values for list containers include type.
type ListContainersIncludeType string

const (
	// ListContainersIncludeMetadata ...
	ListContainersIncludeMetadata ListContainersIncludeType = "metadata"
	// ListContainersIncludeNone represents an empty ListContainersIncludeType.
	ListContainersIncludeNone ListContainersIncludeType = ""
)

// PublicAccessType enumerates the values for public access type.
type PublicAccessType string

const (
	// PublicAccessBlob ...
	PublicAccessBlob PublicAccessType = "blob"
	// PublicAccessContainer ...
	PublicAccessContainer PublicAccessType = "container"
	// PublicAccessNone represents an empty PublicAccessType.
	PublicAccessNone PublicAccessType = ""
)

// PossiblePublicAccessTypeValues returns an array of possible values for the PublicAccessType const type.
func PossiblePublicAccessTypeValues() []PublicAccessType {
	return []PublicAccessType{PublicAccessBlob, PublicAccessContainer, PublicAccessNone}
}

// ContainerItem - An Azure Storage container
type ContainerItem struct {
	// XMLName is used for marshalling and is subject to removal in a future release.
	XMLName    xml.Name            `xml:"Container"`
	Name       string              `xml:"Name"`
	Properties ContainerProperties `xml:"Properties"`
	Metadata   Metadata            `xml:"Metadata"`
}

// ContainerProperties - Properties of a container
type ContainerProperties struct {
	// TODO: LastModified time.Time `xml:"Last-Modified"`
	Etag ETag `xml:"Etag"`
	// LeaseStatus - Possible values include: 'LeaseStatusLocked', 'LeaseStatusUnlocked', 'LeaseStatusNone'
	LeaseStatus LeaseStatusType `xml:"LeaseStatus"`
	// LeaseState - Possible values include: 'LeaseStateAvailable', 'LeaseStateLeased', 'LeaseStateExpired', 'LeaseStateBreaking', 'LeaseStateBroken', 'LeaseStateNone'
	LeaseState LeaseStateType `xml:"LeaseState"`
	// LeaseDuration - Possible values include: 'LeaseDurationInfinite', 'LeaseDurationFixed', 'LeaseDurationNone'
	LeaseDuration LeaseDurationType `xml:"LeaseDuration"`
	// PublicAccess - Possible values include: 'PublicAccessContainer', 'PublicAccessBlob', 'PublicAccessNone'
	PublicAccess          PublicAccessType `xml:"PublicAccess"`
	HasImmutabilityPolicy *bool            `xml:"HasImmutabilityPolicy"`
	HasLegalHold          *bool            `xml:"HasLegalHold"`
}

type ListContainersOptions struct {
	Prefix     *string
	Marker     *string
	Maxresults *int32
	Include    ListContainersIncludeType
	Timeout    *int32
	RequestID  *string
}

// ListContainersPage - An enumeration of containers
type ListContainersPage struct {
	// XMLName is used for marshalling and is subject to removal in a future release.
	XMLName         xml.Name        `xml:"EnumerationResults"`
	ServiceEndpoint string          `xml:"ServiceEndpoint,attr"`
	Prefix          *string         `xml:"Prefix"`
	Marker          *string         `xml:"Marker"`
	MaxResults      *int32          `xml:"MaxResults"`
	ContainerItems  []ContainerItem `xml:"Containers>Container"`
	NextMarker      *string         `xml:"NextMarker"`
}
