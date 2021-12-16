//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/eng/tools/profileBuilder

package translatortext

import original "github.com/Azure/azure-sdk-for-go/services/preview/cognitiveservices/v1.0_preview.1/translatortext"

type Code = original.Code

const (
	InternalServerError Code = original.InternalServerError
	InvalidArgument     Code = original.InvalidArgument
	InvalidRequest      Code = original.InvalidRequest
	RequestRateTooHigh  Code = original.RequestRateTooHigh
	ResourceNotFound    Code = original.ResourceNotFound
	ServiceUnavailable  Code = original.ServiceUnavailable
	Unauthorized        Code = original.Unauthorized
)

type Status = original.Status

const (
	Cancelled  Status = original.Cancelled
	Cancelling Status = original.Cancelling
	Failed     Status = original.Failed
	NotStarted Status = original.NotStarted
	Running    Status = original.Running
	Succeeded  Status = original.Succeeded
)

type Status1 = original.Status1

const (
	Status1Cancelled  Status1 = original.Status1Cancelled
	Status1Cancelling Status1 = original.Status1Cancelling
	Status1Failed     Status1 = original.Status1Failed
	Status1NotStarted Status1 = original.Status1NotStarted
	Status1Running    Status1 = original.Status1Running
	Status1Succeeded  Status1 = original.Status1Succeeded
)

type StorageSource = original.StorageSource

const (
	AzureBlob StorageSource = original.AzureBlob
)

type StorageSource1 = original.StorageSource1

const (
	StorageSource1AzureBlob StorageSource1 = original.StorageSource1AzureBlob
)

type StorageType = original.StorageType

const (
	File   StorageType = original.File
	Folder StorageType = original.Folder
)

type BaseClient = original.BaseClient
type BatchRequest = original.BatchRequest
type BatchStatusDetail = original.BatchStatusDetail
type BatchStatusResponse = original.BatchStatusResponse
type BatchSubmissionRequest = original.BatchSubmissionRequest
type DocumentFilter = original.DocumentFilter
type DocumentStatusDetail = original.DocumentStatusDetail
type DocumentStatusResponse = original.DocumentStatusResponse
type ErrorResponseV2 = original.ErrorResponseV2
type ErrorV2 = original.ErrorV2
type FileFormat = original.FileFormat
type FileFormatListResult = original.FileFormatListResult
type Glossary = original.Glossary
type InnerErrorV2 = original.InnerErrorV2
type SourceInput = original.SourceInput
type StatusSummary = original.StatusSummary
type StorageSourceListResult = original.StorageSourceListResult
type TargetInput = original.TargetInput
type TranslationClient = original.TranslationClient

func New() BaseClient {
	return original.New()
}
func NewTranslationClient() TranslationClient {
	return original.NewTranslationClient()
}
func NewWithoutDefaults() BaseClient {
	return original.NewWithoutDefaults()
}
func PossibleCodeValues() []Code {
	return original.PossibleCodeValues()
}
func PossibleStatus1Values() []Status1 {
	return original.PossibleStatus1Values()
}
func PossibleStatusValues() []Status {
	return original.PossibleStatusValues()
}
func PossibleStorageSource1Values() []StorageSource1 {
	return original.PossibleStorageSource1Values()
}
func PossibleStorageSourceValues() []StorageSource {
	return original.PossibleStorageSourceValues()
}
func PossibleStorageTypeValues() []StorageType {
	return original.PossibleStorageTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
