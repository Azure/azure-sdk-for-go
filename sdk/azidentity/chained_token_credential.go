// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ChainedTokenCredential provides a TokenCredential implementation which chains multiple TokenCredential implementations to be tried in order
// until one of the GetToken methods returns a non-default AccessToken.
type ChainedTokenCredential struct {
	sources []azcore.TokenCredential
}

// NewChainedTokenCredential creates an instance of ChainedTokenCredential with the specified TokenCredential sources.
func NewChainedTokenCredential(sources ...azcore.TokenCredential) (*ChainedTokenCredential, error) {
	if len(sources) == 0 {
		return nil, &CredentialUnavailableError{CredentialType: "Chained Token Credential", Message: "Length of sources cannot be 0"}
	}
	for _, source := range sources {
		if source == nil { // cannot have a nil credential in the chain or else the application will panic when GetToken() is called on nil
			return nil, &CredentialUnavailableError{CredentialType: "Chained Token Credential", Message: "Sources cannot contain a nil TokenCredential"}
		}
	}
	return &ChainedTokenCredential{sources: sources}, nil
}

// GetToken sequentially calls TokenCredential.GetToken on all the specified sources, returning the first non default AccessToken.
func (c *ChainedTokenCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (token *azcore.AccessToken, err error) {
	var errList []*CredentialUnavailableError
	for _, cred := range c.sources { // loop through all of the credentials provided in sources
		token, err = cred.GetToken(ctx, opts) // make a GetToken request for the current credential in the loop
		var credErr *CredentialUnavailableError
		if errors.As(err, &credErr) { // check if we received a CredentialUnavailableError
			errList = append(errList, credErr) // if we did receive a CredentialUnavailableError then we append it to our error slice and continue looping for a good credential
		} else if err != nil { // if we receive some other type of error then we must stop looping and process the error accordingly
			var authenticationFailed *AuthenticationFailedError
			if errors.As(err, &authenticationFailed) { // if the error is an AuthenticationFailedError we return the error related to the invalid credential and append all of the other error messages received prior to this point
				authErr := &AuthenticationFailedError{msg: "Received an AuthenticationFailedError, there is an invalid credential in the chain. " + createChainedErrorMessage(errList), inner: err}
				log := azcore.Log()
				msg := fmt.Sprintf("Azure Identity => ERROR in GetToken() call for %T: %s", cred, authErr.Error())
				log.Write(LogCredential, msg)
				return nil, authErr
			}
			return nil, fmt.Errorf("Received an unexpected error: %w", err) // if we receive some other error type this is unexpected and we simple return the unexpected error
		} else {
			log := azcore.Log()
			msg := fmt.Sprintf("Azure Identity => GetToken() result for %T: SUCCESS", cred)
			log.Write(LogCredential, msg)
			if log.Should(LogCredentialVerbose) {
				vmsg := fmt.Sprintf("Azure Identity => Scopes: [%s]", strings.Join(opts.Scopes, ", "))
				log.Write(LogCredentialVerbose, vmsg)
			}
			return token, nil // if we did not receive an error then we return the token
		}
	}
	return nil, &CredentialUnavailableError{CredentialType: "Chained Token Credential", Message: createChainedErrorMessage(errList)} // if we reach this point it means that all of the credentials in the chain returned CredentialUnavailableErrors
}

// AuthenticationPolicy implements the azcore.Credential interface on ChainedTokenCredential.
func (c *ChainedTokenCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}

// Helper function used to chain the error messages of the CredentialUnavailableError slice
func createChainedErrorMessage(errList []*CredentialUnavailableError) string {
	msg := ""
	for _, err := range errList {
		msg += err.Error()
	}

	return msg
}
