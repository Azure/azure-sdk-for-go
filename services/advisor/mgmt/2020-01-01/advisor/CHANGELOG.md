# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. SetObject

### Signature Changes

#### Funcs

1. RecommendationMetadataClient.Get
	- Returns
		- From: SetObject, error
		- To: MetadataEntity, error
1. RecommendationMetadataClient.GetResponder
	- Returns
		- From: SetObject, error
		- To: MetadataEntity, error
1. SuppressionsClient.Get
	- Returns
		- From: SetObject, error
		- To: SuppressionContract, error
1. SuppressionsClient.GetResponder
	- Returns
		- From: SetObject, error
		- To: SuppressionContract, error

## Additive Changes

### New Funcs

1. ResourceMetadata.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Struct Fields

1. MetadataEntity.autorest.Response
1. RecommendationProperties.Actions
1. RecommendationProperties.Description
1. RecommendationProperties.ExposedMetadataProperties
1. RecommendationProperties.Label
1. RecommendationProperties.LearnMoreLink
1. RecommendationProperties.PotentialBenefits
1. RecommendationProperties.Remediation
1. ResourceMetadata.Action
1. ResourceMetadata.Plural
1. ResourceMetadata.Singular
