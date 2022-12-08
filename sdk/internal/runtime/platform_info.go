// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"fmt"
	"os"
	"runtime"
)

// PlatformInfo is the Go version and OS, formatted properly for insertion
// into a User-Agent string.
// NOTE: the ONLY function that should write to this variable is this func
var PlatformInfo = func() string {
	operatingSystem := runtime.GOOS // Default OS string
	switch operatingSystem {
	case "windows":
		operatingSystem = os.Getenv("OS") // Get more specific OS information
	case "linux": // accept default OS info
	case "freebsd": //  accept default OS info
	}
	return fmt.Sprintf("(%s; %s)", runtime.Version(), operatingSystem)
}()
