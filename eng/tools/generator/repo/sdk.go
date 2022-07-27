// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/link"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/packages/track1"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type SDKRepository interface {
	WorkTree
	RepositoryWithChangelog
	CreateReleaseBranch(releaseBranchName string) error
	AddReleaseCommit(rpName, namespaceName, specHash, version string) error
}

func OpenSDKRepository(path string) (SDKRepository, error) {
	wt, err := NewWorkTree(path)
	if err != nil {
		return nil, err
	}

	return &sdkRepository{
		WorkTree: wt,
	}, nil
}

func CloneSDKRepository(repoUrl, commitID string) (SDKRepository, error) {
	repoBasePath := filepath.Join(os.TempDir(), "generator_sdk")
	if _, err := os.Stat(repoBasePath); err == nil {
		os.RemoveAll(repoBasePath)
	}
	if err := os.Mkdir(repoBasePath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create tmp folder for generation: %+v", err)
	}

	wt, err := CloneWorkTree(fmt.Sprintf("%s.git", repoUrl), repoBasePath)
	if err != nil {
		return nil, err
	}

	err = wt.Checkout(&CheckoutOptions{
		Hash: plumbing.NewHash(commitID),
	})
	if err != nil {
		return nil, err
	}

	return &sdkRepository{
		WorkTree: wt,
	}, nil
}

type sdkRepository struct {
	WorkTree
}

func (s *sdkRepository) AddReleaseCommit(rpName, namespaceName, specHash, version string) error {
	log.Printf("Add release package and commit")
	if err := s.Add(fmt.Sprintf("sdk/resourcemanager/%s/%s", rpName, namespaceName)); err != nil {
		return fmt.Errorf("failed to add 'profiles': %+v", err)
	}

	message := fmt.Sprintf("[Release] sdk/resourcemanager/%s/%s/%s generation from spec commit: %s", rpName, namespaceName, version, specHash)
	if err := s.Commit(message); err != nil {
		if IsNothingToCommit(err) {
			log.Printf("There is nothing to commit. Message: %s", message)
			return nil
		}
		return fmt.Errorf("failed to commit changes: %+v", err)
	}

	return nil
}

func (s *sdkRepository) CreateReleaseBranch(releaseBranchName string) error {
	log.Printf("Checking out to %s", plumbing.NewBranchReferenceName(releaseBranchName))
	return s.Checkout(&CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(releaseBranchName),
		Create: true,
	})
}

type RepoContent map[string]exports.Content

// Print prints the RepoContent to a Writer as JSON string
func (r *RepoContent) Print(o io.Writer) error {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %v", err)
	}
	_, err = o.Write(b)
	return err
}

func GetRepoContent(sdkRoot string) (RepoContent, error) {
	// we must list over the services directory, otherwise it would walk into the .git directory and panic out
	pkgs, err := track1.List(common.ServicesPath(sdkRoot))
	if err != nil {
		return nil, err
	}

	r, err := getExportsForPackages(pkgs, sdkRoot)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// returns repoContent based on the provided slice of package directories
func getExportsForPackages(pkgs []track1.Package, root string) (RepoContent, error) {
	exps := RepoContent{}
	for _, pkg := range pkgs {
		relativePath, err := filepath.Rel(root, pkg.FullPath())
		if err != nil {
			return nil, err
		}
		relativePath = utils.NormalizePath(relativePath)
		if _, ok := exps[relativePath]; ok {
			return nil, fmt.Errorf("duplicate package: %s", pkg.Path())
		}
		exp, err := exports.Get(pkg.FullPath())
		if err != nil {
			return nil, err
		}
		exps[relativePath] = exp
	}
	return exps, nil
}

type RepositoryWithChangelog interface {
	ReportForCommit(commit string) (RepoContent, error)
}

func GetLatestVersion(wt SDKRepository) (*semver.Version, error) {
	b, err := ioutil.ReadFile(common.VersionGoPath(wt.Root()))
	if err != nil {
		return nil, err
	}
	return GetVersion(string(b))
}

func AddCommit(repo SDKRepository, newVersion string) error {
	changelogFile := common.ChangelogPath(repo.Root())
	versionFile := common.VersionGoPath(repo.Root())
	// add changelog and version
	if err := repo.Add(changelogFile); err != nil {
		return fmt.Errorf("failed to add `%s`: %+v", changelogFile, err)
	}
	if err := repo.Add(versionFile); err != nil {
		return fmt.Errorf("failed to add '%s': %+v", versionFile, err)
	}

	if err := repo.Commit(newVersion); err != nil {
		return err
	}

	return nil
}

func (s *sdkRepository) ReportForCommit(commit string) (RepoContent, error) {
	if commit != "" {
		// store the head ref before checkout
		ref, err := s.Head()
		if err != nil {
			return nil, err
		}
		// check out to the commit
		if err := s.Checkout(&CheckoutOptions{
			Hash: plumbing.NewHash(commit),
		}); err != nil {
			return nil, fmt.Errorf("failed to checkout to commit '%s': %+v", commit, err)
		}
		// defer check out back to initial commit or branch
		//defer s.checkoutBack(ref)
		defer func() {
			_ = s.checkoutBack(ref)
		}()
	}

	return GetRepoContent(s.Root())
}

func (s *sdkRepository) checkoutBack(ref *plumbing.Reference) error {
	opt := CheckoutOptions{}
	if ref.Name().IsBranch() {
		opt.Branch = ref.Name()
	} else {
		opt.Hash = ref.Hash()
	}
	return s.Checkout(&opt)
}

// GetRemoteUserName https://github.com/githubName/azure-sdk-for-go.git
func GetRemoteUserName(remote *git.Remote) string {
	if len(remote.Config().URLs) == 0 {
		return ""
	}
	before, _, found := strings.Cut(remote.Config().URLs[0], link.SDKRepo)
	if !found {
		return ""
	}
	_, after, found := strings.Cut(before, "github.com")
	if !found {
		return ""
	}
	return strings.Trim(after, "/")
}

func GetForkRemote(repo WorkTree) (forkRemote *git.Remote, err error) {
	localRemotes, err := repo.Remotes()
	if err != nil {
		return nil, errors.New("local fork remote not set")
	}
	for _, r := range localRemotes {
		if r.Config().Name == "fork" {
			forkRemote = r
		}
	}
	if forkRemote == nil {
		return nil, fmt.Errorf("under %s not set remote fork", link.SDKRepo)
	}
	return
}

func ReleaseTitle(branchName string) string {
	s := strings.Split(branchName, "-")

	inclines := strings.Split(s[0], "/")
	var t1 string
	if len(inclines) > 0 {
		t1 = inclines[len(inclines)-1]
	} else {
		t1 = s[0]
	}

	t1 = strings.Title(t1)
	title := fmt.Sprintf("[%v] ", t1)
	t := []string{"sdk", "resourcemanager"}
	if len(s) == 5 {
		t = append(t, s[1:len(s)-1]...)
	} else if len(s) == 6 { // beta
		t = append(t, s[1:len(s)-2]...)
		t[len(t)-1] = fmt.Sprintf("%s-%s", t[len(t)-1], s[len(s)-2])
	} else {
		return ""
	}
	t2 := strings.Join(t, "/")
	return title + t2
}
