//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatabox

import "encoding/json"

func unmarshalCommonJobDetailsClassification(rawMsg json.RawMessage) (CommonJobDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CommonJobDetailsClassification
	switch m["jobDetailsType"] {
	case string(ClassDiscriminatorDataBox):
		b = &JobDetails{}
	case string(ClassDiscriminatorDataBoxCustomerDisk):
		b = &CustomerDiskJobDetails{}
	case string(ClassDiscriminatorDataBoxDisk):
		b = &DiskJobDetails{}
	case string(ClassDiscriminatorDataBoxHeavy):
		b = &HeavyJobDetails{}
	default:
		b = &CommonJobDetails{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalCommonJobSecretsClassification(rawMsg json.RawMessage) (CommonJobSecretsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CommonJobSecretsClassification
	switch m["jobSecretsType"] {
	case string(ClassDiscriminatorDataBox):
		b = &JobSecrets{}
	case string(ClassDiscriminatorDataBoxCustomerDisk):
		b = &CustomerDiskJobSecrets{}
	case string(ClassDiscriminatorDataBoxDisk):
		b = &DiskJobSecrets{}
	case string(ClassDiscriminatorDataBoxHeavy):
		b = &HeavyJobSecrets{}
	default:
		b = &CommonJobSecrets{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalCommonScheduleAvailabilityRequestClassification(rawMsg json.RawMessage) (CommonScheduleAvailabilityRequestClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CommonScheduleAvailabilityRequestClassification
	switch m["skuName"] {
	case string(SKUNameDataBox):
		b = &ScheduleAvailabilityRequest{}
	case string(SKUNameDataBoxDisk):
		b = &DiskScheduleAvailabilityRequest{}
	case string(SKUNameDataBoxHeavy):
		b = &HeavyScheduleAvailabilityRequest{}
	default:
		b = &CommonScheduleAvailabilityRequest{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalCopyLogDetailsClassification(rawMsg json.RawMessage) (CopyLogDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b CopyLogDetailsClassification
	switch m["copyLogDetailsType"] {
	case string(ClassDiscriminatorDataBox):
		b = &AccountCopyLogDetails{}
	case string(ClassDiscriminatorDataBoxCustomerDisk):
		b = &CustomerDiskCopyLogDetails{}
	case string(ClassDiscriminatorDataBoxDisk):
		b = &DiskCopyLogDetails{}
	case string(ClassDiscriminatorDataBoxHeavy):
		b = &HeavyAccountCopyLogDetails{}
	default:
		b = &CopyLogDetails{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalCopyLogDetailsClassificationArray(rawMsg json.RawMessage) ([]CopyLogDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]CopyLogDetailsClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalCopyLogDetailsClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalDataAccountDetailsClassification(rawMsg json.RawMessage) (DataAccountDetailsClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b DataAccountDetailsClassification
	switch m["dataAccountType"] {
	case string(DataAccountTypeManagedDisk):
		b = &ManagedDiskDetails{}
	case string(DataAccountTypeStorageAccount):
		b = &StorageAccountDetails{}
	default:
		b = &DataAccountDetails{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalDatacenterAddressResponseClassification(rawMsg json.RawMessage) (DatacenterAddressResponseClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b DatacenterAddressResponseClassification
	switch m["datacenterAddressType"] {
	case string(DatacenterAddressTypeDatacenterAddressInstruction):
		b = &DatacenterAddressInstructionResponse{}
	case string(DatacenterAddressTypeDatacenterAddressLocation):
		b = &DatacenterAddressLocationResponse{}
	default:
		b = &DatacenterAddressResponse{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalValidationInputRequestClassification(rawMsg json.RawMessage) (ValidationInputRequestClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ValidationInputRequestClassification
	switch m["validationType"] {
	case string(ValidationInputDiscriminatorValidateAddress):
		b = &ValidateAddress{}
	case string(ValidationInputDiscriminatorValidateCreateOrderLimit):
		b = &CreateOrderLimitForSubscriptionValidationRequest{}
	case string(ValidationInputDiscriminatorValidateDataTransferDetails):
		b = &DataTransferDetailsValidationRequest{}
	case string(ValidationInputDiscriminatorValidatePreferences):
		b = &PreferencesValidationRequest{}
	case string(ValidationInputDiscriminatorValidateSKUAvailability):
		b = &SKUAvailabilityValidationRequest{}
	case string(ValidationInputDiscriminatorValidateSubscriptionIsAllowedToCreateJob):
		b = &SubscriptionIsAllowedToCreateJobValidationRequest{}
	default:
		b = &ValidationInputRequest{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalValidationInputRequestClassificationArray(rawMsg json.RawMessage) ([]ValidationInputRequestClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ValidationInputRequestClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalValidationInputRequestClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}

func unmarshalValidationInputResponseClassification(rawMsg json.RawMessage) (ValidationInputResponseClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b ValidationInputResponseClassification
	switch m["validationType"] {
	case string(ValidationInputDiscriminatorValidateAddress):
		b = &AddressValidationProperties{}
	case string(ValidationInputDiscriminatorValidateCreateOrderLimit):
		b = &CreateOrderLimitForSubscriptionValidationResponseProperties{}
	case string(ValidationInputDiscriminatorValidateDataTransferDetails):
		b = &DataTransferDetailsValidationResponseProperties{}
	case string(ValidationInputDiscriminatorValidatePreferences):
		b = &PreferencesValidationResponseProperties{}
	case string(ValidationInputDiscriminatorValidateSKUAvailability):
		b = &SKUAvailabilityValidationResponseProperties{}
	case string(ValidationInputDiscriminatorValidateSubscriptionIsAllowedToCreateJob):
		b = &SubscriptionIsAllowedToCreateJobValidationResponseProperties{}
	default:
		b = &ValidationInputResponse{}
	}
	return b, json.Unmarshal(rawMsg, b)
}

func unmarshalValidationInputResponseClassificationArray(rawMsg json.RawMessage) ([]ValidationInputResponseClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(rawMsg, &rawMessages); err != nil {
		return nil, err
	}
	fArray := make([]ValidationInputResponseClassification, len(rawMessages))
	for index, rawMessage := range rawMessages {
		f, err := unmarshalValidationInputResponseClassification(rawMessage)
		if err != nil {
			return nil, err
		}
		fArray[index] = f
	}
	return fArray, nil
}
