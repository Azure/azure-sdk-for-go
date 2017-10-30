package azblob

import (
	"context"

	"github.com/Azure/azure-pipeline-go/pipeline"
)

// Credential represent any credential type; it is used to create a credential policy Factory.
type Credential interface {
	pipeline.Factory
	credentialMarker()
}

// NewAnonymousCredential creates an anonymous credential for use with HTTP(S)
// requests that read blobs from public containers or for use with Shared Access
// Signatures (SAS).
func NewAnonymousCredential() Credential {
	return &anonymousCredentialPolicyFactory{}
}

// anonymousCredentialPolicyFactory is the credential's policy factory.
type anonymousCredentialPolicyFactory struct {
}

// New creates a credential policy object.
func (f *anonymousCredentialPolicyFactory) New(node pipeline.Node) pipeline.Policy {
	return &anonymousCredentialPolicy{node: node}
}

// credentialMarker is a package-internal method that exists just to satisfy the Credential interface.
func (*anonymousCredentialPolicyFactory) credentialMarker() {}

// anonymousCredentialPolicy is the credential's policy object.
type anonymousCredentialPolicy struct {
	node pipeline.Node
}

// Do implements the credential's policy interface.
func (p anonymousCredentialPolicy) Do(ctx context.Context, request pipeline.Request) (pipeline.Response, error) {
	// For anonymous credentials, this is effectively a no-op
	return p.node.Do(ctx, request)
}
