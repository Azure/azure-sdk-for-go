package azblob

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// TelemetryOptions configures the telemetry policy's behavior.
type TelemetryOptions struct {
	// Value is a string prepended to each request's User-Agent and sent to the service.
	// The service records the user-agent in logs for diagnostics and tracking of client requests.
	Value string
}

// NewTelemetryPolicyFactory creates a factory that can create telemetry policy objects
// which add telemetry information to outgoing HTTP requests.
func NewTelemetryPolicyFactory(o TelemetryOptions) pipeline.Factory {
	return &telemetryPolicyFactory{serviceVersion: serviceLibVersion, options: o}
}

// telemetryPolicyFactory struct
type telemetryPolicyFactory struct {
	serviceVersion string
	options        TelemetryOptions
}

// New creates a telemetryPolicy object.
func (f *telemetryPolicyFactory) New(node pipeline.Node) pipeline.Policy {
	return &telemetryPolicy{node: node, factory: f}
}

// telemetryPolicy ...
type telemetryPolicy struct {
	node    pipeline.Node
	factory *telemetryPolicyFactory
}

func (p *telemetryPolicy) Do(ctx context.Context, request pipeline.Request) (pipeline.Response, error) {
	request = request.Copy() // Don't mutate the incoming request object's headers
	b := &bytes.Buffer{}
	b.WriteString(p.factory.options.Value)
	if b.Len() > 0 {
		b.WriteRune(' ')
	}
	fmt.Fprintf(b, "Azure-Storage/%s %s", p.factory.serviceVersion, platformInfo)
	request.Header.Set("User-Agent", b.String())
	return p.node.Do(ctx, request)
}

// NOTE: the ONLY function that should read OR write to this variable is platformInfo
var platformInfo = initPlatformInfo()

func initPlatformInfo() string {
	// Azure-Storage/version (runtime; os type and version)”
	// Azure-Storage/1.4.0 (NODE-VERSION v4.5.0; Windows_NT 10.0.14393)'
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		operatingSystem = os.Getenv("OS")
	case "linux": // ...
	case "freebsd": // ...
	}
	return fmt.Sprintf("(%s; %s)", runtime.Version(), operatingSystem)
}
