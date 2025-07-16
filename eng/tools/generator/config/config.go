// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package config

import (
	"encoding/json"
)

type Config struct {
	Track2Requests   Track2ReleaseRequests   `json:"track2Requests,omitempty"`
	TypeSpecRequests TypeSpecReleaseRequests `json:"typespecRequests,omitempty"`
}

func (c Config) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}
