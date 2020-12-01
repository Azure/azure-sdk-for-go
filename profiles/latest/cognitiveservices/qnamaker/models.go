// +build go1.9

// Copyright 2020 Microsoft Corporation
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

package qnamaker

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v4.0/qnamaker"

type EnvironmentType = original.EnvironmentType

const (
	Prod EnvironmentType = original.Prod
	Test EnvironmentType = original.Test
)

type ErrorCodeType = original.ErrorCodeType

const (
	BadArgument       ErrorCodeType = original.BadArgument
	EndpointKeysError ErrorCodeType = original.EndpointKeysError
	ExtractionFailure ErrorCodeType = original.ExtractionFailure
	Forbidden         ErrorCodeType = original.Forbidden
	KbNotFound        ErrorCodeType = original.KbNotFound
	NotFound          ErrorCodeType = original.NotFound
	OperationNotFound ErrorCodeType = original.OperationNotFound
	QnaRuntimeError   ErrorCodeType = original.QnaRuntimeError
	QuotaExceeded     ErrorCodeType = original.QuotaExceeded
	ServiceError      ErrorCodeType = original.ServiceError
	SKULimitExceeded  ErrorCodeType = original.SKULimitExceeded
	Unauthorized      ErrorCodeType = original.Unauthorized
	Unspecified       ErrorCodeType = original.Unspecified
	ValidationFailure ErrorCodeType = original.ValidationFailure
)

type OperationStateType = original.OperationStateType

const (
	Failed     OperationStateType = original.Failed
	NotStarted OperationStateType = original.NotStarted
	Running    OperationStateType = original.Running
	Succeeded  OperationStateType = original.Succeeded
)

type ActiveLearningSettingsDTO = original.ActiveLearningSettingsDTO
type AlterationsClient = original.AlterationsClient
type AlterationsDTO = original.AlterationsDTO
type BaseClient = original.BaseClient
type ContextDTO = original.ContextDTO
type CreateKbDTO = original.CreateKbDTO
type CreateKbInputDTO = original.CreateKbInputDTO
type DeleteKbContentsDTO = original.DeleteKbContentsDTO
type EndpointKeysClient = original.EndpointKeysClient
type EndpointKeysDTO = original.EndpointKeysDTO
type EndpointSettingsClient = original.EndpointSettingsClient
type EndpointSettingsDTO = original.EndpointSettingsDTO
type EndpointSettingsDTOActiveLearning = original.EndpointSettingsDTOActiveLearning
type Error = original.Error
type ErrorResponse = original.ErrorResponse
type ErrorResponseError = original.ErrorResponseError
type FileDTO = original.FileDTO
type InnerErrorModel = original.InnerErrorModel
type KnowledgebaseClient = original.KnowledgebaseClient
type KnowledgebaseDTO = original.KnowledgebaseDTO
type KnowledgebasesDTO = original.KnowledgebasesDTO
type MetadataDTO = original.MetadataDTO
type Operation = original.Operation
type OperationsClient = original.OperationsClient
type PromptDTO = original.PromptDTO
type PromptDTOQna = original.PromptDTOQna
type QnADTO = original.QnADTO
type QnADTOContext = original.QnADTOContext
type QnADocumentsDTO = original.QnADocumentsDTO
type ReplaceKbDTO = original.ReplaceKbDTO
type UpdateContextDTO = original.UpdateContextDTO
type UpdateKbContentsDTO = original.UpdateKbContentsDTO
type UpdateKbOperationDTO = original.UpdateKbOperationDTO
type UpdateKbOperationDTOAdd = original.UpdateKbOperationDTOAdd
type UpdateKbOperationDTODelete = original.UpdateKbOperationDTODelete
type UpdateKbOperationDTOUpdate = original.UpdateKbOperationDTOUpdate
type UpdateMetadataDTO = original.UpdateMetadataDTO
type UpdateQnaDTO = original.UpdateQnaDTO
type UpdateQnaDTOContext = original.UpdateQnaDTOContext
type UpdateQnaDTOMetadata = original.UpdateQnaDTOMetadata
type UpdateQnaDTOQuestions = original.UpdateQnaDTOQuestions
type UpdateQuestionsDTO = original.UpdateQuestionsDTO
type WordAlterationsDTO = original.WordAlterationsDTO

func New(endpoint string) BaseClient {
	return original.New(endpoint)
}
func NewAlterationsClient(endpoint string) AlterationsClient {
	return original.NewAlterationsClient(endpoint)
}
func NewEndpointKeysClient(endpoint string) EndpointKeysClient {
	return original.NewEndpointKeysClient(endpoint)
}
func NewEndpointSettingsClient(endpoint string) EndpointSettingsClient {
	return original.NewEndpointSettingsClient(endpoint)
}
func NewKnowledgebaseClient(endpoint string) KnowledgebaseClient {
	return original.NewKnowledgebaseClient(endpoint)
}
func NewOperationsClient(endpoint string) OperationsClient {
	return original.NewOperationsClient(endpoint)
}
func NewWithoutDefaults(endpoint string) BaseClient {
	return original.NewWithoutDefaults(endpoint)
}
func PossibleEnvironmentTypeValues() []EnvironmentType {
	return original.PossibleEnvironmentTypeValues()
}
func PossibleErrorCodeTypeValues() []ErrorCodeType {
	return original.PossibleErrorCodeTypeValues()
}
func PossibleOperationStateTypeValues() []OperationStateType {
	return original.PossibleOperationStateTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}
