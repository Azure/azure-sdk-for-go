package autorest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Task describes a generation task
type Task struct {
	// AbsReadmeMd absolute path of the readme.md file to generate
	AbsReadmeMd string
}

// Execute executes the autorest task, and then invoke the after scripts
// the error returned will be TaskError
func (t *Task) Execute(options Options) error {
	if err := t.executeAutorest(options.AutorestArguments); err != nil {
		return err
	}

	if err := t.executeAfterScript(options.AfterScripts); err != nil {
		return err
	}

	return nil
}

func (t *Task) executeAutorest(options []string) error {
	arguments := append(options, t.AbsReadmeMd)
	c := exec.Command("autorest", arguments...)
	log.Printf("Executing autorest with parameters: %+v", arguments)

	stdout, _ := c.StdoutPipe()
	stderr, _ := c.StderrPipe()
	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(bufio.ScanLines)
	outScanner := bufio.NewScanner(stdout)
	outScanner.Split(bufio.ScanLines)

	if err := c.Start(); err != nil {
		return &TaskError{
			AbsReadmeMd: t.AbsReadmeMd,
			Script:      "autorest",
			Message:     err.Error(),
		}
	}

	go func() {
		for errScanner.Scan() {
			line := errScanner.Text()
			fmt.Fprintln(os.Stderr, "[AUTOREST] "+line)
		}
	}()

	go func() {
		for outScanner.Scan() {
			line := outScanner.Text()
			fmt.Fprintln(os.Stdout, "[AUTOREST] "+line)
		}
	}()

	if err := c.Wait(); err != nil {
		return &TaskError{
			AbsReadmeMd: t.AbsReadmeMd,
			Script:      "autorest",
			Message:     err.Error(),
		}
	}
	return nil
}

func (t *Task) executeAfterScript(afterScripts []string) error {
	for _, script := range afterScripts {
		log.Printf("Executing after scripts %s...", script)
		arguments := strings.Split(script, " ")
		c := exec.Command(arguments[0], arguments[1:]...)
		output, err := c.CombinedOutput()
		if err != nil {
			return &TaskError{
				AbsReadmeMd: t.AbsReadmeMd,
				Script:      script,
				Message:     string(output),
			}
		}
	}

	return nil
}

// Options describes the options used in an autorest task
type Options struct {
	// AutorestArguments are the optional flags for the autorest tool
	AutorestArguments []string
	// AfterScripts are the scripts that need to be run after the SDK is generated
	AfterScripts []string
}

// NewOptionsFrom returns a new options from a io.Reader
func NewOptionsFrom(reader io.Reader) (*Options, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result Options
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// String ...
func (o Options) String() string {
	b, _ := json.MarshalIndent(o, "", "  ")
	return string(b)
}

// TaskError the error returned during an autorest task
type TaskError struct {
	// AbsReadmeMd relative path of the readme.md file to generate
	AbsReadmeMd string
	// Script the script running when the error is thrown
	Script string
	// Message the error message
	Message string
}

func (r *TaskError) Error() string {
	return fmt.Sprintf("autorest task failed for readme file '%s' during '%s': %s", r.AbsReadmeMd, r.Script, r.Message)
}
