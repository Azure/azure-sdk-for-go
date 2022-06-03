// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package test is an internal package to handle common test setup
package test

import (
	"context"
	"flag"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	common "github.com/Azure/azure-amqp-common-go/v3"
	mgmt "github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	rm "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	"github.com/Azure/go-autorest/autorest/azure"
	azauth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
	debug       = flag.Bool("debug", false, "output debug level logging")
)

const (
	defaultTimeout = 1 * time.Minute
)

type (
	// BaseSuite encapsulates a end to end test of Event Hubs with build up and tear down of all EH resources
	BaseSuite struct {
		suite.Suite
		SubscriptionID    string
		Namespace         string
		ResourceGroupName string
		Location          string
		Env               azure.Environment
		TagID             string
		closer            io.Closer
	}

	// HubMgmtOption represents an option for configuring an Event Hub.
	HubMgmtOption func(model *mgmt.Model) error
	// NamespaceMgmtOption represents an option for configuring a Namespace
	NamespaceMgmtOption func(ns *mgmt.EHNamespace) error
)

func init() {
	rand.Seed(time.Now().Unix())
	loadEnv()
}

// SetupSuite constructs the test suite from the environment and
func (suite *BaseSuite) SetupSuite() {
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	suite.SubscriptionID = MustGetEnv("AZURE_SUBSCRIPTION_ID")

	if suite.SubscriptionID == "" {
		log.Printf("No AZURE_SUBSCRIPTION_ID variable, skipping test")
		suite.T().Skip()
	}

	suite.Namespace = MustGetEnv("EVENTHUB_NAMESPACE")
	suite.ResourceGroupName = MustGetEnv("TEST_EVENTHUB_RESOURCE_GROUP")
	suite.Location = MustGetEnv("TEST_EVENTHUB_LOCATION")
	envName := os.Getenv("AZURE_ENVIRONMENT")
	suite.TagID = RandomString("tag", 5)

	if envName == "" {
		suite.Env = azure.PublicCloud
	} else {
		var err error
		env, err := azure.EnvironmentFromName(envName)
		if !suite.NoError(err) {
			suite.FailNow("could not find env name")
		}
		suite.Env = env
	}

	if !suite.NoError(suite.ensureProvisioned(mgmt.SkuTierStandard)) {
		suite.FailNow("failed provisioning")
	}

	//if !suite.NoError(suite.setupTracing()) {
	//	suite.FailNow("failed to setup tracing")
	//}
}

// TearDownSuite might one day destroy all of the resources in the suite, but I'm not sure we want to do that just yet...
func (suite *BaseSuite) TearDownSuite() {
	// maybe tear down all existing resource??
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	suite.deleteAllTaggedEventHubs(ctx)
	if suite.closer != nil {
		suite.NoError(suite.closer.Close())
	}
}

// RandomHub creates a hub with a random'ish name
func (suite *BaseSuite) RandomHub(opts ...HubMgmtOption) (*mgmt.Model, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout*2)
	defer cancel()

	name := suite.RandomName("goehtest", 6)
	model, err := suite.ensureEventHub(ctx, name, opts...)
	suite.Require().NoError(err)
	suite.Require().NotNil(model)
	suite.Require().NotNil(model.PartitionIds)
	suite.Require().Len(*model.PartitionIds, 4)
	time.Sleep(250 * time.Millisecond) // introduce a bit of a delay before using the hub
	return model, func() {
		if model != nil {
			suite.DeleteEventHub(*model.Name)
		}
	}
}

// EnsureEventHub creates an Event Hub if it doesn't exist
func (suite *BaseSuite) ensureEventHub(ctx context.Context, name string, opts ...HubMgmtOption) (*mgmt.Model, error) {
	client := suite.getEventHubMgmtClient()
	hub, err := client.Get(ctx, suite.ResourceGroupName, suite.Namespace, name)

	if err != nil {
		newHub := &mgmt.Model{
			Name: &name,
			Properties: &mgmt.Properties{
				PartitionCount: common.PtrInt64(4),
			},
		}

		for _, opt := range opts {
			err = opt(newHub)
			if err != nil {
				return nil, err
			}
		}

		var lastErr error
		deadline, _ := ctx.Deadline()
		for time.Now().Before(deadline) {
			hub, err = suite.tryHubCreate(ctx, client, name, newHub)
			if err == nil {
				lastErr = nil
				break
			}
			lastErr = err
		}

		if lastErr != nil {
			return nil, lastErr
		}
	}
	return &hub, nil
}

func (suite *BaseSuite) tryHubCreate(ctx context.Context, client *mgmt.EventHubsClient, name string, hub *mgmt.Model) (mgmt.Model, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	_, err := client.CreateOrUpdate(ctx, suite.ResourceGroupName, suite.Namespace, name, *hub)
	if err != nil {
		return mgmt.Model{}, err
	}

	return client.Get(ctx, suite.ResourceGroupName, suite.Namespace, name)
}

// DeleteEventHub deletes an Event Hub within the given Namespace
func (suite *BaseSuite) DeleteEventHub(name string) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	client := suite.getEventHubMgmtClient()
	_, err := client.Delete(ctx, suite.ResourceGroupName, suite.Namespace, name)
	suite.NoError(err)
}

func (suite *BaseSuite) deleteAllTaggedEventHubs(ctx context.Context) {
	client := suite.getEventHubMgmtClient()
	res, err := client.ListByNamespace(ctx, suite.ResourceGroupName, suite.Namespace, to.Int32Ptr(0), to.Int32Ptr(20))
	if err != nil {
		suite.T().Log("error listing namespaces")
		suite.T().Error(err)
	}

	for res.NotDone() {
		for _, val := range res.Values() {
			if strings.Contains(*val.Name, suite.TagID) {
				for i := 0; i < 5; i++ {
					if _, err := client.Delete(ctx, suite.ResourceGroupName, suite.Namespace, *val.Name); err != nil {
						suite.T().Logf("error deleting %q", *val.Name)
						suite.T().Error(err)
						time.Sleep(3 * time.Second)
					} else {
						break
					}
				}
			} else if !strings.HasPrefix(*val.Name, "examplehub_") {
				suite.T().Logf("%q does not contain %q, so it won't be deleted.", *val.Name, suite.TagID)
			}
		}
		suite.NoError(res.Next())
	}
}

func (suite *BaseSuite) ensureProvisioned(tier mgmt.SkuTier) error {
	_, err := ensureResourceGroup(context.Background(), suite.SubscriptionID, suite.ResourceGroupName, suite.Location, suite.Env)
	if err != nil {
		return err
	}

	_, err = suite.ensureNamespace()
	return err
}

// ensureResourceGroup creates a Azure Resource Group if it does not already exist
func ensureResourceGroup(ctx context.Context, subscriptionID, name, location string, env azure.Environment) (*rm.Group, error) {
	groupClient := getRmGroupClientWithToken(subscriptionID, env)
	group, err := groupClient.Get(ctx, name)
	if group.Response.Response == nil {
		// tcp dial error or something else where the response was not populated
		return nil, err
	}

	if group.StatusCode == http.StatusNotFound {
		group, err = groupClient.CreateOrUpdate(ctx, name, rm.Group{Location: common.PtrString(location)})
		if err != nil {
			return nil, err
		}
	} else if group.StatusCode >= 400 {
		return nil, err
	}

	return &group, nil
}

// ensureNamespace creates a Azure Event Hub Namespace if it does not already exist
func ensureNamespace(ctx context.Context, subscriptionID, rg, name, location string, env azure.Environment, opts ...NamespaceMgmtOption) (*mgmt.EHNamespace, error) {
	_, err := ensureResourceGroup(ctx, subscriptionID, rg, location, env)
	if err != nil {
		return nil, err
	}

	client := getNamespaceMgmtClientWithToken(subscriptionID, env)
	namespace, err := client.Get(ctx, rg, name)
	if err != nil {
		return nil, err
	}

	if namespace.StatusCode == 404 {
		newNamespace := &mgmt.EHNamespace{
			Name: &name,

			Sku: &mgmt.Sku{
				Name:     mgmt.Basic,
				Tier:     mgmt.SkuTierBasic,
				Capacity: common.PtrInt32(1),
			},
			EHNamespaceProperties: &mgmt.EHNamespaceProperties{
				IsAutoInflateEnabled:   common.PtrBool(false),
				MaximumThroughputUnits: common.PtrInt32(1),
			},
		}

		for _, opt := range opts {
			err = opt(newNamespace)
			if err != nil {
				return nil, err
			}
		}

		nsFuture, err := client.CreateOrUpdate(ctx, rg, name, *newNamespace)
		if err != nil {
			return nil, err
		}

		err = nsFuture.WaitForCompletionRef(ctx, client.Client)
		if err != nil {
			return nil, err
		}

		namespace, err = nsFuture.Result(*client)
		if err != nil {
			return nil, err
		}
	} else if namespace.StatusCode >= 400 {
		return nil, err
	}

	return &namespace, nil
}

func (suite *BaseSuite) getEventHubMgmtClient() *mgmt.EventHubsClient {
	client := mgmt.NewEventHubsClientWithBaseURI(suite.Env.ResourceManagerEndpoint, suite.SubscriptionID)
	a, err := azauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	client.Authorizer = a
	return &client
}

func (suite *BaseSuite) ensureNamespace() (*mgmt.EHNamespace, error) {
	ns, err := ensureNamespace(context.Background(), suite.SubscriptionID, suite.ResourceGroupName, suite.Namespace, suite.Location, suite.Env)
	if err != nil {
		return nil, err
	}
	return ns, err
}

func getNamespaceMgmtClientWithToken(subscriptionID string, env azure.Environment) *mgmt.NamespacesClient {
	client := mgmt.NewNamespacesClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID)
	a, err := azauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	client.Authorizer = a
	return &client
}

func getRmGroupClientWithToken(subscriptionID string, env azure.Environment) *rm.GroupsClient {
	groupsClient := rm.NewGroupsClientWithBaseURI(env.ResourceManagerEndpoint, subscriptionID)
	a, err := azauth.NewAuthorizerFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	groupsClient.Authorizer = a
	return &groupsClient
}

//func (suite *BaseSuite) setupTracing() error {
//	if os.Getenv("TRACING") != "true" {
//		return nil
//	}
//	exporter, err := jaeger.NewExporter(jaeger.Options{
//		AgentEndpoint: "localhost:6831",
//		Process: jaeger.Process{
//			ServiceName: "eh-tests",
//		},
//	})
//	if err != nil {
//		return err
//	}
//	trace.RegisterExporter(exporter)
//	return nil
//}

// MustGetEnv will panic or return the env var for a given string key
func MustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Printf("Env variable '" + key + "' required for integration tests.")
		return ""
	}
	return v
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

func loadEnv() {
	lookForMe := []string{".env", "../.env", "../../.env"}
	var reader io.ReadCloser
	for _, env := range lookForMe {
		r, err := os.Open(env)
		if err == nil {
			reader = r
			break
		}
	}

	if reader == nil {
		log.Printf("no .env files were found in %v, no integration tests will be run", lookForMe)
		return
	}

	defer func() {
		if err := reader.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	envMap, err := godotenv.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}

	for key, val := range envMap {
		if err := os.Setenv(key, val); err != nil {
			log.Fatal(err)
		}
	}
}
