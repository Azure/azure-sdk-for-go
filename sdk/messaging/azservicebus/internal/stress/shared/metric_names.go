// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package shared

type Metric string

// These names are modeled off of the metrics from Java
// https://github.com/Azure/azure-sdk-for-java/tree/main/sdk/servicebus/azure-messaging-servicebus/src/main/java/com/azure/messaging/servicebus/implementation/instrumentation/ServiceBusMeter.java
//
// and from our standard for attributes:
// https://gist.github.com/lmolkova/e4215c0f44a49ef824983382762e6b92
const (
	MetricConnectionLost Metric = "messaging.servicebus.connectionlost"
	MetricMessagesSent   Metric = "messaging.servicebus.messages.sent"

	// metrics related to Service Bus sessions (NOT amqp sessions)
	MetricSessionAccept    Metric = "messaging.servicebus.session.accept"
	MetricSessionTimeoutMS Metric = "messaging.servicebus.session.timeout"

	MetricSettlementRequestDuration Metric = "messaging.servicebus.settlement.request.duration"
	MetricReceiveLag                Metric = "messaging.servicebus.receiver.lag"

	MetricAMQPSendDuration Metric = "messaging.az.amqp.producer.send.duration"
	// - amqp.delivery_state (ex: 'accepted')

	MetricAMQPMgmtRequestDuration Metric = "messaging.az.amqp.management.request.duration"
	// - amqp.status_code (ex: 'accepted')

	MetricAMQPSettlementRequestDuration Metric = "messaging.servicebus.settlement.request.duration"
	MetricAMQPSettlementSequenceNum     Metric = "messaging.servicebus.settlement.sequence_number"

	// OTel attributes (should come in naturally when we swap over to using OTel throughout)
	// otel.status_code

	// TODO: I've made these up entirely.
	MetricMessageReceived Metric = "messaging.servicebus.messages.received"
	MetricMessagePeeked   Metric = "messaging.servicebus.messages.peeked"
	MetricCloseDuration   Metric = "messaging.servicebus.close.duration"
	MetricLockRenew       Metric = "messaging.servicebus.lockrenew.duration" // TODO: separate for session vs message lock?
)

const (
	AttrAMQPDeliveryState string = "amqp.delivery_state"
	AttrAMQPStatusCode    string = "amqp.status_code"

	// TODO: I made these up entirely
	AttrMessageCount string = "amqp.message_count"
)

// these metrics are specific to stress tests and wouldn't be in customer code.
const (
	MetricStressSuccessfulCancels = "stress.cancels"
)
