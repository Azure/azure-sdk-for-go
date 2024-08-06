// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// execute `go generate` command and fetch result
func ExecuteGoGenerate(path string) error {
	cmd := exec.Command("go", "generate")
	cmd.Dir = path

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	var stdoutBuffer bytes.Buffer
	if _, err = io.Copy(&stdoutBuffer, stdoutPipe); err != nil {
		return err
	}

	var stderrBuffer bytes.Buffer
	if _, err = io.Copy(&stderrBuffer, stderrPipe); err != nil {
		return err
	}

	err = cmd.Wait()

	fmt.Println(stdoutBuffer.String())
	if stdoutBuffer.Len() > 0 {
		if strings.Contains(stdoutBuffer.String(), "error   |") {
			// find first error message until last
			errMsgs := stdoutBuffer.Bytes()
			index := regexp.MustCompile(`error   \|`).FindIndex(errMsgs)
			if len(index) == 2 {
				return fmt.Errorf("failed to execute `go generate`:\n%s", string(errMsgs[index[0]:]))
			}
		}
	}

	if err != nil || stderrBuffer.Len() > 0 {
		if stderrBuffer.Len() > 0 {
			// filter go downloading log
			// https://github.com/golang/go/blob/1f0c044d60211e435dc58844127544dd3ecb6a41/src/cmd/go/internal/modfetch/fetch.go#L201
			lines := strings.Split(stderrBuffer.String(), "\n")
			newLines := make([]string, 0, len(lines))
			for _, line := range lines {
				l := strings.TrimSpace(line)
				if len(l) == 0 {
					continue
				}
				if !strings.HasPrefix(l, "go: downloading") {
					newLines = append(newLines, line)
				}
			}

			if len(newLines) > 0 {
				newErrMsg := strings.Join(newLines, "\n")
				fmt.Println(newErrMsg)
				return fmt.Errorf("failed to execute `go generate`:\n%s", newErrMsg)
			}

			return nil
		}

		return fmt.Errorf("failed to execute `go generate`:\n%+v", err)
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
	cmd.Stdout = os.Stdout

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, stderr); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil || buf.Len() > 0 {
		if buf.Len() > 0 {
			log.Println(buf.String())

			// filter npm notice log
			lines := strings.Split(buf.String(), "\n")
			newErrInfo := make([]string, 0, len(lines))
			for _, line := range lines {
				if !strings.Contains(line, "npm notice") {
					newErrInfo = append(newErrInfo, line)
				}
			}

			return fmt.Errorf("failed to execute `tsp-client %s`\n%s", strings.Join(args, " "), strings.Join(newErrInfo, "\n"))
		}

		return fmt.Errorf("failed to execute `tsp-client %s`\n%+v", strings.Join(args, " "), err)
	}

	return nil
}

func ExecuteTypeSpecGenerate(ctx *GenerateContext, emitOptions string, tspClientOptions []string) error {
	tspConfigAbs, err := filepath.Abs(ctx.TypeSpecConfig.Path)
	if err != nil {
		return err
	}

	args := []string{
		"init",
		"--tsp-config", tspConfigAbs,
		"--commit", ctx.SpecCommitHash,
		"--repo", ctx.SpecRepoURL[len("https://github.com/"):],
		"--local-spec-repo", filepath.Dir(tspConfigAbs),
		"--emitter-options", emitOptions,
	}

	if len(tspClientOptions) > 0 {
		args = append(args, tspClientOptions...)
	}

	return ExecuteTspClient(ctx.SDKPath, args...)
}
