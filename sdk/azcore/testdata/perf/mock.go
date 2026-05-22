// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const defaultItemCount = 25

// Types below are meant to reflect Storage Blob, which is XML but JSON is supported for comparison.

type list struct {
	Name      *string             `json:"name" xml:"Name"`
	Container *listItemsContainer `json:"items" xml:"Items"`
	Next      *string             `json:"next" xml:"Next"`
}

type listItemsContainer struct {
	Items []*listItems `json:"items" xml:"Item"`
}

type listItems struct {
	Name       *string             `json:"name" xml:"Name"`
	Properties *listItemProperties `json:"properties" xml:"Properties"`
}

type listItemProperties struct {
	ETag         *azcore.ETag `json:"etag" xml:"Etag"`
	CreationTime *time.Time   `json:"creationTime" xml:"Creation-Time"`
	LastModified *time.Time   `json:"lastModified" xml:"Last-Modified"`
	ContentMD5   []byte       `json:"contentMD5" xml:"Content-MD5"`
}

func (l listItemProperties) MarshalJSON() ([]byte, error) {
	type alias listItemProperties
	aux := &struct {
		*alias
		ContentMD5   *string          `xml:"Content-MD5"`
		CreationTime *dateTimeRFC1123 `xml:"Creation-Time"`
		LastModified *dateTimeRFC1123 `xml:"Last-Modified"`
	}{
		alias:        (*alias)(&l),
		CreationTime: (*dateTimeRFC1123)(l.CreationTime),
		LastModified: (*dateTimeRFC1123)(l.LastModified),
	}
	if l.ContentMD5 != nil {
		encodedContentMD5 := runtime.EncodeByteArray(l.ContentMD5, runtime.Base64StdFormat)
		aux.ContentMD5 = &encodedContentMD5
	}
	return json.Marshal(aux)
}

func (l *listItemProperties) UnmarshalJSON(b []byte) error {
	type alias listItemProperties
	aux := &struct {
		*alias
		ContentMD5   *string          `xml:"Content-MD5"`
		CreationTime *dateTimeRFC1123 `xml:"Creation-Time"`
		LastModified *dateTimeRFC1123 `xml:"Last-Modified"`
	}{
		alias: (*alias)(l),
	}
	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}
	if aux.ContentMD5 != nil {
		if err := runtime.DecodeByteArray(*aux.ContentMD5, &l.ContentMD5, runtime.Base64StdFormat); err != nil {
			return err
		}
	}
	if aux.CreationTime != nil && !(*time.Time)(aux.CreationTime).IsZero() {
		l.CreationTime = (*time.Time)(aux.CreationTime)
	}
	if aux.LastModified != nil && !(*time.Time)(aux.LastModified).IsZero() {
		l.LastModified = (*time.Time)(aux.LastModified)
	}
	return nil
}

func (l listItemProperties) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type alias listItemProperties
	aux := &struct {
		*alias
		ContentMD5   *string          `xml:"Content-MD5"`
		CreationTime *dateTimeRFC1123 `xml:"Creation-Time"`
		LastModified *dateTimeRFC1123 `xml:"Last-Modified"`
	}{
		alias:        (*alias)(&l),
		CreationTime: (*dateTimeRFC1123)(l.CreationTime),
		LastModified: (*dateTimeRFC1123)(l.LastModified),
	}
	if l.ContentMD5 != nil {
		encodedContentMD5 := runtime.EncodeByteArray(l.ContentMD5, runtime.Base64StdFormat)
		aux.ContentMD5 = &encodedContentMD5
	}
	return enc.EncodeElement(aux, start)
}

func (l *listItemProperties) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type alias listItemProperties
	aux := &struct {
		*alias
		ContentMD5   *string          `xml:"Content-MD5"`
		CreationTime *dateTimeRFC1123 `xml:"Creation-Time"`
		LastModified *dateTimeRFC1123 `xml:"Last-Modified"`
	}{
		alias: (*alias)(l),
	}
	if err := dec.DecodeElement(aux, &start); err != nil {
		return err
	}
	if aux.ContentMD5 != nil {
		if err := runtime.DecodeByteArray(*aux.ContentMD5, &l.ContentMD5, runtime.Base64StdFormat); err != nil {
			return err
		}
	}
	if aux.CreationTime != nil && !(*time.Time)(aux.CreationTime).IsZero() {
		l.CreationTime = (*time.Time)(aux.CreationTime)
	}
	if aux.LastModified != nil && !(*time.Time)(aux.LastModified).IsZero() {
		l.LastModified = (*time.Time)(aux.LastModified)
	}
	return nil
}

const dateTimeRFC1123JSON = `"` + time.RFC1123 + `"`

type dateTimeRFC1123 time.Time

func (t dateTimeRFC1123) MarshalJSON() ([]byte, error) {
	b := []byte(time.Time(t).Format(dateTimeRFC1123JSON))
	return b, nil
}

func (t *dateTimeRFC1123) UnmarshalJSON(data []byte) error {
	p, err := time.Parse(dateTimeRFC1123JSON, strings.ToUpper(string(data)))
	*t = dateTimeRFC1123(p)
	return err
}

func (t dateTimeRFC1123) MarshalText() ([]byte, error) {
	b := []byte(time.Time(t).Format(time.RFC1123))
	return b, nil
}

func (t *dateTimeRFC1123) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	p, err := time.Parse(time.RFC1123, string(data))
	*t = dateTimeRFC1123(p)
	return err
}

func (t dateTimeRFC1123) String() string {
	return time.Time(t).Format(time.RFC1123)
}
