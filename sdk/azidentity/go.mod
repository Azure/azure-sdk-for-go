module github.com/Azure/azure-sdk-for-go/sdk/azidentity

go 1.16

replace github.com/Azure/azure-sdk-for-go/sdk/azcore => ../azcore

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.20.0
	github.com/Azure/azure-sdk-for-go/sdk/internal v0.8.2
	github.com/AzureAD/microsoft-authentication-library-for-go v0.4.0
	github.com/davecgh/go-spew v1.1.1 // indirect
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
