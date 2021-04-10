// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

// UUID is a struct wrapper around the github.com/Azure/azure-sdk-for-go/sdk/internal/uuid type.
// This wrapper type assists in the proper serialization of uuid types to the Table service.
type UUID struct {
	u uuid.UUID
}

// New returns a new uuid using RFC 4122 algorithm.
func NewTableUUID() UUID {
	return UUID{u: uuid.New()}
}

// String returns an unparsed version of the generated UUID sequence.
func (u UUID) String() string {
	return u.String()
}

// Parse parses a string formatted as "003020100-0504-0706-0809-0a0b0c0d0e0f"
// or "{03020100-0504-0706-0809-0a0b0c0d0e0f}" into a UUID.
func Parse(uuidStr string) UUID {
	return UUID{uuid.Parse(uuidStr)}
}

func (u UUID) bytes() []byte {
	return u.u[:]
}
