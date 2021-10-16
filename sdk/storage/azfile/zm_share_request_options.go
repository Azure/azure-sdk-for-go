package azfile

type CreateShareOptions struct {
	// Specifies the access tier of the share.
	AccessTier *ShareAccessTier
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string
	// Specifies the maximum size of the share, in gigabytes.
	Quota *int32
}

func (o *CreateShareOptions) format() *ShareCreateOptions {
	if o == nil {
		return nil
	}

	return &ShareCreateOptions{
		AccessTier: o.AccessTier,
		Metadata:   o.Metadata,
		Quota:      o.Quota,
	}
}

type CreateShareSnapshotOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string
}

func (o *CreateShareSnapshotOptions) format() *ShareCreateSnapshotOptions {
	if o == nil {
		return nil
	}
	return &ShareCreateSnapshotOptions{Metadata: o.Metadata}
}

type DeleteShareOptions struct {
	// Specifies the option include to delete the base share and all of its snapshots.
	DeleteSnapshots *DeleteSnapshotsOptionType
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *DeleteShareOptions) format() (*ShareDeleteOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &ShareDeleteOptions{
		DeleteSnapshots: o.DeleteSnapshots,
		Sharesnapshot:   o.ShareSnapshot,
	}, o.LeaseAccessConditions
}

type GetSharePropertiesOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot         *string
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetSharePropertiesOptions) format() (*ShareGetPropertiesOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &ShareGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}, o.LeaseAccessConditions
}

type SetSharePropertiesOptions struct {
	// Specifies the access tier of the share.
	AccessTier *ShareAccessTier
	// Specifies the maximum size of the share, in gigabytes.
	Quota                 *int32
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetSharePropertiesOptions) format() (*ShareSetPropertiesOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &ShareSetPropertiesOptions{
		AccessTier: o.AccessTier,
		Quota:      o.Quota,
	}, o.LeaseAccessConditions
}

type SetShareMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string

	leaseAccessConditions *LeaseAccessConditions
}

func (o *SetShareMetadataOptions) format() (*ShareSetMetadataOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &ShareSetMetadataOptions{
		Metadata: o.Metadata,
	}, o.leaseAccessConditions
}

type GetShareAccessPolicyOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetShareAccessPolicyOptions) format() *LeaseAccessConditions {
	if o == nil {
		return nil
	}

	return o.LeaseAccessConditions
}

type CreateSharePermissionOptions struct {
}

type GetSharePermissionOptions struct {
}

// SetShareAccessPolicyOptions contains the optional parameters for the Share.SetAccessPolicy method.
type SetShareAccessPolicyOptions struct {
	// The ACL for the share.
	ShareACL []*SignedIdentifier

	LeaseAccessConditions *LeaseAccessConditions
}

func (o *SetShareAccessPolicyOptions) format() (*ShareSetAccessPolicyOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return &ShareSetAccessPolicyOptions{
		ShareACL: o.ShareACL,
	}, o.LeaseAccessConditions
}

type GetShareStatisticsOptions struct {
	LeaseAccessConditions *LeaseAccessConditions
}

func (o *GetShareStatisticsOptions) format() (*ShareGetStatisticsOptions, *LeaseAccessConditions) {
	if o == nil {
		return nil, nil
	}

	return nil, o.LeaseAccessConditions
}
