// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// SyncToken contains data used in the Sync-Token header.
// See [Azure App Configuration documentation] for more information on sync tokens.
//
// [Azure App Configuration documentation]: https://learn.microsoft.com/azure/azure-app-configuration/rest-api-consistency
// Exported as azappconfig.SyncToken
type SyncToken string

// the following are NOT exported

// SyncTokenValues contains the parts of a SyncToken.
type SyncTokenValues struct {
	ID      string
	Value   string
	Version int64
}

// ParseSyncToken parses the provided SyncToken.
func ParseSyncToken(syncToken SyncToken) ([]SyncTokenValues, error) {
	rawToken := strings.TrimSpace(string(syncToken))
	if rawToken == "" {
		return nil, errors.New("syncToken can't be empty")
	}

	tokenParts := strings.Split(rawToken, ",")
	syncTokens := make([]SyncTokenValues, len(tokenParts))

	// token format is "<id>=<value>;sn=<sn>" and can contain multiple, comman-delimited values
	for i, token := range tokenParts {
		items := strings.Split(token, ";")
		if len(items) != 2 {
			return nil, fmt.Errorf("invalid token %s", token)
		}

		// items[0] contains "<id>=<value>"
		// note that <value> is a base-64 encoded string, so don't try to split on '='
		assignmentIndex := strings.Index(items[0], "=")
		if assignmentIndex < 0 {
			return nil, fmt.Errorf("unexpected token format %s", items[0])
		}
		tokenID := strings.TrimSpace(items[0][:assignmentIndex])
		tokenValue := strings.TrimSpace(items[0][assignmentIndex+1:])

		// items[1] contains "sn=<sn>"
		// parse the version number after the equals sign
		assignmentIndex = strings.Index(items[1], "=")
		if assignmentIndex < 0 {
			return nil, fmt.Errorf("unexpected token version format %s", items[1])
		}
		tokenVersion, err := strconv.ParseInt(strings.TrimSpace(items[1][assignmentIndex+1:]), 10, 64)
		if err != nil {
			return nil, err
		}

		syncTokens[i] = SyncTokenValues{
			ID:      tokenID,
			Value:   tokenValue,
			Version: tokenVersion,
		}
	}

	return syncTokens, nil
}
