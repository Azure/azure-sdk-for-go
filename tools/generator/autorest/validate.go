package autorest

import "github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"

// MetadataValidateFunc is a function that validates a metadata is legal or not
type MetadataValidateFunc func(tag string, metadata model.Metadata) error

// ValidateMetadata validates the tag and metadata using the given validators
func ValidateMetadata(validators []MetadataValidateFunc, tag string, metadata model.Metadata) []error {
	if len(validators) == 0 {
		return nil
	}

	var errors []error
	for _, validator := range validators {
		if validator == nil {
			continue
		}
		if err := validator(tag, metadata); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
