package servicebus

import (
	"context"

	mgmt "github.com/Azure/azure-sdk-for-go/services/servicebus/mgmt/2017-04-01/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"
)

func (sb *serviceBus) GetNamespace(ctx context.Context) (*mgmt.SBNamespace, error) {
	client := sb.getNamespaceMgmtClient()
	res, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	findNamespace := func(namespaces []mgmt.SBNamespace) (*mgmt.SBNamespace, bool) {
		for _, ns := range namespaces {
			if *ns.Name == sb.namespace {
				return &ns, true
			}
		}
		return nil, false
	}

	ns, ok := findNamespace(res.Values())
	if ok {
		return ns, nil
	}

	for res.NotDone() {
		err := res.Next()
		if err != nil {
			return nil, err
		}
		ns, ok := findNamespace(res.Values())
		if ok {
			return ns, nil
		}
	}
	return nil, errors.New("could not find namespace")
}

func (sb *serviceBus) getNamespaceMgmtClient() mgmt.NamespacesClient {
	client := mgmt.NewNamespacesClientWithBaseURI(sb.environment.ResourceManagerEndpoint, sb.subscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(sb.armToken)
	return client
}
