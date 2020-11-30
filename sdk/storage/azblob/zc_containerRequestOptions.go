package azblob

type CreateContainerOptions struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	Access *PublicAccessType

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata *map[string]string

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

type GetPropertiesOptionsContainer struct {
	ContainerGetPropertiesOptions *ContainerGetPropertiesOptions
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetPropertiesOptionsContainer) pointers() (*ContainerGetPropertiesOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerGetPropertiesOptions, o.LeaseAccessConditions
}

type GetAccessPolicyOptions struct {
	ContainerGetAccessPolicyOptions *ContainerGetAccessPolicyOptions
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetAccessPolicyOptions) pointers() (*ContainerGetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerGetAccessPolicyOptions, o.LeaseAccessConditions
}

type SetAccessPolicyOptions struct {
	// At least Access and ContainerAcl must be specified
	ContainerSetAccessPolicyOptions ContainerSetAccessPolicyOptions
	ContainerAccessConditions *ContainerAccessConditions
}

type AcquireLeaseOptionsContainer struct {
	// At least Access and ContainerAcl must be specified
	ContainerSetAccessPolicyOptions *ContainerAcquireLeaseOptions
	ModifiedAccessConditions *ModifiedAccessConditions
}

type RenewLeaseOptionsContainer struct {
	ContainerRenewLeaseOptions *ContainerRenewLeaseOptions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *RenewLeaseOptionsContainer) pointers() (*ContainerRenewLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerRenewLeaseOptions, o.ModifiedAccessConditions
}

type ReleaseLeaseOptionsContainer struct {
	ContainerReleaseLeaseOptions *ContainerReleaseLeaseOptions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ReleaseLeaseOptionsContainer) pointers() (*ContainerReleaseLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerReleaseLeaseOptions, o.ModifiedAccessConditions
}

type BreakLeaseOptionsContainer struct {
	ContainerBreakLeaseOptions *ContainerBreakLeaseOptions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *BreakLeaseOptionsContainer) pointers() (*ContainerBreakLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerBreakLeaseOptions, o.ModifiedAccessConditions
}

type ChangeLeaseOptionsContainer struct {
	ContainerChangeLeaseOptions *ContainerChangeLeaseOptions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *ChangeLeaseOptionsContainer) pointers() (*ContainerChangeLeaseOptions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return o.ContainerChangeLeaseOptions, o.ModifiedAccessConditions
}