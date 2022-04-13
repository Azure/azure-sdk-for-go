//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

// CreateContainerOptions provides set of configurations for CreateContainer operation
type CreateContainerOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Specifies the encryption scope settings to set on the container.
	cpkScope *ContainerCpkScopeInfo
}

func (o *CreateContainerOptions) pointers() (*ContainerCreateOptions, *ContainerCpkScopeInfo) {
	if o == nil {
		return nil, nil
	}

	basicOptions := ContainerCreateOptions{
		Access:   o.Access,
		Metadata: o.Metadata,
	}

	return &basicOptions, o.cpkScope
}

// DeleteContainerOptions provides set of configurations for DeleteContainer operation
type DeleteContainerOptions struct {
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *DeleteContainerOptions) pointers() (*ContainerDeleteOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	return nil, o.LeaseAccessConditions, o.ModifiedAccessConditions
}

// GetPropertiesContainerOptions provides set of configurations for GetPropertiesContainer operation
type GetPropertiesContainerOptions struct {
	ContainerGetPropertiesOptions *ContainerGetPropertiesOptions
	LeaseAccessConditions         *LeaseAccessConditions
}

func (o *GetPropertiesContainerOptions) pointers() (*ContainerGetPropertiesOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerGetPropertiesOptions, o.LeaseAccessConditions
}

// GetAccessPolicyOptions provides set of configurations for GetAccessPolicy operation
type GetAccessPolicyOptions struct {
	ContainerGetAccessPolicyOptions *ContainerGetAccessPolicyOptions
	LeaseAccessConditions           *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) pointers() (*ContainerGetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerGetAccessPolicyOptions, o.LeaseAccessConditions
}

// SetAccessPolicyOptions provides set of configurations for SetAccessPolicy operation
type SetAccessPolicyOptions struct {
	// At least Access and ContainerACL must be specified
	ContainerSetAccessPolicyOptions ContainerSetAccessPolicyOptions
	AccessConditions                *ContainerAccessConditions
}

func (o *SetAccessPolicyOptions) pointers() (ContainerSetAccessPolicyOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return ContainerSetAccessPolicyOptions{}, nil, nil
	}
	mac, lac := o.AccessConditions.pointers()
	return o.ContainerSetAccessPolicyOptions, lac, mac
}

// SetMetadataContainerOptions provides set of configurations for SetMetadataContainer operation
type SetMetadataContainerOptions struct {
	Metadata                 map[string]string
	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *SetMetadataContainerOptions) pointers() (*ContainerSetMetadataOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	options := ContainerSetMetadataOptions{Metadata: o.Metadata}
	return &options, o.LeaseAccessConditions, o.ModifiedAccessConditions
}
