module github.com/Azure/azure-sdk-for-go/sdk/storage/blob/2019-07-07/azblob

go 1.13

require (
    github.com/Azure/azure-sdk-for-go/sdk/azcore v0.5.0
    github.com/Azure/azure-sdk-for-go/sdk/azidentity v0.0.0-00010101000000-000000000000
)

replace (
    github.com/Azure/azure-sdk-for-go/sdk/azidentity => ../../../../azidentity
)