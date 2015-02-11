package azure

import (
	"encoding/xml"
	"fmt"
)

const (
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
}

//NewAnonymouseClient creates a new azure.Client with no credentials set.
func NewAnonymousClient() *Client {
	return &Client{}
}

//NewClientFromPublishSettingsFile creates a new azure.Client and imports the publish
//settings from the specified file path.
func NewClientFromPublishSettingsFile(publishSettingsFilePath string) (*Client, error) {
	client := &Client{}
	err := client.importPublishSettingsFile(publishSettingsFilePath)
	if err != nil {
		return nil, err
	}
	return client, nil
}

//NewClientFromPublishSettingsFile creates a new azure.Client and imports the publish
//settings from the specified file path.
func NewClientFromPublishSettings(subscriptionId string, managementCertificatePath string) (*Client, error) {
	client := &Client{}
	err := client.importPublishSettings(subscriptionId, managementCertificatePath)
	if err != nil {
		return nil, err
	}
	return client, nil
}
