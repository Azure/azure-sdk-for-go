package blockblob

import "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"

const (
	AccessTierArchive AccessTier = generated.AccessTierArchive
	AccessTierCool    AccessTier = generated.AccessTierCool
	AccessTierHot     AccessTier = generated.AccessTierHot
	AccessTierP10     AccessTier = generated.AccessTierP10
	AccessTierP15     AccessTier = generated.AccessTierP15
	AccessTierP20     AccessTier = generated.AccessTierP20
	AccessTierP30     AccessTier = generated.AccessTierP30
	AccessTierP4      AccessTier = generated.AccessTierP4
	AccessTierP40     AccessTier = generated.AccessTierP40
	AccessTierP50     AccessTier = generated.AccessTierP50
	AccessTierP6      AccessTier = generated.AccessTierP6
	AccessTierP60     AccessTier = generated.AccessTierP60
	AccessTierP70     AccessTier = generated.AccessTierP70
	AccessTierP80     AccessTier = generated.AccessTierP80
	AccessTierPremium AccessTier = generated.AccessTierPremium
)

type BlockListType = generated.BlockListType

const (
	BlockListTypeCommitted   BlockListType = "committed"
	BlockListTypeUncommitted BlockListType = "uncommitted"
	BlockListTypeAll         BlockListType = "all"
)

// PossibleBlockListTypeValues returns the possible values for the BlockListType const type.
func PossibleBlockListTypeValues() []BlockListType {
	return []BlockListType{
		BlockListTypeCommitted,
		BlockListTypeUncommitted,
		BlockListTypeAll,
	}
}

const (
	ImmutabilityPolicyModeMutable  ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeMutable
	ImmutabilityPolicyModeUnlocked ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeUnlocked
	ImmutabilityPolicyModeLocked   ImmutabilityPolicyMode = generated.BlobImmutabilityPolicyModeLocked
)

// PossibleBlobImmutabilityPolicyModeValues returns the possible values for the BlobImmutabilityPolicyMode const type.
func PossibleBlobImmutabilityPolicyModeValues() []ImmutabilityPolicyMode {
	return []ImmutabilityPolicyMode{
		ImmutabilityPolicyModeMutable,
		ImmutabilityPolicyModeUnlocked,
		ImmutabilityPolicyModeLocked,
	}
}
