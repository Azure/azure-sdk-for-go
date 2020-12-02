package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
)

type autorestContext struct {
	absReadme      string
	metadataOutput string
	options        model.Options
}

func (c autorestContext) generate() error {
	g := autorest.NewGeneratorFromOptions(c.options).WithReadme(c.absReadme).WithMetadataOutput(c.metadataOutput)

	stdout, _ := g.StdoutPipe()
	stderr, _ := g.StderrPipe()
	defer stdout.Close()
	defer stderr.Close()
	outScanner := bufio.NewScanner(stdout)
	outScanner.Split(bufio.ScanLines)
	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(bufio.ScanLines)

	if err := g.Start(); err != nil {
		return err
	}

	go printWithPrefixTo(outScanner, os.Stdout, "[AUTOREST] ")
	go printWithPrefixTo(errScanner, os.Stderr, "[AUTOREST] ")

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func printWithPrefixTo(scanner *bufio.Scanner, writer io.Writer, prefix string) error {
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := fmt.Fprintln(writer, fmt.Sprintf("%s%s", prefix, line)); err != nil {
			return err
		}
	}
	return nil
}
