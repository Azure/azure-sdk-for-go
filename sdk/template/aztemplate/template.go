//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztemplate

type TemplateClient struct{}

func NewTemplateClient() (TemplateClient, error) {

	return TemplateClient{}, nil

}

func ClientVersion() string {

	return packageinfo.Version

}
