// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"encoding/xml"
	"net/url"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime/datetime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

type TransactionalContentSetter interface {
	SetMD5([]byte)
	SetCRC64([]byte)
	SetStructuredBody(bodyType string, contentLength int64)
}

func (f *FileClientUploadRangeOptions) SetMD5(v []byte) {
	f.ContentMD5 = v
}

func (f *FileClientUploadRangeOptions) SetCRC64([]byte) {
	// no-op: Azure Files does not yet support transactional CRC64 headers on UploadRange
}

func (f *FileClientUploadRangeOptions) SetStructuredBody(bodyType string, contentLength int64) {
	f.StructuredBodyType = &bodyType
	f.StructuredContentLength = &contentLength
}

func (f *FileClientCreateOptions) SetMD5(v []byte) {
	f.ContentMD5 = v
}

func (f *FileClientCreateOptions) SetCRC64([]byte) {
	// no-op: Azure Files does not support transactional CRC64 headers on Create
}

func (f *FileClientCreateOptions) SetStructuredBody(bodyType string, contentLength int64) {
	f.StructuredBodyType = &bodyType
	f.StructuredContentLength = &contentLength
}

type SourceContentSetter interface {
	SetSourceContentCRC64(v []byte)
	// add SetSourceContentMD5() when Azure File service starts supporting it.
}

func (f *FileClientUploadRangeFromURLOptions) SetSourceContentCRC64(v []byte) {
	f.SourceContentCRC64 = v
}

type StringEncoded struct {
	Content *string `xml:",chardata"`
	Encoded *bool   `xml:"Encoded,attr"`
}

// Custom MarshalXML/UnmarshalXML functions for types that need special handling.

// MarshalXML implements the xml.Marshaller interface for type Handle.
func (h Handle) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type alias Handle
	aux := &struct {
		*alias
		LastReconnectTime *datetime.RFC1123 `xml:"LastReconnectTime"`
		OpenTime          *datetime.RFC1123 `xml:"OpenTime"`
	}{
		alias:             (*alias)(&h),
		LastReconnectTime: (*datetime.RFC1123)(h.LastReconnectTime),
		OpenTime:          (*datetime.RFC1123)(h.OpenTime),
	}
	return enc.EncodeElement(aux, start)
}

// UnmarshalXML implements the xml.Unmarshaller interface for type Handle.
func (h *Handle) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias Handle
	aux := &struct {
		*alias
		Path              *StringEncoded    `xml:"Path"`
		LastReconnectTime *datetime.RFC1123 `xml:"LastReconnectTime"`
		OpenTime          *datetime.RFC1123 `xml:"OpenTime"`
	}{
		alias: (*alias)(h),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	h.LastReconnectTime = (*time.Time)(aux.LastReconnectTime)
	h.OpenTime = (*time.Time)(aux.OpenTime)
	if aux.Path != nil {
		if aux.Path.Encoded != nil && *aux.Path.Encoded {
			name, err := url.QueryUnescape(*aux.Path.Content)
			if err != nil {
				return err
			}
			h.Path = to.Ptr(string(name))
		} else {
			h.Path = aux.Path.Content
		}
	}
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type Directory.
func (d *Directory) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias Directory
	aux := &struct {
		*alias
		Name *StringEncoded `xml:"Name"`
	}{
		alias: (*alias)(d),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.Name != nil {
		d.Name = aux.Name.Content
		// A malformed encoding does not fail the page; the name is returned as sent by the service.
		if aux.Name.Encoded != nil && *aux.Name.Encoded && aux.Name.Content != nil {
			if name, err := url.QueryUnescape(*aux.Name.Content); err == nil {
				d.Name = &name
			}
		}
	}
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type File.
func (f *File) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias File
	aux := &struct {
		*alias
		Name *StringEncoded `xml:"Name"`
	}{
		alias: (*alias)(f),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.Name != nil {
		f.Name = aux.Name.Content
		// A malformed encoding does not fail the page; the name is returned as sent by the service.
		if aux.Name.Encoded != nil && *aux.Name.Encoded && aux.Name.Content != nil {
			if name, err := url.QueryUnescape(*aux.Name.Content); err == nil {
				f.Name = &name
			}
		}
	}
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type HardLink.
// The service percent-encodes a FileName that cannot be represented in XML and marks it with Encoded="true".
func (h *HardLink) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias HardLink
	aux := &struct {
		*alias
		FileName *StringEncoded `xml:"FileName"`
	}{
		alias: (*alias)(h),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.FileName != nil {
		if aux.FileName.Encoded != nil && *aux.FileName.Encoded && aux.FileName.Content != nil {
			name, err := url.QueryUnescape(*aux.FileName.Content)
			if err != nil {
				return err
			}
			h.FileName = to.Ptr(name)
		} else {
			h.FileName = aux.FileName.Content
		}
	}
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type ListFilesAndDirectoriesSegmentResponse.
func (l *ListFilesAndDirectoriesSegmentResponse) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias ListFilesAndDirectoriesSegmentResponse
	aux := &struct {
		*alias
		Prefix *StringEncoded `xml:"Prefix"`
	}{
		alias: (*alias)(l),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.Prefix != nil {
		l.Prefix = aux.Prefix.Content
		// A malformed encoding does not fail the page; the prefix is returned as sent by the service.
		if aux.Prefix.Encoded != nil && *aux.Prefix.Encoded && aux.Prefix.Content != nil {
			if prefix, err := url.QueryUnescape(*aux.Prefix.Content); err == nil {
				l.Prefix = &prefix
			}
		}
	}
	return nil
}
