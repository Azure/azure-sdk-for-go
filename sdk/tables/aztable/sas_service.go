// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// TableSASSignatureValues is used to generate a Shared Access Signature (SAS) for an Azure Storage container or blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/constructing-a-service-sas
type TableSASSignatureValues struct {
	Version           string      `param:"sv"`  // If not specified, this defaults to SASVersion
	Protocol          SASProtocol `param:"spr"` // See the SASProtocol* constants
	StartTime         time.Time   `param:"st"`  // Not specified if IsZero
	ExpiryTime        time.Time   `param:"se"`  // Not specified if IsZero
	Permissions       string      `param:"sp"`  // Create by initializing a ContainerSASPermissions or TableSASPermissions and then call String()
	IPRange           IPRange     `param:"sip"`
	Identifier        string      `param:"si"`
	TableName         string      `param:"tn"`
	StartPartitionKey string      `param:"spk"`
	StartRowKey       string      `param:"srk"`
	EndPartitionKey   string      `param:"epk"`
	EndRowKey         string      `param:"erk"`
}

// NewSASQueryParameters uses an account's StorageAccountCredential to sign this signature values to produce
// the proper SAS query parameters.
// See: StorageAccountCredential. Compatible with both UserDelegationCredential and SharedKeyCredential
func (v TableSASSignatureValues) NewSASQueryParameters(credential *SharedKeyCredential) (SASQueryParameters, error) {
	resource := ""

	if v.Version != "" {
		//Make sure the permission characters are in the correct order
		perms := &TableSASPermissions{}
		if err := perms.Parse(v.Permissions); err != nil {
			return SASQueryParameters{}, err
		}
		v.Permissions = perms.String()
	} else if v.TableName == "" {
		// Make sure the permission characters are in the correct order
		perms := &TableSASPermissions{}
		if err := perms.Parse(v.Permissions); err != nil {
			return SASQueryParameters{}, err
		}
		v.Permissions = perms.String()
	} else {
		// Make sure the permission characters are in the correct order
		perms := &TableSASPermissions{}
		if err := perms.Parse(v.Permissions); err != nil {
			return SASQueryParameters{}, err
		}
		v.Permissions = perms.String()
	}
	if v.Version == "" {
		v.Version = SASVersion
	}
	startTime, expiryTime := FormatTimesForSASSigning(v.StartTime, v.ExpiryTime)

	signedIdentifier := v.Identifier

	p := SASQueryParameters{
		// Common SAS parameters
		version:     v.Version,
		protocol:    v.Protocol,
		startTime:   v.StartTime,
		expiryTime:  v.ExpiryTime,
		permissions: v.Permissions,
		ipRange:     v.IPRange,
		tableName:   v.TableName,

		// Table SAS parameters
		resource:   resource,
		identifier: v.Identifier,
	}

	canonicalName := "/" + "table" + "/" + credential.AccountName() + "/" + v.TableName

	// String to sign: http://msdn.microsoft.com/en-us/library/azure/dn140255.aspx
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

	signature, err := credential.ComputeHMACSHA256(stringToSign)
	p.signature = signature
	return p, err
}

// The TableSASPermissions type simplifies creating the permissions string for an Azure Storage blob SAS.
// Initialize an instance of this type and then call its String method to set TableSASSignatureValues's Permissions field.
type TableSASPermissions struct {
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
func (p TableSASPermissions) String() string {
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
func (p *TableSASPermissions) Parse(s string) error {
	*p = TableSASPermissions{} // Clear the flags
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
