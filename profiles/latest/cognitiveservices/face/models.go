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

package face

import original "github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.0/face"

type Client = original.Client
type AccessoryType = original.AccessoryType

const (
	Glasses  AccessoryType = original.Glasses
	HeadWear AccessoryType = original.HeadWear
	Mask     AccessoryType = original.Mask
)

type AttributeType = original.AttributeType

const (
	AttributeTypeAccessories AttributeType = original.AttributeTypeAccessories
	AttributeTypeAge         AttributeType = original.AttributeTypeAge
	AttributeTypeBlur        AttributeType = original.AttributeTypeBlur
	AttributeTypeEmotion     AttributeType = original.AttributeTypeEmotion
	AttributeTypeExposure    AttributeType = original.AttributeTypeExposure
	AttributeTypeFacialHair  AttributeType = original.AttributeTypeFacialHair
	AttributeTypeGender      AttributeType = original.AttributeTypeGender
	AttributeTypeGlasses     AttributeType = original.AttributeTypeGlasses
	AttributeTypeHair        AttributeType = original.AttributeTypeHair
	AttributeTypeHeadPose    AttributeType = original.AttributeTypeHeadPose
	AttributeTypeMakeup      AttributeType = original.AttributeTypeMakeup
	AttributeTypeNoise       AttributeType = original.AttributeTypeNoise
	AttributeTypeOcclusion   AttributeType = original.AttributeTypeOcclusion
	AttributeTypeSmile       AttributeType = original.AttributeTypeSmile
)

type AzureRegions = original.AzureRegions

const (
	Australiaeast  AzureRegions = original.Australiaeast
	Brazilsouth    AzureRegions = original.Brazilsouth
	Eastasia       AzureRegions = original.Eastasia
	Eastus         AzureRegions = original.Eastus
	Eastus2        AzureRegions = original.Eastus2
	Northeurope    AzureRegions = original.Northeurope
	Southcentralus AzureRegions = original.Southcentralus
	Southeastasia  AzureRegions = original.Southeastasia
	Westcentralus  AzureRegions = original.Westcentralus
	Westeurope     AzureRegions = original.Westeurope
	Westus         AzureRegions = original.Westus
	Westus2        AzureRegions = original.Westus2
)

type BlurLevel = original.BlurLevel

const (
	High   BlurLevel = original.High
	Low    BlurLevel = original.Low
	Medium BlurLevel = original.Medium
)

type ExposureLevel = original.ExposureLevel

const (
	GoodExposure  ExposureLevel = original.GoodExposure
	OverExposure  ExposureLevel = original.OverExposure
	UnderExposure ExposureLevel = original.UnderExposure
)

type FindSimilarMatchMode = original.FindSimilarMatchMode

const (
	MatchFace   FindSimilarMatchMode = original.MatchFace
	MatchPerson FindSimilarMatchMode = original.MatchPerson
)

type Gender = original.Gender

const (
	Female     Gender = original.Female
	Genderless Gender = original.Genderless
	Male       Gender = original.Male
)

type GlassesType = original.GlassesType

const (
	NoGlasses       GlassesType = original.NoGlasses
	ReadingGlasses  GlassesType = original.ReadingGlasses
	Sunglasses      GlassesType = original.Sunglasses
	SwimmingGoggles GlassesType = original.SwimmingGoggles
)

type HairColorType = original.HairColorType

const (
	Black   HairColorType = original.Black
	Blond   HairColorType = original.Blond
	Brown   HairColorType = original.Brown
	Gray    HairColorType = original.Gray
	Other   HairColorType = original.Other
	Red     HairColorType = original.Red
	Unknown HairColorType = original.Unknown
	White   HairColorType = original.White
)

type NoiseLevel = original.NoiseLevel

const (
	NoiseLevelHigh   NoiseLevel = original.NoiseLevelHigh
	NoiseLevelLow    NoiseLevel = original.NoiseLevelLow
	NoiseLevelMedium NoiseLevel = original.NoiseLevelMedium
)

type TrainingStatusType = original.TrainingStatusType

const (
	Failed     TrainingStatusType = original.Failed
	Nonstarted TrainingStatusType = original.Nonstarted
	Running    TrainingStatusType = original.Running
	Succeeded  TrainingStatusType = original.Succeeded
)

type Accessory = original.Accessory
type APIError = original.APIError
type Attributes = original.Attributes
type Blur = original.Blur
type Coordinate = original.Coordinate
type DetectedFace = original.DetectedFace
type Emotion = original.Emotion
type Error = original.Error
type Exposure = original.Exposure
type FacialHair = original.FacialHair
type FindSimilarRequest = original.FindSimilarRequest
type GroupRequest = original.GroupRequest
type GroupResult = original.GroupResult
type Hair = original.Hair
type HairColor = original.HairColor
type HeadPose = original.HeadPose
type IdentifyCandidate = original.IdentifyCandidate
type IdentifyRequest = original.IdentifyRequest
type IdentifyResult = original.IdentifyResult
type ImageURL = original.ImageURL
type Landmarks = original.Landmarks
type List = original.List
type ListDetectedFace = original.ListDetectedFace
type ListIdentifyResult = original.ListIdentifyResult
type ListList = original.ListList
type ListPerson = original.ListPerson
type ListPersonGroup = original.ListPersonGroup
type ListSimilarFace = original.ListSimilarFace
type Makeup = original.Makeup
type NameAndUserDataContract = original.NameAndUserDataContract
type Noise = original.Noise
type Occlusion = original.Occlusion
type PersistedFace = original.PersistedFace
type Person = original.Person
type PersonGroup = original.PersonGroup
type Rectangle = original.Rectangle
type SimilarFace = original.SimilarFace
type TrainingStatus = original.TrainingStatus
type UpdatePersonFaceRequest = original.UpdatePersonFaceRequest
type VerifyFaceToFaceRequest = original.VerifyFaceToFaceRequest
type VerifyFaceToPersonRequest = original.VerifyFaceToPersonRequest
type VerifyResult = original.VerifyResult
type ListClient = original.ListClient
type BaseClient = original.BaseClient
type PersonGroupPersonClient = original.PersonGroupPersonClient
type PersonGroupClient = original.PersonGroupClient

func NewPersonGroupClient(azureRegion AzureRegions) PersonGroupClient {
	return original.NewPersonGroupClient(azureRegion)
}
func NewClient(azureRegion AzureRegions) Client {
	return original.NewClient(azureRegion)
}
func PossibleAccessoryTypeValues() []AccessoryType {
	return original.PossibleAccessoryTypeValues()
}
func PossibleAttributeTypeValues() []AttributeType {
	return original.PossibleAttributeTypeValues()
}
func PossibleAzureRegionsValues() []AzureRegions {
	return original.PossibleAzureRegionsValues()
}
func PossibleBlurLevelValues() []BlurLevel {
	return original.PossibleBlurLevelValues()
}
func PossibleExposureLevelValues() []ExposureLevel {
	return original.PossibleExposureLevelValues()
}
func PossibleFindSimilarMatchModeValues() []FindSimilarMatchMode {
	return original.PossibleFindSimilarMatchModeValues()
}
func PossibleGenderValues() []Gender {
	return original.PossibleGenderValues()
}
func PossibleGlassesTypeValues() []GlassesType {
	return original.PossibleGlassesTypeValues()
}
func PossibleHairColorTypeValues() []HairColorType {
	return original.PossibleHairColorTypeValues()
}
func PossibleNoiseLevelValues() []NoiseLevel {
	return original.PossibleNoiseLevelValues()
}
func PossibleTrainingStatusTypeValues() []TrainingStatusType {
	return original.PossibleTrainingStatusTypeValues()
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}
func NewListClient(azureRegion AzureRegions) ListClient {
	return original.NewListClient(azureRegion)
}
func New(azureRegion AzureRegions) BaseClient {
	return original.New(azureRegion)
}
func NewWithoutDefaults(azureRegion AzureRegions) BaseClient {
	return original.NewWithoutDefaults(azureRegion)
}
func NewPersonGroupPersonClient(azureRegion AzureRegions) PersonGroupPersonClient {
	return original.NewPersonGroupPersonClient(azureRegion)
}
