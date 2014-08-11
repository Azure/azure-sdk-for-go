package azureSdkForGo

var publishSettings PublishSettings = PublishSettings{}

func ImportPublishSettings(id, cert string) {
	SetPublishSettings(id, cert)
}

func GetPublishSettings() PublishSettings {
	return publishSettings;
}

func SetPublishSettings(id, cert string) {
	publishSettings.SubscriptionID = id
	publishSettings.SubscriptionCert = cert
}

type PublishSettings struct {
	SubscriptionID				string
	SubscriptionCert			string
}
