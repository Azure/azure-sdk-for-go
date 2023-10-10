//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// AddSettingOptions contains the optional parameters for the AddSetting method.
type AddSettingOptions struct {
	// Configuration setting content type.
	ContentType *string

	// Configuration setting label.
	Label *string
}

// DeleteSettingOptions contains the optional parameters for the DeleteSetting method.
type DeleteSettingOptions struct {
	// Configuration setting label.
	Label *string

	// If set, and the configuration setting exists in the configuration store,
	// delete the setting if the passed-in ETag is the same as the setting's ETag in the configuration store.
	OnlyIfUnchanged *azcore.ETag
}

// GetSettingOptions contains the optional parameters for the GetSetting method.
type GetSettingOptions struct {
	// The setting will be retrieved exactly as it existed at the provided time.
	AcceptDateTime *time.Time

	// Configuration setting label.
	Label *string

	// If set, only retrieve the setting from the configuration store if setting has changed
	// since the client last retrieved it with the ETag provided.
	OnlyIfChanged *azcore.ETag
}

// ListRevisionsOptions contains the optional parameters for the NewListRevisionsPager method.
type ListRevisionsOptions struct {
	// placeholder for future options
}

// ListSettingsOptions contains the optional parameters for the NewListSettingsPager method.
type ListSettingsOptions struct {
	// placeholder for future options
}

// SetReadOnlyOptions contains the optional parameters for the SetReadOnly method.
type SetReadOnlyOptions struct {
	// Configuration setting label.
	Label *string

	// If set, and the configuration setting exists in the configuration store, update the setting
	// if the passed-in configuration setting ETag is the same version as the one in the configuration store.
	OnlyIfUnchanged *azcore.ETag
}

// SetSettingOptions contains the optional parameters for the SetSetting method.
type SetSettingOptions struct {
	// Configuration setting content type.
	ContentType *string

	// Configuration setting label.
	Label *string

	// If set, and the configuration setting exists in the configuration store, overwrite the setting
	// if the passed-in ETag is the same version as the one in the configuration store.
	OnlyIfUnchanged *azcore.ETag
}
