module github.com/Azure/azure-sdk-for-go/sdk/storage/azblob

go 1.13

require (
	github.com/Azure/azure-pipeline-go v0.2.3
	github.com/Azure/azure-sdk-for-go/sdk/armcore v0.1.0
	github.com/Azure/azure-sdk-for-go/sdk/azcore v0.9.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.0.0-00010101000000-000000000000
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f
)

replace github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../azidentity
