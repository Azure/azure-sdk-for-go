package azblob

const (
	// ETagNone represents an empty entity tag.
	ETagNone = ""

	// ETagAny matches any entity tag.
	ETagAny = "*"
)

// ContainerAccessConditions identifies container-specific access conditions which you optionally set.
type ContainerAccessConditions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
	LeaseAccessConditions    *LeaseAccessConditions
}

func (ac *ContainerAccessConditions) pointers() (*ModifiedAccessConditions, *LeaseAccessConditions) {
	if ac == nil {
		return nil, nil
	}

	return ac.ModifiedAccessConditions, ac.LeaseAccessConditions
}

// BlobAccessConditions identifies blob-specific access conditions which you optionally set.
type BlobAccessConditions struct {
	ModifiedAccessConditions *ModifiedAccessConditions
	LeaseAccessConditions    *LeaseAccessConditions
}

func (ac *BlobAccessConditions) pointers() (*ModifiedAccessConditions, *LeaseAccessConditions) {
	if ac == nil {
		return nil, nil
	}

	return ac.ModifiedAccessConditions, ac.LeaseAccessConditions
}
