// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"time"
)

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

// AccessRights represents the rights for an authorization rule.
type AccessRights int

const (
	// AccessRightsManage allows management of entities.
	AccessRightsManage AccessRights = 0
	// AccessRightsSend allows sending to entities.
	AccessRightsSend AccessRights = 1
	// AccessRightsListen allows listening to entities.
	AccessRightsListen AccessRights = 2
)

// EntityAvailabilityStatus is the availability status of the entity.
type EntityAvailabilityStatus string

const (
	// EntityAvailabilityStatusAvailable indicates the entity is available.
	EntityAvailabilityStatusAvailable EntityAvailabilityStatus = "Available"

	EntityAvailabilityStatusLimited EntityAvailabilityStatus = "Limited"

	// EntityAvailabilityStatusRenaming indicates the entity is being renamed.
	EntityAvailabilityStatusRenaming EntityAvailabilityStatus = "Renaming"

	// EntityAvailabilityStatusRestoring indicates the entity is being restored.
	EntityAvailabilityStatusRestoring EntityAvailabilityStatus = "Restoring"

	EntityAvailabilityStatusUnknown EntityAvailabilityStatus = "Unknown"
)

// AuthorizationRule encompasses access rights, metadata and a key for authentication.
type AuthorizationRule struct {
	ClaimType    string
	Rights       []AccessRights
	KeyName      string
	CreatedTime  time.Time
	ModifiedTime time.Time
}
