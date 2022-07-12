// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

const (
	sdk_tag_fetch_url = "https://api.github.com/repos/Azure/azure-sdk-for-go/git/refs/tags"
	sdk_remote_url    = "https://github.com/Azure/azure-sdk-for-go.git"
)

func GetAllVersionTags(rpName, namespaceName string) ([]string, error) {
	log.Printf("Fetching all release tags from GitHub for RP: '%s' Package: '%s' ...", rpName, namespaceName)
	client := http.Client{}
	res, err := client.Get(sdk_tag_fetch_url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	result := []map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, tag := range result {
		tagName := tag["ref"].(string)
		if strings.Contains(tagName, "sdk/resourcemanager/"+rpName+"/"+namespaceName) {
			tags = append(tags, tag["ref"].(string))
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(tags)))

	return tags, nil
}

func ContainsPreviewAPIVersion(packagePath string) (bool, error) {
	log.Printf("Judge whether contains preview API version from '%s' ...", packagePath)

	files, err := ioutil.ReadDir(packagePath)
	if err != nil {
		return false, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".go") {
			b, err := ioutil.ReadFile(path.Join(packagePath, file.Name()))
			if err != nil {
				return false, err
			}

			lines := strings.Split(string(b), "\n")
			for _, line := range lines {
				if strings.Contains(line, "\"api-version\"") {
					parts := strings.Split(line, "\"")
					if len(parts) == 5 && strings.Contains(parts[3], "preview") {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

func GetPreviousVersionTag(isCurrentPreview bool, allReleases []string) string {
	if isCurrentPreview {
		// for preview api, always compare with latest release
		return allReleases[0]
	} else {
		// for stable api, always compare with previous stable, if no stable, then latest release
		for _, release := range allReleases {
			if !strings.Contains(release, "beta") {
				return release
			}
		}
		return allReleases[0]
	}
}

func GetExportsFromTag(sdkRepo repo.SDKRepository, packagePath, tag string) (*exports.Content, error) {
	log.Printf("Get exports from specific tag '%s' ...", tag)

	// get current head branch name
	currentRef, err := sdkRepo.Head()
	if err != nil {
		return nil, err
	}

	// add package change
	err = sdkRepo.Add(packagePath)
	if err != nil {
		return nil, err
	}

	// stash current change
	err = sdkRepo.Stash()
	if err != nil {
		return nil, err
	}

	// create remote with center sdk repo
	remoteName := "release_remote"
	_, err = sdkRepo.CreateRemote(&config.RemoteConfig{Name: remoteName, URLs: []string{sdk_remote_url}})
	if err != nil {
		if err != git.ErrRemoteExists {
			return nil, err
		}
	}

	// fetch tag from remote
	err = sdkRepo.Fetch(&git.FetchOptions{RemoteName: remoteName, RefSpecs: []config.RefSpec{config.RefSpec(tag + ":" + tag)}})
	if err != nil {
		if err.Error() != "already up-to-date" {
			return nil, err
		}
	}

	// checkout to the specific tag
	err = sdkRepo.CheckoutTag(strings.TrimPrefix(tag, "ref/tags/"))
	if err != nil {
		return nil, err
	}

	// get exports
	result, err := exports.Get(packagePath)
	if err != nil {
		return nil, err
	}

	// checkout back to head branch
	err = sdkRepo.Checkout(&repo.CheckoutOptions{
		Branch: plumbing.ReferenceName(currentRef.Name()),
	})
	if err != nil {
		return nil, err
	}

	// remove remote
	err = sdkRepo.DeleteRemote(remoteName)
	if err != nil {
		return nil, err
	}

	// restore current change
	err = sdkRepo.StashPop()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func FilterChangelog(changelog *model.Changelog) {
	if changelog.Modified != nil {
		if changelog.Modified.AdditiveChanges != nil {
			removeMarshalUnmarshalFunc(changelog.Modified.AdditiveChanges.Funcs)
		}
		if changelog.Modified.BreakingChanges != nil && changelog.Modified.BreakingChanges.Removed != nil {
			removeMarshalUnmarshalFunc(changelog.Modified.BreakingChanges.Removed.Funcs)
		}
	}
}

func removeMarshalUnmarshalFunc(funcs map[string]exports.Func) {
	for k := range funcs {
		if strings.HasSuffix(k, ".MarshalJSON") || strings.HasSuffix(k, ".UnmarshalJSON") {
			delete(funcs, k)
		}
	}
}
