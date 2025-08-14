// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package request

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/link"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/google/go-github/v62/github"
)

var (
	resultHandlerMap = map[link.Code]resultHandlerFunc{
		link.CodeSuccess:     handleTrack2,
		link.CodeDataPlane:   handleDataPlane,
		link.CodePRNotMerged: handlePRNotMerged,
		link.CodeTypeSpec:    handleTypeSpec,
	}
)

// ParsingOptions ...
type ParsingOptions struct {
	IncludeDataPlaneRequests bool
}

type resultHandlerFunc func(ctx context.Context, client *query.Client, reqIssue ReleaseRequestIssue, result link.ResolveResult) (*Request, error)

// Request represents a parsed SDK release request
type Request struct {
	RequestLink string
	TargetDate  time.Time
	ReadmePath  string
	Tag         string
	Track       Track
}

// Track ...
type Track string

const (
	// Track2 ...
	Track2 Track = "Track2"

	// TypeSpec ...
	TypeSpec Track = "TypeSpec"
)

const (
	linkKeyword        = `(\*\*API spec pull request link\*\*: )|(\*\*Link\*\*: )`
	tagKeyword         = `\*\*Readme Tag\*\*: `
	releaseDateKeyword = `\*\*Target release date\*\*: `
)

type issueError struct {
	issue github.Issue
	err   error
}

// Error ...
func (e *issueError) Error() string {
	return fmt.Sprintf("cannot parse release request from issue %s: %+v", e.issue.GetHTMLURL(), e.err)
}

func initializeHandlers(options ParsingOptions) {
	if options.IncludeDataPlaneRequests {
		resultHandlerMap[link.CodeDataPlane] = handleTrack2
	}
}

// ParseIssue parses the release request issues to release requests
func ParseIssue(ctx context.Context, client *query.Client, issue github.Issue, options ParsingOptions) (*Request, error) {
	initializeHandlers(options)

	reqIssue, err := NewReleaseRequestIssue(issue)
	if err != nil {
		return nil, err
	}
	result, err := ParseReadmeFromLink(ctx, client, *reqIssue)
	if err != nil {
		return nil, err
	}
	handler := resultHandlerMap[result.GetCode()]
	if handler == nil {
		panic(fmt.Sprintf("unhandled code '%s'", result.GetCode()))
	}
	return handler(ctx, client, *reqIssue, result)
}

// ReleaseRequestIssue represents a release request issue
type ReleaseRequestIssue struct {
	IssueLink   string
	TargetLink  string
	Tag         string
	ReleaseDate time.Time
	Labels      []*github.Label
}

// NewReleaseRequestIssue ...
func NewReleaseRequestIssue(issue github.Issue) (*ReleaseRequestIssue, error) {
	body := issue.GetBody()
	contents := getRawContent(strings.Split(body, "\n"), []string{
		linkKeyword, tagKeyword, releaseDateKeyword,
	})

	// get release date
	targetDate := regexp.MustCompile(`\d+\/\d+\/\d+`).FindString(contents[releaseDateKeyword])
	releaseDate, err := time.Parse("1/2/2006", targetDate)
	if err != nil {
		releaseDate = time.Now()
	}
	return &ReleaseRequestIssue{
		IssueLink:   issue.GetHTMLURL(),
		TargetLink:  parseLink(contents[linkKeyword]),
		Tag:         contents[tagKeyword],
		ReleaseDate: releaseDate,
		Labels:      issue.Labels,
	}, nil
}

func getRawContent(lines []string, keywords []string) map[string]string {
	result := make(map[string]string)
	for _, line := range lines {
		for _, keyword := range keywords {
			raw := getContentByPrefix(line, keyword)
			if raw != "" {
				result[keyword] = raw
			}
		}
	}
	return result
}

func getContentByPrefix(line, regexExp string) string {
	regex := regexp.MustCompile(regexExp)
	r := regex.FindStringSubmatch(line)
	if len(r) < 1 {
		return ""
	}
	prefix := r[0]
	if strings.HasPrefix(line, prefix) {
		return strings.TrimSpace(strings.TrimPrefix(line, prefix))
	}
	return ""
}

func parseLink(rawLink string) string {
	regex := regexp.MustCompile(`^\[.+\]\((.+)\)$`)
	r := regex.FindStringSubmatch(rawLink)
	if len(r) < 1 {
		return ""
	}
	return r[1]
}

// ParseReadmeFromLink ...
func ParseReadmeFromLink(ctx context.Context, client *query.Client, reqIssue ReleaseRequestIssue) (link.ResolveResult, error) {
	// check if invalid characters in a url
	regex := regexp.MustCompile(`^[ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789\-._~:/\?#\[\]@!\$&'\(\)\*\+,;=]+$`)
	if !regex.MatchString(reqIssue.TargetLink) {
		return nil, fmt.Errorf("link '%s' contains invalid characters", reqIssue.TargetLink)
	}
	r, err := parseResolver(ctx, client, reqIssue.IssueLink, reqIssue.TargetLink)
	if err != nil {
		return nil, err
	}
	return r.Resolve()
}

func parseResolver(ctx context.Context, client *query.Client, requestLink, releaseLink string) (link.Resolver, error) {
	if !strings.HasPrefix(releaseLink, link.SpecRepoPrefix) {
		return nil, fmt.Errorf("link '%s' is not from '%s'", releaseLink, link.SpecRepoPrefix)
	}
	releaseLink = strings.TrimPrefix(releaseLink, link.SpecRepoPrefix)
	prefix, err := getLinkPrefix(releaseLink)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve link '%s': %+v", releaseLink, err)
	}
	switch prefix {
	case link.PullRequestPrefix:
		return link.NewPullRequestLink(ctx, client, requestLink, releaseLink), nil
	case link.DirectoryPrefix:
		return link.NewDirectoryLink(ctx, client, requestLink, releaseLink), nil
	case link.FilePrefix:
		return link.NewFileLink(ctx, client, requestLink, releaseLink), nil
	case link.CommitPrefix:
		return link.NewCommitLink(ctx, client, requestLink, releaseLink), nil
	default:
		return nil, fmt.Errorf("prefix '%s' of link '%s' not supported yet", prefix, releaseLink)
	}
}

func getLinkPrefix(link string) (string, error) {
	segments := strings.Split(link, "/")
	if len(segments) < 2 {
		return "", fmt.Errorf("cannot determine the prefix of link")
	}
	return segments[0] + "/", nil
}
