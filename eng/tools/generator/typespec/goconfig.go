// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package typespec

import (
	"errors"
	"regexp"
)

type GoEmitterOptions struct {
	Module           string `yaml:"module,omitempty"`
	ContainingModule string `yaml:"containing-module,omitempty"`
	EmitterOutputDir string `yaml:"emitter-output-dir,omitempty"`
	PackageDir       string `yaml:"package-dir,omitempty"`
	ServiceDir       string `yaml:"service-dir,omitempty"`
}

const moduleRegex = `^github.com/Azure/azure-sdk-for-go/sdk/` +
	`(` +
	`resourcemanager/\w+/arm\w+` + // either an ARM package (ie: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus)
	`|` +
	`.+?/az[^/]+` + // or a data plane package (ie, github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces)
	`)`

var (
	ErrModuleEmpty  = errors.New("`module` or `containing-module` is required")
	ErrModuleFormat = errors.New("`module` or `containing-module` must be in the format of github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/xxx/armxxx or github.com/Azure/azure-sdk-for-go/sdk/xxx/azxxx")
)

// Validate checks if the GoEmitterOptions is valid
func (o *GoEmitterOptions) Validate() error {
	if o.Module == "" && o.ContainingModule == "" {
		return ErrModuleEmpty
	}

	regex := regexp.MustCompile(moduleRegex)
	matched := regex.MatchString(o.Module) || regex.MatchString(o.ContainingModule)
	if !matched {
		return ErrModuleFormat
	}

	return nil
}
