// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
)

// AMQPMessage represents the AMQP message, as received from Service Bus.
// For details about these properties, refer to the AMQP specification:
//   https://docs.oasis-open.org/amqp/core/v1.0/os/amqp-core-messaging-v1.0-os.html#section-message-format
//
// Some fields in this struct are typed 'any', which means they will accept AMQP primitives, or in some
// cases slices and maps.
//
// AMQP simple types include:
// - int (any size), uint (any size)
// - float (any size)
// - string
// - bool
// - time.Time
type AMQPMessage struct {
	// ApplicationProperties corresponds to the "application-properties" section of an AMQP message.
	//
	// The values of the map are restricted to AMQP simple types, as listed in the comment for AMQPMessage.
	ApplicationProperties map[string]any

	// Body represents the body of an AMQP message.
	Body AMQPBody

	// DeliveryAnnotations corresponds to the "delivery-annotations" section in an AMQP message.
	//
	// The values of the map are restricted to AMQP simple types, as listed in the comment for AMQPMessage.
	DeliveryAnnotations map[any]any

	// DeliveryTag corresponds to the delivery-tag property of the TRANSFER frame
	// for this message.
	DeliveryTag []byte

	// Footer is the transport footers for this AMQP message.
	//
	// The values of the map are restricted to AMQP simple types, as listed in the comment for AMQPMessage.
	Footer map[any]any

	// Header is the transport headers for this AMQP message.
	Header *AMQPMessageHeader

	// MessageAnnotations corresponds to the message-annotations section of an AMQP message.
	//
	// The values of the map are restricted to AMQP simple types, as listed in the comment for AMQPMessage.
	MessageAnnotations map[any]any

	// Properties corresponds to the properties section of an AMQP message.
	Properties *AMQPMessageProperties

	linkName string

	// inner is the AMQP message we originally received, which contains some hidden
	// data that's needed to settle with go-amqp. We strip out most of the underlying
	// data so it's fairly minimal.
	inner *amqp.Message
}

type AMQPMessageProperties struct {
	// An absolute time when this message is considered to be expired.
	AbsoluteExpiryTime *time.Time

	// The content-encoding property is used as a modifier to the content-type.
	// When present, its value indicates what additional content encodings have been
	// applied to the application-data, and thus what decoding mechanisms need to be
	// applied in order to obtain the media-type referenced by the content-type header
	// field.
	//
	// Content-encoding is primarily used to allow a document to be compressed without
	// losing the identity of its underlying content type.
	//
	// Content-encodings are to be interpreted as per section 3.5 of RFC 2616 [RFC2616].
	// Valid content-encodings are registered at IANA [IANAHTTPPARAMS].
	//
	// The content-encoding MUST NOT be set when the application-data section is other
	// than data. The binary representation of all other application-data section types
	// is defined completely in terms of the AMQP type system.
	//
	// Implementations MUST NOT use the identity encoding. Instead, implementations
	// SHOULD NOT set this property. Implementations SHOULD NOT use the compress encoding,
	// except as to remain compatible with messages originally sent with other protocols,
	// e.g. HTTP or SMTP.
	//
	// Implementations SHOULD NOT specify multiple content-encoding values except as to
	// be compatible with messages originally sent with other protocols, e.g. HTTP or SMTP.
	ContentEncoding *string

	// The RFC-2046 [RFC2046] MIME type for the message's application-data section
	// (body). As per RFC-2046 [RFC2046] this can contain a charset parameter defining
	// the character encoding used: e.g., 'text/plain; charset="utf-8"'.
	//
	// For clarity, as per section 7.2.1 of RFC-2616 [RFC2616], where the content type
	// is unknown the content-type SHOULD NOT be set. This allows the recipient the
	// opportunity to determine the actual type. Where the section is known to be truly
	// opaque binary data, the content-type SHOULD be set to application/octet-stream.
	//
	// When using an application-data section with a section code other than data,
	// content-type SHOULD NOT be set.
	ContentType *string

	// This is a client-specific id that can be used to mark or identify messages
	// between clients.
	// The type of CorrelationID can be a uint64, UUID, []byte, or a string
	CorrelationID any

	// An absolute time when this message was created.
	CreationTime *time.Time

	// Identifies the group the message belongs to.
	GroupID *string

	// The relative position of this message within its group.
	// RFC-1982 sequence number
	GroupSequence *uint32

	// Message-id, if set, uniquely identifies a message within the message system.
	// The message producer is usually responsible for setting the message-id in
	// such a way that it is assured to be globally unique. A broker MAY discard a
	// message as a duplicate if the value of the message-id matches that of a
	// previously received message sent to the same node.
	MessageID any // uint64, UUID, []byte, or string

	// The address of the node to send replies to.
	ReplyTo *string

	// This is a client-specific id that is used so that client can send replies to this
	// message to a specific group.
	ReplyToGroupID *string

	// A common field for summary information about the message content and purpose.
	Subject *string

	// The to field identifies the node that is the intended destination of the message.
	// On any given transfer this might not be the node at the receiving end of the link.
	To *string

	// The identity of the user responsible for producing the message.
	// The client sets this value, and it MAY be authenticated by intermediaries.
	UserID []byte
}

// AMQPBody represents the body of an AMQP message.
// Only one of these fields can be used a a time. They are mutually exclusive.
type AMQPBody struct {
	// Data is encoded/decoded as multiple data sections in the body.
	Data [][]byte

	// Sequence is encoded/decoded as one or more amqp-sequence sections in the body.
	//
	// The values of the slices are are restricted to AMQP simple types, as listed in the comment for AMQPMessage.
	Sequence [][]any

	// Value is encoded/decoded as the amqp-value section in the body.
	//
	// The type of Value can be any of the AMQP simple types, as listed in the comment for AMQPMessage,
	// as well as slices or maps of AMQP simple types.
	Value any
}

// AMQPMessageHeader carries standard delivery details about the transfer
// of a message.
// See https://docs.oasis-open.org/amqp/core/v1.0/os/amqp-core-messaging-v1.0-os.html#type-header
// for more details.
type AMQPMessageHeader struct {
	// DeliveryCount is the number of unsuccessful previous attempts to deliver this message.
	DeliveryCount uint32

	Durable       bool
	FirstAcquirer bool
	Priority      uint8
	TTL           time.Duration // from milliseconds
}

// toAMQPMessage converts between our (azservicebus) AMQP message
// to the underlying message used by go-amqp.
func (am *AMQPMessage) toAMQPMessage() *amqp.Message {
	var header *amqp.MessageHeader

	if am.Header != nil {
		header = &amqp.MessageHeader{
			DeliveryCount: am.Header.DeliveryCount,
			Durable:       am.Header.Durable,
			FirstAcquirer: am.Header.FirstAcquirer,
			Priority:      am.Header.Priority,
			TTL:           am.Header.TTL,
		}
	}

	var properties *amqp.MessageProperties

	if am.Properties != nil {
		properties = &amqp.MessageProperties{
			AbsoluteExpiryTime: am.Properties.AbsoluteExpiryTime,
			ContentEncoding:    am.Properties.ContentEncoding,
			ContentType:        am.Properties.ContentType,
			CorrelationID:      am.Properties.CorrelationID,
			CreationTime:       am.Properties.CreationTime,
			GroupID:            am.Properties.GroupID,
			GroupSequence:      am.Properties.GroupSequence,
			MessageID:          am.Properties.MessageID,
			ReplyTo:            am.Properties.ReplyTo,
			ReplyToGroupID:     am.Properties.ReplyToGroupID,
			Subject:            am.Properties.Subject,
			To:                 am.Properties.To,
			UserID:             am.Properties.UserID,
		}
	} else {
		properties = &amqp.MessageProperties{}
	}

	var footer amqp.Annotations

	if am.Footer != nil {
		footer = (amqp.Annotations)(am.Footer)
	}

	return &amqp.Message{
		Annotations:           copyAnnotations(am.MessageAnnotations),
		ApplicationProperties: am.ApplicationProperties,
		Data:                  am.Body.Data,
		DeliveryAnnotations:   amqp.Annotations(am.DeliveryAnnotations),
		DeliveryTag:           am.DeliveryTag,
		Footer:                footer,
		Header:                header,
		Properties:            properties,
		Sequence:              am.Body.Sequence,
		Value:                 am.Body.Value,
	}
}

func copyAnnotations(src map[any]any) amqp.Annotations {
	if src == nil {
		return amqp.Annotations{}
	}

	dest := amqp.Annotations{}

	for k, v := range src {
		dest[k] = v
	}

	return dest
}

func newAMQPMessage(goAMQPMessage *amqp.Message) *AMQPMessage {
	var header *AMQPMessageHeader

	if goAMQPMessage.Header != nil {
		header = &AMQPMessageHeader{
			DeliveryCount: goAMQPMessage.Header.DeliveryCount,
			Durable:       goAMQPMessage.Header.Durable,
			FirstAcquirer: goAMQPMessage.Header.FirstAcquirer,
			Priority:      goAMQPMessage.Header.Priority,
			TTL:           goAMQPMessage.Header.TTL,
		}
	}

	var properties *AMQPMessageProperties

	if goAMQPMessage.Properties != nil {
		properties = &AMQPMessageProperties{
			AbsoluteExpiryTime: goAMQPMessage.Properties.AbsoluteExpiryTime,
			ContentEncoding:    goAMQPMessage.Properties.ContentEncoding,
			ContentType:        goAMQPMessage.Properties.ContentType,
			CorrelationID:      goAMQPMessage.Properties.CorrelationID,
			CreationTime:       goAMQPMessage.Properties.CreationTime,
			GroupID:            goAMQPMessage.Properties.GroupID,
			GroupSequence:      goAMQPMessage.Properties.GroupSequence,
			MessageID:          goAMQPMessage.Properties.MessageID,
			ReplyTo:            goAMQPMessage.Properties.ReplyTo,
			ReplyToGroupID:     goAMQPMessage.Properties.ReplyToGroupID,
			Subject:            goAMQPMessage.Properties.Subject,
			To:                 goAMQPMessage.Properties.To,
			UserID:             goAMQPMessage.Properties.UserID,
		}
	}

	var footer map[any]any

	if goAMQPMessage.Footer != nil {
		footer = (map[any]any)(goAMQPMessage.Footer)
	}

	return &AMQPMessage{
		MessageAnnotations:    map[any]any(goAMQPMessage.Annotations),
		ApplicationProperties: goAMQPMessage.ApplicationProperties,
		Body: AMQPBody{
			Data:     goAMQPMessage.Data,
			Sequence: goAMQPMessage.Sequence,
			Value:    goAMQPMessage.Value,
		},
		DeliveryAnnotations: map[any]any(goAMQPMessage.DeliveryAnnotations),
		DeliveryTag:         goAMQPMessage.DeliveryTag,
		Footer:              footer,
		Header:              header,
		linkName:            goAMQPMessage.LinkName(),
		Properties:          properties,
		inner:               goAMQPMessage,
	}
}
