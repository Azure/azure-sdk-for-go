package azureSdkForGo

import (
	"fmt"
	"errors"
	"io/ioutil"
	"encoding/xml"
)

var publishSettings PublishSettings = PublishSettings{}

func GetPublishSettings() PublishSettings {
	return publishSettings;
}

func SetPublishSettings(id, cert string) {
	publishSettings.SubscriptionID = id
	publishSettings.SubscriptionCert = cert
}

func ImportPublishSettings(id, cert string) {
	SetPublishSettings(id, cert)
}

func ImportPublishSettingsFile(filePath string) {

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		PrintErrorAndExit(err)
	}

	activeSubscription, subErr := GetActiveSubscription(content)
	if subErr != nil {
		PrintErrorAndExit(subErr)
	}

	fmt.Println(activeSubscription.Id)
	fmt.Println(activeSubscription.Name)
	fmt.Println(activeSubscription.ManagementCertificate)
	fmt.Println(activeSubscription.ServiceManagementUrl)

	certPath := GetSubscriptionCert(activeSubscription)
	ImportPublishSettings(activeSubscription.Id, certPath)
}

func GetSubscriptionCert(subscription Subscription) string {
	//!TODO
	return ""
}

func GetActiveSubscription(content []byte) (Subscription, error) {
	publishData := PublishData{}
	activeSubscription := Subscription{}

	err := xml.Unmarshal(content, &publishData)
	if err != nil {
		PrintErrorAndExit(err)
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

	if len(subscriptions) == 1 {
		activeSubscription.ManagementCertificate = publishProfile.ManagementCertificate
		activeSubscription.ServiceManagementUrl = publishProfile.Url
	}

	return activeSubscription, nil
}

type PublishSettings struct {
	SubscriptionID				string
	SubscriptionCert			string
}

type PublishData struct {
	PublishProfiles []PublishProfile `xml:"PublishProfile"`
}

type PublishProfile struct {
	SchemaVersion string `xml:",attr"`
	PublishMethod string `xml:",attr"`
	Url string `xml:",attr"`
	ManagementCertificate string `xml:",attr"`
	Subscriptions []Subscription `xml:"Subscription"`
}

type Subscription struct {
	ServiceManagementUrl string `xml:",attr"`
	Id string `xml:",attr"`
	Name string `xml:",attr"`
	ManagementCertificate string `xml:",attr"`
}
