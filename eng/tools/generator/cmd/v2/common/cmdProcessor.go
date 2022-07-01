// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// execute `go generate` command and fetch result
func ExecuteGoGenerate(path string) error {
	cmd := exec.Command("go", "generate")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go generate` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go generate` '%s': %+v", string(output), err)
	}
	return nil
}

// execute `pwsh Invoke-MgmtTestgen` command and fetch result
func ExecuteExampleGenerate(path, packagePath string) error {
	cmd := exec.Command("pwsh", "../../../../eng/scripts/Invoke-MgmtTestGen.ps1", "-skipBuild", "-cleanGenerated", "-format", "-tidy", "-generateExample", packagePath)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `pwsh Invoke-MgmtTestgen` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `pwsh Invoke-MgmtTestgen` '%s': %+v", string(output), err)
	}
	return nil
}

// execute `goimports` command and fetch result
func ExecuteGoimports(path string) error {
	cmd := exec.Command("go", "get", "golang.org/x/tools/cmd/goimports")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `go get golang.org/x/tools/cmd/goimports` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `go get golang.org/x/tools/cmd/goimports` '%s': %+v", string(output), err)
	}
	cmd = exec.Command("goimports", "-w", ".")
	cmd.Dir = path
	output, err = cmd.CombinedOutput()
	log.Printf("Result of `goimports` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `goimports` '%s': %+v", string(output), err)
	}
	return nil
}

func ExecuteGitPush(path, remoteName, branchName string) (string, error) {
	refName := fmt.Sprintf(branchName + ":" + branchName)
	push := exec.Command("git", "push", remoteName, refName)
	push.Dir = path
	msg, err := push.CombinedOutput()
	if err != nil {
		return string(msg), err
	}
	return "", nil
}

func ExecuteCreatePullRequest(path, repoOwner, repoName, prOwner, prBranch, prTitle, prBody, authToken string) (string, error) {
	cmd := exec.Command("pwsh", "./eng/common/scripts/Submit-PullRequest.ps1", "-RepoOwner", repoOwner, "-RepoName", repoName, "-BaseBranch", "main", "-PROwner", prOwner, "-PRBranch", prBranch, "-AuthToken", authToken, "-PRTitle", prTitle, "-PRBody", prBody)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `pwsh Submit-PullRequest` execution: \n%s", string(output))
	if err != nil {
		return "", fmt.Errorf("failed to execute `pwsh Submit-PullRequest` '%s': %+v", string(output), err)
	}

	s1 := strings.Split(string(output), "html_url=")
	s2 := strings.Split(s1[len(s1)-1], ";")

	return s2[0], nil
}

func ExecuteAddIssueComment(path, repoOwner, repoName, issueNumber, comment, authToken string) error {
	cmd := exec.Command("pwsh", "./eng/common/scripts/Add-IssueComment.ps1", "-RepoOwner", repoOwner, "-RepoName", repoName, "-IssueNumber", issueNumber, "-Comment", comment, "-AuthToken", authToken)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `pwsh Add-IssueComment` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `pwsh Add-IssueComment` '%s': %+v", string(output), err)
	}
	return nil
}
