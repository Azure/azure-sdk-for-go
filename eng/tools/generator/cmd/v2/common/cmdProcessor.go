// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package common

import (
	"bytes"
	"fmt"
	"io"
	"log"
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

	cmdWaitErr := cmd.Wait()
	if cmdWaitErr == nil {
		return nil
	}

	fmt.Println(stdoutBuffer.String())
	fmt.Println(stderrBuffer.String())

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

	if cmdWaitErr != nil || stderrBuffer.Len() > 0 {
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
				return fmt.Errorf("failed to execute `go generate`:\n%s", strings.Join(newLines, "\n"))
			}

			return nil
		}

		return fmt.Errorf("failed to execute `go generate`:\n%+v", err)
	}

	return nil
}

// execute `pwsh Invoke-MgmtTestgen` command and fetch result
func ExecuteExampleGenerate(path, packagePath string, flags []string) error {
	cmd := exec.Command("pwsh", append([]string{"../../../../eng/scripts/Invoke-MgmtTestgen.ps1", "-skipBuild", "-cleanGenerated", "-format", "-tidy", "-generateExample", packagePath}, flags...)...)
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
	// Use pinned tsp-client from eng/common/tsp-client instead of global npx
	tspClientDir := filepath.Join(path, "eng", "common", "tsp-client")
	args = append([]string{"--prefix", tspClientDir, "exec", "--no", "--", "tsp-client"}, args...)
	cmd := exec.Command("npm", args...)
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

	cmdWaitErr := cmd.Wait()
	fmt.Println(stdoutBuffer.String())
	fmt.Println(stderrBuffer.String())

	if stdoutBuffer.Len() > 0 {
		for _, line := range strings.Split(stdoutBuffer.String(), "\n") {
			if len(strings.TrimSpace(line)) == 0 {
				continue
			}
			if strings.Contains(line, "generation complete") {
				return nil
			}
		}
	}
	if cmdWaitErr != nil || stderrBuffer.Len() > 0 {
		if stderrBuffer.Len() > 0 {
			// filter npm notice & warning log
			newErrMsgs := make([]string, 0)
			for _, line := range strings.Split(stderrBuffer.String(), "\n") {
				if len(strings.TrimSpace(line)) == 0 {
					continue
				}
				if strings.Contains(line, "npm notice") {
					continue
				}
				if strings.Contains(line, "npm warn") {
					continue
				}
				newErrMsgs = append(newErrMsgs, line)
			}

			// filter diagnostic errors
			if len(newErrMsgs) >= 1 &&
				strings.HasPrefix(newErrMsgs[0], "Diagnostics were reported during compilation.") {
				newErrMsgs = newErrMsgs[1:]

				errDiags := getErrorDiagnostics(strings.Split(stdoutBuffer.String(), "\n"))
				temp := make([]string, 0)
				for _, line := range errDiags {
					line := strings.TrimSpace(line)
					if line == "" ||
						strings.Contains(line, "Cleaning up temp directory") ||
						strings.Contains(line, "Skipping cleanup of temp directory:") {
						continue
					}
					temp = append(temp, line)
				}

				if len(temp) > 0 {
					newErrMsgs = append(newErrMsgs, temp...)
				}
			}

			if len(newErrMsgs) > 0 {
				return fmt.Errorf("failed to execute `tsp-client %s`\n%s", strings.Join(args, " "), strings.Join(newErrMsgs, "\n"))
			}

			return nil
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
		"--update-if-exists",
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

type diagnostic struct {
	// error | warning
	kind       string
	start, end int
}

// get all warning and error diagnostics
func diagnostics(lines []string) []diagnostic {
	var kind string
	start := -1
	diagnostics := make([]diagnostic, 0)

	// get all warning and error diagnostics
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if strings.Contains(line, "warning ") {
			if start != -1 {
				diagnostics = append(diagnostics, diagnostic{kind: kind, start: start, end: i - 1})
			}
			start = i
			kind = "warning"
		} else if strings.Contains(line, "error ") {
			if start != -1 {
				diagnostics = append(diagnostics, diagnostic{kind: kind, start: start, end: i - 1})
			}
			start = i
			kind = "error"
		}

		if i == len(lines)-1 && start != -1 {
			diagnostics = append(diagnostics, diagnostic{kind: kind, start: start, end: i})
		}
	}

	return diagnostics
}

func getErrorDiagnostics(lines []string) []string {
	diags := diagnostics(lines)
	if len(diags) == 0 {
		return nil
	}

	result := make([]string, 0)
	for _, diag := range diags {
		if diag.kind == "error" {
			result = append(result, lines[diag.start:diag.end+1]...)
		}
	}

	return result
}
