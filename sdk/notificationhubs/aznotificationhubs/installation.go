//go:build go1.20
// +build go1.20

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznotificationhubs

type (
	// Represents an installation in Azure Notification Hubs.
	Installation struct {
		InstallationId     string                           `json:"installationId"`               // The installation ID
		UserId             string                           `json:"userId,omitempty"`             // The user ID
		LastActiveOn       string                           `json:"lastActiveOn,omitempty"`       // The date the installation was last active
		ExpirationTime     string                           `json:"expirationTime,omitempty"`     // The installation expiration time
		LastUpdate         string                           `json:"lastUpdate,omitempty"`         // The date the installation was last updated
		Platform           string                           `json:"platform"`                     // The platform of the installation
		ExpiredPushChannel bool                             `json:"expiredPushChannel,omitempty"` // Whether the push channel is expired
		Tags               []string                         `json:"tags,omitempty"`               // The tags for the installation
		Templates          map[string]*InstallationTemplate `json:"templates,omitempty"`          // The templates for the installation
		PushChannel        interface{}                      `json:"pushChannel"`                  // The push channel for the installation whcih can be a string or for web push, a map[string]interface{}
	}

	// Represents a template for an installation in Azure Notification Hubs.
	InstallationTemplate struct {
		Body string `json:"body"` // The body of the template
	}

	// Represents a patch operation for an installation in Azure Notification Hubs.
	InstallationPatch struct {
		Op    string `json:"op"`              // The operation to perform
		Path  string `json:"path"`            // The path to the property to patch
		Value string `json:"value,omitempty"` // The value to set the property to
	}
)
