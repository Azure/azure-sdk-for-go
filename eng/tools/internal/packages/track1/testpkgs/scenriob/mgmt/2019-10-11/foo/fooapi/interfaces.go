package fooapi

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/tools/internal/packages/track1/testpkgs/scenrioa/foo"
)

// GatewaysClientAPI ...
type GatewaysClientAPI interface {
	CreateOrUpdate(ctx context.Context, resourceGroupName string, gatewayName string, parameters foo.Gateway) (result foo.Gateway, err error)
}
