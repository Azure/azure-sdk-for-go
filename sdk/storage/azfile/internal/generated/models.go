// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"encoding/xml"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"net/url"
	"time"
)

type TransactionalContentSetter interface {
	SetMD5([]byte)
	// add SetCRC64() when Azure File service starts supporting it.
}

func (f *FileClientUploadRangeOptions) SetMD5(v []byte) {
	f.ContentMD5 = v
}

type SourceContentSetter interface {
	SetSourceContentCRC64(v []byte)
	// add SetSourceContentMD5() when Azure File service starts supporting it.
}

func (f *FileClientUploadRangeFromURLOptions) SetSourceContentCRC64(v []byte) {
	f.SourceContentCRC64 = v
}

// Custom MarshalXML/UnmarshalXML functions for types that need special handling.

// MarshalXML implements the xml.Marshaller interface for type Handle.
func (h Handle) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type alias Handle
	aux := &struct {
		*alias
		LastReconnectTime *dateTimeRFC1123 `xml:"LastReconnectTime"`
		OpenTime          *dateTimeRFC1123 `xml:"OpenTime"`
	}{
		alias:             (*alias)(&h),
		LastReconnectTime: (*dateTimeRFC1123)(h.LastReconnectTime),
		OpenTime:          (*dateTimeRFC1123)(h.OpenTime),
	}
	return enc.EncodeElement(aux, start)
}

// UnmarshalXML implements the xml.Unmarshaller interface for type Handle.
func (h *Handle) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias Handle
	aux := &struct {
		*alias
		Path              *StringEncoded   `xml:"Path"`
		LastReconnectTime *dateTimeRFC1123 `xml:"LastReconnectTime"`
		OpenTime          *dateTimeRFC1123 `xml:"OpenTime"`
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
		if aux.Name.Encoded != nil && *aux.Name.Encoded {
			name, err := url.QueryUnescape(*aux.Name.Content)
			if err != nil {
				return err
			}
			d.Name = to.Ptr(string(name))
		} else {
			d.Name = aux.Name.Content
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
		if aux.Name.Encoded != nil && *aux.Name.Encoded {
			name, err := url.QueryUnescape(*aux.Name.Content)
			if err != nil {
				return err
			}
			f.Name = to.Ptr(string(name))
		} else {
			f.Name = aux.Name.Content
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
		if aux.Prefix.Encoded != nil && *aux.Prefix.Encoded {
			name, err := url.QueryUnescape(*aux.Prefix.Content)
			if err != nil {
				return err
			}
			l.Prefix = to.Ptr(string(name))
		} else {
			l.Prefix = aux.Prefix.Content
		}
	}
	return nil
}
