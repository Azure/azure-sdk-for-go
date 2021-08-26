// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
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

type specRepository struct {
	WorkTree

	lastRef *plumbing.Reference
}

func (s *specRepository) LastHead() *plumbing.Reference {
	return s.lastRef
}
