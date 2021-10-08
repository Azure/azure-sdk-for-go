// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/spf13/pflag"
)

type GlobalFlags struct {
	OptionsPath string
}

func BindGlobalFlags(flags *pflag.FlagSet) {
	flags.Bool("version", false, "Show version number")
	flags.String("options-path", DefaultOptionPath, "Specify the autorest option file, relative to the root of azure-sdk-for-go")
}

func ParseGlobalFlags(flagSet *pflag.FlagSet) GlobalFlags {
	return GlobalFlags{
		OptionsPath: flags.GetString(flagSet, "options-path"),
	}
}
