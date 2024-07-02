// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
)

// Resolver represent a readme path resolver which resolves a link and produces a readme path
type Resolver interface {
	GetReleaseLink() string
	GetRequestLink() string
	Resolve() (ResolveResult, error)
}

// CommitHashLink ...
type CommitHashLink interface {
	GetCommitHash() (string, error)
}

func getCommitRefFromLink(l, prefix string) (string, error) {
	if !strings.HasPrefix(l, prefix) {
		return "", fmt.Errorf("link '%s' does not have prefix '%s'", l, prefix)
	}
	l = strings.TrimPrefix(l, prefix)
	segments := strings.Split(l, "/")
	return segments[0], nil
}

// ResolveResult ...
type ResolveResult interface {
	GetReadme() Readme
	GetCode() Code
}

type Code string

const (
	// CodeSuccess marks the resolve is successful
	CodeSuccess Code = "Success"
	// CodeDataPlane marks the resolved readme belongs to a data plane package
	CodeDataPlane Code = "DataPlaneRequest"
	// CodePRNotMerged marks the resolve succeeds but the requested PR is not merged yet
	CodePRNotMerged Code = "PRNotMerged"

	CodeTypeSpec Code = "TypeSpec"
)

type result struct {
	readme Readme
	code   Code
}

// GetReadme ...
func (r result) GetReadme() Readme {
	return r.readme
}

// GetCode ...
func (r result) GetCode() Code {
	return r.code
}

type linkBase struct {
	ctx         context.Context
	client      *query.Client
	releaseLink string
	requestLink string
}

// GetReleaseLink ...
func (l linkBase) GetReleaseLink() string {
	return l.releaseLink
}

// GetRequestLink ...
func (l linkBase) GetRequestLink() string {
	return l.requestLink
}

func getResult(readme Readme) ResolveResult {
	code := CodeDataPlane
	if readme.IsMgmt() {
		code = CodeSuccess
	} else if readme.IsTsp() {
		code = CodeTypeSpec
	}
	return result{
		readme: readme,
		code:   code,
	}
}
