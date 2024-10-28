// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strings"
	"sync"
)

// ServerFactory is a fake server for instances of the armiotoperations.ClientFactory type.
type ServerFactory struct {
	// BrokerAuthenticationServer contains the fakes for client BrokerAuthenticationClient
	BrokerAuthenticationServer BrokerAuthenticationServer

	// BrokerAuthorizationServer contains the fakes for client BrokerAuthorizationClient
	BrokerAuthorizationServer BrokerAuthorizationServer

	// BrokerServer contains the fakes for client BrokerClient
	BrokerServer BrokerServer

	// BrokerListenerServer contains the fakes for client BrokerListenerClient
	BrokerListenerServer BrokerListenerServer

	// DataflowServer contains the fakes for client DataflowClient
	DataflowServer DataflowServer

	// DataflowEndpointServer contains the fakes for client DataflowEndpointClient
	DataflowEndpointServer DataflowEndpointServer

	// DataflowProfileServer contains the fakes for client DataflowProfileClient
	DataflowProfileServer DataflowProfileServer

	// InstanceServer contains the fakes for client InstanceClient
	InstanceServer InstanceServer

	// OperationsServer contains the fakes for client OperationsClient
	OperationsServer OperationsServer
}

// NewServerFactoryTransport creates a new instance of ServerFactoryTransport with the provided implementation.
// The returned ServerFactoryTransport instance is connected to an instance of armiotoperations.ClientFactory via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerFactoryTransport(srv *ServerFactory) *ServerFactoryTransport {
	return &ServerFactoryTransport{
		srv: srv,
	}
}

// ServerFactoryTransport connects instances of armiotoperations.ClientFactory to instances of ServerFactory.
// Don't use this type directly, use NewServerFactoryTransport instead.
type ServerFactoryTransport struct {
	srv                          *ServerFactory
	trMu                         sync.Mutex
	trBrokerAuthenticationServer *BrokerAuthenticationServerTransport
	trBrokerAuthorizationServer  *BrokerAuthorizationServerTransport
	trBrokerServer               *BrokerServerTransport
	trBrokerListenerServer       *BrokerListenerServerTransport
	trDataflowServer             *DataflowServerTransport
	trDataflowEndpointServer     *DataflowEndpointServerTransport
	trDataflowProfileServer      *DataflowProfileServerTransport
	trInstanceServer             *InstanceServerTransport
	trOperationsServer           *OperationsServerTransport
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
	case "BrokerAuthenticationClient":
		initServer(s, &s.trBrokerAuthenticationServer, func() *BrokerAuthenticationServerTransport {
			return NewBrokerAuthenticationServerTransport(&s.srv.BrokerAuthenticationServer)
		})
		resp, err = s.trBrokerAuthenticationServer.Do(req)
	case "BrokerAuthorizationClient":
		initServer(s, &s.trBrokerAuthorizationServer, func() *BrokerAuthorizationServerTransport {
			return NewBrokerAuthorizationServerTransport(&s.srv.BrokerAuthorizationServer)
		})
		resp, err = s.trBrokerAuthorizationServer.Do(req)
	case "BrokerClient":
		initServer(s, &s.trBrokerServer, func() *BrokerServerTransport { return NewBrokerServerTransport(&s.srv.BrokerServer) })
		resp, err = s.trBrokerServer.Do(req)
	case "BrokerListenerClient":
		initServer(s, &s.trBrokerListenerServer, func() *BrokerListenerServerTransport {
			return NewBrokerListenerServerTransport(&s.srv.BrokerListenerServer)
		})
		resp, err = s.trBrokerListenerServer.Do(req)
	case "DataflowClient":
		initServer(s, &s.trDataflowServer, func() *DataflowServerTransport { return NewDataflowServerTransport(&s.srv.DataflowServer) })
		resp, err = s.trDataflowServer.Do(req)
	case "DataflowEndpointClient":
		initServer(s, &s.trDataflowEndpointServer, func() *DataflowEndpointServerTransport {
			return NewDataflowEndpointServerTransport(&s.srv.DataflowEndpointServer)
		})
		resp, err = s.trDataflowEndpointServer.Do(req)
	case "DataflowProfileClient":
		initServer(s, &s.trDataflowProfileServer, func() *DataflowProfileServerTransport {
			return NewDataflowProfileServerTransport(&s.srv.DataflowProfileServer)
		})
		resp, err = s.trDataflowProfileServer.Do(req)
	case "InstanceClient":
		initServer(s, &s.trInstanceServer, func() *InstanceServerTransport { return NewInstanceServerTransport(&s.srv.InstanceServer) })
		resp, err = s.trInstanceServer.Do(req)
	case "OperationsClient":
		initServer(s, &s.trOperationsServer, func() *OperationsServerTransport { return NewOperationsServerTransport(&s.srv.OperationsServer) })
		resp, err = s.trOperationsServer.Do(req)
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
