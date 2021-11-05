// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func WrapWithQueueEnvelope(qd *QueueDescription, tokenProvider auth.TokenProvider) (*QueueEnvelope, []MiddlewareFunc) {
	qd.ServiceBusSchema = to.StringPtr(serviceBusSchema)

	qe := &QueueEnvelope{
		Entry: &Entry{
			AtomSchema: atomSchema,
		},
		Content: &queueContent{
			Type:             applicationXML,
			QueueDescription: *qd,
		},
	}

	var mw []MiddlewareFunc
	if qd.ForwardTo != nil {
		mw = append(mw, addSupplementalAuthorization(*qd.ForwardTo, tokenProvider))
	}

	if qd.ForwardDeadLetteredMessagesTo != nil {
		mw = append(mw, addDeadLetterSupplementalAuthorization(*qd.ForwardDeadLetteredMessagesTo, tokenProvider))
	}

	return qe, mw
}

func WrapWithTopicEnvelope(td *TopicDescription) *TopicEnvelope {
	td.ServiceBusSchema = to.StringPtr(serviceBusSchema)

	return &TopicEnvelope{
		Entry: &Entry{
			AtomSchema: atomSchema,
		},
		Content: &topicContent{
			Type:             applicationXML,
			TopicDescription: *td,
		},
	}
}

func WrapWithSubscriptionEnvelope(sd *SubscriptionDescription) *SubscriptionEnvelope {
	sd.ServiceBusSchema = to.StringPtr(serviceBusSchema)

	return &SubscriptionEnvelope{
		Entry: &Entry{
			AtomSchema: atomSchema,
		},
		Content: &subscriptionContent{
			Type:                    applicationXML,
			SubscriptionDescription: *sd,
		},
	}
}
