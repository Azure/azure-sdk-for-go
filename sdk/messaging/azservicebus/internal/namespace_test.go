// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/suite"
)

type (
	serviceBusSuite struct {
		test.BaseSuite
	}
)

func TestSB(t *testing.T) {
	if os.Getenv("SERVICEBUS_CONNECTION_STRING") == "" {
		t.Skipf("environment variable SERVICEBUS_CONNECTION_STRING was not set")
	}

	suite.Run(t, new(serviceBusSuite))
}

func (suite *serviceBusSuite) TestCreateNamespaceFromConnectionString() {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING") // `Endpoint=sb://XXXX.servicebus.windows.net/;SharedAccessKeyName=XXXX;SharedAccessKey=XXXX`
	if connStr == "" {
		suite.T().Skipf("environment variable SERVICEBUS_CONNECTION_STRING was not set")
	}

	ns, err := NewNamespace(NamespaceWithConnectionString(connStr))
	if suite.NoError(err) {
		suite.Contains(connStr, ns.Name)
	}
}

func TestNewNamespaceWithAzureEnvironment(t *testing.T) {
	ns, err := NewNamespace(NamespaceWithAzureEnvironment("namespaceName", "AzureGermanCloud"))
	if err != nil {
		t.Fatalf("unexpected error creating namespace: %s", err)
	}
	if ns.Environment != azure.GermanCloud {
		t.Fatalf("expected namespace environment to be %q but was %q", azure.GermanCloud, ns.Environment)
	}
	if !strings.EqualFold(ns.Suffix, azure.GermanCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("expected suffix to be %q but was %q", azure.GermanCloud.ServiceBusEndpointSuffix, ns.Suffix)
	}
	if ns.Name != "namespaceName" {
		t.Fatalf("expected namespace name to be %q but was %q", "namespaceName", ns.Name)
	}
}

// TearDownSuite destroys created resources during the run of the suite
func (suite *serviceBusSuite) TearDownSuite() {
	suite.BaseSuite.TearDownSuite()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	suite.deleteAllTaggedQueues(ctx)
	suite.deleteAllTaggedTopics(ctx)
}

func (suite *serviceBusSuite) deleteAllTaggedQueues(ctx context.Context) {
	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()

	qs, err := qm.List(ctx)
	if err != nil {
		suite.T().Fatal(err)
	}

	for _, q := range qs {
		if strings.HasSuffix(q.Name, suite.TagID) {
			err := qm.Delete(ctx, q.Name)
			if err != nil {
				suite.T().Fatal(err)
			}
		}
	}
}

func (suite *serviceBusSuite) deleteAllTaggedTopics(ctx context.Context) {
	ns := suite.getNewSasInstance()
	tm := ns.NewTopicManager()

	topics, err := tm.List(ctx)
	if err != nil {
		suite.T().Fatal(err)
	}

	for _, topic := range topics {
		if strings.HasSuffix(topic.Name, suite.TagID) {
			err := tm.Delete(ctx, topic.Name)
			if err != nil {
				suite.T().Fatal(err)
			}
		}
	}
}

func (suite *serviceBusSuite) getNewSasInstance(opts ...NamespaceOption) *Namespace {
	ns, err := NewNamespace(append(opts, NamespaceWithConnectionString(suite.ConnStr))...)
	if err != nil {
		suite.T().Fatal(err)
	}
	return ns
}
