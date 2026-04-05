// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

// ConnectionMode specifies the connection mode to use when communicating with Azure Cosmos DB.
// Gateway mode routes all requests through the Azure Cosmos DB gateway using HTTPS.
// Direct mode connects directly to the backend nodes using the RNTBD (binary) protocol for data operations.
type ConnectionMode int

const (
	// ConnectionModeGateway routes all requests through the Azure Cosmos DB gateway using HTTPS.
	// This is the default mode and works with any network configuration.
	// All operations (data and control plane) use HTTP/HTTPS.
	ConnectionModeGateway ConnectionMode = iota

	// ConnectionModeDirect connects directly to Azure Cosmos DB backend nodes for data operations.
	// Uses the RNTBD (Remote Native Transport Binary Direct) protocol over TCP.
	// Provides lower latency and higher throughput compared to Gateway mode.
	// Control plane operations (database/container management) still use HTTPS via the gateway.
	// Requires TCP connectivity to the backend ports (typically 10255).
	ConnectionModeDirect
)

// String returns the string representation of the ConnectionMode.
func (c ConnectionMode) String() string {
	switch c {
	case ConnectionModeGateway:
		return "Gateway"
	case ConnectionModeDirect:
		return "Direct"
	default:
		return "Unknown"
	}
}

// ConnectionModes returns all available connection modes.
func ConnectionModes() []ConnectionMode {
	return []ConnectionMode{ConnectionModeGateway, ConnectionModeDirect}
}

// ToPtr returns a pointer to the ConnectionMode.
func (c ConnectionMode) ToPtr() *ConnectionMode {
	return &c
}
