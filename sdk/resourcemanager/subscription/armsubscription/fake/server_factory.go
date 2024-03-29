//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strings"
	"sync"
)

// ServerFactory is a fake server for instances of the armsubscription.ClientFactory type.
type ServerFactory struct {
	AliasServer          AliasServer
	BillingAccountServer BillingAccountServer
	Server               Server
	OperationsServer     OperationsServer
	PolicyServer         PolicyServer
	SubscriptionsServer  SubscriptionsServer
	TenantsServer        TenantsServer
}

// NewServerFactoryTransport creates a new instance of ServerFactoryTransport with the provided implementation.
// The returned ServerFactoryTransport instance is connected to an instance of armsubscription.ClientFactory via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerFactoryTransport(srv *ServerFactory) *ServerFactoryTransport {
	return &ServerFactoryTransport{
		srv: srv,
	}
}

// ServerFactoryTransport connects instances of armsubscription.ClientFactory to instances of ServerFactory.
// Don't use this type directly, use NewServerFactoryTransport instead.
type ServerFactoryTransport struct {
	srv                    *ServerFactory
	trMu                   sync.Mutex
	trAliasServer          *AliasServerTransport
	trBillingAccountServer *BillingAccountServerTransport
	trServer               *ServerTransport
	trOperationsServer     *OperationsServerTransport
	trPolicyServer         *PolicyServerTransport
	trSubscriptionsServer  *SubscriptionsServerTransport
	trTenantsServer        *TenantsServerTransport
}

// Do implements the policy.Transporter interface for ServerFactoryTransport.
func (s *ServerFactoryTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	client := method[:strings.Index(method, ".")]
	var resp *http.Response
	var err error

	switch client {
	case "AliasClient":
		initServer(s, &s.trAliasServer, func() *AliasServerTransport { return NewAliasServerTransport(&s.srv.AliasServer) })
		resp, err = s.trAliasServer.Do(req)
	case "BillingAccountClient":
		initServer(s, &s.trBillingAccountServer, func() *BillingAccountServerTransport {
			return NewBillingAccountServerTransport(&s.srv.BillingAccountServer)
		})
		resp, err = s.trBillingAccountServer.Do(req)
	case "Client":
		initServer(s, &s.trServer, func() *ServerTransport { return NewServerTransport(&s.srv.Server) })
		resp, err = s.trServer.Do(req)
	case "OperationsClient":
		initServer(s, &s.trOperationsServer, func() *OperationsServerTransport { return NewOperationsServerTransport(&s.srv.OperationsServer) })
		resp, err = s.trOperationsServer.Do(req)
	case "PolicyClient":
		initServer(s, &s.trPolicyServer, func() *PolicyServerTransport { return NewPolicyServerTransport(&s.srv.PolicyServer) })
		resp, err = s.trPolicyServer.Do(req)
	case "SubscriptionsClient":
		initServer(s, &s.trSubscriptionsServer, func() *SubscriptionsServerTransport {
			return NewSubscriptionsServerTransport(&s.srv.SubscriptionsServer)
		})
		resp, err = s.trSubscriptionsServer.Do(req)
	case "TenantsClient":
		initServer(s, &s.trTenantsServer, func() *TenantsServerTransport { return NewTenantsServerTransport(&s.srv.TenantsServer) })
		resp, err = s.trTenantsServer.Do(req)
	default:
		err = fmt.Errorf("unhandled client %s", client)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func initServer[T any](s *ServerFactoryTransport, dst **T, src func() *T) {
	s.trMu.Lock()
	if *dst == nil {
		*dst = src()
	}
	s.trMu.Unlock()
}
