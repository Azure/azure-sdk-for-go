// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package tests

import (
	"sync"
	"sync/atomic"
)

type trackingData struct {
	PartitionID string

	// atomically updated number, from the sender, telling us how many messages have been sent.
	SentCount *int64

	// tracks the received sequence numbers, with the idea that there could be duplicates
	// as partitions are stolen/rearranged. Mostly we care about # of unique events received
	// but in the near future we'll also start checking duplication as well.
	mu                      *sync.Mutex
	receivedSequenceNumbers *map[int64]int
}

type trackedCounts struct {
	Received int64
	Sent     int64
}

func newTrackingData(partitionID string) trackingData {
	var sentCount int64

	return trackingData{
		PartitionID:             partitionID,
		SentCount:               &sentCount,
		mu:                      &sync.Mutex{},
		receivedSequenceNumbers: &map[int64]int{},
	}
}

func (td trackingData) Inc(seqNumber int64) {
	td.mu.Lock()
	defer td.mu.Unlock()

	(*td.receivedSequenceNumbers)[seqNumber]++
}

func (td trackingData) Stats() trackedCounts {
	// not transactional...
	sent := atomic.LoadInt64(td.SentCount)
	received := int64(len(*td.receivedSequenceNumbers))

	return trackedCounts{
		Received: received,
		Sent:     sent,
	}
}
