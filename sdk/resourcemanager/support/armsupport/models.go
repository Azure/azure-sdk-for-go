//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsupport

import "time"

// CheckNameAvailabilityInput - Input of CheckNameAvailability API.
type CheckNameAvailabilityInput struct {
	// REQUIRED; The resource name to validate.
	Name *string

	// REQUIRED; The type of resource.
	Type *Type
}

// CheckNameAvailabilityOutput - Output of check name availability API.
type CheckNameAvailabilityOutput struct {
	// READ-ONLY; The detailed error message describing why the name is not available.
	Message *string

	// READ-ONLY; Indicates whether the name is available.
	NameAvailable *bool

	// READ-ONLY; The reason why the name is not available.
	Reason *string
}

// CommunicationDetails - Object that represents a Communication resource.
type CommunicationDetails struct {
	// Properties of the resource.
	Properties *CommunicationDetailsProperties

	// READ-ONLY; Id of the resource.
	ID *string

	// READ-ONLY; Name of the resource.
	Name *string

	// READ-ONLY; Type of the resource 'Microsoft.Support/communications'.
	Type *string
}

// CommunicationDetailsProperties - Describes the properties of a communication resource.
type CommunicationDetailsProperties struct {
	// REQUIRED; Body of the communication.
	Body *string

	// REQUIRED; Subject of the communication.
	Subject *string

	// Email address of the sender. This property is required if called by a service principal.
	Sender *string

	// READ-ONLY; Direction of communication.
	CommunicationDirection *CommunicationDirection

	// READ-ONLY; Communication type.
	CommunicationType *CommunicationType

	// READ-ONLY; Time in UTC (ISO 8601 format) when the communication was created.
	CreatedDate *time.Time
}

// CommunicationsListResult - Collection of Communication resources.
type CommunicationsListResult struct {
	// The URI to fetch the next page of Communication resources.
	NextLink *string

	// List of Communication resources.
	Value []*CommunicationDetails
}

// ContactProfile - Contact information associated with the support ticket.
type ContactProfile struct {
	// REQUIRED; Country of the user. This is the ISO 3166-1 alpha-3 code.
	Country *string

	// REQUIRED; First name.
	FirstName *string

	// REQUIRED; Last name.
	LastName *string

	// REQUIRED; Preferred contact method.
	PreferredContactMethod *PreferredContactMethod

	// REQUIRED; Preferred language of support from Azure. Support languages vary based on the severity you choose for your support
// ticket. Learn more at Azure Severity and responsiveness
// [https://azure.microsoft.com/support/plans/response]. Use the standard language-country code. Valid values are 'en-us'
// for English, 'zh-hans' for Chinese, 'es-es' for Spanish, 'fr-fr' for French,
// 'ja-jp' for Japanese, 'ko-kr' for Korean, 'ru-ru' for Russian, 'pt-br' for Portuguese, 'it-it' for Italian, 'zh-tw' for
// Chinese and 'de-de' for German.
	PreferredSupportLanguage *string

	// REQUIRED; Time zone of the user. This is the name of the time zone from Microsoft Time Zone Index Values [https://support.microsoft.com/help/973627/microsoft-time-zone-index-values].
	PreferredTimeZone *string

	// REQUIRED; Primary email address.
	PrimaryEmailAddress *string

	// Additional email addresses listed will be copied on any correspondence about the support ticket.
	AdditionalEmailAddresses []*string

	// Phone number. This is required if preferred contact method is phone.
	PhoneNumber *string
}

// Engineer - Support engineer information.
type Engineer struct {
	// READ-ONLY; Email address of the Azure Support engineer assigned to the support ticket.
	EmailAddress *string
}

// ExceptionResponse - The API error.
type ExceptionResponse struct {
	// The API error details.
	Error *ServiceError
}

// Operation - The operation supported by Microsoft Support resource provider.
type Operation struct {
	// The object that describes the operation.
	Display *OperationDisplay

	// READ-ONLY; Operation name: {provider}/{resource}/{operation}.
	Name *string
}

// OperationDisplay - The object that describes the operation.
type OperationDisplay struct {
	// READ-ONLY; The description of the operation.
	Description *string

	// READ-ONLY; The action that users can perform, based on their permission level.
	Operation *string

	// READ-ONLY; Service provider: Microsoft Support.
	Provider *string

	// READ-ONLY; Resource on which the operation is performed.
	Resource *string
}

// OperationsListResult - The list of operations supported by Microsoft Support resource provider.
type OperationsListResult struct {
	// The list of operations supported by Microsoft Support resource provider.
	Value []*Operation
}

// ProblemClassification resource object.
type ProblemClassification struct {
	// Properties of the resource.
	Properties *ProblemClassificationProperties

	// READ-ONLY; Id of the resource.
	ID *string

	// READ-ONLY; Name of the resource.
	Name *string

	// READ-ONLY; Type of the resource 'Microsoft.Support/problemClassification'.
	Type *string
}

// ProblemClassificationProperties - Details about a problem classification available for an Azure service.
type ProblemClassificationProperties struct {
	// Localized name of problem classification.
	DisplayName *string
}

// ProblemClassificationsListResult - Collection of ProblemClassification resources.
type ProblemClassificationsListResult struct {
	// List of ProblemClassification resources.
	Value []*ProblemClassification
}

// QuotaChangeRequest - This property is required for providing the region and new quota limits.
type QuotaChangeRequest struct {
	// Payload of the quota increase request.
	Payload *string

	// Region for which the quota increase request is being made.
	Region *string
}

// QuotaTicketDetails - Additional set of information required for quota increase support ticket for certain quota types,
// e.g.: Virtual machine cores. Get complete details about Quota payload support request along with
// examples at Support quota request [https://aka.ms/supportrpquotarequestpayload].
type QuotaTicketDetails struct {
	// Required for certain quota types when there is a sub type, such as Batch, for which you are requesting a quota increase.
	QuotaChangeRequestSubType *string

	// Quota change request version.
	QuotaChangeRequestVersion *string

	// This property is required for providing the region and new quota limits.
	QuotaChangeRequests []*QuotaChangeRequest
}

// Service - Object that represents a Service resource.
type Service struct {
	// Properties of the resource.
	Properties *ServiceProperties

	// READ-ONLY; Id of the resource.
	ID *string

	// READ-ONLY; Name of the resource.
	Name *string

	// READ-ONLY; Type of the resource 'Microsoft.Support/services'.
	Type *string
}

// ServiceError - The API error details.
type ServiceError struct {
	// The error code.
	Code *string

	// The error message.
	Message *string

	// The target of the error.
	Target *string

	// READ-ONLY; The list of error details.
	Details []*ServiceErrorDetail
}

// ServiceErrorDetail - The error details.
type ServiceErrorDetail struct {
	// The target of the error.
	Target *string

	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The error message.
	Message *string
}

// ServiceLevelAgreement - Service Level Agreement details for a support ticket.
type ServiceLevelAgreement struct {
	// READ-ONLY; Time in UTC (ISO 8601 format) when the service level agreement expires.
	ExpirationTime *time.Time

	// READ-ONLY; Service Level Agreement in minutes.
	SLAMinutes *int32

	// READ-ONLY; Time in UTC (ISO 8601 format) when the service level agreement starts.
	StartTime *time.Time
}

// ServiceProperties - Details about an Azure service available for support ticket creation.
type ServiceProperties struct {
	// Localized name of the Azure service.
	DisplayName *string

	// ARM Resource types.
	ResourceTypes []*string
}

// ServicesListResult - Collection of Service resources.
type ServicesListResult struct {
	// List of Service resources.
	Value []*Service
}

// TechnicalTicketDetails - Additional information for technical support ticket.
type TechnicalTicketDetails struct {
	// This is the resource Id of the Azure service resource (For example: A virtual machine resource or an HDInsight resource)
// for which the support ticket is created.
	ResourceID *string
}

// TicketDetails - Object that represents SupportTicketDetails resource.
type TicketDetails struct {
	// Properties of the resource.
	Properties *TicketDetailsProperties

	// READ-ONLY; Id of the resource.
	ID *string

	// READ-ONLY; Name of the resource.
	Name *string

	// READ-ONLY; Type of the resource 'Microsoft.Support/supportTickets'.
	Type *string
}

// TicketDetailsProperties - Describes the properties of a support ticket.
type TicketDetailsProperties struct {
	// REQUIRED; Contact information of the user requesting to create a support ticket.
	ContactDetails *ContactProfile

	// REQUIRED; Detailed description of the question or issue.
	Description *string

	// REQUIRED; Each Azure service has its own set of issue categories, also known as problem classification. This parameter
// is the unique Id for the type of problem you are experiencing.
	ProblemClassificationID *string

	// REQUIRED; This is the resource Id of the Azure service resource associated with the support ticket.
	ServiceID *string

	// REQUIRED; A value that indicates the urgency of the case, which in turn determines the response time according to the service
// level agreement of the technical support plan you have with Azure. Note: 'Highest
// critical impact', also known as the 'Emergency - Severe impact' level in the Azure portal is reserved only for our Premium
// customers.
	Severity *SeverityLevel

	// REQUIRED; Title of the support ticket.
	Title *string

	// Time in UTC (ISO 8601 format) when the problem started.
	ProblemStartTime *time.Time

	// Additional ticket details associated with a quota support ticket request.
	QuotaTicketDetails *QuotaTicketDetails

	// Indicates if this requires a 24x7 response from Azure.
	Require24X7Response *bool

	// Service Level Agreement information for this support ticket.
	ServiceLevelAgreement *ServiceLevelAgreement

	// Information about the support engineer working on this support ticket.
	SupportEngineer *Engineer

	// System generated support ticket Id that is unique.
	SupportTicketID *string

	// Additional ticket details associated with a technical support ticket request.
	TechnicalTicketDetails *TechnicalTicketDetails

	// READ-ONLY; Time in UTC (ISO 8601 format) when the support ticket was created.
	CreatedDate *time.Time

	// READ-ONLY; Enrollment Id associated with the support ticket.
	EnrollmentID *string

	// READ-ONLY; Time in UTC (ISO 8601 format) when the support ticket was last modified.
	ModifiedDate *time.Time

	// READ-ONLY; Localized name of problem classification.
	ProblemClassificationDisplayName *string

	// READ-ONLY; Localized name of the Azure service.
	ServiceDisplayName *string

	// READ-ONLY; Status of the support ticket.
	Status *string

	// READ-ONLY; Support plan type associated with the support ticket.
	SupportPlanType *string
}

// TicketsListResult - Object that represents a collection of SupportTicket resources.
type TicketsListResult struct {
	// The URI to fetch the next page of SupportTicket resources.
	NextLink *string

	// List of SupportTicket resources.
	Value []*TicketDetails
}

// UpdateContactProfile - Contact information associated with the support ticket.
type UpdateContactProfile struct {
	// Email addresses listed will be copied on any correspondence about the support ticket.
	AdditionalEmailAddresses []*string

	// Country of the user. This is the ISO 3166-1 alpha-3 code.
	Country *string

	// First name.
	FirstName *string

	// Last name.
	LastName *string

	// Phone number. This is required if preferred contact method is phone.
	PhoneNumber *string

	// Preferred contact method.
	PreferredContactMethod *PreferredContactMethod

	// Preferred language of support from Azure. Support languages vary based on the severity you choose for your support ticket.
// Learn more at Azure Severity and responsiveness
// [https://azure.microsoft.com/support/plans/response/]. Use the standard language-country code. Valid values are 'en-us'
// for English, 'zh-hans' for Chinese, 'es-es' for Spanish, 'fr-fr' for French,
// 'ja-jp' for Japanese, 'ko-kr' for Korean, 'ru-ru' for Russian, 'pt-br' for Portuguese, 'it-it' for Italian, 'zh-tw' for
// Chinese and 'de-de' for German.
	PreferredSupportLanguage *string

	// Time zone of the user. This is the name of the time zone from Microsoft Time Zone Index Values [https://support.microsoft.com/help/973627/microsoft-time-zone-index-values].
	PreferredTimeZone *string

	// Primary email address.
	PrimaryEmailAddress *string
}

// UpdateSupportTicket - Updates severity, ticket status, and contact details in the support ticket.
type UpdateSupportTicket struct {
	// Contact details to be updated on the support ticket.
	ContactDetails *UpdateContactProfile

	// Severity level.
	Severity *SeverityLevel

	// Status to be updated on the ticket.
	Status *Status
}

