package autorest

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
)

type Generator struct {
	arguments []string
	cmd       *exec.Cmd
}

func FromOptions(o model.Options) *Generator {
	return &Generator{
		arguments: o.Arguments(),
	}
}

func (g *Generator) WithOption(option model.Option) *Generator {
	g.arguments = append(g.arguments, option.Format())
	return g
}

func (g *Generator) WithTag(tag string) *Generator {
	return g.WithOption(model.NewKeyValueOption("tag", tag))
}

func (g *Generator) WithMultiApi() *Generator {
	return g.WithOption(model.NewFlagOption("multiapi"))
}

func (g *Generator) WithMetadataOutput(output string) *Generator {
	return g.WithOption(model.NewKeyValueOption("metadata-output-folder", output))
}

func (g *Generator) WithReadme(readme string) *Generator {
	return g.WithOption(model.NewArgument(readme))
}

func (g *Generator) Generate() error {
	g.buildCommand()
	c := exec.Command("autorest", g.arguments...)
	o, err := c.CombinedOutput()
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

func (g *Generator) StdoutPipe() (io.ReadCloser, error) {
	g.buildCommand()
	return g.cmd.StdoutPipe()
}

func (g *Generator) StderrPipe() (io.ReadCloser, error) {
	g.buildCommand()
	return g.cmd.StderrPipe()
}

type GenerateError struct {
	Arguments []string
	Message   string
}

func (e *GenerateError) Error() string {
	return fmt.Sprintf("autorest error with arguments '%s': \n%s", e.Arguments, e.Message)
}
