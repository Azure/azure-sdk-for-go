//go:build go1.9
// +build go1.9

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package customimagesearch

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/customimagesearch"

const (
	DefaultEndpoint = original.DefaultEndpoint
)

type ErrorCode = original.ErrorCode

const (
	InsufficientAuthorization ErrorCode = original.InsufficientAuthorization
	InvalidAuthorization      ErrorCode = original.InvalidAuthorization
	InvalidRequest            ErrorCode = original.InvalidRequest
	None                      ErrorCode = original.None
	RateLimitExceeded         ErrorCode = original.RateLimitExceeded
	ServerError               ErrorCode = original.ServerError
)

type ErrorSubCode = original.ErrorSubCode

const (
	AuthorizationDisabled   ErrorSubCode = original.AuthorizationDisabled
	AuthorizationExpired    ErrorSubCode = original.AuthorizationExpired
	AuthorizationMissing    ErrorSubCode = original.AuthorizationMissing
	AuthorizationRedundancy ErrorSubCode = original.AuthorizationRedundancy
	Blocked                 ErrorSubCode = original.Blocked
	HTTPNotAllowed          ErrorSubCode = original.HTTPNotAllowed
	NotImplemented          ErrorSubCode = original.NotImplemented
	ParameterInvalidValue   ErrorSubCode = original.ParameterInvalidValue
	ParameterMissing        ErrorSubCode = original.ParameterMissing
	ResourceError           ErrorSubCode = original.ResourceError
	UnexpectedError         ErrorSubCode = original.UnexpectedError
)

type Freshness = original.Freshness

const (
	Day   Freshness = original.Day
	Month Freshness = original.Month
	Week  Freshness = original.Week
)

type ImageAspect = original.ImageAspect

const (
	All    ImageAspect = original.All
	Square ImageAspect = original.Square
	Tall   ImageAspect = original.Tall
	Wide   ImageAspect = original.Wide
)

type ImageColor = original.ImageColor

const (
	Black      ImageColor = original.Black
	Blue       ImageColor = original.Blue
	Brown      ImageColor = original.Brown
	ColorOnly  ImageColor = original.ColorOnly
	Gray       ImageColor = original.Gray
	Green      ImageColor = original.Green
	Monochrome ImageColor = original.Monochrome
	Orange     ImageColor = original.Orange
	Pink       ImageColor = original.Pink
	Purple     ImageColor = original.Purple
	Red        ImageColor = original.Red
	Teal       ImageColor = original.Teal
	White      ImageColor = original.White
	Yellow     ImageColor = original.Yellow
)

type ImageContent = original.ImageContent

const (
	Face     ImageContent = original.Face
	Portrait ImageContent = original.Portrait
)

type ImageLicense = original.ImageLicense

const (
	ImageLicenseAll                ImageLicense = original.ImageLicenseAll
	ImageLicenseAny                ImageLicense = original.ImageLicenseAny
	ImageLicenseModify             ImageLicense = original.ImageLicenseModify
	ImageLicenseModifyCommercially ImageLicense = original.ImageLicenseModifyCommercially
	ImageLicensePublic             ImageLicense = original.ImageLicensePublic
	ImageLicenseShare              ImageLicense = original.ImageLicenseShare
	ImageLicenseShareCommercially  ImageLicense = original.ImageLicenseShareCommercially
)

type ImageSize = original.ImageSize

const (
	ImageSizeAll       ImageSize = original.ImageSizeAll
	ImageSizeLarge     ImageSize = original.ImageSizeLarge
	ImageSizeMedium    ImageSize = original.ImageSizeMedium
	ImageSizeSmall     ImageSize = original.ImageSizeSmall
	ImageSizeWallpaper ImageSize = original.ImageSizeWallpaper
)

type ImageType = original.ImageType

const (
	AnimatedGif ImageType = original.AnimatedGif
	Clipart     ImageType = original.Clipart
	Line        ImageType = original.Line
	Photo       ImageType = original.Photo
	Shopping    ImageType = original.Shopping
	Transparent ImageType = original.Transparent
)

type SafeSearch = original.SafeSearch

const (
	Moderate SafeSearch = original.Moderate
	Off      SafeSearch = original.Off
	Strict   SafeSearch = original.Strict
)

type Type = original.Type

const (
	TypeAnswer              Type = original.TypeAnswer
	TypeCreativeWork        Type = original.TypeCreativeWork
	TypeErrorResponse       Type = original.TypeErrorResponse
	TypeIdentifiable        Type = original.TypeIdentifiable
	TypeImageObject         Type = original.TypeImageObject
	TypeImages              Type = original.TypeImages
	TypeMediaObject         Type = original.TypeMediaObject
	TypeResponse            Type = original.TypeResponse
	TypeResponseBase        Type = original.TypeResponseBase
	TypeSearchResultsAnswer Type = original.TypeSearchResultsAnswer
	TypeThing               Type = original.TypeThing
	TypeWebPage             Type = original.TypeWebPage
)

type Answer = original.Answer
type BaseClient = original.BaseClient
type BasicAnswer = original.BasicAnswer
type BasicCreativeWork = original.BasicCreativeWork
type BasicIdentifiable = original.BasicIdentifiable
type BasicMediaObject = original.BasicMediaObject
type BasicResponse = original.BasicResponse
type BasicResponseBase = original.BasicResponseBase
type BasicSearchResultsAnswer = original.BasicSearchResultsAnswer
type BasicThing = original.BasicThing
type CreativeWork = original.CreativeWork
type CustomInstanceClient = original.CustomInstanceClient
type Error = original.Error
type ErrorResponse = original.ErrorResponse
type Identifiable = original.Identifiable
type ImageObject = original.ImageObject
type Images = original.Images
type MediaObject = original.MediaObject
type Query = original.Query
type Response = original.Response
type ResponseBase = original.ResponseBase
type SearchResultsAnswer = original.SearchResultsAnswer
type Thing = original.Thing
type WebPage = original.WebPage

func New() BaseClient {
	return original.New()
}
func NewCustomInstanceClient() CustomInstanceClient {
	return original.NewCustomInstanceClient()
}
func NewWithoutDefaults(endpoint string) BaseClient {
	return original.NewWithoutDefaults(endpoint)
}
func PossibleErrorCodeValues() []ErrorCode {
	return original.PossibleErrorCodeValues()
}
func PossibleErrorSubCodeValues() []ErrorSubCode {
	return original.PossibleErrorSubCodeValues()
}
func PossibleFreshnessValues() []Freshness {
	return original.PossibleFreshnessValues()
}
func PossibleImageAspectValues() []ImageAspect {
	return original.PossibleImageAspectValues()
}
func PossibleImageColorValues() []ImageColor {
	return original.PossibleImageColorValues()
}
func PossibleImageContentValues() []ImageContent {
	return original.PossibleImageContentValues()
}
func PossibleImageLicenseValues() []ImageLicense {
	return original.PossibleImageLicenseValues()
}
func PossibleImageSizeValues() []ImageSize {
	return original.PossibleImageSizeValues()
}
func PossibleImageTypeValues() []ImageType {
	return original.PossibleImageTypeValues()
}
func PossibleSafeSearchValues() []SafeSearch {
	return original.PossibleSafeSearchValues()
}
func PossibleTypeValues() []Type {
	return original.PossibleTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/preview"
}
func Version() string {
	return original.Version()
}
