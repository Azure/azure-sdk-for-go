package servicebus

import (
	"time"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/mitchellh/mapstructure"
	"pack.ag/amqp"
)

func (suite *serviceBusSuite) TestMapStructureEncode() {
	sp := new(SystemProperties)
	m, err := encodeStructureToMap(sp)
	if suite.NoError(err) {
		suite.Len(m, 0)
	}

	now := time.Now()
	pID := int16(1)
	sp.LockedUntil = &now
	m, err = encodeStructureToMap(sp)
	if suite.NoError(err) {
		suite.Equal(now, m["x-opt-locked-until"])
		suite.Len(m, 1)
	}

	sp.PartitionKey = to.StringPtr("foo")
	sp.PartitionID = &pID
	sp.SequenceNumber = to.Int64Ptr(42)
	sp.EnqueuedTime = &now
	sp.EnqueuedSequenceNumber = to.Int64Ptr(43)
	sp.DeadLetterSource = to.StringPtr("bar")
	sp.ScheduledEnqueueTime = &now
	sp.ViaPartitionKey = to.StringPtr("via")

	m, err = encodeStructureToMap(sp)
	if suite.NoError(err) {
		var sp2 SystemProperties
		err = mapstructure.Decode(&m, &sp2)
		if suite.NoError(err) {
			suite.Equal(sp, &sp2)
		}
	}
}

func (suite *serviceBusSuite) TestMessageToAMQPMessage() {
	sequence := uint32(1234)
	d := 30 * time.Second
	until := time.Now().Add(d)
	pID := int16(12)
	id, err := uuid.NewV4()
	suite.NoError(err)
	msg := Message{
		ContentType:    "application/json",
		CorrelationID:  "1",
		Data:           []byte("foo"),
		GroupID:        to.StringPtr("12"),
		GroupSequence:  &sequence,
		ID:             "123",
		Label:          "subject",
		ReplyTo:        "replyTo",
		ReplyToGroupID: "replyToGroupID",
		To:             "to",
		TTL:            &d,
		LockToken:      &id,
		SystemProperties: &SystemProperties{
			LockedUntil:            &until,
			SequenceNumber:         to.Int64Ptr(1),
			PartitionID:            &pID,
			PartitionKey:           to.StringPtr("key"),
			EnqueuedTime:           &until,
			DeadLetterSource:       to.StringPtr("deadLetterSource"),
			ScheduledEnqueueTime:   &until,
			EnqueuedSequenceNumber: to.Int64Ptr(1),
			ViaPartitionKey:        to.StringPtr("via"),
		},
		UserProperties: map[string]interface{}{
			"test": "foo",
		},
	}
	aMsg, err := msg.toMsg()
	if suite.NoError(err) {
		suite.Equal(msg.ID, aMsg.Properties.MessageID, "message id")
		suite.Equal(*msg.GroupID, aMsg.Properties.GroupID, "groupID")
		suite.Equal(*msg.GroupSequence, aMsg.Properties.GroupSequence, "GroupSequence")
		suite.Equal(msg.CorrelationID, aMsg.Properties.CorrelationID, "CorrelationID")
		suite.Equal(msg.ContentType, aMsg.Properties.ContentType, "ContentType")
		suite.Equal(msg.Data, aMsg.Data[0], "Data")
		suite.Equal(msg.Label, aMsg.Properties.Subject, "Label")
		suite.Equal(msg.ReplyTo, aMsg.Properties.ReplyTo, "ReplyTo")
		suite.Equal(msg.ReplyToGroupID, aMsg.Properties.ReplyToGroupID, "ReplyToGroupID")
		suite.Equal(msg.To, aMsg.Properties.To, "To")
		suite.Equal(*msg.TTL, aMsg.Header.TTL, "TTL")

		suite.Equal(*msg.LockToken, aMsg.DeliveryAnnotations["x-opt-lock-token"])

		sysPropMap, err := encodeStructureToMap(msg.SystemProperties)
		if suite.NoError(err) {
			for key, val := range sysPropMap {
				suite.Equal(val, aMsg.Annotations[key], key)
			}
		}
	}
}

func (suite *serviceBusSuite) TestAMQPMessageToMessage() {
	d := 30 * time.Second
	until := time.Now().Add(d)
	pID := int16(12)
	id, err := uuid.NewV4()
	suite.NoError(err)
	aMsg := &amqp.Message{
		Properties: &amqp.MessageProperties{
			MessageID:          "messageID",
			To:                 "to",
			Subject:            "subject",
			ReplyTo:            "replyTo",
			ReplyToGroupID:     "replyToGroupID",
			CorrelationID:      "correlationID",
			ContentType:        "contentType",
			ContentEncoding:    "contentEncoding",
			AbsoluteExpiryTime: until,
			CreationTime:       until,
			GroupID:            "groupID",
			GroupSequence:      uint32(1),
		},
		DeliveryAnnotations: amqp.Annotations{
			"x-opt-lock-token": amqp.UUID(id),
		},
		Annotations: amqp.Annotations{
			"x-opt-locked-until":            until,
			"x-opt-sequence-number":         int64(1),
			"x-opt-partition-id":            pID,
			"x-opt-partition-key":           "key",
			"x-opt-enqueued-time":           until,
			"x-opt-deadletter-source":       "deadLetterSource",
			"x-opt-scheduled-enqueue-time":  until,
			"x-opt-enqueue-sequence-number": int64(1),
			"x-opt-via-partition-key":       "via",
		},
		ApplicationProperties: map[string]interface{}{
			"test": "foo",
		},
		Header: &amqp.MessageHeader{
			TTL: d,
		},
		Data: [][]byte{[]byte("foo")},
	}

	msg, err := messageFromAMQPMessage(aMsg)
	if suite.NoError(err) {
		suite.Equal(msg.ID, aMsg.Properties.MessageID, "messageID")
		suite.Equal(*msg.GroupSequence, aMsg.Properties.GroupSequence, "groupSequence")
		suite.Equal(*msg.GroupID, aMsg.Properties.GroupID, "groupID")
		suite.Equal(msg.ContentType, aMsg.Properties.ContentType, "contentType")
		suite.Equal(msg.CorrelationID, aMsg.Properties.CorrelationID, "correlation")
		suite.Equal(msg.ReplyToGroupID, aMsg.Properties.ReplyToGroupID, "replyToGroupID")
		suite.Equal(msg.ReplyTo, aMsg.Properties.ReplyTo, "replyTo")
		suite.Equal(*msg.TTL, aMsg.Header.TTL, "ttl")
		suite.Equal(msg.Label, aMsg.Properties.Subject, "subject")
		suite.Equal(msg.To, aMsg.Properties.To, "to")
		suite.Equal(msg.Data, aMsg.Data[0], "data")

		suite.EqualValues(aMsg.DeliveryAnnotations["x-opt-lock-token"], *msg.LockToken, "lockToken")

		sysPropMap, err := encodeStructureToMap(msg.SystemProperties)
		if suite.NoError(err) {
			for key, val := range sysPropMap {
				suite.Equal(val, aMsg.Annotations[key], key)
			}
		}

	}
}
