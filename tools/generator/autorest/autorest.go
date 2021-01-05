package autorest

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
)

// Generator collects all the related context of an autorest generation
type Generator struct {
	arguments []string
	cmd       *exec.Cmd
}

// NewGeneratorFromOptions returns a new Generator with the given model.Options
func NewGeneratorFromOptions(o model.Options) *Generator {
	return &Generator{
		arguments: o.Arguments(),
	}
}

// WithOption appends an model.Option to the argument list of the autorest generation
func (g *Generator) WithOption(option model.Option) *Generator {
	g.arguments = append(g.arguments, option.Format())
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
			Arguments: g.arguments,
			Message:   string(o),
		}
	}
	return nil
}

func (g *Generator) buildCommand() {
	if g.cmd != nil {
		return
	}
	g.cmd = exec.Command("autorest", g.arguments...)
}

// Start starts the generation
func (g *Generator) Start() error {
	g.buildCommand()
	if err := g.cmd.Start(); err != nil {
		return &GenerateError{
			Arguments: g.arguments,
			Message:   err.Error(),
		}
	}
	return nil
}

// Wait waits for the generation to complete
func (g *Generator) Wait() error {
	g.buildCommand()
	if err := g.cmd.Wait(); err != nil {
		return &GenerateError{
			Arguments: g.arguments,
			Message:   err.Error(),
		}
	}
	return nil
}

// Run starts and waits the generation
func (g *Generator) Run() error {
	g.buildCommand()
	if err := g.cmd.Run(); err != nil {
		return &GenerateError{
			Arguments: g.arguments,
			Message:   err.Error(),
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
	Arguments []string
	Message   string
}

// Error ...
func (e *GenerateError) Error() string {
	return fmt.Sprintf("autorest error with arguments '%s': \n%s", e.Arguments, e.Message)
}
