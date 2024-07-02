// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
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
func ExecuteExampleGenerate(path, packagePath, flag string) error {
	cmd := exec.Command("pwsh", "../../../../eng/scripts/Invoke-MgmtTestgen.ps1", "-skipBuild", "-cleanGenerated", "-format", "-tidy", "-generateExample", packagePath, flag)
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
	cmd := exec.Command("git", "push", remoteName, refName)
	cmd.Dir = path
	msg, err := cmd.CombinedOutput()
	if err != nil {
		return string(msg), err
	}
	return "", nil
}

func ExecuteCreatePullRequest(path, repoOwner, repoName, prOwner, prBranch, prTitle, prBody, authToken, prLabels string) (string, error) {
	cmd := exec.Command("pwsh", "./eng/common/scripts/Submit-PullRequest.ps1", "-RepoOwner", repoOwner, "-RepoName", repoName, "-BaseBranch", "main", "-PROwner", prOwner, "-PRBranch", prBranch, "-AuthToken", authToken, "-PRTitle", prTitle, "-PRBody", prBody, "-PRLabels", prLabels)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `pwsh Submit-PullRequest` execution: \n%s", string(output))
	if err != nil {
		return "", fmt.Errorf("failed to execute `pwsh Submit-PullRequest` '%s': %+v", string(output), err)
	}

	match := regexp.MustCompile(`html_url[ ]*:.*`).FindString(string(output))
	_, after, _ := strings.Cut(match, ":")
	index := strings.Index(after, "https")
	return after[index:], nil
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

func ExecuteAddIssueLabels(path, repoOwner, repoName, issueNumber, authToken string, labels []string) error {
	var l string
	if len(labels) == 1 {
		l = labels[0]
	} else {
		l = strings.Join(labels, ",")
	}
	cmd := exec.Command("pwsh", "./eng/common/scripts/Add-IssueLabels.ps1", "-RepoOwner", repoOwner, "-RepoName", repoName, "-IssueNumber", issueNumber, "-Labels", l, "-AuthToken", authToken)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `pwsh Add-IssueLabels` execution: \n%s", string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `pwsh Add-IssueLabels` '%s': %+v", string(output), err)
	}
	return nil
}

func ExecuteGo(dir string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute `go %s` '%s': %+v", strings.Join(args, " "), string(combinedOutput), err)
	}

	return nil
}

func ExecuteGoFmt(dir string, args ...string) error {
	cmd := exec.Command("gofmt", args...)
	cmd.Dir = dir
	combinedOutput, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute `gofmt %s` '%s': %+v", strings.Join(args, " "), string(combinedOutput), err)
	}

	return nil
}

// execute tsp-client command
func ExecuteTspClient(path string, args ...string) error {
	cmd := exec.Command("tsp-client", args...)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	log.Printf("Result of `tsp-client %s` execution: \n%s", strings.Join(args, " "), string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `tsp-client %s` '%s': %+v", strings.Join(args, " "), string(output), err)
	}
	if strings.Contains(string(output), "error:") {
		return fmt.Errorf("failed to execute `tsp-client %s` '%s'", strings.Join(args, " "), string(output))
	}
	return nil
}

func ExecuteTypeSpecGenerate(path, tspConfigPath, specCommit, specRepo, tspDir, emitOptions string) error {

	return ExecuteTspClient(
		path,
		"init",
		"--tsp-config", tspConfigPath,
		"--commit", specCommit,
		"--repo", specRepo,
		"--local-spec-repo", tspDir,
		"--emitter-options", emitOptions,
	)
}
