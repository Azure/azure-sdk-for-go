// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bytes"
	"crypto/sha1"
	"encoding/pem"
	"errors"
	"fmt"
)

// fingerprint type wraps a byte slice that contains the corresponding SHA-1 fingerprint for the client's certificate
type fingerprint []byte

// String represents the fingerprint digest as a series of
// colon-delimited hexadecimal octets.
func (f fingerprint) String() string {
	var buf bytes.Buffer
	for i, b := range f {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02x", b)
	}
	return buf.String()
}

// spkiFingerprint calculates the fingerprint of the certificate based on it's Subject Public Key Info with the SHA-1
// signing algorithm.
func spkiFingerprint(pemBytes []byte) (fingerprint, error) {
	// Get first block of PEM file
	data, rest := pem.Decode([]byte(pemBytes))
	const certificateBlock = "CERTIFICATE"
	if data.Type != certificateBlock {
		for len(rest) > 0 {
			data, rest = pem.Decode(rest)
			if data.Type == certificateBlock {
				// Sign the CERTIFICATE block with SHA1
				h := sha1.New()
				h.Write(data.Bytes)

				return fingerprint(h.Sum(nil)), nil
			}
		}
		return nil, errors.New("Cannot find CERTIFICATE in file")
	}
	h := sha1.New()
	h.Write(data.Bytes)

	return fingerprint(h.Sum(nil)), nil
}
