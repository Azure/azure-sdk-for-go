//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/internal/generated"
)

// SetExpiryType defines values for ExpiryType
type SetExpiryType interface {
	Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.PathClientSetExpiryOptions)
	notPubliclyImplementable()
}

// SetExpiryTypeAbsolute defines the absolute time for the blob expiry
type SetExpiryTypeAbsolute time.Time

// SetExpiryTypeRelativeToNow defines the duration relative to now for the blob expiry
type SetExpiryTypeRelativeToNow time.Duration

// SetExpiryTypeRelativeToCreation defines the duration relative to creation for the blob expiry
type SetExpiryTypeRelativeToCreation time.Duration

// SetExpiryTypeNever defines that the blob will be set to never expire
type SetExpiryTypeNever struct {
	// empty struct since NeverExpire expiry type does not require expiry time
}

// SetExpiryOptions contains the optional parameters for the Client.SetExpiry method.
type SetExpiryOptions struct {
	// placeholder for future options
}

func (e SetExpiryTypeAbsolute) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.PathClientSetExpiryOptions) {
	return generated.ExpiryOptionsAbsolute, &generated.PathClientSetExpiryOptions{
		ExpiresOn: to.Ptr(time.Time(e).UTC().Format(http.TimeFormat)),
	}
}

func (e SetExpiryTypeAbsolute) notPubliclyImplementable() {}

func (e SetExpiryTypeRelativeToNow) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.PathClientSetExpiryOptions) {
	return generated.ExpiryOptionsRelativeToNow, &generated.PathClientSetExpiryOptions{
		ExpiresOn: to.Ptr(strconv.FormatInt(time.Duration(e).Milliseconds(), 10)),
	}
}

func (e SetExpiryTypeRelativeToNow) notPubliclyImplementable() {}

func (e SetExpiryTypeRelativeToCreation) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.PathClientSetExpiryOptions) {
	return generated.ExpiryOptionsRelativeToCreation, &generated.PathClientSetExpiryOptions{
		ExpiresOn: to.Ptr(strconv.FormatInt(time.Duration(e).Milliseconds(), 10)),
	}
}

func (e SetExpiryTypeRelativeToCreation) notPubliclyImplementable() {}

func (e SetExpiryTypeNever) Format(o *SetExpiryOptions) (generated.ExpiryOptions, *generated.PathClientSetExpiryOptions) {
	return generated.ExpiryOptionsNeverExpire, &generated.PathClientSetExpiryOptions{}
}

func (e SetExpiryTypeNever) notPubliclyImplementable() {}
