package armconnectedcache_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/connectedcache/armconnectedcache"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/assert"
)

// C:\Users\v-liujudy\go\pkg\mod\github.com\!azure\azure-sdk-for-go\sdk\resourcemanager
const (
	SubscriptionID    = "faa080af-c1d8-40ad-9cce-e1a450ca5b57"
	ResourceGroupName = "judytest04"
	ResourceLocation  = "eastus2"
	ResourceName      = "judy-test001"
	PathToPackage     = "sdk/resourcemanager/sql/armsql/testdata"
)

func TestCreateOrUpdateGroupOfSubscriptionId(t *testing.T) {
	subsriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	assert.NotEmpty(t, subsriptionId)
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	assert.Nil(t, err)
	assert.NotNil(t, cred)
	fmt.Println("get default credential")

	// new client factory
	clientFactory, err := armresources.NewClientFactory(subsriptionId, cred, nil)
	assert.Nil(t, err)
	assert.NotNil(t, clientFactory)
	client := clientFactory.NewResourceGroupsClient()
	assert.NotNil(t, client)
	ctx := context.Background()

	createOrUpdateGroupResponse, err := client.CreateOrUpdate(ctx, ResourceGroupName, armresources.ResourceGroup{
		Location: to.Ptr(ResourceLocation),
	}, nil)
	assert.Nil(t, err)
	assert.NotNil(t, createOrUpdateGroupResponse)
	fmt.Println("create resource group client")

	// check whether create new group successfully
	checkGroupExistResponse, err := client.CheckExistence(ctx, ResourceGroupName, nil)
	assert.Nil(t, err)
	assert.True(t, checkGroupExistResponse.Success)
	fmt.Println("create new resource group ", ResourceGroupName, " of ", subsriptionId, "successfully")
}

func TestCreateOrUpdateIspCustomersClient(t *testing.T) {
	// make sure that the group has been created
	TestCreateOrUpdateGroupOfSubscriptionId(t)

	// cred := GetAzureCredentail()
	// assert.NotNil(t, cred)
	// // create new msql resource under group
	// connectedCacheClientFactory, err := armconnectedcache.NewClientFactory(SubscriptionID, cred, nil)
	// assert.Nil(t, err)
	// assert.NotNil(t, connectedCacheClientFactory)

	// serverclient, err := connectedCacheClientFactory.NewCacheNodesOperationsClient().Get()
	// assert.NotNil(t, serverclient)

	subsriptionId := os.Getenv("AZURE_SUBSCRIPTION_ID")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	clientFactory, err := armconnectedcache.NewClientFactory(subsriptionId, cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := clientFactory.NewIspCustomersClient().BeginCreateOrUpdate(ctx, ResourceGroupName, ResourceName, armconnectedcache.IspCustomerResource{
		Location: to.Ptr("westus"),
		Properties: &armconnectedcache.CustomerProperty{
			Customer: &armconnectedcache.CustomerEntity{
				FullyQualifiedResourceID: to.Ptr("uqsbtgae"),
				CustomerName:             to.Ptr("mkpzynfqihnjfdbaqbqwyhd"),
				ContactEmail:             to.Ptr("xquos"),
				ContactPhone:             to.Ptr("vue"),
				ContactName:              to.Ptr("wxyqjoyoscmvimgwhpitxky"),
				IsEntitled:               to.Ptr(true),
				ReleaseVersion:           to.Ptr[int32](20),
				ClientTenantID:           to.Ptr("fproidkpgvpdnac"),
				IsEnterpriseManaged:      to.Ptr(true),
				ShouldMigrate:            to.Ptr(true),
				ResendSignupCode:         to.Ptr(true),
				VerifySignupCode:         to.Ptr(true),
				VerifySignupPhrase:       to.Ptr("tprjvttkgmrqlsyicnidhm"),
			},
			AdditionalCustomerProperties: &armconnectedcache.AdditionalCustomerProperties{
				CustomerEmail:                 to.Ptr("zdjgibsidydyzm"),
				CustomerTransitAsn:            to.Ptr("habgklnxqzmozqpazoyejwiphezpi"),
				CustomerAsn:                   to.Ptr("hgrelgnrtdkleisnepfolu"),
				CustomerEntitlementSKUID:      to.Ptr("b"),
				CustomerEntitlementSKUGUID:    to.Ptr("rvzmdpxyflgqetvpwupnfaxsweiiz"),
				CustomerEntitlementSKUName:    to.Ptr("waaqfijr"),
				CustomerEntitlementExpiration: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-01-30T00:54:04.773Z"); return t }()),
				OptionalProperty1:             to.Ptr("qhmwxza"),
				OptionalProperty2:             to.Ptr("l"),
				OptionalProperty3:             to.Ptr("mblwwvbie"),
				OptionalProperty4:             to.Ptr("vzuek"),
				OptionalProperty5:             to.Ptr("fzjodscdfcdr"),
			},
			Error: &armconnectedcache.ErrorDetail{},
		},
		Tags: map[string]*string{
			"key1878": to.Ptr("warz"),
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	fmt.Println("created resourceId:", res.ID)

	res1, err1 := clientFactory.NewIspCustomersClient().Get(ctx, ResourceGroupName, ResourceName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err1)
	}
	fmt.Println("resourceId:", res1.ID)
}

func GetAzureCredentail() (cred *azidentity.DefaultAzureCredential) {
	// get default credential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Println("get default credentail fail, error:", err.Error())
		return
	}
	fmt.Println("get default credential")
	return cred
}

func CreateAzureResourceGroup(groupName string) (group *armresources.ResourceGroup, err error) {
	cred := GetAzureCredentail()
	clientFactory, err := armresources.NewClientFactory(SubscriptionID, cred, nil)
	if err != nil {
		fmt.Println("create resource client factory fail, error:", err.Error())
		return
	}
	client := clientFactory.NewResourceGroupsClient()
	createOrUpdateGroupResponse, err := client.CreateOrUpdate(context.TODO(), groupName, armresources.ResourceGroup{
		Location: to.Ptr(ResourceLocation),
	}, nil)
	if err != nil {
		fmt.Println("create resource group fail, error:", err.Error())
		return
	}
	group = &createOrUpdateGroupResponse.ResourceGroup
	return
}
