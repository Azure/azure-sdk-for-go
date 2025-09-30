// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package issue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/request"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config/validate"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/flags"
	"github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	repoName  = "sdk-release-request"
	repoOwner = "Azure"
)

// Command returns the issue command
func Command() *cobra.Command {
	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Fetch and parse the release request issues to get the configuration of the release",
		Long: `This command fetches the release request from https://github.com/Azure/sdk-release-request/issues?q=is%3Aissue+is%3Aopen+label%3AGo
and produces the configuration to stdout which other command in this tool can consume.

In order to query the issues from GitHub, you need to provide some authentication information to this command.
You can either use populate the personal access token by assigning the flag '-t', or you can use your account and 
password (also otp if needed).

WARNING: This command is still working in progress. The current version of this command cannot handle the request of
a data plane RP at all (an error will be thrown out). Use with caution.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := query.Info{
				UserInfo: query.UserInfo{
					Username: flags.GetString(cmd.Flags(), "username"),
					Password: flags.GetString(cmd.Flags(), "password"),
					Otp:      flags.GetString(cmd.Flags(), "otp"),
				},
				Token: flags.GetString(cmd.Flags(), "token"),
			}
			ctx := context.Background()
			cmdCtx := &commandContext{
				ctx:    ctx,
				client: query.Login(ctx, info),
				flags:  ParseFlags(cmd.Flags()),
			}
			return cmdCtx.execute()
		},
	}

	BindFlags(issueCmd.Flags())

	return issueCmd
}

// BindFlags binds the flags to this command
func BindFlags(flagSet *pflag.FlagSet) {
	flagSet.StringP("username", "u", "", "Specify username of github account")
	flagSet.StringP("password", "p", "", "Specify the password of github account")
	flagSet.String("otp", "", "Specify the two-factor authentication code")
	flagSet.StringP("token", "t", "", "Specify the personal access token")
	flagSet.Bool("include-data-plane", false, "Specify whether we include the requests from data plane RPs")
	flagSet.BoolP("skip-validate", "l", false, "Skip the validate for readme files and tags.")
	flagSet.IntSlice("request-issues", []int{}, "Specify the release request IDs to parse.")
}

// ParseFlags parses the flags to a Flags struct
func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		IncludeDataPlaneRequests: flags.GetBool(flagSet, "include-data-plane"),
		SkipValidate:             flags.GetBool(flagSet, "skip-validate"),
		ReleaseRequestIDs:        flags.GetIntSlice(flagSet, "request-issues"),
	}
}

// Flags ...
type Flags struct {
	IncludeDataPlaneRequests bool
	SkipValidate             bool
	ReleaseRequestIDs        []int
}

type commandContext struct {
	ctx    context.Context
	client *query.Client

	flags Flags
}

func (c *commandContext) execute() error {
	issues, err := c.listIssues()
	if err != nil {
		return err
	}
	requests, reqErr := c.parseIssues(issues)
	if reqErr != nil {
		log.Printf("[ERROR] We are getting errors during parsing the release requests: %+v", reqErr)
	}
	log.Printf("Successfully parsed %d request(s)", len(requests))
	cfg, err := c.buildConfig(requests)
	if err != nil {
		return err
	}
	// validate the config
	if err := c.validateConfig(*cfg); err != nil {
		log.Printf("validate config fail:error(%s)", err.Error())
	}
	// output the config to stdout after filtering out some invalid request,  so that the user could always get a usable config
	// write config to stdout
	b, err := json.MarshalIndent(*cfg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func (c *commandContext) listIssues() ([]*github.Issue, error) {
	if len(c.flags.ReleaseRequestIDs) == 0 {
		return c.listOpenIssues()
	}

	return c.listSpecifiedIssues(c.flags.ReleaseRequestIDs)
}

func (c *commandContext) listOpenIssues() ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		Labels: []string{"Go"},
		ListOptions: github.ListOptions{
			PerPage: 10,
		},
	}
	var issues []*github.Issue
	for {
		r, resp, err := c.client.Issues.ListByRepo(c.ctx, repoOwner, repoName, opt)
		if err != nil {
			return nil, err
		}
		issues = append(issues, r...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return issues, nil
}

func (c *commandContext) listSpecifiedIssues(ids []int) ([]*github.Issue, error) {
	var issues []*github.Issue
	for _, id := range ids {
		issue, _, err := c.client.Issues.Get(c.ctx, repoOwner, repoName, id)
		if err != nil {
			return nil, err
		}

		if !isGoReleaseRequest(issue) {
			return nil, fmt.Errorf("release request '%s' is not a Go SDK release request", issue.GetHTMLURL())
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

func issueHasLabel(issue *github.Issue, label IssueLabel) bool {
	if issue == nil {
		return false
	}

	for _, l := range issue.Labels {
		if IssueLabel(l.GetName()) == label {
			return true
		}
	}

	return false
}

type IssueLabel string

const (
	GoLabel                       IssueLabel = "Go"
	AutoLinkLabel                 IssueLabel = "auto-link"
	PRreadyLabel                  IssueLabel = "PRready"
	InconsistentTagLabel          IssueLabel = "Inconsistent tag"
	SdkReleasedByServiceTeamLabel IssueLabel = "SDK released by service owner"
	HoldOnLabel                   IssueLabel = "HoldOn"
)

func isGoReleaseRequest(issue *github.Issue) bool {
	return issueHasLabel(issue, GoLabel)
}

func isAutoLink(issue *github.Issue) bool {
	return issueHasLabel(issue, AutoLinkLabel)
}

func isPRReady(issue *github.Issue) bool {
	return issueHasLabel(issue, PRreadyLabel)
}

func isInconsistentTag(issue *github.Issue) bool {
	return issueHasLabel(issue, InconsistentTagLabel)
}

func isSdkReleasedByServiceTeam(issue *github.Issue) bool {
	return issueHasLabel(issue, SdkReleasedByServiceTeamLabel)
}

func isHoldOn(issue *github.Issue) bool {
	return issueHasLabel(issue, HoldOnLabel)
}

func (c *commandContext) parseIssues(issues []*github.Issue) ([]request.Request, error) {
	var requests []request.Request
	var errResult error
	for _, issue := range issues {
		if issue == nil {
			continue
		}
		if isHoldOn(issue) {
			continue
		}
		if isPRReady(issue) {
			continue
		}
		if !isAutoLink(issue) {
			continue
		}
		if isInconsistentTag(issue) && !isSdkReleasedByServiceTeam(issue) {
			log.Printf("[ERROR] %s Readme tag is inconsistent with default tag\n", issue.GetHTMLURL())
			errResult = errors.Join(errResult, fmt.Errorf("%s: readme tag is inconsistent with default tag", issue.GetHTMLURL()))
			continue
		}

		log.Printf("Parsing issue %s (%s)", issue.GetHTMLURL(), issue.GetTitle())
		req, err := request.ParseIssue(c.ctx, c.client, *issue, request.ParsingOptions{
			IncludeDataPlaneRequests: c.flags.IncludeDataPlaneRequests,
		})
		if err != nil {
			log.Printf("[ERROR] Cannot parse release request %s: %+v", issue.GetHTMLURL(), err)
			errResult = errors.Join(errResult, err)
			continue
		}
		if req == nil {
			continue
		}
		requests = append(requests, *req)
	}
	return requests, errResult
}

func (c *commandContext) buildConfig(requests []request.Request) (*config.Config, error) {
	track2Requests := config.Track2ReleaseRequests{}
	typespecRequests := config.TypeSpecReleaseRequests{}

	for _, req := range requests {
		switch req.Track {
		case request.Track2:
			track2Requests.Add(req.ReadmePath, config.Track2Request{
				ReleaseRequestInfo: config.ReleaseRequestInfo{
					TargetDate:  timePtr(req.TargetDate),
					RequestLink: req.RequestLink,
				},
				PackageFlag: req.Tag, // TODO -- we need a better place to put this in the request
			})
		case request.TypeSpec:
			typespecRequests.Add(req.ReadmePath, config.Track2Request{
				ReleaseRequestInfo: config.ReleaseRequestInfo{
					TargetDate:  timePtr(req.TargetDate),
					RequestLink: req.RequestLink,
				},
				PackageFlag: req.Tag, // TODO -- we need a better place to put this in the request
			})
		default:
			panic("unhandled track " + req.Track)
		}
	}
	return &config.Config{
		Track2Requests:   track2Requests,
		TypeSpecRequests: typespecRequests,
	}, nil
}

func (c *commandContext) validateConfig(cfg config.Config) error {
	if c.flags.SkipValidate {
		return nil
	}
	log.Printf("Validating the generated config...")
	validator := validate.NewRemoteValidator(c.ctx, c.client)
	return validator.Validate(cfg)
}

func timePtr(t time.Time) *time.Time {
	return &t
}
