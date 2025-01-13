// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package request

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/link"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
)

func handleTrack2(_ context.Context, _ *query.Client, reqIssue ReleaseRequestIssue, result link.ResolveResult) (*Request, error) {
	return &Request{
		RequestLink: reqIssue.IssueLink,
		TargetDate:  reqIssue.ReleaseDate,
		ReadmePath:  string(result.GetReadme()),
		Tag:         reqIssue.Tag,
		Track:       Track2,
	}, nil
}

func handleDataPlane(_ context.Context, _ *query.Client, reqIssue ReleaseRequestIssue, result link.ResolveResult) (*Request, error) {
	log.Printf("[WARNING] Release request %s is requesting a release from a data-plane readme file `%s`, treat this as a track 2 request by default", reqIssue.IssueLink, result.GetReadme())
	return &Request{
		RequestLink: reqIssue.IssueLink,
		TargetDate:  reqIssue.ReleaseDate,
		ReadmePath:  string(result.GetReadme()),
		Tag:         reqIssue.Tag,
		Track:       Track2,
	}, nil
}

func handlePRNotMerged(_ context.Context, _ *query.Client, reqIssue ReleaseRequestIssue, result link.ResolveResult) (*Request, error) {
	log.Printf("[WARNING] Release request %s is requesting a release from a non-merged PR `%s`, discard this request", reqIssue.IssueLink, reqIssue.TargetLink)
	// TODO -- add comment and close this issue
	return nil, nil
}

func handleTypeSpec(_ context.Context, _ *query.Client, reqIssue ReleaseRequestIssue, result link.ResolveResult) (*Request, error) {
	return &Request{
		RequestLink: reqIssue.IssueLink,
		TargetDate:  reqIssue.ReleaseDate,
		ReadmePath:  string(result.GetReadme()),
		Tag:         reqIssue.Tag,
		Track:       TypeSpec,
	}, nil
}
