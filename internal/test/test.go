package test

import (
	"math/rand"
	"os"
	"time"

	"context"
	rm "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	sbmgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/suite"
)

const (
	location          = "eastus"
	resourceGroupName = "test"
)

type (
	// BaseSuite encapsulates a end to end test of Service Bus with build up and tear down of all SB resources
	BaseSuite struct {
		suite.Suite
		TenantID       string
		SubscriptionID string
		ClientID       string
		ClientSecret   string
		ConnStr        string
		Namespace      string
		Token          *adal.ServicePrincipalToken
		Environment    azure.Environment
		TagID          string
	}
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

func init() {
	rand.Seed(time.Now().Unix())
}

// SetupSuite prepares the test suite and provisions a standard Service Bus Namespace
func (suite *BaseSuite) SetupSuite() {
	suite.TenantID = mustGetEnv("AZURE_TENANT_ID")
	suite.SubscriptionID = mustGetEnv("AZURE_SUBSCRIPTION_ID")
	suite.ClientID = mustGetEnv("AZURE_CLIENT_ID")
	suite.ClientSecret = mustGetEnv("AZURE_CLIENT_SECRET")
	suite.Namespace = mustGetEnv("SERVICEBUS_NAMESPACE")
	suite.ConnStr = mustGetEnv("SERVICEBUS_CONNECTION_STRING")
	suite.Token = suite.servicePrincipalToken()
	suite.Environment = azure.PublicCloud
	suite.TagID = RandomString("tag", 10)

	err := suite.ensureProvisioned(sbmgmt.SkuTierStandard)
	if err != nil {
		suite.T().Fatal(err)
	}
}

// TearDownSuite destroys created resources during the run of the suite
func (suite *BaseSuite) TearDownSuite() {
	// tear down queues and subscriptions maybe??
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}

func (suite *BaseSuite) servicePrincipalToken() *adal.ServicePrincipalToken {

	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, suite.TenantID)
	if err != nil {
		suite.T().Fatal(err)
	}

	tokenProvider, err := adal.NewServicePrincipalToken(*oauthConfig,
		suite.ClientID,
		suite.ClientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		suite.T().Fatal(err)
	}

	return tokenProvider
}

func (suite *BaseSuite) getRmGroupClient() *rm.GroupsClient {
	groupsClient := rm.NewGroupsClient(suite.SubscriptionID)
	groupsClient.Authorizer = autorest.NewBearerAuthorizer(suite.Token)
	return &groupsClient
}

func (suite *BaseSuite) getNamespaceClient() *sbmgmt.NamespacesClient {
	nsClient := sbmgmt.NewNamespacesClient(suite.SubscriptionID)
	nsClient.Authorizer = autorest.NewBearerAuthorizer(suite.Token)
	return &nsClient
}

func (suite *BaseSuite) ensureProvisioned(tier sbmgmt.SkuTier) error {
	groupsClient := suite.getRmGroupClient()
	location := location
	_, err := groupsClient.CreateOrUpdate(context.Background(), resourceGroupName, rm.Group{Location: &location})
	if err != nil {
		return err
	}

	nsClient := suite.getNamespaceClient()
	_, err = nsClient.Get(context.Background(), resourceGroupName, suite.Namespace)
	if err != nil {
		ns := sbmgmt.SBNamespace{
			Sku: &sbmgmt.SBSku{
				Name: sbmgmt.SkuName(tier),
				Tier: tier,
			},
			Location: &location,
		}
		res, err := nsClient.CreateOrUpdate(context.Background(), resourceGroupName, suite.Namespace, ns)
		if err != nil {
			return err
		}

		return res.WaitForCompletion(context.Background(), nsClient.Client)
	}

	return nil
}

// RandomName generates a random Event Hub name tagged with the suite id
func (suite *BaseSuite) RandomName(prefix string, length int) string {
	return RandomString(prefix, length) + "-" + suite.TagID
}

// RandomString generates a random string with prefix
func RandomString(prefix string, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}
