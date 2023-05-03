//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"encoding/base64"
	"encoding/xml"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// Custom UnmarshalXML functions for types that need special handling.
// This is necessary since there is a break in the generated BlobItem return type for listing.

// UnmarshalXML implements the xml.Unmarshaller interface for type BlobPrefix.
func (b *BlobPrefix) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias BlobPrefix
	aux := &struct {
		*alias
		BlobName *BlobName `xml:"Name"`
	}{
		alias: (*alias)(b),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.BlobName != nil {
		if aux.BlobName.Encoded != nil && *aux.BlobName.Encoded {
			name, err := base64.StdEncoding.DecodeString(*aux.BlobName.Content)
			if err != nil {
				return err
			}
			b.Name = to.Ptr(string(name))
		} else {
			b.Name = aux.BlobName.Content
		}
	}
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type BlobItem.
func (b *BlobItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias BlobItem
	aux := &struct {
		*alias
		BlobName   *BlobName            `xml:"Name"`
		Metadata   additionalProperties `xml:"Metadata"`
		OrMetadata additionalProperties `xml:"OrMetadata"`
	}{
		alias: (*alias)(b),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	b.Metadata = (map[string]*string)(aux.Metadata)
	b.OrMetadata = (map[string]*string)(aux.OrMetadata)
	if aux.BlobName != nil {
		if aux.BlobName.Encoded != nil && *aux.BlobName.Encoded {
			name, err := base64.StdEncoding.DecodeString(*aux.BlobName.Content)
			if err != nil {
				return err
			}
			b.Name = to.Ptr(string(name))
		} else {
			b.Name = aux.BlobName.Content
		}
	}
	return nil
}
