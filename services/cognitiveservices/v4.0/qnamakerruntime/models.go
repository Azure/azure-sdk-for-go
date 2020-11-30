package qnamakerruntime

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
)

// The package's fully qualified name.
const fqdn = "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v4.0/qnamakerruntime"

// ContextDTO context associated with Qna.
type ContextDTO struct {
	// IsContextOnly - To mark if a prompt is relevant only with a previous question or not.
	// true - Do not include this QnA as search result for queries without context
	// false - ignores context and includes this QnA in search result
	IsContextOnly *bool `json:"isContextOnly,omitempty"`
	// Prompts - List of prompts associated with the answer.
	Prompts *[]PromptDTO `json:"prompts,omitempty"`
}

// Error the error object. As per Microsoft One API guidelines -
// https://github.com/Microsoft/api-guidelines/blob/vNext/Guidelines.md#7102-error-condition-responses.
type Error struct {
	// Code - One of a server-defined set of error codes. Possible values include: 'BadArgument', 'Forbidden', 'NotFound', 'KbNotFound', 'Unauthorized', 'Unspecified', 'EndpointKeysError', 'QuotaExceeded', 'QnaRuntimeError', 'SKULimitExceeded', 'OperationNotFound', 'ServiceError', 'ValidationFailure', 'ExtractionFailure'
	Code ErrorCodeType `json:"code,omitempty"`
	// Message - A human-readable representation of the error.
	Message *string `json:"message,omitempty"`
	// Target - The target of the error.
	Target *string `json:"target,omitempty"`
	// Details - An array of details about specific errors that led to this reported error.
	Details *[]Error `json:"details,omitempty"`
	// InnerError - An object containing more specific information than the current object about the error.
	InnerError *InnerErrorModel `json:"innerError,omitempty"`
}

// ErrorResponse error response. As per Microsoft One API guidelines -
// https://github.com/Microsoft/api-guidelines/blob/vNext/Guidelines.md#7102-error-condition-responses.
type ErrorResponse struct {
	// Error - The error object.
	Error *ErrorResponseError `json:"error,omitempty"`
}

// ErrorResponseError the error object.
type ErrorResponseError struct {
	// Code - One of a server-defined set of error codes. Possible values include: 'BadArgument', 'Forbidden', 'NotFound', 'KbNotFound', 'Unauthorized', 'Unspecified', 'EndpointKeysError', 'QuotaExceeded', 'QnaRuntimeError', 'SKULimitExceeded', 'OperationNotFound', 'ServiceError', 'ValidationFailure', 'ExtractionFailure'
	Code ErrorCodeType `json:"code,omitempty"`
	// Message - A human-readable representation of the error.
	Message *string `json:"message,omitempty"`
	// Target - The target of the error.
	Target *string `json:"target,omitempty"`
	// Details - An array of details about specific errors that led to this reported error.
	Details *[]Error `json:"details,omitempty"`
	// InnerError - An object containing more specific information than the current object about the error.
	InnerError *InnerErrorModel `json:"innerError,omitempty"`
}

// FeedbackRecordDTO active learning feedback record.
type FeedbackRecordDTO struct {
	// UserID - Unique identifier for the user.
	UserID *string `json:"userId,omitempty"`
	// UserQuestion - The suggested question being provided as feedback.
	UserQuestion *string `json:"userQuestion,omitempty"`
	// QnaID - The qnaId for which the suggested question is provided as feedback.
	QnaID *int32 `json:"qnaId,omitempty"`
}

// FeedbackRecordsDTO active learning feedback records.
type FeedbackRecordsDTO struct {
	// FeedbackRecords - List of feedback records.
	FeedbackRecords *[]FeedbackRecordDTO `json:"feedbackRecords,omitempty"`
}

// InnerErrorModel an object containing more specific information about the error. As per Microsoft One API
// guidelines -
// https://github.com/Microsoft/api-guidelines/blob/vNext/Guidelines.md#7102-error-condition-responses.
type InnerErrorModel struct {
	// Code - A more specific error code than was provided by the containing error.
	Code *string `json:"code,omitempty"`
	// InnerError - An object containing more specific information than the current object about the error.
	InnerError *InnerErrorModel `json:"innerError,omitempty"`
}

// MetadataDTO name - value pair of metadata.
type MetadataDTO struct {
	// Name - Metadata name.
	Name *string `json:"name,omitempty"`
	// Value - Metadata value.
	Value *string `json:"value,omitempty"`
}

// PromptDTO prompt for an answer.
type PromptDTO struct {
	// DisplayOrder - Index of the prompt - used in ordering of the prompts
	DisplayOrder *int32 `json:"displayOrder,omitempty"`
	// QnaID - Qna id corresponding to the prompt - if QnaId is present, QnADTO object is ignored.
	QnaID *int32 `json:"qnaId,omitempty"`
	// Qna - QnADTO - Either QnaId or QnADTO needs to be present in a PromptDTO object
	Qna *PromptDTOQna `json:"qna,omitempty"`
	// DisplayText - Text displayed to represent a follow up question prompt
	DisplayText *string `json:"displayText,omitempty"`
}

// PromptDTOQna qnADTO - Either QnaId or QnADTO needs to be present in a PromptDTO object
type PromptDTOQna struct {
	// ID - Unique id for the Q-A.
	ID *int32 `json:"id,omitempty"`
	// Answer - Answer text
	Answer *string `json:"answer,omitempty"`
	// Source - Source from which Q-A was indexed. eg. https://docs.microsoft.com/en-us/azure/cognitive-services/QnAMaker/FAQs
	Source *string `json:"source,omitempty"`
	// Questions - List of questions associated with the answer.
	Questions *[]string `json:"questions,omitempty"`
	// Metadata - List of metadata associated with the answer.
	Metadata *[]MetadataDTO `json:"metadata,omitempty"`
	// Context - Context of a QnA
	Context *QnADTOContext `json:"context,omitempty"`
}

// QnADTO q-A object.
type QnADTO struct {
	// ID - Unique id for the Q-A.
	ID *int32 `json:"id,omitempty"`
	// Answer - Answer text
	Answer *string `json:"answer,omitempty"`
	// Source - Source from which Q-A was indexed. eg. https://docs.microsoft.com/en-us/azure/cognitive-services/QnAMaker/FAQs
	Source *string `json:"source,omitempty"`
	// Questions - List of questions associated with the answer.
	Questions *[]string `json:"questions,omitempty"`
	// Metadata - List of metadata associated with the answer.
	Metadata *[]MetadataDTO `json:"metadata,omitempty"`
	// Context - Context of a QnA
	Context *QnADTOContext `json:"context,omitempty"`
}

// QnADTOContext context of a QnA
type QnADTOContext struct {
	// IsContextOnly - To mark if a prompt is relevant only with a previous question or not.
	// true - Do not include this QnA as search result for queries without context
	// false - ignores context and includes this QnA in search result
	IsContextOnly *bool `json:"isContextOnly,omitempty"`
	// Prompts - List of prompts associated with the answer.
	Prompts *[]PromptDTO `json:"prompts,omitempty"`
}

// QnASearchResult represents Search Result.
type QnASearchResult struct {
	// Questions - List of questions.
	Questions *[]string `json:"questions,omitempty"`
	// Answer - Answer.
	Answer *string `json:"answer,omitempty"`
	// Score - Search result score.
	Score *float64 `json:"score,omitempty"`
	// ID - Id of the QnA result.
	ID *int32 `json:"id,omitempty"`
	// Source - Source of QnA result.
	Source *string `json:"source,omitempty"`
	// Metadata - List of metadata.
	Metadata *[]MetadataDTO `json:"metadata,omitempty"`
	// Context - Context object of the QnA
	Context *QnASearchResultContext `json:"context,omitempty"`
}

// QnASearchResultContext context object of the QnA
type QnASearchResultContext struct {
	// IsContextOnly - To mark if a prompt is relevant only with a previous question or not.
	// true - Do not include this QnA as search result for queries without context
	// false - ignores context and includes this QnA in search result
	IsContextOnly *bool `json:"isContextOnly,omitempty"`
	// Prompts - List of prompts associated with the answer.
	Prompts *[]PromptDTO `json:"prompts,omitempty"`
}

// QnASearchResultList represents List of Question Answers.
type QnASearchResultList struct {
	autorest.Response `json:"-"`
	// Answers - Represents Search Result list.
	Answers *[]QnASearchResult `json:"answers,omitempty"`
}

// QueryContextDTO context object with previous QnA's information.
type QueryContextDTO struct {
	// PreviousQnaID - Previous QnA Id - qnaId of the top result.
	PreviousQnaID *string `json:"previousQnaId,omitempty"`
	// PreviousUserQuery - Previous user query.
	PreviousUserQuery *string `json:"previousUserQuery,omitempty"`
}

// QueryDTO POST body schema to query the knowledgebase.
type QueryDTO struct {
	// QnaID - Exact qnaId to fetch from the knowledgebase, this field takes priority over question.
	QnaID *string `json:"qnaId,omitempty"`
	// Question - User question to query against the knowledge base.
	Question *string `json:"question,omitempty"`
	// Top - Max number of answers to be returned for the question.
	Top *int32 `json:"top,omitempty"`
	// UserID - Unique identifier for the user. Optional parameter for telemetry. For more information, refer <a href="http://aka.ms/qnamaker-analytics#user-traffic" target="blank">Analytics and Telemetry</a>.
	UserID *string `json:"userId,omitempty"`
	// IsTest - Query against the test index.
	IsTest *bool `json:"isTest,omitempty"`
	// ScoreThreshold - Threshold for answers returned based on score.
	ScoreThreshold *float64 `json:"scoreThreshold,omitempty"`
	// Context - Context object with previous QnA's information.
	Context *QueryDTOContext `json:"context,omitempty"`
	// RankerType - Optional field. Set to 'QuestionOnly' for using a question only Ranker.
	RankerType *string `json:"rankerType,omitempty"`
	// StrictFilters - Find only answers that contain these metadata.
	StrictFilters *[]MetadataDTO `json:"strictFilters,omitempty"`
	// StrictFiltersCompoundOperationType - Optional field. Set to OR for using OR as Operation for Strict Filters. Possible values include: 'AND', 'OR'
	StrictFiltersCompoundOperationType StrictFiltersCompoundOperationType `json:"strictFiltersCompoundOperationType,omitempty"`
}

// QueryDTOContext context object with previous QnA's information.
type QueryDTOContext struct {
	// PreviousQnaID - Previous QnA Id - qnaId of the top result.
	PreviousQnaID *string `json:"previousQnaId,omitempty"`
	// PreviousUserQuery - Previous user query.
	PreviousUserQuery *string `json:"previousUserQuery,omitempty"`
}
