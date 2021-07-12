// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package flags

import (
	"log"

	"github.com/spf13/pflag"
)

func GetBool(f *pflag.FlagSet, name string) bool {
	b, err := f.GetBool(name)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func GetString(f *pflag.FlagSet, name string) string {
	s, err := f.GetString(name)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func GetStringSlice(f *pflag.FlagSet, name string) []string {
	s, err := f.GetStringSlice(name)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func GetInt(f *pflag.FlagSet, name string) int {
	i, err := f.GetInt(name)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func GetIntSlice(f *pflag.FlagSet, name string) []int {
	s, err := f.GetIntSlice(name)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
