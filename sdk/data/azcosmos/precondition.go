// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// Precondition is an immutable ETag condition for an item operation.
// Its zero value adds no condition. Use IfMatch or IfNoneMatch to construct a value.
// Service support for each condition depends on the operation.
type Precondition struct {
	kind preconditionKind
	etag azcore.ETag
}

// IfMatch constructs a condition requiring the item to match etag.
// Empty ETags and ETags containing NUL are rejected. The ETag is passed unchanged.
func IfMatch(etag azcore.ETag) (Precondition, error) {
	if err := validateETag("IfMatch", &etag); err != nil {
		return Precondition{}, err
	}
	return Precondition{kind: preconditionKindIfMatch, etag: etag}, nil
}

// IfNoneMatch constructs a condition requiring the item not to match etag.
// Empty ETags and ETags containing NUL are rejected. The ETag is passed unchanged.
func IfNoneMatch(etag azcore.ETag) (Precondition, error) {
	if err := validateETag("IfNoneMatch", &etag); err != nil {
		return Precondition{}, err
	}
	return Precondition{kind: preconditionKindIfNoneMatch, etag: etag}, nil
}

func (p Precondition) apply(req *itemRequest) {
	req.preconditionKind = p.kind
	req.preconditionETag = string(p.etag)
}

func validateETag(name string, etag *azcore.ETag) error {
	if etag == nil {
		return nil
	}
	if *etag == "" {
		return fmt.Errorf("azcosmos: %s must not be empty", name)
	}
	if strings.IndexByte(string(*etag), 0) >= 0 {
		return fmt.Errorf("azcosmos: %s must not contain a NUL byte", name)
	}
	return nil
}
