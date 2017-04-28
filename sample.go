package main

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

var (
	groupClient resources.GroupsClient
)

func init() {
	subscriptionID := getEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	tenantID := getEnvVarOrExit("AZURE_TENANT_ID")

	oauthConfig, err := azure.PublicCloud.OAuthConfigForTenant(tenantID)
	onErrorFail(err, "OAuthConfigForTenant failed")

	clientID := getEnvVarOrExit("AZURE_CLIENT_ID")
	clientSecret := getEnvVarOrExit("AZURE_CLIENT_SECRET")
	spToken, err := azure.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, azure.PublicCloud.ResourceManagerEndpoint)
	onErrorFail(err, "NewServicePrincipalToken failed")

	createClients(subscriptionID, spToken)
}

func main() {
	list, err := groupClient.List("", to.Int32Ptr(1))
	if err != nil {
		panic(err)
	}
	fmt.Println("List len: %d", len(*list.Value))
	list, err = list.Next()
	if err != nil {
		panic(err)
	}
	fmt.Println("List len: %d", len(*list.Value))

	anotherList, err := groupClient.ListComplete("", to.Int32Ptr(1), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(anotherList)
}

// getEnvVarOrExit returns the value of specified environment variable or terminates if it's not defined.
func getEnvVarOrExit(varName string) string {
	value := os.Getenv(varName)
	if value == "" {
		fmt.Printf("Missing environment variable %s\n", varName)
		os.Exit(1)
	}

	return value
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string) {
	if err != nil {
		fmt.Printf("%s: %s\n", message, err)
		os.Exit(1)
	}
}

func createClients(subscriptionID string, spToken *azure.ServicePrincipalToken) {
	groupClient = resources.NewGroupsClient(subscriptionID)
	groupClient.Authorizer = spToken
}
