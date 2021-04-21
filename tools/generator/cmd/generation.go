// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"bufio"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
)

type autorestContext struct {
	generator *autorest.Generator
}

func (c autorestContext) generate() error {
	stdout, _ := c.generator.StdoutPipe()
	stderr, _ := c.generator.StderrPipe()
	defer stdout.Close()
	defer stderr.Close()
	outScanner := bufio.NewScanner(stdout)
	outScanner.Split(bufio.ScanLines)
	errScanner := bufio.NewScanner(stderr)
	errScanner.Split(bufio.ScanLines)

	if err := c.generator.Start(); err != nil {
		return err
	}

	go printWithPrefixTo(outScanner, os.Stdout, "[AUTOREST] ")
	go printWithPrefixTo(errScanner, os.Stderr, "[AUTOREST] ")

	if err := c.generator.Wait(); err != nil {
		return err
	}

	return nil
}

func (c autorestContext) autorestArguments() []model.Option {
	return c.generator.Arguments()
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
