package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v32/github"
	"strconv"
	"strings"
)

var confirmComment = "Hi @%s, Please confirm the SDK package change: %s."

func CreatePullRequest(ctx context.Context, client *github.Client, owner, repo, fork, branch, body string) (*github.PullRequest, *github.Response, error) {
	if branch == "" {
		return nil, nil, errors.New("branch name is nil")
	}

	newPR := &github.NewPullRequest{
		Title:               github.String(ReleaseTitle(branch)),
		Head:                github.String(fork + ":" + branch),
		Base:                github.String("main"),
		Body:                github.String(body),
		MaintainerCanModify: github.Bool(true),
	}

	pr, resp, err := client.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		return nil, resp, fmt.Errorf("create pull request error: %v", err)
	}
	return pr, resp, nil
}

func AddIssueComment(ctx context.Context, client *github.Client, owner, repo, issue, prUrl string) (*github.IssueComment, error) {
	s := strings.Split(issue, "/")
	prNumber, err := strconv.Atoi(s[len(s)-1])
	if err != nil {
		return nil, fmt.Errorf("issue link invalid format: %v", err)
	}

	issueInfo, _, err := client.Issues.Get(ctx, owner, repo, prNumber)
	if err != nil {
		return nil, err
	}

	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf(confirmComment, *issueInfo.User.Login, prUrl)),
	}
	issueComment, _, err := client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return nil, err
	}

	return issueComment, nil
}
