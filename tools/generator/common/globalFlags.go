package common

import (
	"github.com/Azure/azure-sdk-for-go/tools/generator/flags"
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
	"github.com/spf13/pflag"
)

type GlobalFlags struct {
	OptionsPath string
}

func BindGlobalFlags(flags *pflag.FlagSet) {
	flags.Bool("version", false, "Show version number")
	flags.String("options-path", sdk.DefaultOptionsPath, "Specify the autorest option file, relative to the root of azure-sdk-for-go")
}

func ParseGlobalFlags(flagSet *pflag.FlagSet) GlobalFlags {
	return GlobalFlags{
		OptionsPath: flags.GetString(flagSet, "options-path"),
	}
}
