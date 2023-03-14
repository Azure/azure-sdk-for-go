package shared

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

func NewStorageChallengePolicy(cred azcore.TokenCredential) policy.Policy {
	return runtime.NewBearerTokenPolicy(cred, []string{TokenScope}, nil)
}
