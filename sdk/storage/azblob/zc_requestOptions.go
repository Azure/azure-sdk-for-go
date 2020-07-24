package azblob

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
		Metadata: nil, // TODO
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

type DeleteBlobOptions struct {
	// Required if the blob has associated snapshots. Specify one of the following two options: include: Delete the base blob
	// and all of its snapshots. only: Delete only the blob's snapshots and not the blob itself
	DeleteSnapshots *DeleteSnapshotsOptionType

	LeaseAccessConditions    *LeaseAccessConditions
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *DeleteBlobOptions) pointers() (*BlobDeleteOptions, *LeaseAccessConditions, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil
	}

	basics := BlobDeleteOptions{
		DeleteSnapshots: o.DeleteSnapshots,
	}

	return &basics, o.LeaseAccessConditions, o.ModifiedAccessConditions
}

type UploadBlockBlobOptions struct {
	// Optional. Used to set blob tags in various blob operations.
	BlobTagsString *string

	// Optional. Specifies a user-defined name-value pair associated with the blob.
	Metadata map[string]string

	// Optional. Indicates the tier to be set on the blob.
	Tier *AccessTier

	// Specify the transactional md5 for the body, to be validated by the service.
	TransactionalContentMd5 *[]byte

	BlobHttpHeaders          *BlobHttpHeaders
	LeaseAccessConditions    *LeaseAccessConditions
	CpkInfo                  *CpkInfo
	CpkScopeInfo             *CpkScopeInfo
	ModifiedAccessConditions *ModifiedAccessConditions
}

func (o *UploadBlockBlobOptions) pointers() (*BlockBlobUploadOptions, *BlobHttpHeaders, *LeaseAccessConditions,
	*CpkInfo, *CpkScopeInfo, *ModifiedAccessConditions) {
	if o == nil {
		return nil, nil, nil, nil, nil, nil
	}

	basics := BlockBlobUploadOptions{
		BlobTagsString:          o.BlobTagsString,
		Metadata:                nil, // TODO
		Tier:                    o.Tier,
		TransactionalContentMd5: o.TransactionalContentMd5,
	}
	return &basics, o.BlobHttpHeaders, o.LeaseAccessConditions, o.CpkInfo, o.CpkScopeInfo, o.ModifiedAccessConditions
}
