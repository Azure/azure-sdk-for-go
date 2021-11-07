// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package admin

// EntityStatus represents the current status of the entity.
type EntityStatus string

const (
	// EntityStatusActive indicates an entity can be used for sending and receiving.
	EntityStatusActive EntityStatus = "Active"
	// EntityStatusDisabled indicates an entity cannot be used for sending or receiving.
	EntityStatusDisabled EntityStatus = "Disabled"
	// EntityStatusSendDisabled indicates that an entity cannot be used for sending.
	EntityStatusSendDisabled EntityStatus = "SendDisabled"
	// EntityStatusReceiveDisabled indicates that an entity cannot be used for receiving.
	EntityStatusReceiveDisabled EntityStatus = "ReceiveDisabled"
)
