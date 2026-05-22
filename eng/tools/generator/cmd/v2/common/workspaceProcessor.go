// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
)

var (
	commitIDRegex = regexp.MustCompile("^[0-9a-f]{40}$")
)

func GetSDKRepo(sdkRepoParam, sdkRepoURL string) (repo.SDKRepository, error) {
	var err error
	var sdkRepo repo.SDKRepository
	// create sdk git repo ref
	if commitIDRegex.Match([]byte(sdkRepoParam)) {
		if sdkRepo, err = repo.CloneSDKRepository(sdkRepoURL, sdkRepoParam); err != nil {
			return nil, fmt.Errorf("failed to get sdk repo: %+v", err)
		}
	} else {
		path, err := filepath.Abs(sdkRepoParam)
		if err != nil {
			return nil, fmt.Errorf("failed to get the directory of azure-sdk-for-go: %v", err)
		}

		if sdkRepo, err = repo.OpenSDKRepository(path); err != nil {
			return nil, fmt.Errorf("failed to get sdk repo: %+v", err)
		}
	}
	return sdkRepo, nil
}

func GetSpecCommit(specRepoParam string) (string, error) {
	specCommitHash := ""
	// create spec git repo ref
	if commitIDRegex.Match([]byte(specRepoParam)) {
		specCommitHash = specRepoParam
	} else {
		path, err := filepath.Abs(specRepoParam)
		if err != nil {
			return "", fmt.Errorf("failed to get the directory of azure-rest-api-specs: %v", err)
		}
		specRepo, err := repo.OpenSpecRepository(path)
		if err != nil {
			return "", fmt.Errorf("failed to get spec repo: %+v", err)
		}
		specHeader, err := specRepo.Head()
		if err != nil {
			return "", fmt.Errorf("failed to get HEAD ref of azure-rest-api-specs: %+v", err)
		}
		specCommitHash = specHeader.Hash().String()
	}

	return specCommitHash, nil
}
