// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/tools/internal/ioext"
	"github.com/Azure/azure-sdk-for-go/tools/internal/repo"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
)

func printf(format string, a ...interface{}) {
	if !quietFlag {
		fmt.Printf(format, a...)
	}
}

func println(a ...interface{}) {
	if !quietFlag {
		fmt.Println(a...)
	}
}

func dprintf(format string, a ...interface{}) {
	if debugFlag {
		printf(format, a...)
	}
}

func dprintln(a ...interface{}) {
	if debugFlag {
		println(a...)
	}
}

func vprintf(format string, a ...interface{}) {
	if verboseFlag {
		printf(format, a...)
	}
}

func vprintln(a ...interface{}) {
	if verboseFlag {
		println(a...)
	}
}

func processArgsAndClone(args []string) (cln repo.WorkingTree, err error) {
	if onlyAdditiveChangesFlag && onlyBreakingChangesFlag {
		err = errors.New("flags 'additions' and 'breakingchanges' are mutually exclusive")
		return
	}

	// there should be at minimum two args, a directory and a
	// sequence of commits, i.e. "d:\foo 1,2,3".  else a directory
	// and two commits, i.e. "d:\foo 1 2" or "d:\foo 1 2,3"
	if len(args) < 2 {
		err = errors.New("not enough args were supplied")
		return
	}

	// here args[1] should be a comma-delimited list of commits
	if len(args) == 2 && strings.Index(args[1], ",") < 0 {
		err = errors.New("expected a comma-delimited list of commits")
		return
	}

	dir := args[0]
	dir, err = filepath.Abs(dir)
	if err != nil {
		err = fmt.Errorf("failed to convert path '%s' to absolute path: %v", dir, err)
		return
	}

	src, err := repo.Get(dir)
	if err != nil {
		err = fmt.Errorf("failed to get repository: %v", err)
		return
	}

	tempRepoDir := filepath.Join(os.TempDir(), fmt.Sprintf("apidiff-%v", time.Now().Unix()))
	if copyRepoFlag {
		vprintf("copying '%s' into '%s'...\n", src.Root(), tempRepoDir)
		err = ioext.CopyDir(src.Root(), tempRepoDir)
		if err != nil {
			err = fmt.Errorf("failed to copy repo: %v", err)
			return
		}
		cln, err = repo.Get(tempRepoDir)
		if err != nil {
			err = fmt.Errorf("failed to get copied repo: %v", err)
			return
		}
	} else {
		vprintf("cloning '%s' into '%s'...\n", src.Root(), tempRepoDir)
		cln, err = src.Clone(tempRepoDir)
		if err != nil {
			err = fmt.Errorf("failed to clone repository: %v", err)
			return
		}
	}

	// fix up pkgDir to the clone
	args[0] = strings.Replace(dir, src.Root(), cln.Root(), 1)

	return
}

type reportGenFunc func(dir string, cln repo.WorkingTree, baseCommit, targetCommit string) error

func generateReports(args []string, cln repo.WorkingTree, fn reportGenFunc) error {
	defer func() {
		// delete clone
		vprintln("cleaning up clone")
		err := os.RemoveAll(cln.Root())
		if err != nil {
			vprintf("failed to delete temp repo: %v\n", err)
		}
	}()

	var commits []string

	// if the commits are specified as 1 2,3,4 then it means that commit 1 is
	// always the base commit and to compare it against each target commit in
	// the sequence.  however if it's specifed as 1,2,3,4 then the base commit
	// is relative to the iteration, i.e. compare 1-2, 2-3, 3-4.
	fixedBase := true
	if len(args) == 3 {
		commits = make([]string, 2, 2)
		commits[0] = args[1]
		commits[1] = args[2]
	} else {
		commits = strings.Split(args[1], ",")
		fixedBase = false
	}

	for i := 0; i+1 < len(commits); i++ {
		baseCommit := commits[i]
		if fixedBase {
			baseCommit = commits[0]
		}
		targetCommit := commits[i+1]

		err := fn(args[0], cln, baseCommit, targetCommit)
		if err != nil {
			return err
		}
	}
	return nil
}

// compares report status with the desired report options (breaking/additions)
// to determine if the program should terminate with a non-zero exit code.
func evalReportStatus(r report.Status) {
	if onlyBreakingChangesFlag && r.HasBreakingChanges() {
		os.Exit(1)
	}
	if onlyAdditiveChangesFlag && !r.HasAdditiveChanges() {
		os.Exit(1)
	}
}

// PrintReport prints the report to stdout
func PrintReport(r report.Status) error {
	if r.IsEmpty() {
		println("no changes were found")
		return nil
	}

	if !suppressReport {
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal report: %v", err)
		}
		println(string(b))
	}
	return nil
}
