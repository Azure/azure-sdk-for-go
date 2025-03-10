// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tracing

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/go-amqp"
)

// messageWrapper implements the TextMapCarrier interface for sender side
type messageWrapper amqp.Message

// messageCarrierAdapter wraps a Message so that it implements the propagation.TextMapCarrier interface
func messageCarrierAdapter(message amqp.Message) tracing.Carrier {
	mw := messageWrapper(message)
	return tracing.NewCarrier(tracing.CarrierImpl{
		Get:  mw.Get,
		Set:  mw.Set,
		Keys: mw.Keys,
	})
}

func (mw messageWrapper) Set(key string, value string) {
	if mw.ApplicationProperties == nil {
		mw.ApplicationProperties = make(map[string]any)
	}
	mw.ApplicationProperties[key] = value
}

func (mw messageWrapper) Get(key string) string {
	if mw.ApplicationProperties == nil || mw.ApplicationProperties[key] == nil {
		return ""
	}
	return mw.ApplicationProperties[key].(string)
}

func (mw messageWrapper) Keys() []string {
	keys := make([]string, 0, len(mw.ApplicationProperties))
	for k := range mw.ApplicationProperties {
		keys = append(keys, k)
	}
	return keys
}
