// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package atom

import (
	"encoding/xml"

	"github.com/Azure/go-autorest/autorest/date"
)

// Queues
type (
	// QueueEntity is the Azure Service Bus description of a Queue for management activities
	QueueEntity struct {
		*QueueDescription
		*Entity
	}

	// QueueFeed is a specialized feed containing QueueEntries
	QueueFeed struct {
		*Feed
		Entries []QueueEnvelope `xml:"entry"`
	}

	// QueueEnvelope is a specialized Queue feed entry
	QueueEnvelope struct {
		*Entry
		Content *queueContent `xml:"content"`
	}

	// queueContent is a specialized Queue body for an Atom entry
	queueContent struct {
		XMLName          xml.Name         `xml:"content"`
		Type             string           `xml:"type,attr"`
		QueueDescription QueueDescription `xml:"QueueDescription"`
	}

	// QueueDescription is the content type for Queue management requests
	QueueDescription struct {
		XMLName xml.Name `xml:"QueueDescription"`
		BaseEntityDescription
		LockDuration                        *string       `xml:"LockDuration,omitempty"`               // LockDuration - ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1 minute.
		MaxSizeInMegabytes                  *int32        `xml:"MaxSizeInMegabytes,omitempty"`         // MaxSizeInMegabytes - The maximum size of the queue in megabytes, which is the size of memory allocated for the queue. Default is 1024.
		RequiresDuplicateDetection          *bool         `xml:"RequiresDuplicateDetection,omitempty"` // RequiresDuplicateDetection - A value indicating if this queue requires duplicate detection.
		RequiresSession                     *bool         `xml:"RequiresSession,omitempty"`
		DefaultMessageTimeToLive            *string       `xml:"DefaultMessageTimeToLive,omitempty"`            // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
		DeadLetteringOnMessageExpiration    *bool         `xml:"DeadLetteringOnMessageExpiration,omitempty"`    // DeadLetteringOnMessageExpiration - A value that indicates whether this queue has dead letter support when a message expires.
		DuplicateDetectionHistoryTimeWindow *string       `xml:"DuplicateDetectionHistoryTimeWindow,omitempty"` // DuplicateDetectionHistoryTimeWindow - ISO 8601 timeSpan structure that defines the duration of the duplicate detection history. The default value is 10 minutes.
		MaxDeliveryCount                    *int32        `xml:"MaxDeliveryCount,omitempty"`                    // MaxDeliveryCount - The maximum delivery count. A message is automatically deadlettered after this number of deliveries. default value is 10.
		EnableBatchedOperations             *bool         `xml:"EnableBatchedOperations,omitempty"`             // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		SizeInBytes                         *int64        `xml:"SizeInBytes,omitempty"`                         // SizeInBytes - The size of the queue, in bytes.
		MessageCount                        *int64        `xml:"MessageCount,omitempty"`                        // MessageCount - The number of messages in the queue.
		IsAnonymousAccessible               *bool         `xml:"IsAnonymousAccessible,omitempty"`
		Status                              *EntityStatus `xml:"Status,omitempty"`
		AccessedAt                          *date.Time    `xml:"AccessedAt,omitempty"`
		CreatedAt                           *date.Time    `xml:"CreatedAt,omitempty"`
		UpdatedAt                           *date.Time    `xml:"UpdatedAt,omitempty"`
		SupportOrdering                     *bool         `xml:"SupportOrdering,omitempty"`
		AutoDeleteOnIdle                    *string       `xml:"AutoDeleteOnIdle,omitempty"`
		EnablePartitioning                  *bool         `xml:"EnablePartitioning,omitempty"`
		EnableExpress                       *bool         `xml:"EnableExpress,omitempty"`
		CountDetails                        *CountDetails `xml:"CountDetails,omitempty"`
		ForwardTo                           *string       `xml:"ForwardTo,omitempty"`
		ForwardDeadLetteredMessagesTo       *string       `xml:"ForwardDeadLetteredMessagesTo,omitempty"` // ForwardDeadLetteredMessagesTo - absolute URI of the entity to forward dead letter messages
	}
)

// Topics
type (
	// TopicEntity is the Azure Service Bus description of a Topic for management activities
	TopicEntity struct {
		*TopicDescription
		*Entity
	}

	// TopicEnvelope is a specialized Topic feed entry
	TopicEnvelope struct {
		*Entry
		Content *topicContent `xml:"content"`
	}

	// topicContent is a specialized Topic body for an Atom entry
	topicContent struct {
		XMLName          xml.Name         `xml:"content"`
		Type             string           `xml:"type,attr"`
		TopicDescription TopicDescription `xml:"TopicDescription"`
	}

	// TopicFeed is a specialized feed containing Topic Entries
	TopicFeed struct {
		*Feed
		Entries []TopicEnvelope `xml:"entry"`
	}
)

// Subscriptions (and rules)
type (
	// FilterDescriber can transform itself into a FilterDescription
	FilterDescriber interface {
		ToFilterDescription() FilterDescription
	}

	// ActionDescriber can transform itself into a ActionDescription
	ActionDescriber interface {
		ToActionDescription() ActionDescription
	}

	// RuleDescription is the content type for Subscription Rule management requests
	RuleDescription struct {
		XMLName xml.Name `xml:"RuleDescription"`
		BaseEntityDescription
		CreatedAt *date.Time         `xml:"CreatedAt,omitempty"`
		Filter    FilterDescription  `xml:"Filter"`
		Action    *ActionDescription `xml:"Action,omitempty"`
	}
	// DefaultRuleDescription is the content type for Subscription Rule management requests
	DefaultRuleDescription struct {
		XMLName xml.Name          `xml:"DefaultRuleDescription"`
		Filter  FilterDescription `xml:"Filter"`
		Name    *string           `xml:"Name,omitempty"`
	}

	// FilterDescription describes a filter which can be applied to a subscription to filter messages from the topic.
	//
	// Subscribers can define which messages they want to receive from a topic. These messages are specified in the
	// form of one or more named subscription rules. Each rule consists of a condition that selects particular messages
	// and an action that annotates the selected message. For each matching rule condition, the subscription produces a
	// copy of the message, which may be differently annotated for each matching rule.
	//
	// Each newly created topic subscription has an initial default subscription rule. If you don't explicitly specify a
	// filter condition for the rule, the applied filter is the true filter that enables all messages to be selected
	// into the subscription. The default rule has no associated annotation action.
	FilterDescription struct {
		XMLName xml.Name `xml:"Filter"`
		CorrelationFilter
		Type               string  `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
		SQLExpression      *string `xml:"SqlExpression,omitempty"`
		CompatibilityLevel int     `xml:"CompatibilityLevel,omitempty"`
	}

	// ActionDescription describes an action upon a message that matches a filter
	//
	// With SQL filter conditions, you can define an action that can annotate the message by adding, removing, or
	// replacing properties and their values. The action uses a SQL-like expression that loosely leans on the SQL
	// UPDATE statement syntax. The action is performed on the message after it has been matched and before the message
	// is selected into the subscription. The changes to the message properties are private to the message copied into
	// the subscription.
	ActionDescription struct {
		Type                  string `xml:"http://www.w3.org/2001/XMLSchema-instance type,attr"`
		SQLExpression         string `xml:"SqlExpression"`
		RequiresPreprocessing bool   `xml:"RequiresPreprocessing"`
		CompatibilityLevel    int    `xml:"CompatibilityLevel,omitempty"`
	}

	// RuleEntity is the Azure Service Bus description of a Subscription Rule for management activities
	RuleEntity struct {
		*RuleDescription
		*Entity
	}

	// ruleContent is a specialized Subscription body for an Atom entry
	ruleContent struct {
		XMLName         xml.Name        `xml:"content"`
		Type            string          `xml:"type,attr"`
		RuleDescription RuleDescription `xml:"RuleDescription"`
	}

	RuleEnvelope struct {
		*Entry
		Content *ruleContent `xml:"content"`
	}

	// SubscriptionDescription is the content type for Subscription management requests
	SubscriptionDescription struct {
		XMLName xml.Name `xml:"SubscriptionDescription"`
		BaseEntityDescription
		LockDuration                              *string                 `xml:"LockDuration,omitempty"` // LockDuration - ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The maximum value for LockDuration is 5 minutes; the default value is 1 minute.
		RequiresSession                           *bool                   `xml:"RequiresSession,omitempty"`
		DefaultMessageTimeToLive                  *string                 `xml:"DefaultMessageTimeToLive,omitempty"`         // DefaultMessageTimeToLive - ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself.
		DeadLetteringOnMessageExpiration          *bool                   `xml:"DeadLetteringOnMessageExpiration,omitempty"` // DeadLetteringOnMessageExpiration - A value that indicates whether this queue has dead letter support when a message expires.
		DeadLetteringOnFilterEvaluationExceptions *bool                   `xml:"DeadLetteringOnFilterEvaluationExceptions,omitempty"`
		DefaultRuleDescription                    *DefaultRuleDescription `xml:"DefaultRuleDescription,omitempty"`
		MaxDeliveryCount                          *int32                  `xml:"MaxDeliveryCount,omitempty"`        // MaxDeliveryCount - The maximum delivery count. A message is automatically deadlettered after this number of deliveries. default value is 10.
		MessageCount                              *int64                  `xml:"MessageCount,omitempty"`            // MessageCount - The number of messages in the queue.
		EnableBatchedOperations                   *bool                   `xml:"EnableBatchedOperations,omitempty"` // EnableBatchedOperations - Value that indicates whether server-side batched operations are enabled.
		Status                                    *EntityStatus           `xml:"Status,omitempty"`
		ForwardTo                                 *string                 `xml:"ForwardTo,omitempty"` // ForwardTo - absolute URI of the entity to forward messages
		UserMetadata                              *string                 `xml:"UserMetadata,omitempty"`
		ForwardDeadLetteredMessagesTo             *string                 `xml:"ForwardDeadLetteredMessagesTo,omitempty"` // ForwardDeadLetteredMessagesTo - absolute URI of the entity to forward dead letter messages
		AutoDeleteOnIdle                          *string                 `xml:"AutoDeleteOnIdle,omitempty"`
		CreatedAt                                 *date.Time              `xml:"CreatedAt,omitempty"`
		UpdatedAt                                 *date.Time              `xml:"UpdatedAt,omitempty"`
		AccessedAt                                *date.Time              `xml:"AccessedAt,omitempty"`
		CountDetails                              *CountDetails           `xml:"CountDetails,omitempty"`
	}

	// SubscriptionEntity is the Azure Service Bus description of a topic Subscription for management activities
	SubscriptionEntity struct {
		*SubscriptionDescription
		*Entity
	}

	// SubscriptionFeed is a specialized feed containing Topic Subscriptions
	SubscriptionFeed struct {
		*Feed
		Entries []SubscriptionEnvelope `xml:"entry"`
	}

	// subscriptionEntryContent is a specialized Topic feed Subscription
	SubscriptionEnvelope struct {
		*Entry
		Content *subscriptionContent `xml:"content"`
	}

	// subscriptionContent is a specialized Subscription body for an Atom entry
	subscriptionContent struct {
		XMLName                 xml.Name                `xml:"content"`
		Type                    string                  `xml:"type,attr"`
		SubscriptionDescription SubscriptionDescription `xml:"SubscriptionDescription"`
	}

	// Entity is represents the most basic form of an Azure Service Bus entity.
	Entity struct {
		Name string
		ID   string
	}
)
