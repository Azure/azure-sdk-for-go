package management

import (
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/Azure/go-pkcs12"
)

func (client *Client) importPublishSettings(id string, certPath string) error {
	if id == "" {
		return fmt.Errorf(errParamNotSpecified, "id")
	}
	if certPath == "" {
		return fmt.Errorf(errParamNotSpecified, "certPath")
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}

	client.publishSettings.SubscriptionID = id
	client.publishSettings.SubscriptionCert = cert
	client.publishSettings.SubscriptionKey = cert

	return nil
}

func (client *Client) importPublishSettingsFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf(errParamNotSpecified, "filePath")
	}

	publishSettingsContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	activeSubscription, err := getActiveSubscription(publishSettingsContent)
	if err != nil {
		return err
	}

	cert, err := getSubscriptionCert(activeSubscription)
	if err != nil {
		return err
	}

	client.publishSettings.SubscriptionID = activeSubscription.Id
	client.publishSettings.SubscriptionCert = cert
	client.publishSettings.SubscriptionKey = cert
	return nil
}

func getSubscriptionCert(subscription subscription) ([]byte, error) {
	certPassword := ""

	pfxCert, err := base64.StdEncoding.DecodeString(subscription.ManagementCertificate)
	if err != nil {
		return nil, err
	}

	pemBlocks, err := pkcs12.ConvertToPEM(pfxCert, certPassword)
	if err != nil {
		return nil, err
	}

	var certData []byte
	for _, block := range pemBlocks {
		certData = append(certData, pem.EncodeToMemory(block)...)
	}

	return certData, nil
}

func getActiveSubscription(publishSettingsContent []byte) (subscription, error) {
	publishData := publishData{}
	activeSubscription := subscription{}

	err := xml.Unmarshal(publishSettingsContent, &publishData)
	if err != nil {
		return activeSubscription, err
	}

	if len(publishData.PublishProfiles) == 0 {
		err = errors.New("No publish profiles were found")
		return activeSubscription, err
	}

	publishProfile := publishData.PublishProfiles[0]
	subscriptions := publishProfile.Subscriptions
	if len(subscriptions) == 0 {
		err = errors.New("No subscriptions were found")
		return activeSubscription, err
	}

	activeSubscription = subscriptions[0]

	if activeSubscription.ManagementCertificate == "" {
		activeSubscription.ManagementCertificate = publishProfile.ManagementCertificate
		activeSubscription.ServiceManagementUrl = publishProfile.Url
	}

	return activeSubscription, nil
}

type publishSettings struct {
	XMLName          xml.Name `xml:"PublishSettings"`
	SubscriptionID   string
	SubscriptionCert []byte
	SubscriptionKey  []byte
}

type publishData struct {
	XMLName         xml.Name         `xml:"PublishData"`
	PublishProfiles []publishProfile `xml:"PublishProfile"`
}

type publishProfile struct {
	XMLName               xml.Name       `xml:"PublishProfile"`
	SchemaVersion         string         `xml:",attr"`
	PublishMethod         string         `xml:",attr"`
	Url                   string         `xml:",attr"`
	ManagementCertificate string         `xml:",attr"`
	Subscriptions         []subscription `xml:"Subscription"`
}

type subscription struct {
	XMLName               xml.Name `xml:"Subscription"`
	ServiceManagementUrl  string   `xml:",attr"`
	Id                    string   `xml:",attr"`
	Name                  string   `xml:",attr"`
	ManagementCertificate string   `xml:",attr"`
}
