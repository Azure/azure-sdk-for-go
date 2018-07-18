// +build go1.9

// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package alertsmanagement

import original "github.com/Azure/azure-sdk-for-go/services/alertsmanagement/mgmt/2018-05-05-preview/alertsmanagement"

type AlertsClient = original.AlertsClient

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient
type AlertModificationEvent = original.AlertModificationEvent

const (
	AlertCreated           AlertModificationEvent = original.AlertCreated
	MonitorConditionChange AlertModificationEvent = original.MonitorConditionChange
	StateChange            AlertModificationEvent = original.StateChange
)

type AlertsSortByFields = original.AlertsSortByFields

const (
	AlertsSortByFieldsAlertState           AlertsSortByFields = original.AlertsSortByFieldsAlertState
	AlertsSortByFieldsLastModifiedDateTime AlertsSortByFields = original.AlertsSortByFieldsLastModifiedDateTime
	AlertsSortByFieldsMonitorCondition     AlertsSortByFields = original.AlertsSortByFieldsMonitorCondition
	AlertsSortByFieldsName                 AlertsSortByFields = original.AlertsSortByFieldsName
	AlertsSortByFieldsSeverity             AlertsSortByFields = original.AlertsSortByFieldsSeverity
	AlertsSortByFieldsStartDateTime        AlertsSortByFields = original.AlertsSortByFieldsStartDateTime
	AlertsSortByFieldsTargetResource       AlertsSortByFields = original.AlertsSortByFieldsTargetResource
	AlertsSortByFieldsTargetResourceGroup  AlertsSortByFields = original.AlertsSortByFieldsTargetResourceGroup
	AlertsSortByFieldsTargetResourceName   AlertsSortByFields = original.AlertsSortByFieldsTargetResourceName
	AlertsSortByFieldsTargetResourceType   AlertsSortByFields = original.AlertsSortByFieldsTargetResourceType
)

type AlertState = original.AlertState

const (
	AlertStateAcknowledged AlertState = original.AlertStateAcknowledged
	AlertStateClosed       AlertState = original.AlertStateClosed
	AlertStateNew          AlertState = original.AlertStateNew
)

type APIVersion = original.APIVersion

const (
	TwoZeroOneEightHyphenMinusZeroFiveHyphenMinusZeroFiveHyphenMinuspreview     APIVersion = original.TwoZeroOneEightHyphenMinusZeroFiveHyphenMinusZeroFiveHyphenMinuspreview
	TwoZeroOneSevenHyphenMinusOneOneHyphenMinusOneFiveHyphenMinusprivatepreview APIVersion = original.TwoZeroOneSevenHyphenMinusOneOneHyphenMinusOneFiveHyphenMinusprivatepreview
)

type MonitorCondition = original.MonitorCondition

const (
	Fired    MonitorCondition = original.Fired
	Resolved MonitorCondition = original.Resolved
)

type MonitorService = original.MonitorService

const (
	ActivityLogAdministrative MonitorService = original.ActivityLogAdministrative
	ActivityLogAutoscale      MonitorService = original.ActivityLogAutoscale
	ActivityLogPolicy         MonitorService = original.ActivityLogPolicy
	ActivityLogRecommendation MonitorService = original.ActivityLogRecommendation
	ActivityLogSecurity       MonitorService = original.ActivityLogSecurity
	ApplicationInsights       MonitorService = original.ApplicationInsights
	InfrastructureInsights    MonitorService = original.InfrastructureInsights
	LogAnalytics              MonitorService = original.LogAnalytics
	Nagios                    MonitorService = original.Nagios
	Platform                  MonitorService = original.Platform
	SCOM                      MonitorService = original.SCOM
	ServiceHealth             MonitorService = original.ServiceHealth
	SmartDetector             MonitorService = original.SmartDetector
	Zabbix                    MonitorService = original.Zabbix
)

type Severity = original.Severity

const (
	Sev0 Severity = original.Sev0
	Sev1 Severity = original.Sev1
	Sev2 Severity = original.Sev2
	Sev3 Severity = original.Sev3
	Sev4 Severity = original.Sev4
)

type SignalType = original.SignalType

const (
	Log     SignalType = original.Log
	Metric  SignalType = original.Metric
	Unknown SignalType = original.Unknown
)

type SmartGroupModificationEvent = original.SmartGroupModificationEvent

const (
	SmartGroupModificationEventAlertAdded        SmartGroupModificationEvent = original.SmartGroupModificationEventAlertAdded
	SmartGroupModificationEventAlertRemoved      SmartGroupModificationEvent = original.SmartGroupModificationEventAlertRemoved
	SmartGroupModificationEventSmartGroupCreated SmartGroupModificationEvent = original.SmartGroupModificationEventSmartGroupCreated
	SmartGroupModificationEventStateChange       SmartGroupModificationEvent = original.SmartGroupModificationEventStateChange
)

type SmartGroupsSortByFields = original.SmartGroupsSortByFields

const (
	SmartGroupsSortByFieldsAlertsCount          SmartGroupsSortByFields = original.SmartGroupsSortByFieldsAlertsCount
	SmartGroupsSortByFieldsLastModifiedDateTime SmartGroupsSortByFields = original.SmartGroupsSortByFieldsLastModifiedDateTime
	SmartGroupsSortByFieldsSeverity             SmartGroupsSortByFields = original.SmartGroupsSortByFieldsSeverity
	SmartGroupsSortByFieldsStartDateTime        SmartGroupsSortByFields = original.SmartGroupsSortByFieldsStartDateTime
	SmartGroupsSortByFieldsState                SmartGroupsSortByFields = original.SmartGroupsSortByFieldsState
)

type State = original.State

const (
	StateAcknowledged State = original.StateAcknowledged
	StateClosed       State = original.StateClosed
	StateNew          State = original.StateNew
)

type TimeRange = original.TimeRange

const (
	Oned       TimeRange = original.Oned
	Oneh       TimeRange = original.Oneh
	Sevend     TimeRange = original.Sevend
	ThreeZerod TimeRange = original.ThreeZerod
)

type Alert = original.Alert
type AlertModification = original.AlertModification
type AlertModificationItem = original.AlertModificationItem
type AlertModificationProperties = original.AlertModificationProperties
type AlertProperties = original.AlertProperties
type AlertsList = original.AlertsList
type AlertsListIterator = original.AlertsListIterator
type AlertsListPage = original.AlertsListPage
type AlertsSummary = original.AlertsSummary
type AlertsSummaryByMonitorCondition = original.AlertsSummaryByMonitorCondition
type AlertsSummaryByMonitorService = original.AlertsSummaryByMonitorService
type AlertsSummaryBySeverityAndMonitorCondition = original.AlertsSummaryBySeverityAndMonitorCondition
type AlertsSummaryBySeverityAndMonitorConditionSev0 = original.AlertsSummaryBySeverityAndMonitorConditionSev0
type AlertsSummaryBySeverityAndMonitorConditionSev1 = original.AlertsSummaryBySeverityAndMonitorConditionSev1
type AlertsSummaryBySeverityAndMonitorConditionSev2 = original.AlertsSummaryBySeverityAndMonitorConditionSev2
type AlertsSummaryBySeverityAndMonitorConditionSev3 = original.AlertsSummaryBySeverityAndMonitorConditionSev3
type AlertsSummaryBySeverityAndMonitorConditionSev4 = original.AlertsSummaryBySeverityAndMonitorConditionSev4
type AlertsSummaryByState = original.AlertsSummaryByState
type AlertsSummaryProperties = original.AlertsSummaryProperties
type AlertsSummaryPropertiesSummaryByMonitorService = original.AlertsSummaryPropertiesSummaryByMonitorService
type AlertsSummaryPropertiesSummaryBySeverity = original.AlertsSummaryPropertiesSummaryBySeverity
type AlertsSummaryPropertiesSummaryBySeverityAndMonitorCondition = original.AlertsSummaryPropertiesSummaryBySeverityAndMonitorCondition
type AlertsSummaryPropertiesSummaryBySeveritySev0 = original.AlertsSummaryPropertiesSummaryBySeveritySev0
type AlertsSummaryPropertiesSummaryBySeveritySev1 = original.AlertsSummaryPropertiesSummaryBySeveritySev1
type AlertsSummaryPropertiesSummaryBySeveritySev2 = original.AlertsSummaryPropertiesSummaryBySeveritySev2
type AlertsSummaryPropertiesSummaryBySeveritySev3 = original.AlertsSummaryPropertiesSummaryBySeveritySev3
type AlertsSummaryPropertiesSummaryBySeveritySev4 = original.AlertsSummaryPropertiesSummaryBySeveritySev4
type AlertsSummaryPropertiesSummaryByState = original.AlertsSummaryPropertiesSummaryByState
type ErrorResponse = original.ErrorResponse
type ErrorResponseBody = original.ErrorResponseBody
type Operation = original.Operation
type OperationDisplay = original.OperationDisplay
type OperationsList = original.OperationsList
type OperationsListIterator = original.OperationsListIterator
type OperationsListPage = original.OperationsListPage
type Resource = original.Resource
type SmartGroup = original.SmartGroup
type SmartGroupAggregatedProperty = original.SmartGroupAggregatedProperty
type SmartGroupModification = original.SmartGroupModification
type SmartGroupModificationItem = original.SmartGroupModificationItem
type SmartGroupModificationProperties = original.SmartGroupModificationProperties
type SmartGroupProperties = original.SmartGroupProperties
type SmartGroupsList = original.SmartGroupsList
type OperationsClient = original.OperationsClient
type SmartGroupsClient = original.SmartGroupsClient

func NewAlertsClient(subscriptionID string, monitorService1 MonitorService) AlertsClient {
	return original.NewAlertsClient(subscriptionID, monitorService1)
}
func NewAlertsClientWithBaseURI(baseURI string, subscriptionID string, monitorService1 MonitorService) AlertsClient {
	return original.NewAlertsClientWithBaseURI(baseURI, subscriptionID, monitorService1)
}
func New(subscriptionID string, monitorService1 MonitorService) BaseClient {
	return original.New(subscriptionID, monitorService1)
}
func NewWithBaseURI(baseURI string, subscriptionID string, monitorService1 MonitorService) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID, monitorService1)
}
func PossibleAlertModificationEventValues() []AlertModificationEvent {
	return original.PossibleAlertModificationEventValues()
}
func PossibleAlertsSortByFieldsValues() []AlertsSortByFields {
	return original.PossibleAlertsSortByFieldsValues()
}
func PossibleAlertStateValues() []AlertState {
	return original.PossibleAlertStateValues()
}
func PossibleAPIVersionValues() []APIVersion {
	return original.PossibleAPIVersionValues()
}
func PossibleMonitorConditionValues() []MonitorCondition {
	return original.PossibleMonitorConditionValues()
}
func PossibleMonitorServiceValues() []MonitorService {
	return original.PossibleMonitorServiceValues()
}
func PossibleSeverityValues() []Severity {
	return original.PossibleSeverityValues()
}
func PossibleSignalTypeValues() []SignalType {
	return original.PossibleSignalTypeValues()
}
func PossibleSmartGroupModificationEventValues() []SmartGroupModificationEvent {
	return original.PossibleSmartGroupModificationEventValues()
}
func PossibleSmartGroupsSortByFieldsValues() []SmartGroupsSortByFields {
	return original.PossibleSmartGroupsSortByFieldsValues()
}
func PossibleStateValues() []State {
	return original.PossibleStateValues()
}
func PossibleTimeRangeValues() []TimeRange {
	return original.PossibleTimeRangeValues()
}
func NewOperationsClient(subscriptionID string, monitorService1 MonitorService) OperationsClient {
	return original.NewOperationsClient(subscriptionID, monitorService1)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string, monitorService1 MonitorService) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID, monitorService1)
}
func NewSmartGroupsClient(subscriptionID string, monitorService1 MonitorService) SmartGroupsClient {
	return original.NewSmartGroupsClient(subscriptionID, monitorService1)
}
func NewSmartGroupsClientWithBaseURI(baseURI string, subscriptionID string, monitorService1 MonitorService) SmartGroupsClient {
	return original.NewSmartGroupsClientWithBaseURI(baseURI, subscriptionID, monitorService1)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
