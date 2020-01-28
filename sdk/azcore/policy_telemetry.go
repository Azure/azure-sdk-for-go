// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcore

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
)

// TelemetryOptions configures the telemetry policy's behavior.
type TelemetryOptions struct {
	// Value is a string prepended to each request's User-Agent and sent to the service.
	// The service records the user-agent in logs for diagnostics and tracking of client requests.
	Value string
}

type telemetryPolicy struct {
	telemetryValue string
}

// NewTelemetryPolicy creates a telemetry policy object that adds telemetry information to outgoing HTTP requests.
func NewTelemetryPolicy(o TelemetryOptions) Policy {
	b := &bytes.Buffer{}
	b.WriteString(o.Value)
	if b.Len() > 0 {
		b.WriteRune(' ')
	}
	b.WriteString(platformInfo)
	return &telemetryPolicy{telemetryValue: b.String()}
}

func (p telemetryPolicy) Do(ctx context.Context, req *Request) (*Response, error) {
	req.Request.Header.Set(HeaderUserAgent, p.telemetryValue)
	return req.Next(ctx)
}

// NOTE: the ONLY function that should write to this variable is this func
var platformInfo = func() string {
	operatingSystem := runtime.GOOS // Default OS string
	switch operatingSystem {
	case "windows":
		operatingSystem = os.Getenv("OS") // Get more specific OS information
	case "linux": // accept default OS info
	case "freebsd": //  accept default OS info
	}
	return fmt.Sprintf("(%s; %s)", runtime.Version(), operatingSystem)
}()
