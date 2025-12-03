// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

type WorkTree interface {
	Root() string
	Add(path string) error
	Commit(message string) error
	Checkout(opt *CheckoutOptions) error
	CheckoutTag(tag string) error
	CreateBranch(branch *Branch) error
	DeleteBranch(name string) error
	CherryPick(commit string) error
	Stash() error
	StashPop() error
	Head() (*plumbing.Reference, error)
	Tags() (storer.ReferenceIter, error)
	Remotes() ([]*git.Remote, error)
	DeleteRemote(name string) error
	CreateRemote(c *config.RemoteConfig) (*git.Remote, error)
	Fetch(o *git.FetchOptions) error
}

type CheckoutOptions git.CheckoutOptions
type Branch config.Branch

type repository struct {
	*git.Repository

	wt   *git.Worktree
	root string
}

func NewWorkTree(path string) (WorkTree, error) {
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot open '%s': %+v", path, err)
	}
	wt, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("cannot get the work tree of '%s': %+v", path, err)
	}
	return &repository{
		Repository: r,
		wt:         wt,
		root:       wt.Filesystem.Root(),
	}, nil
}

func CloneWorkTree(repoURL, workingPath string) (WorkTree, error) {
	r, err := git.PlainClone(workingPath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot clone '%s' to '%s': %+v", repoURL, workingPath, err)
	}
	wt, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("cannot get the work tree of '%s': %+v", workingPath, err)
	}

	return &repository{
		Repository: r,
		wt:         wt,
		root:       wt.Filesystem.Root(),
	}, nil
}

func (r *repository) Root() string {
	return r.root
}

// TODO -- go-git has some performance issue during Add, therefore we use the git command as a workaround
func (r *repository) Add(path string) error {
	cmd := exec.Command("git", "add", path)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// TODO -- go-git has some performance and permission issue during Commit, therefore we use the git command as a workaround
func (r *repository) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		m := string(output)
		if strings.Contains(m, "nothing added to commit") {
			return &NothingToCommit{
				message: m,
			}
		}
		return fmt.Errorf("%s", m)
	}
	return nil
}

// TODO -- go-git has some issue on the Checkout command, it will keep the CRLF changes after switching branches in stage
func (r *repository) Checkout(opt *CheckoutOptions) error {
	if len(opt.Branch) > 0 {
		return r.checkoutBranch(opt.Branch.Short(), opt.Create)
	}
	if !opt.Hash.IsZero() {
		return r.checkoutHash(opt.Hash.String())
	}
	return fmt.Errorf("must set one of hash or branch")
}

func (r *repository) checkoutBranch(branch string, create bool) error {
	var cmd *exec.Cmd
	if create {
		cmd = exec.Command("git", "checkout", "-b", branch)
	} else {
		cmd = exec.Command("git", "checkout", branch)
	}
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func (r *repository) checkoutHash(hash string) error {
	cmd := exec.Command("git", "checkout", hash)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func (r *repository) CheckoutTag(tag string) error {
	cmd := exec.Command("git", "checkout", tag)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func (r *repository) CreateBranch(branch *Branch) error {
	cmd := exec.Command("git", "branch", branch.Name)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// TODO -- we cannot delete the branch that is not created using go-git, therefore we use the git command as a workaround
func (r *repository) DeleteBranch(name string) error {
	cmd := exec.Command("git", "branch", "-D", name)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// TODO -- go-git now does not support cherry-pick (or I did not find this?), therefore we use the git command as a workaround
func (r *repository) CherryPick(commit string) error {
	cmd := exec.Command("git", "cherry-pick", commit)
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// TODO -- go-git now does not support stash (or I did not find this?), therefore we use the git command as a workaround
func (r *repository) Stash() error {
	cmd := exec.Command("git", "stash")
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

// TODO -- go-git now does not support stash (or I did not find this?), therefore we use the git command as a workaround
func (r *repository) StashPop() error {
	cmd := exec.Command("git", "stash", "pop")
	cmd.Dir = r.root
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

type NothingToCommit struct {
	message string
}

func (n *NothingToCommit) Error() string {
	return n.message
}

func IsNothingToCommit(err error) bool {
	_, ok := err.(*NothingToCommit)
	return ok
}
