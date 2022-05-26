// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"runtime"
	"strings"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-amqp-common-go/v3/cbs"
	"github.com/Azure/azure-amqp-common-go/v3/conn"
	"github.com/Azure/azure-amqp-common-go/v3/sas"
	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/azure"
	"golang.org/x/net/websocket"
)

type (
	namespace struct {
		name          string
		tokenProvider auth.TokenProvider
		host          string
		useWebSocket  bool
	}

	// namespaceOption provides structure for configuring a new Event Hub namespace
	namespaceOption func(h *namespace) error
)

// newNamespaceWithConnectionString configures a namespace with the information provided in a Service Bus connection string
func namespaceWithConnectionString(connStr string) namespaceOption {
	return func(ns *namespace) error {
		parsed, err := conn.ParsedConnectionFromStr(connStr)
		if err != nil {
			return err
		}
		ns.name = parsed.Namespace
		ns.host = parsed.Host
		provider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(parsed.KeyName, parsed.Key))
		if err != nil {
			return err
		}
		ns.tokenProvider = provider
		return nil
	}
}

func namespaceWithAzureEnvironment(name string, tokenProvider auth.TokenProvider, env azure.Environment) namespaceOption {
	return func(ns *namespace) error {
		ns.name = name
		ns.tokenProvider = tokenProvider
		ns.host = "amqps://" + ns.name + "." + env.ServiceBusEndpointSuffix
		return nil
	}
}

// newNamespace creates a new namespace configured through NamespaceOption(s)
func newNamespace(opts ...namespaceOption) (*namespace, error) {
	ns := &namespace{}

	for _, opt := range opts {
		err := opt(ns)
		if err != nil {
			return nil, err
		}
	}

	return ns, nil
}

func (ns *namespace) newConnection() (*amqp.Client, error) {
	host := ns.getAmqpsHostURI()

	defaultConnOptions := []amqp.ConnOption{
		amqp.ConnSASLAnonymous(),
		amqp.ConnProperty("product", "MSGolangClient"),
		amqp.ConnProperty("version", Version),
		amqp.ConnProperty("platform", runtime.GOOS),
		amqp.ConnProperty("framework", runtime.Version()),
		amqp.ConnProperty("user-agent", rootUserAgent),
	}

	if ns.useWebSocket {
		trimmedHost := strings.TrimPrefix(ns.host, "amqps://")
		wssConn, err := websocket.Dial("wss://"+trimmedHost+"/$servicebus/websocket", "amqp", "http://localhost/")
		if err != nil {
			return nil, err
		}

		wssConn.PayloadType = websocket.BinaryFrame
		return amqp.New(wssConn, append(defaultConnOptions, amqp.ConnServerHostname(trimmedHost))...)
	}

	return amqp.Dial(host, defaultConnOptions...)
}

func (ns *namespace) negotiateClaim(ctx context.Context, conn *amqp.Client, entityPath string) error {
	span, ctx := ns.startSpanFromContext(ctx, "eh.namespace.negotiateClaim")
	defer span.End()

	audience := ns.getEntityAudience(entityPath)
	return cbs.NegotiateClaim(ctx, audience, conn, ns.tokenProvider)
}

func (ns *namespace) getAmqpsHostURI() string {
	return ns.host + "/"
}

func (ns *namespace) getAmqpHostURI() string {
	return strings.Replace(ns.getAmqpsHostURI(), "amqps", "amqp", 1)
}

func (ns *namespace) getEntityAudience(entityPath string) string {
	return ns.getAmqpsHostURI() + entityPath
}

func (ns *namespace) getHTTPSHostURI() string {
	return strings.Replace(ns.getAmqpsHostURI(), "amqps", "https", 1)
}
