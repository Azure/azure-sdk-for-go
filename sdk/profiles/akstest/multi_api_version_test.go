package akstest

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armresource201805 "github.com/Azure/azure-sdk-for-go/sdk/profiles/aksprofile1/resourcemanager/resources/armresources"
	armresource202010 "github.com/Azure/azure-sdk-for-go/sdk/profiles/aksprofile2/resourcemanager/resources/armresources"
)

func TestResources(t *testing.T) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %+v", err)
	}
	ctx := context.Background()

	client1, err := armresource201805.NewResourceGroupsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %+v", err)
	}
	pager1 := client1.NewListPager(nil)
	if pager1.More() {
		result, err := pager1.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to fetch page: %+v", err)
		}
		for _, v := range result.Value {
			log.Printf("Get resource group: %s", *v.Name)
		}
	}

	client2, err := armresource202010.NewDeploymentsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %+v", err)
	}
	pager2 := client2.NewListAtSubscriptionScopePager(nil)
	if pager2.More() {
		result, err := pager2.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to fetch page: %+v", err)
		}
		for _, v := range result.Value {
			log.Printf("Get deployment: %s", *v.Name)
		}
	}
}
