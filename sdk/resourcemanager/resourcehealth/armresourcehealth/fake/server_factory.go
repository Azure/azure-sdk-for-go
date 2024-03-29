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

// ServerFactory is a fake server for instances of the armresourcehealth.ClientFactory type.
type ServerFactory struct {
	AvailabilityStatusesServer              AvailabilityStatusesServer
	ChildAvailabilityStatusesServer         ChildAvailabilityStatusesServer
	ChildResourcesServer                    ChildResourcesServer
	EmergingIssuesServer                    EmergingIssuesServer
	EventServer                             EventServer
	EventsServer                            EventsServer
	ImpactedResourcesServer                 ImpactedResourcesServer
	MetadataServer                          MetadataServer
	OperationsServer                        OperationsServer
	SecurityAdvisoryImpactedResourcesServer SecurityAdvisoryImpactedResourcesServer
}

// NewServerFactoryTransport creates a new instance of ServerFactoryTransport with the provided implementation.
// The returned ServerFactoryTransport instance is connected to an instance of armresourcehealth.ClientFactory via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerFactoryTransport(srv *ServerFactory) *ServerFactoryTransport {
	return &ServerFactoryTransport{
		srv: srv,
	}
}

// ServerFactoryTransport connects instances of armresourcehealth.ClientFactory to instances of ServerFactory.
// Don't use this type directly, use NewServerFactoryTransport instead.
type ServerFactoryTransport struct {
	srv                                       *ServerFactory
	trMu                                      sync.Mutex
	trAvailabilityStatusesServer              *AvailabilityStatusesServerTransport
	trChildAvailabilityStatusesServer         *ChildAvailabilityStatusesServerTransport
	trChildResourcesServer                    *ChildResourcesServerTransport
	trEmergingIssuesServer                    *EmergingIssuesServerTransport
	trEventServer                             *EventServerTransport
	trEventsServer                            *EventsServerTransport
	trImpactedResourcesServer                 *ImpactedResourcesServerTransport
	trMetadataServer                          *MetadataServerTransport
	trOperationsServer                        *OperationsServerTransport
	trSecurityAdvisoryImpactedResourcesServer *SecurityAdvisoryImpactedResourcesServerTransport
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
	case "AvailabilityStatusesClient":
		initServer(s, &s.trAvailabilityStatusesServer, func() *AvailabilityStatusesServerTransport {
			return NewAvailabilityStatusesServerTransport(&s.srv.AvailabilityStatusesServer)
		})
		resp, err = s.trAvailabilityStatusesServer.Do(req)
	case "ChildAvailabilityStatusesClient":
		initServer(s, &s.trChildAvailabilityStatusesServer, func() *ChildAvailabilityStatusesServerTransport {
			return NewChildAvailabilityStatusesServerTransport(&s.srv.ChildAvailabilityStatusesServer)
		})
		resp, err = s.trChildAvailabilityStatusesServer.Do(req)
	case "ChildResourcesClient":
		initServer(s, &s.trChildResourcesServer, func() *ChildResourcesServerTransport {
			return NewChildResourcesServerTransport(&s.srv.ChildResourcesServer)
		})
		resp, err = s.trChildResourcesServer.Do(req)
	case "EmergingIssuesClient":
		initServer(s, &s.trEmergingIssuesServer, func() *EmergingIssuesServerTransport {
			return NewEmergingIssuesServerTransport(&s.srv.EmergingIssuesServer)
		})
		resp, err = s.trEmergingIssuesServer.Do(req)
	case "EventClient":
		initServer(s, &s.trEventServer, func() *EventServerTransport { return NewEventServerTransport(&s.srv.EventServer) })
		resp, err = s.trEventServer.Do(req)
	case "EventsClient":
		initServer(s, &s.trEventsServer, func() *EventsServerTransport { return NewEventsServerTransport(&s.srv.EventsServer) })
		resp, err = s.trEventsServer.Do(req)
	case "ImpactedResourcesClient":
		initServer(s, &s.trImpactedResourcesServer, func() *ImpactedResourcesServerTransport {
			return NewImpactedResourcesServerTransport(&s.srv.ImpactedResourcesServer)
		})
		resp, err = s.trImpactedResourcesServer.Do(req)
	case "MetadataClient":
		initServer(s, &s.trMetadataServer, func() *MetadataServerTransport { return NewMetadataServerTransport(&s.srv.MetadataServer) })
		resp, err = s.trMetadataServer.Do(req)
	case "OperationsClient":
		initServer(s, &s.trOperationsServer, func() *OperationsServerTransport { return NewOperationsServerTransport(&s.srv.OperationsServer) })
		resp, err = s.trOperationsServer.Do(req)
	case "SecurityAdvisoryImpactedResourcesClient":
		initServer(s, &s.trSecurityAdvisoryImpactedResourcesServer, func() *SecurityAdvisoryImpactedResourcesServerTransport {
			return NewSecurityAdvisoryImpactedResourcesServerTransport(&s.srv.SecurityAdvisoryImpactedResourcesServer)
		})
		resp, err = s.trSecurityAdvisoryImpactedResourcesServer.Do(req)
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
