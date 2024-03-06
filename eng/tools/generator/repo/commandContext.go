// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package repo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5/plumbing"
)

type CommandContext interface {
	SDK() SDKRepository
	Spec() SpecRepository
	CreateReleaseBranch(version string) (string, error)
	CheckExternalChanges()
}

type commandContext struct {
	sdk  SDKRepository
	spec SpecRepository

	checkExternalChanges func(ref, newRef *plumbing.Reference, err error)
}

func (c *commandContext) CheckExternalChanges() {
	ref := c.spec.LastHead()
	newRef, err := c.spec.Head()
	c.checkExternalChanges(ref, newRef, err)
}

func checkExternalChangesPanic(ref, newRef *plumbing.Reference, err error) {
	if err != nil {
		log.Fatalf("Cannot get head ref of azure-rest-api-specs when checking external changes: %+v", err)
		return
	}
	if ref.Hash().String() != newRef.Hash().String() {
		log.Fatalf("External changes detected in azure-rest-api-specs. Command executed on %s, but now on %s", ref.Hash().String(), newRef.Hash().String())
	}
}

func checkExternalChangesWarning(ref, newRef *plumbing.Reference, err error) {
	if err != nil {
		log.Printf("[WARNING] Cannot get head ref of azure-rest-api-specs when checking external changes: %+v", err)
		return
	}
	if ref.Hash().String() != newRef.Hash().String() {
		log.Printf("[WARNING] External changes detected in azure-rest-api-specs. Command executed on %s, but now on %s", ref.Hash().String(), newRef.Hash().String())
	}
}

func (c *commandContext) CreateReleaseBranch(version string) (string, error) {
	// append a time in long to avoid collision of branch names
	releaseBranchName := fmt.Sprintf(releaseBranchNamePattern, version, time.Now().Unix())

	log.Printf("Checking out to %s", plumbing.NewBranchReferenceName(releaseBranchName))
	if err := c.SDK().Checkout(&CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(releaseBranchName),
		Create: true,
	}); err != nil {
		return "", err
	}

	return releaseBranchName, nil

}

func (c *commandContext) SDK() SDKRepository {
	return c.sdk
}

func (c *commandContext) Spec() SpecRepository {
	return c.spec
}

func NewCommandContext(sdkPath, specPath string, panicWhenDetectExternalChanges bool) (CommandContext, error) {
	sdkRepo, err := OpenSDKRepository(sdkPath)
	if err != nil {
		return nil, err
	}
	specRepo, err := OpenSpecRepository(specPath)
	if err != nil {
		return nil, err
	}
	ctx := commandContext{
		sdk:  sdkRepo,
		spec: specRepo,
	}
	if panicWhenDetectExternalChanges {
		ctx.checkExternalChanges = checkExternalChangesPanic
	} else {
		ctx.checkExternalChanges = checkExternalChangesWarning
	}
	return &ctx, nil
}

func TempDir() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("generator-%v", time.Now().Unix()))
}

const (
	releaseBranchNamePattern = "release-%s-%v"
)
