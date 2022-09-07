//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package checkpoints provides a CheckpointStore that uses Azure Blob Storage.
// This checkpoint store can be used with the azeventhubs.Processor type to
// coordinate distributed consumption of events from an event hub.
package checkpoints
