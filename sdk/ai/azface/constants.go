//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface

// ServiceVersion - The version of the Face service to use
type ServiceVersion string

const (
	// ServiceVersionV1_0 - Version 1.0 of the Face service
	ServiceVersionV1_0 ServiceVersion = "v1.0"
)

// PossibleServiceVersionValues returns the possible values for the ServiceVersion const type.
func PossibleServiceVersionValues() []ServiceVersion {
	return []ServiceVersion{
		ServiceVersionV1_0,
	}
}

// DetectionModel - The detection model to use for face detection
type DetectionModel string

const (
	// DetectionModelDetection01 - Detection model 01
	DetectionModelDetection01 DetectionModel = "detection_01"
	// DetectionModelDetection02 - Detection model 02
	DetectionModelDetection02 DetectionModel = "detection_02"
	// DetectionModelDetection03 - Detection model 03
	DetectionModelDetection03 DetectionModel = "detection_03"
)

// PossibleDetectionModelValues returns the possible values for the DetectionModel const type.
func PossibleDetectionModelValues() []DetectionModel {
	return []DetectionModel{
		DetectionModelDetection01,
		DetectionModelDetection02,
		DetectionModelDetection03,
	}
}

// RecognitionModel - The recognition model to use for face recognition
type RecognitionModel string

const (
	// RecognitionModelRecognition01 - Recognition model 01
	RecognitionModelRecognition01 RecognitionModel = "recognition_01"
	// RecognitionModelRecognition02 - Recognition model 02
	RecognitionModelRecognition02 RecognitionModel = "recognition_02"
	// RecognitionModelRecognition03 - Recognition model 03
	RecognitionModelRecognition03 RecognitionModel = "recognition_03"
	// RecognitionModelRecognition04 - Recognition model 04
	RecognitionModelRecognition04 RecognitionModel = "recognition_04"
)

// PossibleRecognitionModelValues returns the possible values for the RecognitionModel const type.
func PossibleRecognitionModelValues() []RecognitionModel {
	return []RecognitionModel{
		RecognitionModelRecognition01,
		RecognitionModelRecognition02,
		RecognitionModelRecognition03,
		RecognitionModelRecognition04,
	}
}