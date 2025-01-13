// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package tools

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/stress/shared"
)

func DeleteUsingRegexp(remainingArgs []string) int {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)

	entityTypeFn := addEntityTypeFlag(fs, queueEntityType, topicEntityType)
	entityRegex := fs.String("re", "", "The regex to match against the entity name")
	clientCreator := shared.AddAuthFlags(fs)

	if err := fs.Parse(remainingArgs); err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		fs.PrintDefaults()
		return 1
	}

	entityType, err := entityTypeFn()

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		fs.PrintDefaults()
		return 1
	}

	_ = shared.LoadEnvironment()

	_, adminClient, err := clientCreator()

	if err != nil {
		fmt.Printf("Failed to create client: %s", err.Error())
		return 1
	}

	re := regexp.MustCompile(*entityRegex)

	ctx, cancel := shared.NewCtrlCContext()
	defer cancel()

	switch entityType {
	case "queue":
		pager := adminClient.NewListQueuesPager(nil)

		var queuesToDelete []string

		for pager.More() {
			page, err := pager.NextPage(ctx)

			if err != nil {
				log.Printf("Failed to get queues: %s", err)
				return 1
			}

			for _, queueProps := range page.Queues {
				if re.MatchString(queueProps.QueueName) {
					queuesToDelete = append(queuesToDelete, queueProps.QueueName)
				}
			}
		}

		log.Printf("Deleting %d queues", len(queuesToDelete))

		for _, queue := range queuesToDelete {
			if _, err := adminClient.DeleteQueue(ctx, queue, nil); err != nil {
				fmt.Printf("Failed to delete %s: %s\n", queue, err.Error())
				return 1
			} else {
				log.Printf("Deleted %s", queue)
			}
		}
	case "topic":
		pager := adminClient.NewListTopicsPager(nil)

		var topicsToDelete []string

		for pager.More() {
			page, err := pager.NextPage(ctx)

			if err != nil {
				fmt.Printf("Failed to get topics: %s\n", err)
				return 1
			}

			for _, topicProps := range page.Topics {
				if re.MatchString(topicProps.TopicName) {
					topicsToDelete = append(topicsToDelete, topicProps.TopicName)
				}
			}
		}

		log.Printf("Deleting %d topics", len(topicsToDelete))

		for _, topic := range topicsToDelete {
			if _, err := adminClient.DeleteTopic(ctx, topic, nil); err != nil {
				fmt.Printf("Failed to delete %s: %s\n", topic, err.Error())
				return 1
			} else {
				log.Printf("Deleted %s", topic)
			}
		}
	}

	return 0
}
