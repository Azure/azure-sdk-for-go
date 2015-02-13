package azure

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

//AzureError represents an error returned by the management API. It has an error
//code (for example, ResourceNotFound) and a descriptive message.
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
	publishSettings publishSettings
	baseURL         string
}

//NewAnonymouseClient creates a new azure.Client with no credentials set.
func NewAnonymousClient() Client {
	return Client{}
}

// NewBasicClient creates a new azure.Client, imports the publish
// settings from the specified file path and uses the default API endpoint.
func NewBasicClient(publishSettingsFilePath string) (Client, error) {
	return NewClientFromPublishSettingsFile(publishSettingsFilePath, defaultAzureManagementURL)
}

// NewClient creates a new azure.Client using the given subscription ID,
// management certificate and custom API endpoint
func NewClient(subscriptionID string, managementCertificatePath string, baseURL string) (Client, error) {
	return NewClientFromPublishSettings(subscriptionID, managementCertificatePath, baseURL)
}

// NewClientFromPublishSettingsFile creates a new azure.Client using the
// given publishsettings file path and custom API endpoint
func NewClientFromPublishSettingsFile(publishSettingsFilePath string, baseURL string) (Client, error) {
	var client Client
	if publishSettingsFilePath == "" {
		return client, errors.New("azure: publishsettings file path required")
	} else if baseURL == "" {
		return client, errors.New("azure: base URL required")
	}

	err := client.importPublishSettingsFile(publishSettingsFilePath)
	if err != nil {
		return client, err
	}

	return Client{
		baseURL: baseURL,
	}, nil
}

// NewClientFromPublishSettings creates a new azure.Client using the
// given subscription ID, management certificate path and custom API endpoint
func NewClientFromPublishSettings(subscriptionID string, managementCertificatePath string, baseURL string) (Client, error) {
	var client Client
	if subscriptionID == "" {
		return client, errors.New("azure: subscription ID required")
	} else if managementCertificatePath == "" {
		return client, errors.New("azure: management certificate path required")
	} else if baseURL == "" {
		return client, errors.New("azure: base URL required")
	}

	err := client.importPublishSettings(subscriptionID, managementCertificatePath)
	if err != nil {
		return client, err
	}

	return Client{
		baseURL: baseURL,
	}, nil
}
