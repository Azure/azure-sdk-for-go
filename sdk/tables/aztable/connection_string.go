// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var ErrConnectionString = errors.New("connection string is either blank or malformed. The expected connection string should contain key value pairs separated by semicolons. For example 'DefaultEndpointsProtocol=https;AccountName=<accountName>;AccountKey=<accountKey>;EndpointSuffix=core.windows.net'")

// NewTableClientFromConnectionString creates a new TableClient struct from a connection string. The connection
// string must contain either an account name and account key or an account name and a shared access signature.
func NewTableClientFromConnectionString(tableName string, connectionString string, options *TableClientOptions) (*TableClient, error) {
	if options == nil {
		options = &TableClientOptions{}
	}
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return NewTableClient(tableName, endpoint, azcore.AnonymousCredential(), nil)
	}
	return NewTableClient(tableName, endpoint, credential, options)
}

// NewTableServiceClientFromConnectionString creates a new TableServiceClient struct from a connection string. The connection
// string must contain either an account name and account key or an account name and a shared access signature.
func NewTableServiceClientFromConnectionString(connectionString string, options *TableClientOptions) (*TableServiceClient, error) {
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	if credential == nil {
		return NewTableServiceClient(endpoint, azcore.AnonymousCredential(), nil)
	}
	return NewTableServiceClient(endpoint, credential, options)
}

func convertConnStrToMap(connStr string) (map[string]string, error) {
	ret := make(map[string]string)
	connStr = strings.TrimRight(connStr, ";")

	splitString := strings.Split(connStr, ";")
	if len(splitString) == 0 {
		return ret, ErrConnectionString
	}
	for _, stringPart := range splitString {
		parts := strings.Split(stringPart, "=")
		if len(parts) != 2 {
			return ret, ErrConnectionString
		}
		ret[parts[0]] = parts[1]
	}
	return ret, nil
}

// parseConnectionString parses a connection string into a service URL and a SharedKeyCredential or a service url with the
// SharedAccessSignature combined.
func parseConnectionString(connStr string) (string, azcore.Credential, error) {
	var serviceURL string
	var cred azcore.Credential

	defaultScheme := "https"
	defaultSuffix := "core.windows.net"

	connStrMap, err := convertConnStrToMap(connStr)
	if err != nil {
		return "", nil, err
	}

	accountName, ok := connStrMap["AccountName"]
	if !ok {
		return "", nil, ErrConnectionString
	}
	accountKey, ok := connStrMap["AccountKey"]
	if !ok {

		// if accountName == "" || accountKey == "" {
		// Try sharedaccesssignature
		sharedAccessSignature, ok := connStrMap["SharedAccessSignature"]
		if !ok {
			return "", nil, ErrConnectionString
		}
		return fmt.Sprintf("%v://%v.table.%v/?%v", defaultScheme, accountName, defaultSuffix, sharedAccessSignature), nil, nil
		// }

		// return "", nil, ErrConnectionString
	}

	protocol, ok := connStrMap["DefaultEndpointsProtocol"]
	if !ok {
		protocol = defaultScheme
	}

	suffix, ok := connStrMap["EndpointSuffix"]
	if !ok {
		suffix = defaultSuffix
	}

	tableEndpoint, ok := connStrMap["TableEndpoint"]
	if ok {
		cred, err = NewSharedKeyCredential(accountName, accountKey)
		return tableEndpoint, cred, err
	}
	serviceURL = fmt.Sprintf("%v://%v.table.%v", protocol, accountName, suffix)

	cred, err = NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", nil, err
	}

	return serviceURL, cred, nil
}
