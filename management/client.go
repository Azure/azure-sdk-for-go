package management

import (
	"errors"
	"time"
)

const (
	defaultAzureManagementURL             = "https://management.core.windows.net"
	defaultOperationPollInterval          = time.Second * 30
	errPublishSettingsConfiguration       = "PublishSettingsFilePath is set. Consequently ManagementCertificatePath and SubscriptionId must not be set."
	errManagementCertificateConfiguration = "Both ManagementCertificatePath and SubscriptionId should be set, and PublishSettingsFilePath must not be set."
	errParamNotSpecified                  = "Parameter %s is not specified."
)

// Client provides a client to the Azure API.
type Client struct {
	managementURL   string
	publishSettings publishSettings
	pollInterval    time.Duration
}

// ClientConfig provides a configuration for use by a Client
type ClientConfig struct {
	ManagementURL         string
	OperationPollInterval time.Duration
}

// NewAnonymousClient creates a new azure.Client with no credentials set.
func NewAnonymousClient() Client {
	return Client{}
}

func defaultConfig() ClientConfig {
	return ClientConfig{
		ManagementURL:         defaultAzureManagementURL,
		OperationPollInterval: defaultOperationPollInterval,
	}
}

// NewClient creates a new Client using the given subscription ID and
// management certificate
func NewClient(subscriptionID string, managementCert []byte) (Client, error) {
	return NewClientFromConfig(subscriptionID, managementCert, defaultConfig())
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
	} else if config.ManagementURL == "" {
		return client, errors.New("azure: base URL required")
	} else if config.OperationPollInterval <= 0 {
		return client, errors.New("azure: operation polling interval must be a positive duration")
	}

	publishSettings := publishSettings{
		SubscriptionID:   subscriptionID,
		SubscriptionCert: managementCert,
		SubscriptionKey:  managementCert,
	}

	return Client{
		managementURL:   config.ManagementURL,
		publishSettings: publishSettings,
		pollInterval:    config.OperationPollInterval,
	}, nil
}
