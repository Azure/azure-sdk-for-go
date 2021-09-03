// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
)

// Generator collects all the related context of an autorest generation
type Generator struct {
	options model.Options
	cmd     *exec.Cmd
}

// NewGeneratorFromOptions returns a new Generator with the given model.Options
func NewGeneratorFromOptions(o model.Options) *Generator {
	return &Generator{
		options: o,
	}
}

// WithOption appends an model.Option to the argument list of the autorest generation
func (g *Generator) WithOption(option model.Option) *Generator {
	g.options = g.options.MergeOptions(option)
	return g
}

// WithTag appends a tag option to the autorest argument list
func (g *Generator) WithTag(tag string) *Generator {
	return g.WithOption(model.NewKeyValueOption("tag", tag))
}

// WithMultiAPI appends a multiapi flag to the autorest argument list
func (g *Generator) WithMultiAPI() *Generator {
	return g.WithOption(model.NewFlagOption("multiapi"))
}

// WithMetadataOutput appends a `metadata-output-folder` option to the autorest argument list
func (g *Generator) WithMetadataOutput(output string) *Generator {
	return g.WithOption(model.NewKeyValueOption("metadata-output-folder", output))
}

// WithReadme appends a readme argument
func (g *Generator) WithReadme(readme string) *Generator {
	return g.WithOption(model.NewArgument(readme))
}

// Generate executes the autorest generation. The error will be of type *GenerateError
func (g *Generator) Generate() error {
	g.buildCommand()
	o, err := g.cmd.CombinedOutput()
	if err != nil {
		return &GenerateError{
			Options: g.options,
			Message: string(o),
		}
	}
	return nil
}

func (g *Generator) buildCommand() {
	if g.cmd != nil {
		return
	}
	arguments := make([]string, len(g.options.Arguments()))
	for i, o := range g.options.Arguments() {
		arguments[i] = o.Format()
	}
	g.cmd = exec.Command("autorest", arguments...)
}

// Arguments returns the arguments which are using in the autorest command ('autorest' itself excluded)
func (g *Generator) Arguments() []model.Option {
	return g.options.Arguments()
}

// Start starts the generation
func (g *Generator) Start() error {
	g.buildCommand()
	if err := g.cmd.Start(); err != nil {
		return &GenerateError{
			Options: g.options,
			Message: err.Error(),
		}
	}
	return nil
}

// Wait waits for the generation to complete
func (g *Generator) Wait() error {
	g.buildCommand()
	if err := g.cmd.Wait(); err != nil {
		return &GenerateError{
			Options: g.options,
			Message: err.Error(),
		}
	}
	return nil
}

// Run starts and waits the generation
func (g *Generator) Run() error {
	g.buildCommand()
	if err := g.cmd.Run(); err != nil {
		return &GenerateError{
			Options: g.options,
			Message: err.Error(),
		}
	}
	return nil
}

// StdoutPipe returns the stdout pipeline of the command
func (g *Generator) StdoutPipe() (io.ReadCloser, error) {
	g.buildCommand()
	return g.cmd.StdoutPipe()
}

// StderrPipe returns the stderr pipeline of the command
func (g *Generator) StderrPipe() (io.ReadCloser, error) {
	g.buildCommand()
	return g.cmd.StderrPipe()
}

// GenerateError ...
type GenerateError struct {
	Options model.Options
	Message string
}

// Error ...
func (e *GenerateError) Error() string {
	arguments := make([]string, len(e.Options.Arguments()))
	for i, o := range e.Options.Arguments() {
		arguments[i] = o.Format()
	}
	return fmt.Sprintf("autorest error with arguments '%s': \n%s", strings.Join(arguments, ", "), e.Message)
}
