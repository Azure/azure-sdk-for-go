// Package imagesearch implements the Azure ARM Imagesearch service API version 1.0.
//
// The Image Search API lets you send a search query to Bing and get back a list of relevant images. This section
// provides technical details about the query parameters and headers that you use to request images and the JSON
// response objects that contain them. For examples that show how to make requests, see [Searching the Web for
// Images](https://docs.microsoft.com/azure/cognitive-services/bing-image-search/search-the-web).
package imagesearch

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
)

const (
	// DefaultEndpoint is the default value for endpoint
	DefaultEndpoint = "https://api.cognitive.microsoft.com"
)

// BaseClient is the base client for Imagesearch.
type BaseClient struct {
	autorest.Client
	Endpoint string
}

// New creates an instance of the BaseClient client.
func New() BaseClient {
	return NewWithoutDefaults(DefaultEndpoint)
}

// NewWithoutDefaults creates an instance of the BaseClient client.
func NewWithoutDefaults(endpoint string) BaseClient {
	return BaseClient{
		Client:   autorest.NewClientWithUserAgent(UserAgent()),
		Endpoint: endpoint,
	}
}
