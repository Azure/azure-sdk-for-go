// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armplatformvalidation

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	// ExecutionPlanAPIVersion is the API version emitted in execution plan configurations.
	ExecutionPlanAPIVersion = "microsoft.validate/executionPlan.v0"

	// ExecutionPlanKind is the kind emitted in execution plan configurations.
	ExecutionPlanKind = "ExecutionPlan"
)

// ExecutionPlanConfiguration describes an execution plan configuration for
// ValidationExecutionPlanProperties.PlanConfigurationJSON.
type ExecutionPlanConfiguration struct {
	// Name is the execution plan metadata name.
	Name string

	// CertificationPackageReference describes the package to validate.
	CertificationPackageReference CertificationPackageReference

	// Steps contains the validation tests to execute.
	Steps []ValidationStep
}

// CertificationPackageReference describes the certification package used by an execution plan.
type CertificationPackageReference struct {
	// OSType is the operating system type of the package.
	OSType string

	// VMGenerationType is the virtual machine generation of the package.
	VMGenerationType string

	// ArchitectureType is the processor architecture of the package.
	ArchitectureType string

	// RecommendedVMSizes contains the virtual machine sizes recommended for validation.
	RecommendedVMSizes []string

	// StorageProfile describes the package's disk images.
	StorageProfile CertificationPackageStorageProfile

	// AdditionalProperties contains package metadata not represented by the other fields.
	AdditionalProperties map[string]any
}

// CertificationPackageStorageProfile describes the disk images in a certification package.
type CertificationPackageStorageProfile struct {
	// OSDiskImage is the operating system disk image.
	OSDiskImage CertificationPackageDiskImage

	// DataDiskImages contains any data disk images.
	DataDiskImages []CertificationPackageDiskImage
}

// CertificationPackageDiskImage identifies a disk image by its source VHD URI.
type CertificationPackageDiskImage struct {
	// SourceVHDURI is the URI of the source VHD.
	SourceVHDURI string
}

// ValidationStep describes a test invocation in an execution plan.
type ValidationStep struct {
	// Name identifies the step within the execution plan.
	Name string

	// TestRef is the complete validation test reference.
	TestRef string

	// Inputs contains test-specific input values.
	Inputs map[string]any
}

// NewValidationStep creates a test step using testRef exactly as supplied.
func NewValidationStep(name, testRef string) (ValidationStep, error) {
	step := ValidationStep{
		Name:    name,
		TestRef: testRef,
	}
	if err := step.validate(); err != nil {
		return ValidationStep{}, err
	}
	return step, nil
}

// ToJSON validates and serializes the configuration for use as
// ValidationExecutionPlanProperties.PlanConfigurationJSON.
func (c *ExecutionPlanConfiguration) ToJSON() (string, error) {
	if err := c.validate(); err != nil {
		return "", err
	}

	additionalProperties := c.CertificationPackageReference.AdditionalProperties
	if additionalProperties == nil {
		additionalProperties = map[string]any{}
	}
	dataDiskImages := c.CertificationPackageReference.StorageProfile.DataDiskImages
	if dataDiskImages == nil {
		dataDiskImages = []CertificationPackageDiskImage{}
	}

	configuration := executionPlanConfigurationJSON{
		APIVersion: ExecutionPlanAPIVersion,
		Kind:       ExecutionPlanKind,
		Metadata: executionPlanMetadataJSON{
			Name: c.Name,
		},
		Parameters: executionPlanParametersJSON{
			CertificationPackageReference: certificationPackageReferenceJSON{
				OSType:               c.CertificationPackageReference.OSType,
				VMGenerationType:     c.CertificationPackageReference.VMGenerationType,
				ArchitectureType:     c.CertificationPackageReference.ArchitectureType,
				RecommendedVMSizes:   c.CertificationPackageReference.RecommendedVMSizes,
				AdditionalProperties: additionalProperties,
				StorageProfile: certificationPackageStorageProfileJSON{
					OSDiskImage:    c.CertificationPackageReference.StorageProfile.OSDiskImage,
					DataDiskImages: dataDiskImages,
				},
			},
		},
		Authoring: executionPlanAuthoringJSON{
			Steps: c.Steps,
		},
	}

	data, err := json.Marshal(configuration)
	if err != nil {
		return "", fmt.Errorf("marshalling execution plan configuration: %w", err)
	}
	return string(data), nil
}

func (c *ExecutionPlanConfiguration) validate() error {
	if c == nil {
		return errors.New("execution plan configuration is nil")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("execution plan configuration name is required")
	}

	reference := c.CertificationPackageReference
	if strings.TrimSpace(reference.OSType) == "" {
		return errors.New("certification package reference OS type is required")
	}
	if strings.TrimSpace(reference.VMGenerationType) == "" {
		return errors.New("certification package reference VM generation type is required")
	}
	if strings.TrimSpace(reference.ArchitectureType) == "" {
		return errors.New("certification package reference architecture type is required")
	}
	if len(reference.RecommendedVMSizes) == 0 {
		return errors.New("certification package reference must include at least one recommended VM size")
	}
	for i, size := range reference.RecommendedVMSizes {
		if strings.TrimSpace(size) == "" {
			return fmt.Errorf("certification package reference recommended VM size at index %d is required", i)
		}
	}
	if err := reference.StorageProfile.OSDiskImage.validate("OS disk image"); err != nil {
		return err
	}
	for i, image := range reference.StorageProfile.DataDiskImages {
		if err := image.validate(fmt.Sprintf("data disk image at index %d", i)); err != nil {
			return err
		}
	}
	if len(c.Steps) == 0 {
		return errors.New("execution plan configuration must include at least one validation step")
	}
	for i := range c.Steps {
		if err := c.Steps[i].validate(); err != nil {
			return fmt.Errorf("validation step at index %d: %w", i, err)
		}
	}
	return nil
}

func (d CertificationPackageDiskImage) validate(field string) error {
	if strings.TrimSpace(d.SourceVHDURI) == "" {
		return fmt.Errorf("%s source VHD URI is required", field)
	}
	return nil
}

func (s ValidationStep) validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(s.TestRef) == "" {
		return errors.New("testRef is required")
	}
	return nil
}

type executionPlanConfigurationJSON struct {
	APIVersion string                      `json:"apiVersion"`
	Kind       string                      `json:"kind"`
	Metadata   executionPlanMetadataJSON   `json:"metadata"`
	Parameters executionPlanParametersJSON `json:"parameters"`
	Authoring  executionPlanAuthoringJSON  `json:"authoring"`
}

type executionPlanMetadataJSON struct {
	Name string `json:"name"`
}

type executionPlanParametersJSON struct {
	CertificationPackageReference certificationPackageReferenceJSON `json:"certificationPackageReference"`
}

type certificationPackageReferenceJSON struct {
	OSType               string                                 `json:"osType"`
	VMGenerationType     string                                 `json:"vmGenerationType"`
	ArchitectureType     string                                 `json:"architectureType"`
	RecommendedVMSizes   []string                               `json:"recommendedVMSizes"`
	StorageProfile       certificationPackageStorageProfileJSON `json:"storageProfile"`
	AdditionalProperties map[string]any                         `json:"additionalProperties"`
}

type certificationPackageStorageProfileJSON struct {
	OSDiskImage    CertificationPackageDiskImage   `json:"osDiskImage"`
	DataDiskImages []CertificationPackageDiskImage `json:"dataDiskImages"`
}

type executionPlanAuthoringJSON struct {
	Steps []ValidationStep `json:"steps"`
}

func (d CertificationPackageDiskImage) MarshalJSON() ([]byte, error) {
	type diskImageJSON struct {
		SourceVHDURI string `json:"sourceVhdUri"`
	}
	return json.Marshal(diskImageJSON{SourceVHDURI: d.SourceVHDURI})
}

func (s ValidationStep) MarshalJSON() ([]byte, error) {
	type validationStepJSON struct {
		Name    string         `json:"name"`
		Type    string         `json:"type"`
		TestRef string         `json:"testRef"`
		Inputs  map[string]any `json:"inputs"`
	}
	step := validationStepJSON{
		Name:    s.Name,
		Type:    "test",
		TestRef: s.TestRef,
		Inputs:  s.Inputs,
	}
	if s.Inputs != nil {
		return json.Marshal(step)
	}
	return json.Marshal(struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		TestRef string `json:"testRef"`
	}{
		Name:    step.Name,
		Type:    "test",
		TestRef: step.TestRef,
	})
}
