package management

import (
	"encoding/xml"
	"errors"
	"fmt"
)

const (
	defaultAzureManagementURL             = "https://management.core.windows.net"
	errPublishSettingsConfiguration       = "PublishSettingsFilePath is set. Consequently ManagementCertificatePath and SubscriptionId must not be set."
	errManagementCertificateConfiguration = "Both ManagementCertificatePath and SubscriptionId should be set, and PublishSettingsFilePath must not be set."
	errParamNotSpecified                  = "Parameter %s is not specified."
)

// AzureError represents an error returned by the management API. It has an error
// code (for example, ResourceNotFound) and a descriptive message.
type AzureError struct {
	XMLName xml.Name `xml:"Error"`
	Code    string
	Message string
}

//Error implements the error interface for the AzureError type.
func (e *AzureError) Error() string {
	return fmt.Sprintf("Error response from Azure. Code: %s, Message: %s", e.Code, e.Message)
}

// Client provides a client to the Azure API.
type Client struct {
	managementURL   string
	publishSettings publishSettings
}

// ClientConfig provides a configuration for use by a Client
type ClientConfig struct {
	ManagementURL string
}

// NewAnonymousClient creates a new azure.Client with no credentials set.
func NewAnonymousClient() Client {
	return Client{}
}

// NewClient creates a new Client using the given subscription ID and
// management certificate
func NewClient(subscriptionID string, managementCert []byte) (Client, error) {
	return makeClient(subscriptionID, managementCert, defaultAzureManagementURL)
}

// NewClientFromConfig creates a new Client using a given ClientConfig
func NewClientFromConfig(subscriptionID string, managementCert []byte, config ClientConfig) (Client, error) {
	return makeClient(subscriptionID, managementCert, config.ManagementURL)
}

func makeClient(subscriptionID string, managementCert []byte, managementURL string) (Client, error) {
	var client Client
	if subscriptionID == "" {
		return client, errors.New("azure: subscription ID required")
	} else if len(managementCert) == 0 {
		return client, errors.New("azure: management certificate required")
	} else if managementURL == "" {
		return client, errors.New("azure: base URL required")
	}

	publishSettings := publishSettings{
		SubscriptionID:   subscriptionID,
		SubscriptionCert: managementCert,
		SubscriptionKey:  managementCert,
	}

	return Client{
		managementURL:   managementURL,
		publishSettings: publishSettings,
	}, nil
}
