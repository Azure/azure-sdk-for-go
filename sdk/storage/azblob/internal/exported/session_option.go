package exported

// SessionMode specifies how session-based authentication is handled.
type SessionMode string

const ( // SessionModeDefault is the default mode where sessions are disabled.
	SessionModeDefault SessionMode = SessionModeOff // TODO : I dont think this is idiomatic in Go. Should this be ""?
	// SessionModeOff explicitly disables session-based authentication.
	SessionModeOff SessionMode = "off"
	// SessionModeSingleContainer enables session-based authentication for a single container.
	SessionModeSingleContainer SessionMode = "singlecontainer"
)

// PossibleSessionModeValues returns a slice of possible values for SessionMode.
func PossibleSessionModeValues() []SessionMode {
	return []SessionMode{
		SessionModeDefault,
		SessionModeOff,
		SessionModeSingleContainer,
	}
}

// SessionOptions configures session-based authentication behavior.
type SessionOptions struct {
	// Mode specifies the session authentication mode.
	Mode SessionMode

	// AccountName is the storage account name.
	AccountName string
	// ContainerName is the container name for SingleContainer mode.
	ContainerName string
}
