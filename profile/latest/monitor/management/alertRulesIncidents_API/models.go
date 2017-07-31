package alertrulesincidentsapi

import (
	 original "github.com/Azure/azure-sdk-for-go/service/monitor/management/2016-03-01/alertRulesIncidents_API"
)

type (
	 ErrorResponse = original.ErrorResponse
	 Incident = original.Incident
	 IncidentListResult = original.IncidentListResult
	 AlertRuleIncidentsClient = original.AlertRuleIncidentsClient
	 ManagementClient = original.ManagementClient
)

const (
	 DefaultBaseURI = original.DefaultBaseURI
)
