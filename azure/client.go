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

// Config is used to configure the creation of a client
type Config struct {
	//ManagementCertificatePath is the path to a file containing the management
	//certificate for the subscription on which this client will operate. If using
	//this method of authentication, SubscriptionId must also be set, and
	//PublishSettingsFilePath must not be set.
	ManagementCertificatePath string

	//SubscriptionId is the name of the subscription on which the subscription will
	//operate. If using this method of authentication, ManagementCertificatePath must
	//also be set, and PublishSettingsFilePath must not be set.
	SubscriptionId string

	//PublishSettingsFilePath is the path to an Azure PublishSettings file containing
	//the management settings for a subscription. If using this method of authentication,
	//neither ManagementCertificatePath or SubscriptionId may be set.
	PublishSettingsFilePath string

	//Anonymous flags whether any authentication options will be specified for the
	//client. If Anonymous is set to true, this will override any other authentication
	//options specified.
	Anonymous bool
}

//AzureError represents an error returned by the management API. It has an error
//code (for example, ResourceNotFound) and a descriptive message.
type AzureError struct {
	XMLName xml.Name `xml:"Error"`
	Code    string
	Message string
}

//Error implements the error interface for the AzureError type.
func (e *AzureError) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s", e.Code, e.Message)
}

// Client provides a client to the Azure API.
type Client struct {
	publishSettings publishSettings
}

// NewClient returns a new client.
func NewClient(config *Config) (*Client, error) {
	client := &Client{}

	//If the client is anonymous there's no need to load auth options.
	if config.Anonymous {
		return client, nil
	}

	//Set the publish settings for the client according to the configuration
	if len(config.PublishSettingsFilePath) > 0 {
		if len(config.ManagementCertificatePath) > 0 || len(config.SubscriptionId) > 0 {
			return nil, fmt.Errorf(errPublishSettingsConfiguration)
		}
		err := client.importPublishSettingsFile(config.PublishSettingsFilePath)
		if err != nil {
			return nil, err
		}
	} else {
		if len(config.ManagementCertificatePath) == 0 || len(config.SubscriptionId) == 0 {
			return nil, fmt.Errorf(errManagementCertificateConfiguration)
		}
		err := client.importPublishSettings(config.SubscriptionId, config.ManagementCertificatePath)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}
