// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package textanalytics

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.1/textanalytics"

type BaseClient = original.BaseClient
type DetectedLanguage = original.DetectedLanguage
type DocumentStatistics = original.DocumentStatistics
type EntitiesBatchResult = original.EntitiesBatchResult
type EntitiesBatchResultItem = original.EntitiesBatchResultItem
type EntityRecord = original.EntityRecord
type ErrorRecord = original.ErrorRecord
type ErrorResponse = original.ErrorResponse
type InternalError = original.InternalError
type KeyPhraseBatchResult = original.KeyPhraseBatchResult
type KeyPhraseBatchResultItem = original.KeyPhraseBatchResultItem
type LanguageBatchInput = original.LanguageBatchInput
type LanguageBatchResult = original.LanguageBatchResult
type LanguageBatchResultItem = original.LanguageBatchResultItem
type LanguageInput = original.LanguageInput
type MatchRecord = original.MatchRecord
type MultiLanguageBatchInput = original.MultiLanguageBatchInput
type MultiLanguageInput = original.MultiLanguageInput
type RequestStatistics = original.RequestStatistics
type SentimentBatchResult = original.SentimentBatchResult
type SentimentBatchResultItem = original.SentimentBatchResultItem

func New(endpoint string) BaseClient {
	return original.New(endpoint)
}
func NewWithoutDefaults(endpoint string) BaseClient {
	return original.NewWithoutDefaults(endpoint)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
