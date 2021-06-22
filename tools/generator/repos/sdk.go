package repos

import (
	"fmt"
	"io/ioutil"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest_ext/changelog_ext"
	"github.com/Azure/azure-sdk-for-go/tools/generator/sdk"
	"github.com/Masterminds/semver"
	"github.com/go-git/go-git/v5/plumbing"
)

type RepositoryWithChangelog interface {
	ReportForCommit(commit string) (changelog_ext.RepoContent, error)
}

type SDKRepository interface {
	WorkTree
	RepositoryWithChangelog
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

type sdkRepository struct {
	WorkTree
}

func GetLatestVersion(wt SDKRepository) (*semver.Version, error) {
	b, err := ioutil.ReadFile(sdk.VersionGoPath(wt.Root()))
	if err != nil {
		return nil, err
	}
	return sdk.GetVersion(string(b))
}

func AddCommit(repo SDKRepository, newVersion string) error {
	changelogFile := sdk.ChangelogPath(repo.Root())
	versionFile := sdk.VersionGoPath(repo.Root())
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

func (s *sdkRepository) ReportForCommit(commit string) (changelog_ext.RepoContent, error) {
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
		defer s.checkoutBack(ref)
	}

	return changelog_ext.GetRepoContent(s.Root())
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
