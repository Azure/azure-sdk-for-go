//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

type Snapshot struct {
	// REQUIRED; A list of filters used to filter the key-values included in the snapshot.
	Filters []KeyValueFilter

	// The composition type describes how the key-values within the snapshot are composed. The 'key' composition type ensures
	// there are no two key-values containing the same key. The 'key_label' composition
	// type ensures there are no two key-values containing the same key and label.
	CompositionType *CompositionType

	// The amount of time, in seconds, that a snapshot will remain in the archived state before expiring. This property is only
	// writable during the creation of a snapshot. If not specified, the default
	// lifetime of key-value revisions will be used.
	RetentionPeriod *int64

	// The tags of the snapshot.
	Tags map[string]*string

	// READ-ONLY; The time that the snapshot was created.
	Created *time.Time

	// READ-ONLY; A value representing the current state of the snapshot.
	ETag *azcore.ETag

	// READ-ONLY; The time that the snapshot will expire.
	Expires *time.Time

	// READ-ONLY; The amount of key-values in the snapshot.
	ItemsCount *int64

	// READ-ONLY; The name of the snapshot.
	Name *string

	// READ-ONLY; The size in bytes of the snapshot.
	Size *int64

	// READ-ONLY; The current status of the snapshot.
	Status *SnapshotStatus
}
