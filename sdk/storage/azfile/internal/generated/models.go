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

// FileType - The type of a listed file system item. It is not defined by the service; it is derived from the XML element
// in which the item was returned by the List Directories and Files operation.
type FileType string

const (
	// FileTypeRegular - Regular file.
	FileTypeRegular FileType = "Regular"
	// FileTypeDirectory - Directory.
	FileTypeDirectory FileType = "Directory"
	// FileTypeSymLink - Symbolic link.
	FileTypeSymLink FileType = "SymLink"
	// FileTypeBlockDevice - Block device.
	FileTypeBlockDevice FileType = "BlockDevice"
	// FileTypeCharacterDevice - Character device.
	FileTypeCharacterDevice FileType = "CharacterDevice"
	// FileTypeSocket - Socket.
	FileTypeSocket FileType = "Socket"
	// FileTypeFifo - FIFO.
	FileTypeFifo FileType = "Fifo"
)

// PossibleFileTypeValues returns the possible values for the FileType const type.
func PossibleFileTypeValues() []FileType {
	return []FileType{
		FileTypeRegular,
		FileTypeDirectory,
		FileTypeSymLink,
		FileTypeBlockDevice,
		FileTypeCharacterDevice,
		FileTypeSocket,
		FileTypeFifo,
	}
}

// FileType returns FileTypeDirectory. It is intended to be used for NFS shares.
func (d *Directory) FileType() FileType { return FileTypeDirectory }

// FileType returns FileTypeRegular. It is intended to be used for NFS shares.
func (f *File) FileType() FileType { return FileTypeRegular }

// FileType returns FileTypeSymLink.
func (s *SymLinkItem) FileType() FileType { return FileTypeSymLink }

// FileType returns FileTypeBlockDevice.
func (b *BlockDeviceItem) FileType() FileType { return FileTypeBlockDevice }

// FileType returns FileTypeCharacterDevice.
func (c *CharDeviceItem) FileType() FileType { return FileTypeCharacterDevice }

// FileType returns FileTypeFifo.
func (f *FifoItem) FileType() FileType { return FileTypeFifo }

// FileType returns FileTypeSocket.
func (s *SocketItem) FileType() FileType { return FileTypeSocket }

// decodeStringEncoded URL-decodes s.Content in place if the service returned it encoded, and sets s.Encoded to false.
// If s.Content is not a valid URL encoding, it is left as returned by the service with s.Encoded still true, so that a
// single malformed name does not fail the whole list response.
func decodeStringEncoded(s *StringEncoded) {
	if s == nil || s.Content == nil || s.Encoded == nil || !*s.Encoded {
		return
	}
	content, err := url.QueryUnescape(*s.Content)
	if err != nil {
		return
	}
	s.Content = &content
	s.Encoded = to.Ptr(false)
}

// UnmarshalXML implements the xml.Unmarshaller interface for type SymLinkItem.
func (s *SymLinkItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias SymLinkItem
	if err := dec.DecodeElement((*alias)(s), &start); err != nil {
		return err
	}
	decodeStringEncoded(s.Name)
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type BlockDeviceItem.
func (b *BlockDeviceItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias BlockDeviceItem
	if err := dec.DecodeElement((*alias)(b), &start); err != nil {
		return err
	}
	decodeStringEncoded(b.Name)
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type CharDeviceItem.
func (c *CharDeviceItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias CharDeviceItem
	if err := dec.DecodeElement((*alias)(c), &start); err != nil {
		return err
	}
	decodeStringEncoded(c.Name)
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type FifoItem.
func (f *FifoItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias FifoItem
	if err := dec.DecodeElement((*alias)(f), &start); err != nil {
		return err
	}
	decodeStringEncoded(f.Name)
	return nil
}

// UnmarshalXML implements the xml.Unmarshaller interface for type SocketItem.
func (s *SocketItem) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias SocketItem
	if err := dec.DecodeElement((*alias)(s), &start); err != nil {
		return err
	}
	decodeStringEncoded(s.Name)
	return nil
}
