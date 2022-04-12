//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azfile

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// ---------------------------------------------------------------------------------------------------------------------

type ShareCreateOptions struct {
	// Specifies the access tier of the share.
	AccessTier *ShareAccessTier
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string
	// Specifies the maximum size of the share, in gigabytes.
	Quota *int32
}

func (o *ShareCreateOptions) format() *shareClientCreateOptions {
	if o == nil {
		return nil
	}

	return &shareClientCreateOptions{
		AccessTier: o.AccessTier,
		Metadata:   o.Metadata,
		Quota:      o.Quota,
	}
}

type ShareCreateResponse struct {
	shareClientCreateResponse
}

func toShareCreateResponse(resp shareClientCreateResponse) ShareCreateResponse {
	return ShareCreateResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareCreateSnapshotOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string
}

func (o *ShareCreateSnapshotOptions) format() *shareClientCreateSnapshotOptions {
	if o == nil {
		return nil
	}
	return &shareClientCreateSnapshotOptions{Metadata: o.Metadata}
}

type ShareCreateSnapshotResponse struct {
	shareClientCreateSnapshotResponse
}

func toShareCreateSnapshotResponse(resp shareClientCreateSnapshotResponse) ShareCreateSnapshotResponse {
	return ShareCreateSnapshotResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareDeleteOptions struct {
	// Specifies the option include to delete the base share and all of its snapshots.
	DeleteSnapshots *DeleteSnapshotsOptionType
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ShareDeleteOptions) format() (*shareClientDeleteOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &shareClientDeleteOptions{
		DeleteSnapshots: o.DeleteSnapshots,
		Sharesnapshot:   o.ShareSnapshot,
	}, o.LeaseAccessConditions
}

type ShareDeleteResponse struct {
	shareClientDeleteResponse
}

func toShareDeleteResponse(resp shareClientDeleteResponse) ShareDeleteResponse {
	return ShareDeleteResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareGetPropertiesOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ShareGetPropertiesOptions) format() (*shareClientGetPropertiesOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &shareClientGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}, o.LeaseAccessConditions
}

type ShareGetPropertiesResponse struct {
	shareClientGetPropertiesResponse
}

func toShareGetPropertiesResponse(resp shareClientGetPropertiesResponse) ShareGetPropertiesResponse {
	return ShareGetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareSetPropertiesOptions struct {
	// Specifies the access tier of the share.
	AccessTier *ShareAccessTier
	// Specifies the maximum size of the share, in gigabytes.
	Quota *int32

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ShareSetPropertiesOptions) format() (*shareClientSetPropertiesOptions, *LeaseAccessConditions, error) {
	if o == nil {
		return nil, nil, nil
	}

	if o.Quota != nil && *o.Quota < 0 {
		return nil, nil, errors.New("validation failed: share quote cannot be negative")
	}

	return &shareClientSetPropertiesOptions{
		AccessTier: o.AccessTier,
		Quota:      o.Quota,
	}, o.LeaseAccessConditions, nil
}

type ShareSetPropertiesResponse struct {
	shareClientSetPropertiesResponse
}

func toShareSetPropertiesResponse(resp shareClientSetPropertiesResponse) ShareSetPropertiesResponse {
	return ShareSetPropertiesResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareSetMetadataOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ShareSetMetadataOptions) format(metadata map[string]string) (*shareClientSetMetadataOptions, *LeaseAccessConditions, error) {
	if metadata == nil {
		return nil, nil, errors.New("metadata cannot be nil")
	}

	if o == nil {
		return &shareClientSetMetadataOptions{Metadata: metadata}, nil, nil
	}

	return &shareClientSetMetadataOptions{Metadata: metadata}, o.LeaseAccessConditions, nil
}

type ShareSetMetadataResponse struct {
	shareClientSetMetadataResponse
}

func toShareSetMetadataResponse(resp shareClientSetMetadataResponse) ShareSetMetadataResponse {
	return ShareSetMetadataResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareGetAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ShareGetAccessPolicyOptions) format() (*shareClientGetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

type ShareGetAccessPolicyResponse struct {
	shareClientGetAccessPolicyResponse
}

func toShareGetAccessPolicyResponse(resp shareClientGetAccessPolicyResponse) ShareGetAccessPolicyResponse {
	return ShareGetAccessPolicyResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareCreatePermissionOptions struct {
}

func (o *ShareCreatePermissionOptions) format(sharePermission string) (SharePermission, *shareClientCreatePermissionOptions) {
	return SharePermission{Permission: to.Ptr(sharePermission)}, nil
}

type ShareCreatePermissionResponse struct {
	shareClientCreatePermissionResponse
}

func toShareCreatePermissionResponse(resp shareClientCreatePermissionResponse) ShareCreatePermissionResponse {
	return ShareCreatePermissionResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareGetPermissionOptions struct {
}

func (o *ShareGetPermissionOptions) format() *shareClientGetPermissionOptions {
	return nil
}

type ShareGetPermissionResponse struct {
	shareClientGetPermissionResponse
}

func toShareGetPermissionResponse(resp shareClientGetPermissionResponse) ShareGetPermissionResponse {
	return ShareGetPermissionResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

// SetShareAccessPolicyOptions contains the optional parameters for the Share.SetAccessPolicy method.
type SetShareAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetShareAccessPolicyOptions) format(shareACLs []*SignedIdentifier) (*shareClientSetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return &shareClientSetAccessPolicyOptions{ShareACL: shareACLs}, nil
	}
	return &shareClientSetAccessPolicyOptions{ShareACL: shareACLs}, o.LeaseAccessConditions
}

type ShareSetAccessPolicyResponse struct {
	shareClientSetAccessPolicyResponse
}

func toShareSetAccessPolicyResponse(resp shareClientSetAccessPolicyResponse) ShareSetAccessPolicyResponse {
	return ShareSetAccessPolicyResponse{resp}
}

// ---------------------------------------------------------------------------------------------------------------------

type ShareGetStatisticsOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *ShareGetStatisticsOptions) format() (*shareClientGetStatisticsOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}

type ShareGetStatisticsResponse struct {
	shareClientGetStatisticsResponse
}

func toShareGetStatisticsResponse(resp shareClientGetStatisticsResponse) ShareGetStatisticsResponse {
	return ShareGetStatisticsResponse{resp}
}
