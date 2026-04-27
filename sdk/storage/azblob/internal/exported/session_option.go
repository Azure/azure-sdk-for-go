// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

// SessionMode specifies how session-based authentication is handled.
type SessionMode string

const (
	// SessionModeDefault is the default mode where sessions are disabled.
	SessionModeDefault SessionMode = ""
	// SessionModeOff explicitly disables session-based authentication.
	SessionModeOff SessionMode = "off"
	// SessionModeSingleSpecifiedContainer enables session-based authentication for a single container.
	SessionModeSingleSpecifiedContainer SessionMode = "singlespecifiedcontainer"
)

// PossibleSessionModeValues returns a slice of possible values for SessionMode.
func PossibleSessionModeValues() []SessionMode {
	return []SessionMode{
		SessionModeDefault,
		SessionModeOff,
		SessionModeSingleSpecifiedContainer,
	}
}

// SessionOptions configures session-based authentication behavior.
type SessionOptions struct {
	// Mode specifies the session authentication mode.
	Mode SessionMode

	// AccountName is the storage account name.
	AccountName string
	// ContainerName is the container name for SingleSpecifiedContainer mode.
	ContainerName string
}
