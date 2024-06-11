// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// SASSignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Table instance.
// For more information, see https://learn.microsoft.com/rest/api/storageservices/constructing-a-service-sas
type SASSignatureValues struct {
	Version           string      // If not specified, this defaults to SASVersion
	Protocol          SASProtocol // See the SASProtocol* constants
	StartTime         time.Time   // Not specified if IsZero
	ExpiryTime        time.Time   // Not specified if IsZero
	Permissions       string      // Create by initializing a ContainerSASPermissions or TableSASPermissions and then call String()
	IPRange           IPRange
	Identifier        string
	TableName         string
	StartPartitionKey string
	StartRowKey       string
	EndPartitionKey   string
	EndRowKey         string
}

// Sign uses an account's SharedKeyCredential to sign this signature values to produce
// the proper SAS string.
func (v SASSignatureValues) Sign(credential *SharedKeyCredential) (string, error) {
	resource := ""

	// Make sure the permission characters are in the correct order
	perms := &SASPermissions{}
	if err := perms.Parse(v.Permissions); err != nil {
		return "", err
	}
	v.Permissions = perms.String()

	if v.Version == "" {
		v.Version = SASVersion
	}
	startTime, expiryTime := FormatTimesForSASSigning(v.StartTime, v.ExpiryTime)

	signedIdentifier := v.Identifier

	lowerCaseTableName := strings.ToLower(v.TableName)

	p := SASQueryParameters{
		// Common SAS parameters
		version:     v.Version,
		protocol:    v.Protocol,
		startTime:   v.StartTime,
		expiryTime:  v.ExpiryTime,
		permissions: v.Permissions,
		ipRange:     v.IPRange,
		tableName:   lowerCaseTableName,

		// Table SAS parameters
		resource:   resource,
		identifier: v.Identifier,
	}

	canonicalName := "/" + "table" + "/" + credential.AccountName() + "/" + lowerCaseTableName

	// String to sign: https://learn.microsoft.com/rest/api/storageservices/create-service-sas
	stringToSign := strings.Join([]string{
		v.Permissions,
		startTime,
		expiryTime,
		canonicalName,
		signedIdentifier,
		v.IPRange.String(),
		string(v.Protocol),
		v.Version,
		v.StartPartitionKey,
		v.StartRowKey,
		v.EndPartitionKey,
		v.EndRowKey,
	},
		"\n",
	)

	signature, err := credential.computeHMACSHA256(stringToSign)
	p.signature = signature
	return p.Encode(), err
}

// SASPermissions simplifies creating the permissions string for an Azure Table.
// Initialize an instance of this type and then call its String method to set TableSASSignatureValues's Permissions field.
type SASPermissions struct {
	Read              bool
	Add               bool
	Update            bool
	Delete            bool
	StartPartitionKey string
	StartRowKey       string
	EndPartitionKey   string
	EndRowKey         string
}

// String produces the SAS permissions string for an Azure Storage blob.
// Call this method to set TableSASSignatureValues's Permissions field.
func (p SASPermissions) String() string {
	var b bytes.Buffer
	if p.Read {
		b.WriteRune('r')
	}
	if p.Add {
		b.WriteRune('a')
	}
	if p.Update {
		b.WriteRune('u')
	}
	if p.Delete {
		b.WriteRune('d')
	}
	return b.String()
}

// Parse initializes the TableSASPermissions's fields from a string.
func (p *SASPermissions) Parse(s string) error {
	*p = SASPermissions{} // Clear the flags
	for _, r := range s {
		switch r {
		case 'r':
			p.Read = true
		case 'a':
			p.Add = true
		case 'u':
			p.Update = true
		case 'd':
			p.Delete = true
		default:
			return fmt.Errorf("invalid permission: '%v'", r)
		}
	}
	return nil
}
