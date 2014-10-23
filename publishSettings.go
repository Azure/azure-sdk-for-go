package azureSdkForGo

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
)

var settings publishSettings = publishSettings{}

func GetPublishSettings() publishSettings {
	return settings
}

func setPublishSettings(id string, cert []byte, key []byte) {
	settings.SubscriptionID = id
	settings.SubscriptionCert = cert
	settings.SubscriptionKey = key
}

func ImportPublishSettings(id string, certPath string) error {
	if len(id) == 0 {
		return fmt.Errorf(ParamNotSpecifiedError, "id")
	}
	if len(certPath) == 0 {
		return fmt.Errorf(ParamNotSpecifiedError, "certPath")
	}

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return err
	}

	setPublishSettings(id, cert, cert)
	return nil
}

func ImportPublishSettingsFile(filePath string) error {
	if len(filePath) == 0 {
		return fmt.Errorf(ParamNotSpecifiedError, "filePath")
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

	setPublishSettings(activeSubscription.Id, cert, cert)
	return nil
}

func getSubscriptionCert(subscription subscription) ([]byte, error) {
	certPassword := ""

	azureDir, err := getAzureDir()
	if err != nil {
		return nil, err
	}

	pfxCertPath := path.Join(azureDir, "cert.pfx")
	pemCertPath := path.Join(azureDir, "cert.pem")

	err = createPfxCert(subscription.ManagementCertificate, pfxCertPath)
	if err != nil {
		return nil, err
	}

	ExecuteCommand(fmt.Sprintf("openssl pkcs12 -in %s -out %s -nodes -passin pass:%s", pfxCertPath, pemCertPath, certPassword))

	pemCert, readErr := ioutil.ReadFile(pemCertPath)
	if readErr != nil {
		return nil, readErr
	}

	os.Remove(pfxCertPath)
	os.Remove(pemCertPath)

	return pemCert, nil
}

func createPfxCert(managementCert string, pfxCertPath string) error {

	pfxCert, err := base64.StdEncoding.DecodeString(managementCert)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pfxCertPath, pfxCert, 0644)
	if err != nil {
		return err
	}

	return nil
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

	if len(activeSubscription.ManagementCertificate) == 0 {
		activeSubscription.ManagementCertificate = publishProfile.ManagementCertificate
		activeSubscription.ServiceManagementUrl = publishProfile.Url
	}

	return activeSubscription, nil
}

func getAzureDir() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	azureDir := path.Join(homeDir, ".azure")
	//Create azure dir if does not exists
	if _, err := os.Stat(azureDir); os.IsNotExist(err) {
		mkdirErr := os.Mkdir(azureDir, 0644)
		if mkdirErr != nil {
			return "", mkdirErr
		}
	}

	return azureDir, nil
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
