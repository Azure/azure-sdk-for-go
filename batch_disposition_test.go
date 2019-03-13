package servicebus

import (
	"context"
	"testing"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBatchDispositionIterator(t *testing.T) {
	count := 20
	fetched := 0
	lockIDs := []*uuid.UUID{}

	for i := count; i > 0; i-- {
		lockIDs = append(lockIDs, &uuid.UUID{})
	}

	bdi := &BatchDispositionIterator{
		LockTokenIDs: lockIDs,
	}

	assert.Equal(t, 0, bdi.cursor)

	for !bdi.Done() {
		if uuid := bdi.Next(); uuid != nil {
			fetched++
		}
	}
	assert.Equal(t, count, fetched)
}

func TestBatchDispositionUnsupportedStatus(t *testing.T) {
	status := MessageStatus(suspendedDisposition)
	id := uuid.UUID{}
	bdi := BatchDispositionIterator{
		LockTokenIDs: []*uuid.UUID{
			&id,
		},
		Status: status,
	}

	subscription := Subscription{}
	err := subscription.SendBatchDisposition(context.Background(), bdi)
	assert.EqualErrorf(t, err, "unsupported bulk disposition status \"suspended\"", err.Error())
}
