//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package internal

import (
	"errors"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

const (
	TokenScope = "https://webpubsub.azure.com/.default"
)

const (
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
	HeaderXmsVersion        = "x-ms-version"
	HeaderXmsRequestID      = "x-ms-request-id"
)

func NewPipeline(plOpts runtime.PipelineOptions, options *policy.ClientOptions) runtime.Pipeline {
	return runtime.NewPipeline(ModuleName+".Client", ModuleVersion, plOpts, options)
}

var errConnectionString = errors.New("connection string is either blank or malformed. The expected connection string " +
	"should contain key value pairs separated by semicolons. For example 'Endpoint=<endpoint>;AccessKey=<key>;'")

type ParsedConnectionString struct {
	Endpoint  string
	AccessKey string
}

func ParseConnectionString(connectionString string) (ParsedConnectionString, error) {
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
		connStrMap[strings.ToLower(parts[0])] = parts[1]
	}

	endpoint, ok := connStrMap["endpoint"]
	if !ok {
		return ParsedConnectionString{}, errConnectionString
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return ParsedConnectionString{}, errors.New("endpoint is not a valid URL")
	}

	port, has_port := connStrMap["port"]
	if has_port {
		parsedURL.Host = parsedURL.Hostname() + ":" + port
		endpoint = parsedURL.String()
	}

	if !strings.HasSuffix(endpoint, "/") {
		// add a trailing slash to be consistent with the portal
		endpoint += "/"
	}

	key, ok := connStrMap["accesskey"]
	if !ok {
		return ParsedConnectionString{}, errConnectionString
	}

	return ParsedConnectionString{
		Endpoint:  endpoint,
		AccessKey: key,
	}, nil
}
