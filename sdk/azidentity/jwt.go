// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

// headerJWT type contains the fields necessary to create a JSON Web Token including the x5t field which must contain a x.509 certificate thumbprint
type headerJWT struct {
	Typ string `json:"typ"`
	Alg string `json:"alg"`
	X5t string `json:"x5t"`
}

// payloadJWT type contains all fields that are necessary when creating a JSON Web Token payload section
type payloadJWT struct {
	JTI string `json:"jti"`
	AUD string `json:"aud"`
	ISS string `json:"iss"`
	SUB string `json:"sub"`
	NBF int64  `json:"nbf"`
	EXP int64  `json:"exp"`
}

// createClientAssertionJWT build the JWT header, payload and signature,
// then returns a string for the JWT assertion
func createClientAssertionJWT(clientID string, audience string, clientCertificate string) (string, error) {
	cert, err := os.Open(clientCertificate)
	defer cert.Close()
	if err != nil { // here the error from os.Open is descriptive enough to be returned directly
		return "", err
	}
	pemFileInfo, err := cert.Stat()
	if err != nil {
		return "", fmt.Errorf("Getting certificate file info: %w", err)
	}
	size := pemFileInfo.Size()

	pemBytes := make([]byte, size)
	buffer := bufio.NewReader(cert)
	_, err = buffer.Read(pemBytes)
	if err != nil {
		return "", fmt.Errorf("Read PEM file bytes: %w", err)
	}
	fingerprint, err := spkiFingerprint(pemBytes)
	if err != nil {
		return "", err
	}
	headerData := headerJWT{
		Typ: "JWT",
		Alg: "RS256",
		X5t: base64.RawURLEncoding.EncodeToString(fingerprint),
	}

	headerJSON, err := json.Marshal(headerData)
	if err != nil {
		return "", fmt.Errorf("Marshal headerJWT: %w", err)
	}
	header := base64.RawURLEncoding.EncodeToString(headerJSON)

	payloadData := payloadJWT{
		JTI: uuid.New().String(),
		AUD: audience,
		ISS: clientID,
		SUB: clientID,
		NBF: time.Now().Unix(),
		EXP: time.Now().Add(30 * time.Minute).Unix(),
	}

	payloadJSON, err := json.Marshal(payloadData)
	if err != nil {
		return "", fmt.Errorf("Marshal payloadJWT: %w", err)
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	result := header + "." + payload
	hashed := []byte(result)
	hashedSum := sha256.Sum256(hashed)
	cryptoRand := rand.Reader

	privateKey, err := getPrivateKey(pemBytes)
	if err != nil {
		return "", err
	}

	signed, err := rsa.SignPKCS1v15(cryptoRand, privateKey, crypto.SHA256, hashedSum[:])
	if err != nil {
		return "", err
	}

	signature := base64.RawURLEncoding.EncodeToString(signed)

	return result + "." + signature, nil
}
