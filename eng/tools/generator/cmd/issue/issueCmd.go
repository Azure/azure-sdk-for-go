// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package issue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/tools/generator/cmd/issue/request"
	"github.com/Azure/azure-sdk-for-go/tools/generator/config"
	"github.com/Azure/azure-sdk-for-go/tools/generator/config/validate"
	"github.com/Azure/azure-sdk-for-go/tools/generator/flags"
	"github.com/google/go-github/v32/github"
	"github.com/hashicorp/go-multierror"
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
	flagSet.StringSlice("additional-options", []string{"--enum-prefix"}, "Specify the default additional options for the upcoming new version of SDK.")
}

// ParseFlags parses the flags to a Flags struct
func ParseFlags(flagSet *pflag.FlagSet) Flags {
	return Flags{
		IncludeDataPlaneRequests: flags.GetBool(flagSet, "include-data-plane"),
		SkipValidate:             flags.GetBool(flagSet, "skip-validate"),
		ReleaseRequestIDs:        flags.GetIntSlice(flagSet, "request-issues"),
		AdditionalOptions:        flags.GetStringSlice(flagSet, "additional-options"),
	}
}

// Flags ...
type Flags struct {
	IncludeDataPlaneRequests bool
	SkipValidate             bool
	ReleaseRequestIDs        []int
	AdditionalOptions        []string
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
	// write config to stdout
	b, err := json.MarshalIndent(*cfg, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	// we first output the config to stdout, then validate it so that the user could always get a usable config
	// validate the config
	if err := c.validateConfig(*cfg); err != nil {
		return err
	}
	return reqErr
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

func issueHasLabel(issue *github.Issue, label string) bool {
	if issue == nil {
		return false
	}

	for _, l := range issue.Labels {
		if l.GetName() == label {
			return true
		}
	}

	return false
}

func isGoReleaseRequest(issue *github.Issue) bool {
	return issueHasLabel(issue, "Go")
}

func (c *commandContext) parseIssues(issues []*github.Issue) ([]request.Request, error) {
	var requests []request.Request
	var errResult error
	for _, issue := range issues {
		if issue == nil {
			continue
		}
		log.Printf("Parsing issue %s (%s)", issue.GetHTMLURL(), issue.GetTitle())
		req, err := request.ParseIssue(c.ctx, c.client, *issue, request.ParsingOptions{
			IncludeDataPlaneRequests: c.flags.IncludeDataPlaneRequests,
		})
		if err != nil {
			log.Printf("[ERROR] Cannot parse release request %s: %+v", issue.GetHTMLURL(), err)
			errResult = multierror.Append(errResult, err)
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
	track1Requests := config.Track1ReleaseRequests{}
	track2Requests := config.Track2ReleaseRequests{}
	for _, req := range requests {
		switch req.Track {
		case request.Track1:
			track1Requests.Add(req.ReadmePath, req.Tag, config.ReleaseRequestInfo{
				TargetDate:  timePtr(req.TargetDate),
				RequestLink: req.RequestLink,
			})
		case request.Track2:
			track2Requests.Add(req.ReadmePath, config.Track2Request{
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
		Track1Requests:  track1Requests,
		Track2Requests:  track2Requests,
		AdditionalFlags: c.flags.AdditionalOptions,
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
