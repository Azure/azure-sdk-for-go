// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"context"
	"fmt"
	"io"

	eventhub "github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/persist"
)

type persistRecord struct {
	namespace     string
	name          string
	consumerGroup string
	partitionID   string
	checkpoint    persist.Checkpoint
}

type batchWriter struct {
	persister persist.CheckpointPersister
	writer    io.Writer

	batchSize      int
	batch          []string
	persistRecords []*persistRecord
	flushed        *persistRecord
}

var batchSize = 10

// NewBatchWriter creates an object that can be used as both a `persist.CheckpointPersister` and an Event Hubs Event Handler `batchWriter.HandleEvent`
func NewBatchWriter(persister persist.CheckpointPersister, writer io.Writer) (*batchWriter, error) {
	return &batchWriter{
		persister:      persister,
		writer:         writer,
		batchSize:      batchSize,
		batch:          make([]string, 0, batchSize),
		persistRecords: make([]*persistRecord, 0, batchSize),
	}, nil
}

// Read reads the last checkpoint
func (w batchWriter) Read(namespace, name, consumerGroup, partitionID string) (persist.Checkpoint, error) {
	return w.persister.Read(namespace, name, consumerGroup, partitionID)
}

// Write will write the last checkpoint of the last event flushed and record persist records for future use
func (w *batchWriter) Write(namespace, name, consumerGroup, partitionID string, checkpoint persist.Checkpoint) error {
	var err error
	if w.flushed != nil {
		r := w.flushed
		err = w.persister.Write(r.namespace, r.name, r.consumerGroup, r.partitionID, r.checkpoint)
		if err != nil {
			w.flushed = nil
		}
	}
	w.persistRecords = append(w.persistRecords, &persistRecord{
		namespace:     namespace,
		name:          name,
		consumerGroup: consumerGroup,
		partitionID:   partitionID,
		checkpoint:    checkpoint,
	})
	return err
}

// HandleEvent will handle Event Hubs Events
// If the length of the batch buffer has reached the max batchSize, the buffer will be flushed before appending the new event
// If flush fails and it hasn't made space in the buffer, the flush error will be returned to the caller
func (w *batchWriter) HandleEvent(ctx context.Context, event *eventhub.Event) error {
	if len(w.batch) >= batchSize {
		err := w.Flush(ctx)
		// If we received an error flushing and still don't have room in the buffer return the error
		if err != nil && len(w.batch) >= batchSize {
			return err
		}
	}
	// Append the event to the buffer if we have room for it
	w.batch = append(w.batch, string(event.Data))
	return nil
}

// Flush flushes the buffer to the given io.Writer
// Post-condition:
//   error == nil: buffer has been flushed successfully, buffer has been replaced with a new buffer
//   error != nil: some or no events have been flushed, buffer contains only events that failed to flush
func (w *batchWriter) Flush(ctx context.Context) error {
	for i, s := range w.batch {
		_, err := fmt.Fprintln(w.writer, s)
		if err != nil {
			w.batch = w.batch[i:]
			w.persistRecords = w.persistRecords[i:]
			return err
		}
		w.flushed = w.persistRecords[i]
	}
	w.batch = make([]string, 0, batchSize)
	w.persistRecords = make([]*persistRecord, 0, batchSize)
	return nil
}
