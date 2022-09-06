module akstest

go 1.18

require (
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/profiles/aksprofile1 v0.0.0-20220906082653-6b805d8278b6
	github.com/Azure/azure-sdk-for-go/sdk/profiles/aksprofile2 v0.0.0-20220906082653-6b805d8278b6
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.0.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v0.4.0 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/pkg/browser v0.0.0-20210115035449-ce105d075bb4 // indirect
	golang.org/x/crypto v0.0.0-20220511200225-c6db032c6c88 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/Azure/azure-sdk-for-go/sdk/profiles/aksprofile1/resourcemanager/resources/armresources => ../aksprofile1/resourcemanager/resources/armresources

replace github.com/Azure/azure-sdk-for-go/sdk/profiles/aksprofile2/resourcemanager/resources/armresources => ../aksprofile2/resourcemanager/resources/armresources
