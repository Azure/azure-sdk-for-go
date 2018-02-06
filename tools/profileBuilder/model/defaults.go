package model

import (
	"os"
	"path"
	"strings"
)

// AzureSDKforGoLocation returns the default location for the Azure-SDK-for-Go to reside.
func AzureSDKforGoLocation() string {
	raw := path.Join(
		os.Getenv("GOPATH"),
		"src",
		"github.com",
		"Azure",
		"azure-sdk-for-go",
	)

	return strings.Replace(raw, "\\", "/", -1)
}

// DefaultOutputLocation establishes the location where profiles should be
// output to unless otherwise specified.
func DefaultOutputLocation() string {
	return path.Join(AzureSDKforGoLocation(), "profiles")
}

// DefaultInputRoot establishes the location where we expect to find the packages
// to create aliases for.
func DefaultInputRoot() string {
	return path.Join(AzureSDKforGoLocation(), "services")
}
