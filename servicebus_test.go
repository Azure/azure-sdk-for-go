package servicebus

import (
	"context"
	rm "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	sbmgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"math/rand"
	"os"
	"pack.ag/amqp"
	"testing"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

const (
	RootRuleName      = "RootManageSharedAccessKey"
	WestUS2           = "westus2"
	ResourceGroupName = "sbtest"
)

// ServiceBusSuite encapsulates a end to end test of Service Bus with build up and tear down of all SB resources
type ServiceBusSuite struct {
	suite.Suite
	TenantID       string
	SubscriptionID string
	ClientID       string
	ClientSecret   string
	Namespace      string
	Token          *adal.ServicePrincipalToken
	Environment    azure.Environment
}

func (suite *ServiceBusSuite) SetupSuite() {
	suite.TenantID = mustGetenv("AZURE_TENANT_ID")
	suite.SubscriptionID = mustGetenv("AZURE_SUBSCRIPTION_ID")
	suite.ClientID = mustGetenv("AZURE_CLIENT_ID")
	suite.ClientSecret = mustGetenv("AZURE_CLIENT_SECRET")
	suite.Namespace = mustGetenv("SERVICEBUS_NAMESPACE")
	suite.Token = suite.servicePrincipalToken()
	suite.Environment = azure.PublicCloud

	err := suite.ensureProvisioned()
	if err != nil {
		log.Fatalln(err)
	}
}

func (suite *ServiceBusSuite) TearDownSuite() {
	// tear down queues and subscriptions
}

func (suite *ServiceBusSuite) TestQueue() {
	tests := []func(*testing.T, SenderReceiver, string){testQueueSend}

	spToken := suite.servicePrincipalToken()
	sb, err := NewWithSPToken(spToken, suite.SubscriptionID, ResourceGroupName, suite.Namespace, RootRuleName, suite.Environment)
	if err != nil {
		log.Fatalln(err)
	}
	defer sb.Close()

	queueName := randomName("gosbtest", 10)
	for _, testFunc := range tests {
		testFunc(suite.T(), sb, queueName)
	}
}

func testQueueSend(t *testing.T, sb SenderReceiver, queueName string) {
	err := sb.Send(context.Background(), queueName, &amqp.Message{
		Data: []byte("Hello!"),
	})
	assert.Nil(t, err)
}

func TestServiceBusSuite(t *testing.T) {
	suite.Run(t, new(ServiceBusSuite))
}

func TestServiceBusConstruction(t *testing.T) {
	connStr := os.Getenv("AZURE_SERVICE_BUS_CONN_STR") // `Endpoint=sb://XXXX.servicebus.windows.net/;SharedAccessKeyName=XXXX;SharedAccessKey=XXXX`
	sb, err := NewWithConnectionString(connStr)
	defer sb.Close()
	assert.Nil(t, err)
}

func mustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("Environment variable '" + key + "' required for integration tests.")
	}
	return v
}

func randomName(prefix string, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + "-" + string(b)
}

func (suite *ServiceBusSuite) servicePrincipalToken() *adal.ServicePrincipalToken {

	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, suite.TenantID)
	if err != nil {
		log.Fatalln(err)
	}

	tokenProvider, err := adal.NewServicePrincipalToken(*oauthConfig,
		suite.ClientID,
		suite.ClientSecret,
		azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	return tokenProvider
}

func (suite *ServiceBusSuite) getRmGroupClient() *rm.GroupsClient {
	groupsClient := rm.NewGroupsClient(suite.SubscriptionID)
	groupsClient.Authorizer = autorest.NewBearerAuthorizer(suite.Token)
	return &groupsClient
}

func (suite *ServiceBusSuite) getServiceBusNamespaceClient() *sbmgmt.NamespacesClient {
	nsClient := sbmgmt.NewNamespacesClient(suite.SubscriptionID)
	nsClient.Authorizer = autorest.NewBearerAuthorizer(suite.Token)
	return &nsClient
}

func (suite *ServiceBusSuite) ensureProvisioned() error {
	log.Println("ensuring test resource group is provisioned")
	groupsClient := suite.getRmGroupClient()
	location := WestUS2
	_, err := groupsClient.CreateOrUpdate(context.Background(), ResourceGroupName, rm.Group{Location: &location})
	if err != nil {
		return err
	}

	nsClient := suite.getServiceBusNamespaceClient()
	_, err = nsClient.Get(context.Background(), ResourceGroupName, suite.Namespace)
	if err != nil {
		log.Println("namespace is not there, create it")
		res, err := nsClient.CreateOrUpdate(
			context.Background(),
			ResourceGroupName,
			suite.Namespace,
			sbmgmt.SBNamespace{
				Sku: &sbmgmt.SBSku{
					Name: "Standard",
					Tier: sbmgmt.SkuTierStandard,
				},
				Location: &location,
			})
		if err != nil {
			return err
		}

		log.Println("waiting for namespace to provision")
		return res.WaitForCompletion(context.Background(), nsClient.Client)
	}

	log.Println("namespace was already provisioned")
	return nil
}
