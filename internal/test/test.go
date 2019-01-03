package test

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"context"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/Azure/azure-amqp-common-go/conn"
	rm "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2017-05-10/resources"
	sbmgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
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
		Location       string
		Namespace      string
		ResourceGroup  string
		Token          *adal.ServicePrincipalToken
		Environment    azure.Environment
		TagID          string
		closer         io.Closer
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
	if err := godotenv.Load(); err != nil {
		suite.T().Log(err)
	}

	setFromEnv := func(key string, target *string) {
		v := os.Getenv(key)
		if v == "" {
			suite.FailNowf("Environment variable %q required for integration tests.", key)
		}

		*target = v
	}

	setFromEnv("AZURE_TENANT_ID", &suite.TenantID)
	setFromEnv("AZURE_SUBSCRIPTION_ID", &suite.SubscriptionID)
	setFromEnv("AZURE_CLIENT_ID", &suite.ClientID)
	setFromEnv("AZURE_CLIENT_SECRET", &suite.ClientSecret)
	setFromEnv("SERVICEBUS_CONNECTION_STRING", &suite.ConnStr)
	setFromEnv("TEST_SERVICEBUS_RESOURCE_GROUP", &suite.ResourceGroup)

	// TODO: automatically infer the location from the resource group, if it's not specified.
	// https://github.com/Azure/azure-service-bus-go/issues/40
	setFromEnv("TEST_SERVICEBUS_LOCATION", &suite.Location)

	parsed, err := conn.ParsedConnectionFromStr(suite.ConnStr)
	if !suite.NoError(err) {
		suite.FailNowf("connection string could not be parsed", "Connection String: %q", suite.ConnStr)
	}
	suite.Namespace = parsed.Namespace
	suite.Token = suite.servicePrincipalToken()
	suite.Environment = azure.PublicCloud
	suite.TagID = RandomString("tag", 10)
	suite.setupTracing()

	if !suite.NoError(suite.ensureProvisioned(sbmgmt.SkuTierStandard)) {
		suite.FailNow("failed to ensure provisioned")
	}
}

// TearDownSuite destroys created resources during the run of the suite
func (suite *BaseSuite) TearDownSuite() {
	if suite.closer != nil {
		_ = suite.closer.Close()
	}
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
	_, err := groupsClient.CreateOrUpdate(context.Background(), suite.ResourceGroup, rm.Group{Location: &suite.Location})
	if err != nil {
		return err
	}

	nsClient := suite.getNamespaceClient()
	_, err = nsClient.Get(context.Background(), suite.ResourceGroup, suite.Namespace)
	if err != nil {
		ns := sbmgmt.SBNamespace{
			Sku: &sbmgmt.SBSku{
				Name: sbmgmt.SkuName(tier),
				Tier: tier,
			},
			Location: &suite.Location,
		}
		res, err := nsClient.CreateOrUpdate(context.Background(), suite.ResourceGroup, suite.Namespace, ns)
		if err != nil {
			return err
		}

		return res.WaitForCompletionRef(context.Background(), nsClient.Client)
	}

	return nil
}

func (suite *BaseSuite) setupTracing() error {
	if os.Getenv("TRACING") == "true" {
		// Sample configuration for testing. Use constant sampling to sample every trace
		// and enable LogSpan to log every span via configured Logger.
		cfg := config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LocalAgentHostPort: "0.0.0.0:6831",
			},
		}

		// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
		// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
		// frameworks.
		jLogger := jaegerlog.StdLogger

		closer, err := cfg.InitGlobalTracer(
			"ehtests",
			config.Logger(jLogger),
		)

		suite.closer = closer
		return err
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
