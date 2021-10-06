// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5/plumbing"
)

type SpecRepository interface {
	WorkTree
	LastHead() *plumbing.Reference
}

func OpenSpecRepository(path string) (SpecRepository, error) {
	spec, err := NewWorkTree(path)
	if err != nil {
		return nil, err
	}

	lastRef, err := spec.Head()
	if err != nil {
		return nil, err
	}

	return &specRepository{
		WorkTree: spec,
		lastRef:  lastRef,
	}, nil
}

func CloneSpecRepository(repoUrl, commitID string) (SpecRepository, error) {
	repoBasePath := filepath.Join(os.TempDir(), "generator_spec")
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

	lastRef, err := wt.Head()
	if err != nil {
		return nil, err
	}

	return &specRepository{
		WorkTree: wt,
		lastRef:  lastRef,
	}, nil
}

type specRepository struct {
	WorkTree

	lastRef *plumbing.Reference
}

func (s *specRepository) LastHead() *plumbing.Reference {
	return s.lastRef
}
