//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

// Package checkpoints provides a CheckpointStore using Azure Blob Storage.
//
// CheckpointStore's are generally not used on their own and will be created so they
// can be passed to a [Processor] to coordinate distributed consumption of events from an event hub.
//
// See [example_processor_test.go] for an example that uses the [checkpoints.BlobStore] with
// a [Processor].
//
// [Processor]: https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs#Processor
// [example_processor_test.go]: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/messaging/azeventhubs/example_processor_test.go
package checkpoints
