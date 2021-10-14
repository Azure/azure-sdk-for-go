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
	DeleteSnapshots *string
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *DeleteShareOptions) format() *ShareDeleteOptions {
	if o == nil {
		return nil
	}

	return &ShareDeleteOptions{
		DeleteSnapshots: o.DeleteSnapshots,
		Sharesnapshot:   o.ShareSnapshot,
	}
}

type GetSharePropertiesOptions struct {
	// The snapshot parameter is an opaque DateTime value that, when present, specifies the share snapshot to query.
	ShareSnapshot *string
}

func (o *GetSharePropertiesOptions) format() *ShareGetPropertiesOptions {
	if o == nil {
		return nil
	}

	return &ShareGetPropertiesOptions{
		Sharesnapshot: o.ShareSnapshot,
	}
}

type SetSharePropertiesOptions struct {
	// Specifies the access tier of the share.
	AccessTier *ShareAccessTier
	// Specifies the maximum size of the share, in gigabytes.
	Quota *int32
}

func (o *SetSharePropertiesOptions) format() *ShareSetPropertiesOptions {
	if o == nil {
		return nil
	}

	return &ShareSetPropertiesOptions{
		AccessTier: o.AccessTier,
		Quota:      o.Quota,
	}
}

type SetShareMetadataOptions struct {
	// A name-value pair to associate with a file storage object.
	Metadata map[string]string
}

func (o *SetShareMetadataOptions) format() *ShareSetMetadataOptions {
	if o == nil {
		return nil
	}

	return &ShareSetMetadataOptions{
		Metadata: o.Metadata,
	}
}

type GetShareAccessPolicyOptions struct {
}

type CreateSharePermissionOptions struct {
}

type GetSharePermissionOptions struct {
}

// SetShareAccessPolicyOptions contains the optional parameters for the Share.SetAccessPolicy method.
type SetShareAccessPolicyOptions struct {
	// The ACL for the share.
	ShareACL []*SignedIdentifier
}

func (o *SetShareAccessPolicyOptions) format() *ShareSetAccessPolicyOptions {
	if o == nil {
		return nil
	}

	return &ShareSetAccessPolicyOptions{
		ShareACL: o.ShareACL,
	}
}

type GetShareStatisticsOptions struct {
}
