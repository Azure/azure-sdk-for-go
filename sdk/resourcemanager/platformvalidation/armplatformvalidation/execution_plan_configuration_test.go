// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package armplatformvalidation

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestExecutionPlanConfigurationToJSON(t *testing.T) {
	firstStep, err := NewValidationStep(
		"os-disk-size",
		"/providers/Microsoft.Validate/validationTests/os-disk-size/versions/1.0.0",
	)
	if err != nil {
		t.Fatal(err)
	}
	secondStep, err := NewValidationStep(
		"linux-quality-validation",
		"/providers/Microsoft.Validate/validationTests/linux-quality-validation/versions/1.0.0",
	)
	if err != nil {
		t.Fatal(err)
	}
	secondStep.Inputs = map[string]any{"concurrency": 1}

	configuration := ExecutionPlanConfiguration{
		Name: "contoso-linux-cert",
		CertificationPackageReference: CertificationPackageReference{
			OSType:             "Linux",
			VMGenerationType:   "V1",
			ArchitectureType:   "X64",
			RecommendedVMSizes: []string{"Standard_D4s_v3"},
			StorageProfile: CertificationPackageStorageProfile{
				OSDiskImage: CertificationPackageDiskImage{
					SourceVHDURI: "https://contoso.example/img.vhd",
				},
			},
		},
		Steps: []ValidationStep{firstStep, secondStep},
	}

	actualJSON, err := configuration.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	const expectedJSON = `{
		"apiVersion": "microsoft.validate/executionPlan.v0",
		"kind": "ExecutionPlan",
		"metadata": {"name": "contoso-linux-cert"},
		"parameters": {"certificationPackageReference": {
			"osType": "Linux",
			"vmGenerationType": "V1",
			"architectureType": "X64",
			"recommendedVMSizes": ["Standard_D4s_v3"],
			"storageProfile": {
				"osDiskImage": {"sourceVhdUri": "https://contoso.example/img.vhd"},
				"dataDiskImages": []
			},
			"additionalProperties": {}
		}},
		"authoring": {"steps": [
			{
				"name": "os-disk-size",
				"type": "test",
				"testRef": "/providers/Microsoft.Validate/validationTests/os-disk-size/versions/1.0.0"
			},
			{
				"name": "linux-quality-validation",
				"type": "test",
				"testRef": "/providers/Microsoft.Validate/validationTests/linux-quality-validation/versions/1.0.0",
				"inputs": {"concurrency": 1}
			}
		]}
	}`

	var actual, expected any
	if err := json.Unmarshal([]byte(actualJSON), &actual); err != nil {
		t.Fatalf("unmarshalling actual JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("unmarshalling expected JSON: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("JSON mismatch\nactual: %s\nexpected: %s", actualJSON, expectedJSON)
	}
}

func TestNewValidationStepPreservesTestRef(t *testing.T) {
	const testRef = "  /custom/provider/tests/example/versions/latest  "
	step, err := NewValidationStep("example", testRef)
	if err != nil {
		t.Fatal(err)
	}
	if step.TestRef != testRef {
		t.Fatalf("testRef was changed: got %q, want %q", step.TestRef, testRef)
	}
}

func TestExecutionPlanConfigurationValidation(t *testing.T) {
	validConfiguration := func() ExecutionPlanConfiguration {
		return ExecutionPlanConfiguration{
			Name: "plan",
			CertificationPackageReference: CertificationPackageReference{
				OSType:             "Linux",
				VMGenerationType:   "V1",
				ArchitectureType:   "X64",
				RecommendedVMSizes: []string{"Standard_D4s_v3"},
				StorageProfile: CertificationPackageStorageProfile{
					OSDiskImage: CertificationPackageDiskImage{SourceVHDURI: "https://example.test/os.vhd"},
				},
			},
			Steps: []ValidationStep{{
				Name:    "test",
				TestRef: "/providers/Example/validationTests/test/versions/1.0.0",
			}},
		}
	}

	tests := []struct {
		name        string
		mutate      func(*ExecutionPlanConfiguration)
		errorSubstr string
	}{
		{
			name:        "missing name",
			mutate:      func(c *ExecutionPlanConfiguration) { c.Name = "" },
			errorSubstr: "name is required",
		},
		{
			name:        "missing OS type",
			mutate:      func(c *ExecutionPlanConfiguration) { c.CertificationPackageReference.OSType = "" },
			errorSubstr: "OS type is required",
		},
		{
			name: "missing recommended VM sizes",
			mutate: func(c *ExecutionPlanConfiguration) {
				c.CertificationPackageReference.RecommendedVMSizes = nil
			},
			errorSubstr: "recommended VM size",
		},
		{
			name: "missing OS disk image URI",
			mutate: func(c *ExecutionPlanConfiguration) {
				c.CertificationPackageReference.StorageProfile.OSDiskImage.SourceVHDURI = ""
			},
			errorSubstr: "OS disk image source VHD URI is required",
		},
		{
			name:        "missing steps",
			mutate:      func(c *ExecutionPlanConfiguration) { c.Steps = nil },
			errorSubstr: "at least one validation step",
		},
		{
			name:        "missing testRef",
			mutate:      func(c *ExecutionPlanConfiguration) { c.Steps[0].TestRef = "" },
			errorSubstr: "testRef is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			configuration := validConfiguration()
			test.mutate(&configuration)
			_, err := configuration.ToJSON()
			if err == nil {
				t.Fatal("expected validation error")
			}
			if !strings.Contains(err.Error(), test.errorSubstr) {
				t.Fatalf("error %q does not contain %q", err, test.errorSubstr)
			}
		})
	}
}

func TestExecutionPlanConfigurationNil(t *testing.T) {
	var configuration *ExecutionPlanConfiguration
	_, err := configuration.ToJSON()
	if err == nil || !strings.Contains(err.Error(), "nil") {
		t.Fatalf("expected nil configuration error, got %v", err)
	}
}

func TestNewValidationStepValidation(t *testing.T) {
	tests := []struct {
		name     string
		stepName string
		testRef  string
	}{
		{name: "missing name", testRef: "/complete/test/ref"},
		{name: "missing testRef", stepName: "test"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewValidationStep(test.stepName, test.testRef); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
