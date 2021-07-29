// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var ErrConnectionString = errors.New("connection string is either blank or malformed")

func NewTableClientFromConnectionString(tableName string, connectionString string, options *TableClientOptions) (*TableClient, error) {
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewTableClient(tableName, endpoint, *credential, options)
}

func NewTableServiceClientFrommConnectionString(connectionString string, options *TableClientOptions) (*TableServiceClient, error) {
	endpoint, credential, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewTableServiceClient(endpoint, *credential, options)
}

func convertConnStrToMap(connStr string) (map[string]string, error) {
	ret := make(map[string]string)

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

func parseConnectionString(connStr string) (string, *azcore.Credential, error) {
	var serviceURL string
	var cred azcore.Credential

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
		return "", nil, ErrConnectionString
	}

	if accountName == "" || accountKey == "" {
		// Try sharedaccesssignature
		sharedAccessSignature, ok := connStrMap["sharedaccesssignature"]
		if !ok {
			return serviceURL, nil, ErrConnectionString
		}
		return sharedAccessSignature, nil, errors.New("there is not support for SharedAccessSignature yet")
		// cred = azcore.SharedAccessSignature(sharedAccessSignature)
	}
	defaultProtocol, ok := connStrMap["DefaultEndpointsProtocol"]
	if !ok {
		defaultProtocol = "https"
	}

	endpointSuffix, ok := connStrMap["EndpointSuffix"]
	if !ok {
		endpointSuffix = "core.windows.net"
	}
	serviceURL = fmt.Sprintf("%v://%v.table.%v", defaultProtocol, accountName, endpointSuffix)

	cred, err = NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return "", nil, err
	}

	// primary, okPrimary := pairsMap["tableendpoint"]
	// secondary, okSecondary := pairsMap["tablesecondaryendpoint"]
	// if !okPrimary {
	// 	if okSecondary {
	// 		return serviceURL, cred, errors.New("Connection string specifies only secondary connection")
	// 	}
	// 	if endpointsProtocol, ok := pairsMap["defaultendpointsprotocol"]; ok {
	// 		if accountName, ok := pairsMap["accountname"]; ok {
	// 			if endpointSuffix, ok := pairsMap["endpointsuffix"]; ok {
	// 				primary = fmt.Sprintf("%v://%v.table.%v", endpointsProtocol, accountName, endpointSuffix)
	// 				secondary = fmt.Sprintf("%v-secondary.table.%v", accountName, endpointSuffix)
	// 				okPrimary = true
	// 				okSecondary = true
	// 			}
	// 		}
	// 	}
	// }

	// if !okPrimary {

	// }

	// if serviceURL, ok = pairsMap["tableendpoint"]; !ok {
	// 	return serviceURL, cred, errors.New("Connection string does not specify")
	// }

	return serviceURL, &cred, nil
}
