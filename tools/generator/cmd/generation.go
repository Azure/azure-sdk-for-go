// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

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

	g *autorest.Generator
}

func (c autorestContext) generate() error {
	c.g = autorest.NewGeneratorFromOptions(c.options).WithReadme(c.absReadme).WithMetadataOutput(c.metadataOutput)

	stdout, _ := c.g.StdoutPipe()
	stderr, _ := c.g.StderrPipe()
	defer stdout.Close()
	defer stderr.Close()
	outScanner := bufio.NewScanner(stdout)
	outScanner.Split(bufio.ScanLines)
	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(bufio.ScanLines)

	if err := c.g.Start(); err != nil {
		return err
	}

	go printWithPrefixTo(outScanner, os.Stdout, "[AUTOREST] ")
	go printWithPrefixTo(errScanner, os.Stderr, "[AUTOREST] ")

	if err := c.g.Wait(); err != nil {
		return err
	}

	return nil
}

func (c autorestContext) autorestArguments() []string {
	return c.g.Arguments()
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
