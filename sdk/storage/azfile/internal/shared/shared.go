//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"errors"
	"fmt"
	"strings"
)

const (
	TokenScope              = "https://storage.azure.com/.default"
	StorageAnalyticsVersion = "1.0"

	HeaderAuthorization     = "Authorization"
	HeaderXmsDate           = "x-ms-date"
	HeaderContentLength     = "Content-Length"
	HeaderContentEncoding   = "Content-Encoding"
	HeaderContentLanguage   = "Content-Language"
	HeaderContentType       = "Content-Type"
	HeaderContentMD5        = "Content-MD5"
	HeaderIfModifiedSince   = "If-Modified-Since"
	HeaderIfMatch           = "If-Match"
	HeaderIfNoneMatch       = "If-None-Match"
	HeaderIfUnmodifiedSince = "If-Unmodified-Since"
	HeaderRange             = "Range"
)

func GetClientOptions[T any](o *T) *T {
	if o == nil {
		return new(T)
	}
	return o
}

var errConnectionString = errors.New("connection string is either blank or malformed. The expected connection string " +
	"should contain key value pairs separated by semicolons. For example 'DefaultEndpointsProtocol=https;AccountName=<accountName>;" +
	"AccountKey=<accountKey>;EndpointSuffix=core.windows.net'")

type ParsedConnectionString struct {
	ServiceURL  string
	AccountName string
	AccountKey  string
}

func ParseConnectionString(connectionString string) (ParsedConnectionString, error) {
	const (
		defaultScheme = "https"
		defaultSuffix = "core.windows.net"
	)

	connStrMap := make(map[string]string)
	connectionString = strings.TrimRight(connectionString, ";")

	splitString := strings.Split(connectionString, ";")
	if len(splitString) == 0 {
		return ParsedConnectionString{}, errConnectionString
	}
	for _, stringPart := range splitString {
		parts := strings.SplitN(stringPart, "=", 2)
		if len(parts) != 2 {
			return ParsedConnectionString{}, errConnectionString
		}
		connStrMap[parts[0]] = parts[1]
	}

	accountName, ok := connStrMap["AccountName"]
	if !ok {
		return ParsedConnectionString{}, errors.New("connection string missing AccountName")
	}

	accountKey, ok := connStrMap["AccountKey"]
	if !ok {
		sharedAccessSignature, ok := connStrMap["SharedAccessSignature"]
		if !ok {
			return ParsedConnectionString{}, errors.New("connection string missing AccountKey and SharedAccessSignature")
		}
		return ParsedConnectionString{
			ServiceURL: fmt.Sprintf("%v://%v.file.%v/?%v", defaultScheme, accountName, defaultSuffix, sharedAccessSignature),
		}, nil
	}

	protocol, ok := connStrMap["DefaultEndpointsProtocol"]
	if !ok {
		protocol = defaultScheme
	}

	suffix, ok := connStrMap["EndpointSuffix"]
	if !ok {
		suffix = defaultSuffix
	}

	if fileEndpoint, ok := connStrMap["FileEndpoint"]; ok {
		return ParsedConnectionString{
			ServiceURL:  fileEndpoint,
			AccountName: accountName,
			AccountKey:  accountKey,
		}, nil
	}

	return ParsedConnectionString{
		ServiceURL:  fmt.Sprintf("%v://%v.file.%v", protocol, accountName, suffix),
		AccountName: accountName,
		AccountKey:  accountKey,
	}, nil
}
