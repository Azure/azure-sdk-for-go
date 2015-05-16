// Package management provides the main API client to construct other clients
// and make requests to the Microsoft Azure Service Management REST API.
package management

import (
	"errors"
	"time"
)

const (
	DefaultAzureManagementURL    = "https://management.core.windows.net"
	DefaultOperationPollInterval = time.Second * 30
	DefaultAPIVersion            = "2014-10-01"
	DefaultUserAgent             = "azure-sdk-for-go"

	errPublishSettingsConfiguration       = "PublishSettingsFilePath is set. Consequently ManagementCertificatePath and SubscriptionId must not be set."
	errManagementCertificateConfiguration = "Both ManagementCertificatePath and SubscriptionId should be set, and PublishSettingsFilePath must not be set."
	errParamNotSpecified                  = "Parameter %s is not specified."
)

// Client provides a client to the Azure API.
type Client struct {
	publishSettings publishSettings
	config          ClientConfig
}

// ClientConfig provides a configuration for use by a Client
type ClientConfig struct {
	ManagementURL         string
	OperationPollInterval time.Duration
	UserAgent             string
	APIVersion            string
}

// NewAnonymousClient creates a new azure.Client with no credentials set.
func NewAnonymousClient() Client {
	return Client{}
}

// DefaultConfig returns the default client configuration used to construct
// a client. This value can be used to make modifications on the default API
// configuration.
func DefaultConfig() ClientConfig {
	return ClientConfig{
		ManagementURL:         DefaultAzureManagementURL,
		OperationPollInterval: DefaultOperationPollInterval,
		APIVersion:            DefaultAPIVersion,
		UserAgent:             DefaultUserAgent,
	}
}

// NewClient creates a new Client using the given subscription ID and
// management certificate
func NewClient(subscriptionID string, managementCert []byte) (Client, error) {
	return NewClientFromConfig(subscriptionID, managementCert, DefaultConfig())
}

// NewClientFromConfig creates a new Client using a given ClientConfig
func NewClientFromConfig(subscriptionID string, managementCert []byte, config ClientConfig) (Client, error) {
	return makeClient(subscriptionID, managementCert, config)
}

func makeClient(subscriptionID string, managementCert []byte, config ClientConfig) (Client, error) {
	var client Client

	if subscriptionID == "" {
		return client, errors.New("azure: subscription ID required")
	} else if len(managementCert) == 0 {
		return client, errors.New("azure: management certificate required")
	}

	publishSettings := publishSettings{
		SubscriptionID:   subscriptionID,
		SubscriptionCert: managementCert,
		SubscriptionKey:  managementCert,
	}

	// Validate client configuration
	if config.ManagementURL == "" {
		return client, errors.New("azure: base URL required")
	} else if config.OperationPollInterval <= 0 {
		return client, errors.New("azure: operation polling interval must be a positive duration")
	} else if config.APIVersion == "" {
		return client, errors.New("azure: client configuration must specify an API version")
	} else if config.UserAgent == "" {
		config.UserAgent = DefaultUserAgent
	}

	return Client{
		publishSettings: publishSettings,
		config:          config,
	}, nil
}
