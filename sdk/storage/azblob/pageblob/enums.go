//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	PremiumPageBlobAccessTierP10 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP10
	PremiumPageBlobAccessTierP15 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP15
	PremiumPageBlobAccessTierP20 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP20
	PremiumPageBlobAccessTierP30 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP30
	PremiumPageBlobAccessTierP4  PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP4
	PremiumPageBlobAccessTierP40 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP40
	PremiumPageBlobAccessTierP50 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP50
	PremiumPageBlobAccessTierP6  PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP6
	PremiumPageBlobAccessTierP60 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP60
	PremiumPageBlobAccessTierP70 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP70
	PremiumPageBlobAccessTierP80 PremiumPageBlobAccessTier = generated.PremiumPageBlobAccessTierP80
)

const (
	SequenceNumberActionTypeMax       SequenceNumberActionType = "max"
	SequenceNumberActionTypeUpdate    SequenceNumberActionType = "update"
	SequenceNumberActionTypeIncrement SequenceNumberActionType = "increment"
)

// PossibleSequenceNumberActionTypeValues returns the possible values for the SequenceNumberActionType const type.
func PossibleSequenceNumberActionTypeValues() []SequenceNumberActionType {
	return []SequenceNumberActionType{
		SequenceNumberActionTypeMax,
		SequenceNumberActionTypeUpdate,
		SequenceNumberActionTypeIncrement,
	}
}
