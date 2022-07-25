//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	PremiumPageBlobAccessTierP10 = generated.PremiumPageBlobAccessTierP10
	PremiumPageBlobAccessTierP15 = generated.PremiumPageBlobAccessTierP15
	PremiumPageBlobAccessTierP20 = generated.PremiumPageBlobAccessTierP20
	PremiumPageBlobAccessTierP30 = generated.PremiumPageBlobAccessTierP30
	PremiumPageBlobAccessTierP4  = generated.PremiumPageBlobAccessTierP4
	PremiumPageBlobAccessTierP40 = generated.PremiumPageBlobAccessTierP40
	PremiumPageBlobAccessTierP50 = generated.PremiumPageBlobAccessTierP50
	PremiumPageBlobAccessTierP6  = generated.PremiumPageBlobAccessTierP6
	PremiumPageBlobAccessTierP60 = generated.PremiumPageBlobAccessTierP60
	PremiumPageBlobAccessTierP70 = generated.PremiumPageBlobAccessTierP70
	PremiumPageBlobAccessTierP80 = generated.PremiumPageBlobAccessTierP80
)

const (
	SequenceNumberActionTypeMax       = generated.SequenceNumberActionTypeMax
	SequenceNumberActionTypeUpdate    = generated.SequenceNumberActionTypeUpdate
	SequenceNumberActionTypeIncrement = generated.SequenceNumberActionTypeIncrement
)

// PossibleSequenceNumberActionTypeValues returns the possible values for the SequenceNumberActionType const type.
func PossibleSequenceNumberActionTypeValues() []SequenceNumberActionType {
	return []SequenceNumberActionType{
		SequenceNumberActionTypeMax,
		SequenceNumberActionTypeUpdate,
		SequenceNumberActionTypeIncrement,
	}
}
