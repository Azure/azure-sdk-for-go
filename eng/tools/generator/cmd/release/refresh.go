// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package release

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/refresh"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/go-git/go-git/v5/plumbing"
)

func (c *commandContext) refresh(r *config.RefreshInfo) (*plumbing.Reference, error) {
	// we will not refresh any package if RefreshAll flag is not set, and there are no explicit packages to refresh in the configuration
	if !c.flags.RefreshAll && (r == nil || len(r.Packages) == 0) {
		return c.sdkRef, nil
	}

	// create a temporary branch to hold the generation result
	log.Printf("Creating temporary branch...")
	tempBranchName, err := c.CreateReleaseBranch("refresh")
	if err != nil {
		return nil, err
	}
	log.Printf("Temporary branch '%s' created", tempBranchName)

	// defer checkout back in azure-rest-api-specs
	defer func() {
		if err := c.Spec().Checkout(&repo.CheckoutOptions{
			Branch: c.specRef.Name(),
			Force:  true,
		}); err != nil {
			log.Printf("Error checking out azure-rest-api-specs to %v", c.specRef)
		}
	}()

	// execute the refresh command
	refreshContext := refresh.CommandContext{
		CommandContext: c.CommandContext,
		Flags:          refresh.Flags{},
		RepoContent:    c.repoContent,
	}
	refreshRef, err := refreshContext.Refresh(r)
	if err != nil {
		return nil, fmt.Errorf("failed to execute refresh: %+v", err)
	}

	return refreshRef, nil
}
