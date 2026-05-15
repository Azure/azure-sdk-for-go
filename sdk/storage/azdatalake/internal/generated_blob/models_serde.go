// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated_blob

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

// MarshalXML implements the xml.Marshaller interface for type UserDelegationKey.
func (u UserDelegationKey) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type alias UserDelegationKey
	aux := &struct {
		*alias
		SignedExpiry *datetime.RFC3339 `xml:"SignedExpiry"`
		SignedStart  *datetime.RFC3339 `xml:"SignedStart"`
	}{
		alias:        (*alias)(&u),
		SignedExpiry: (*datetime.RFC3339)(u.SignedExpiry),
		SignedStart:  (*datetime.RFC3339)(u.SignedStart),
	}
	return enc.EncodeElement(aux, start)
}

// UnmarshalXML implements the xml.Unmarshaller interface for type UserDelegationKey.
func (u *UserDelegationKey) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias UserDelegationKey
	aux := &struct {
		*alias
		SignedExpiry *datetime.RFC3339 `xml:"SignedExpiry"`
		SignedStart  *datetime.RFC3339 `xml:"SignedStart"`
	}{
		alias: (*alias)(u),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.SignedExpiry != nil && !(*time.Time)(aux.SignedExpiry).IsZero() {
		u.SignedExpiry = (*time.Time)(aux.SignedExpiry)
	}
	if aux.SignedStart != nil && !(*time.Time)(aux.SignedStart).IsZero() {
		u.SignedStart = (*time.Time)(aux.SignedStart)
	}
	return nil
}

// MarshalXML implements the xml.Marshaller interface for type ListFileSystemsSegmentResponse.
func (l ListFileSystemsSegmentResponse) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type alias ListFileSystemsSegmentResponse
	aux := &struct {
		*alias
		FileSystemItems *[]*FileSystemItem `xml:"Containers>Container"`
	}{
		alias: (*alias)(&l),
	}
	if l.FileSystemItems != nil {
		aux.FileSystemItems = &l.FileSystemItems
	}
	return enc.EncodeElement(aux, start)
}

// UnmarshalXML implements the xml.Unmarshaller interface for type FileSystemItem.
func (c *FileSystemItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias FileSystemItem
	aux := &struct {
		*alias
		Metadata additionalProperties `xml:"Metadata"`
	}{
		alias: (*alias)(c),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	c.Metadata = (map[string]*string)(aux.Metadata)
	return nil
}

// MarshalXML implements the xml.Marshaller interface for type FileSystemProperties.
func (c FileSystemProperties) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type alias FileSystemProperties
	aux := &struct {
		*alias
		DeletedTime  *datetime.RFC1123 `xml:"DeletedTime"`
		LastModified *datetime.RFC1123 `xml:"Last-Modified"`
	}{
		alias:        (*alias)(&c),
		DeletedTime:  (*datetime.RFC1123)(c.DeletedTime),
		LastModified: (*datetime.RFC1123)(c.LastModified),
	}
	return enc.EncodeElement(aux, start)
}

// UnmarshalXML implements the xml.Unmarshaller interface for type FileSystemProperties.
func (c *FileSystemProperties) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias FileSystemProperties
	aux := &struct {
		*alias
		DeletedTime  *datetime.RFC1123 `xml:"DeletedTime"`
		LastModified *datetime.RFC1123 `xml:"Last-Modified"`
	}{
		alias: (*alias)(c),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.DeletedTime != nil && !(*time.Time)(aux.DeletedTime).IsZero() {
		c.DeletedTime = (*time.Time)(aux.DeletedTime)
	}
	if aux.LastModified != nil && !(*time.Time)(aux.LastModified).IsZero() {
		c.LastModified = (*time.Time)(aux.LastModified)
	}
	return nil
}

type additionalProperties map[string]*string

// UnmarshalXML implements the xml.Unmarshaler interface for additionalProperties.
func (ap *additionalProperties) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tokName := ""
	tokValue := ""
	for {
		t, err := d.Token()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			tokName = strings.ToLower(tt.Name.Local)
			tokValue = ""
		case xml.CharData:
			if tokName == "" {
				continue
			}
			tokValue = string(tt)
		case xml.EndElement:
			if tokName == "" {
				continue
			}
			if *ap == nil {
				*ap = additionalProperties{}
			}
			(*ap)[tokName] = to.Ptr(tokValue)
			tokName = ""
		}
	}
	return nil
}
