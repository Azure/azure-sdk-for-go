// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"strings"
	"testing"
)

func TestEnrichTelemetryValue(t *testing.T) {
	cosmosClientOptions := &CosmosClientOptions{}
	telemetryOptions := cosmosClientOptions.enrichTelemetryOptions()
	if !strings.Contains(telemetryOptions.Value, "azsdk-go-azcosmos") {
		t.Errorf("Expected azsdk-go-azcosmos in telemetryOptions.Value, but got %s", telemetryOptions.Value)
	}
}
