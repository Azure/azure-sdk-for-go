package servicebus

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-service-bus-go/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type (
	serviceBusSuite struct {
		test.BaseSuite
	}
)

func TestCreateNamespaceFromConnectionString(t *testing.T) {
	connStr := os.Getenv("SERVICEBUS_CONNECTION_STRING") // `Endpoint=sb://XXXX.servicebus.windows.net/;SharedAccessKeyName=XXXX;SharedAccessKey=XXXX`
	ns, err := NewNamespace(NamespaceWithConnectionString(connStr))
	assert.Nil(t, err)
	assert.Contains(t, connStr, ns.Name)
}

func TestServiceBusSuite(t *testing.T) {
	suite.Run(t, new(serviceBusSuite))
}

// TearDownSuite destroys created resources during the run of the suite
func (suite *serviceBusSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	suite.deleteAllTaggedQueues(ctx)
}

func (suite *serviceBusSuite) deleteAllTaggedQueues(ctx context.Context) {
	ns := suite.getNewSasInstance()
	qm := ns.NewQueueManager()

	feed, err := qm.List(ctx)
	if err != nil {
		suite.T().Fatal(err)
	}

	for _, entry := range feed.Entries {
		if strings.HasSuffix(entry.Title, suite.TagID) {
			err := qm.Delete(ctx, entry.Title)
			if err != nil {
				suite.T().Fatal(err)
			}
		}
	}
}

func (suite *serviceBusSuite) getNewSasInstance() *Namespace {
	ns, err := getNewSasInstance(suite.ConnStr)
	if err != nil {
		suite.T().Fatal(err)
	}
	return ns
}

func getNewSasInstance(connStr string) (*Namespace, error) {
	return NewNamespace(NamespaceWithConnectionString(connStr))
}
