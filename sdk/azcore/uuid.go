// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"fmt"
	"math/rand"
	"time"
)

// The UUID reserved variants.
const (
	reservedNCS       byte = 0x80
	reservedRFC4122   byte = 0x40
	reservedMicrosoft byte = 0x20
	reservedFuture    byte = 0x00
)

func init() {
	rand.Seed(time.Now().Unix())
}

// A UUID representation compliant with specification in RFC 4122 document.
type uuid [16]byte

// NewUUID returns a new uuid using RFC 4122 algorithm.
// It uses math/rand.Read() for obtaining the byte sequence.
func newUUID() uuid {
	u := uuid{}
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	// math/rand.Read() is no-fail so we omit any error checking.
	// NOTE: this takes a process-wide lock
	rand.Read(u[:])
	u[8] = (u[8] | reservedRFC4122) & 0x7F // u.setVariant(ReservedRFC4122)

	var version byte = 4
	u[6] = (u[6] & 0xF) | (version << 4) // u.setVersion(4)
	return u
}

// String returns an unparsed version of the generated UUID sequence.
func (u uuid) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
