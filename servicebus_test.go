package servicebus

import (
	"context"
	"flag"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	rm "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	sbmgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pack.ag/amqp"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
	debug       = flag.Bool("debug", false, "output debug level logging")
)

const (
	RootRuleName      = "RootManageSharedAccessKey"
	Location          = "westus"
	ResourceGroupName = "sbtest"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type (
	// ServiceBusSuite encapsulates a end to end test of Service Bus with build up and tear down of all SB resources
	ServiceBusSuite struct {
		suite.Suite
		TenantID       string
		SubscriptionID string
		ClientID       string
		ClientSecret   string
		Namespace      string
		Token          *adal.ServicePrincipalToken
		Environment    azure.Environment
	}
)

func (suite *ServiceBusSuite) SetupSuite() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	suite.TenantID = mustGetenv("AZURE_TENANT_ID")
	suite.SubscriptionID = mustGetenv("AZURE_SUBSCRIPTION_ID")
	suite.ClientID = mustGetenv("AZURE_CLIENT_ID")
	suite.ClientSecret = mustGetenv("AZURE_CLIENT_SECRET")
	suite.Namespace = "something" //mustGetenv("SERVICEBUS_NAMESPACE")
	suite.Token = suite.servicePrincipalToken()
	suite.Environment = azure.PublicCloud

	err := suite.ensureProvisioned(sbmgmt.SkuTierStandard)
	if err != nil {
		log.Fatalln(err)
	}
}

func (suite *ServiceBusSuite) TearDownSuite() {
	// tear down queues and subscriptions maybe??
}

func (suite *ServiceBusSuite) TestQueue() {
	tests := map[string]func(*testing.T, SenderReceiver, string){
		"SimpleSend":            testQueueSend,
		"SendAndReceiveInOrder": testQueueSendAndReceiveInOrder,
		"DuplicateDetection":    testDuplicateDetection,
	}

	sb := suite.getNewInstance()
	defer func() {
		sb.Close()
	}()

	for name, testFunc := range tests {
		setupTestTeardown := func(t *testing.T) {
			queueName := randomName("gosbtest", 10)
			defer sb.DeleteQueue(context.Background(), queueName)

			window := 3 * time.Minute
			_, err := sb.EnsureQueue(
				context.Background(),
				queueName,
				QueueWithPartitioning(),
				QueueWithDuplicateDetection(nil),
				QueueWithLockDuration(&window))
			if err != nil {
				log.Fatalln(err)
			}

			testFunc(t, sb, queueName)
		}

		suite.T().Run(name, setupTestTeardown)
	}
}

func testQueueSend(t *testing.T, sb SenderReceiver, queueName string) {
	err := sb.Send(context.Background(), queueName, &amqp.Message{
		Data: []byte("Hello!"),
	})
	assert.Nil(t, err)
}

func testQueueSendAndReceiveInOrder(t *testing.T, sb SenderReceiver, queueName string) {
	numMessages := rand.Intn(100) + 20
	var wg sync.WaitGroup
	wg.Add(numMessages + 1)
	messages := make([]string, numMessages)
	for i := 0; i < numMessages; i++ {
		messages[i] = randomName("hello", 10)
	}

	go func() {
		for _, message := range messages {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := sb.Send(ctx, queueName, &amqp.Message{Data: []byte(message)})
			cancel()
			if err != nil {
				log.Fatalln(err)
			}
		}
		defer wg.Done()
	}()

	// ensure in-order processing of messages from the queue
	count := 0
	sb.Receive(queueName, func(ctx context.Context, msg *amqp.Message) error {
		assert.Equal(t, messages[count], string(msg.Data))
		count++
		wg.Done()
		return nil
	})
	wg.Wait()
}

func testDuplicateDetection(t *testing.T, sb SenderReceiver, queueName string) {
	dupID := uuid.NewV4().String()
	messages := []struct {
		ID   string
		Data string
	}{
		{
			ID:   dupID,
			Data: "hello 1!",
		},
		{
			ID:   dupID,
			Data: "hello 1!",
		},
		{
			ID:   uuid.NewV4().String(),
			Data: "hello 2!",
		},
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for _, msg := range messages {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err := sb.Send(ctx, queueName, &amqp.Message{Data: []byte(msg.Data)}, SendWithMessageID(msg.ID))
			cancel()
			if err != nil {
				log.Fatalln(err)
			}
		}
		defer wg.Done()
	}()

	received := make(map[interface{}]string)
	sb.Receive(queueName, func(ctx context.Context, msg *amqp.Message) error {
		// we should get 2 messages discarding the duplicate ID
		received[msg.Properties.MessageID] = string(msg.Data)
		wg.Done()
		return nil
	})
	wg.Wait()
	assert.Equal(t, 2, len(received), "should not have more than 2 messages", received)
}

func TestServiceBusSuite(t *testing.T) {
	suite.Run(t, new(ServiceBusSuite))
}

func TestCreateFromConnectionString(t *testing.T) {
	connStr := os.Getenv("AZURE_SERVICE_BUS_CONN_STR") // `Endpoint=sb://XXXX.servicebus.windows.net/;SharedAccessKeyName=XXXX;SharedAccessKey=XXXX`
	sb, err := NewWithConnectionString(connStr)
	defer sb.Close()
	assert.Nil(t, err)
}

func BenchmarkSend(b *testing.B) {
	sbSuite := &ServiceBusSuite{}
	sbSuite.SetupSuite()
	defer sbSuite.TearDownSuite()

	sb := sbSuite.getNewInstance()
	defer func() {
		err := sb.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	queueName := randomName("gosbbench", 10)
	_, err := sb.EnsureQueue(context.Background(), queueName, nil)
	if err != nil {
		log.Fatalln(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sb.Send(context.Background(), queueName, &amqp.Message{
			Data: []byte("Hello!"),
		})
	}
	b.StopTimer()
	err = sb.DeleteQueue(context.Background(), queueName)
	if err != nil {
		log.Fatalln(err)
	}
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

func (suite *ServiceBusSuite) ensureProvisioned(tier sbmgmt.SkuTier) error {
	groupsClient := suite.getRmGroupClient()
	location := Location
	_, err := groupsClient.CreateOrUpdate(context.Background(), ResourceGroupName, rm.Group{Location: &location})
	if err != nil {
		return err
	}

	nsClient := suite.getServiceBusNamespaceClient()
	_, err = nsClient.Get(context.Background(), ResourceGroupName, suite.Namespace)
	if err != nil {
		ns := sbmgmt.SBNamespace{
			Sku: &sbmgmt.SBSku{
				Name: sbmgmt.SkuName(tier),
				Tier: tier,
			},
			Location: &location,
		}
		res, err := nsClient.CreateOrUpdate(context.Background(), ResourceGroupName, suite.Namespace, ns)
		if err != nil {
			return err
		}

		return res.WaitForCompletion(context.Background(), nsClient.Client)
	}

	return nil
}

func (suite *ServiceBusSuite) getNewInstance() SenderReceiverManager {
	return getNewInstance(suite.TenantID, suite.SubscriptionID, suite.Namespace, suite.ClientID, suite.ClientSecret, suite.Environment)
}

func getNewInstance(tenantID, subscriptionID, namespace, appID, secret string, env azure.Environment) SenderReceiverManager {
	cred := ServicePrincipalCredentials{
		TenantID:      tenantID,
		ApplicationID: appID,
		Secret:        secret,
	}
	sb, err := NewWithServicePrincipal(subscriptionID, namespace, cred, env)
	if err != nil {
		log.Fatalln(err)
	}
	return sb
}
