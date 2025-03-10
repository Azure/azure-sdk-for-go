// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tools

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func Run(remainingArgs []string) {
	onBadCommand := func() {
		fmt.Printf("ERROR: missing tool name\n")
		fmt.Printf("Usage: stress tools (delete|constantupdate|tempqueue|generatesas)\n")
		os.Exit(1)
	}

	if len(remainingArgs) == 0 {
		onBadCommand()
	}

	var ec int

	switch remainingArgs[0] {
	case "delete":
		ec = DeleteUsingRegexp(remainingArgs[1:])
	case "constantupdate":
		ec = ConstantlyUpdateQueue(remainingArgs[1:])
	case "tempqueue":
		ec = CreateTempQueue(remainingArgs[1:])
	case "generatesas":
		ec = GenerateSas(remainingArgs[1:])
	default:
		onBadCommand()
	}

	os.Exit(ec)
}

const queueEntityType = "queue"
const topicEntityType = "topic"

func addEntityTypeFlag(fs *flag.FlagSet, entityTypes ...string) func() (string, error) {
	all := strings.Join(entityTypes, ",")
	entityType := fs.String("type", "", fmt.Sprintf("The type of entity to target (%s)", all))

	return func() (string, error) {
		if *entityType != queueEntityType && *entityType != topicEntityType {
			return "", fmt.Errorf("invalid entity type '%s': must be one of (%s)", *entityType, all)
		}

		return *entityType, nil
	}
}
