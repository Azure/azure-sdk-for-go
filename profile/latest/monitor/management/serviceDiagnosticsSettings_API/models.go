package servicediagnosticssettingsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/management/2016-09-01/serviceDiagnosticsSettings_API"
)

type (
	 ErrorResponse = original.ErrorResponse
	 LogSettings = original.LogSettings
	 MetricSettings = original.MetricSettings
	 Resource = original.Resource
	 RetentionPolicy = original.RetentionPolicy
	 ServiceDiagnosticSettings = original.ServiceDiagnosticSettings
	 ServiceDiagnosticSettingsResource = original.ServiceDiagnosticSettingsResource
	 ServiceDiagnosticSettingsResourcePatch = original.ServiceDiagnosticSettingsResourcePatch
	 ServiceDiagnosticSettingsClient = original.ServiceDiagnosticSettingsClient
	 ManagementClient = original.ManagementClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
