# Unreleased

## Breaking Changes

### Struct Changes

#### Removed Structs

1. ListAgreementTerms

### Signature Changes

#### Funcs

1. MarketplaceAgreementsClient.Cancel
	- Returns
		- From: AgreementTerms, error
		- To: OldAgreementTerms, error
1. MarketplaceAgreementsClient.CancelResponder
	- Returns
		- From: AgreementTerms, error
		- To: OldAgreementTerms, error
1. MarketplaceAgreementsClient.GetAgreement
	- Returns
		- From: AgreementTerms, error
		- To: OldAgreementTerms, error
1. MarketplaceAgreementsClient.GetAgreementResponder
	- Returns
		- From: AgreementTerms, error
		- To: OldAgreementTerms, error
1. MarketplaceAgreementsClient.List
	- Returns
		- From: ListAgreementTerms, error
		- To: AgreementTermsList, error
1. MarketplaceAgreementsClient.ListResponder
	- Returns
		- From: ListAgreementTerms, error
		- To: AgreementTermsList, error
1. MarketplaceAgreementsClient.Sign
	- Returns
		- From: AgreementTerms, error
		- To: OldAgreementTerms, error
1. MarketplaceAgreementsClient.SignResponder
	- Returns
		- From: AgreementTerms, error
		- To: OldAgreementTerms, error

## Additive Changes

### New Constants

1. State.Active
1. State.Canceled

### New Funcs

1. *OldAgreementTerms.UnmarshalJSON([]byte) error
1. OldAgreementTerms.MarshalJSON() ([]byte, error)
1. PossibleStateValues() []State

### Struct Changes

#### New Structs

1. AgreementTermsList
1. OldAgreementProperties
1. OldAgreementTerms

#### New Struct Fields

1. OperationDisplay.Description
