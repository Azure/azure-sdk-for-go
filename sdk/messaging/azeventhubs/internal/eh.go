// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
)

func (l *rpcLink) LinkName() string {
	return l.sender.LinkName()
}

func (ns *Namespace) NewRPCLink(ctx context.Context, managementPath string) (amqpwrap.RPCLink, error) {
	client, _, err := ns.GetAMQPClientImpl(ctx)

	if err != nil {
		return nil, err
	}

	rpcLink, err := NewRPCLink(ctx, RPCLinkArgs{
		Client:   client,
		Address:  managementPath,
		LogEvent: exported.EventProducer,
	})

	if err != nil {
		return nil, err
	}

	return rpcLink, nil
}

func (ns *Namespace) GetTokenForEntity(eventHub string) (*auth.Token, error) {
	audience := ns.GetEntityAudience(eventHub)
	return ns.TokenProvider.GetToken(audience)
}

type NamespaceForManagementOps interface {
	NamespaceForAMQPLinks
	GetTokenForEntity(eventHub string) (*auth.Token, error)
}

// TODO: might just consolidate.
type NamespaceForProducerOrConsumer = NamespaceForManagementOps
