//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azqueue_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue/v2"
	"log"
	"os"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Example_client_NewClient() {
	// this example uses Azure Active Directory (AAD) to authenticate with Azure Queue Storage
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	// https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	fmt.Println(client.URL())
}

func Example_client_NewClientWithSharedKeyCredential() {
	// this example uses a shared key to authenticate with Azure queue Storage
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	accountKey, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_KEY")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_KEY could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	// shared key authentication requires the storage account name and access key
	cred, err := azqueue.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)
	serviceClient, err := azqueue.NewServiceClientWithSharedKeyCredential(serviceURL, cred, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_client_NewClientFromConnectionString() {
	// this example uses a connection string to authenticate with Azure queue Storage
	connectionString, ok := os.LookupEnv("AZURE_STORAGE_CONNECTION_STRING")
	if !ok {
		log.Fatal("the environment variable 'AZURE_STORAGE_CONNECTION_STRING' could not be found")
	}

	serviceClient, err := azqueue.NewServiceClientFromConnectionString(connectionString, nil)
	handleError(err)
	fmt.Println(serviceClient.URL())
}

func Example_client_CreateQueue() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.CreateQueue(context.TODO(), "testqueue", &azqueue.CreateOptions{
		Metadata: map[string]*string{"hello": to.Ptr("world")},
	})
	handleError(err)
	fmt.Println(resp)

	// delete the queue
	_, err = client.DeleteQueue(context.TODO(), "testqueue", nil)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_DeleteQueue() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	opts := &azqueue.DeleteOptions{} // or just pass nil to the method below
	resp, err := client.DeleteQueue(context.TODO(), "testqueue", opts)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_NewListQueuesPager() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	pager := client.NewListQueuesPager(&azqueue.ListQueuesOptions{
		Include: azqueue.ListQueuesInclude{Metadata: true},
	})

	// list pre-existing queues
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		handleError(err) // if err is not nil, break the loop.
		for _, _queue := range resp.Queues {
			fmt.Printf("%v", _queue)
		}
	}
}

func Example_client_Enqueue_DequeueMessage() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.CreateQueue(context.TODO(), "testqueue", &azqueue.CreateOptions{
		Metadata: map[string]*string{"hello": to.Ptr("world")},
	})
	handleError(err)
	fmt.Println(resp)

	opts := &azqueue.EnqueueMessageOptions{TimeToLive: to.Ptr(int32(10))}
	queueClient := client.NewQueueClient("testqueue")
	resp1, err := queueClient.EnqueueMessage(context.Background(), "test content", opts)
	handleError(err)
	fmt.Println(resp1)

	resp2, err := queueClient.DequeueMessage(context.Background(), nil)
	handleError(err)
	// check message content
	fmt.Println(resp2.Messages[0].MessageText)

	// delete the queue
	_, err = client.DeleteQueue(context.TODO(), "testqueue", nil)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_PeekMessages() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.CreateQueue(context.TODO(), "testqueue", &azqueue.CreateOptions{
		Metadata: map[string]*string{"hello": to.Ptr("world")},
	})
	handleError(err)
	fmt.Println(resp)

	queueClient := client.NewQueueClient("testqueue")

	// enqueue 4 messages
	for i := 0; i < 4; i++ {
		resp1, err := queueClient.EnqueueMessage(context.Background(), "test content", nil)
		handleError(err)
		fmt.Println(resp1)
	}
	// only check 3 messages
	opts := &azqueue.PeekMessagesOptions{NumberOfMessages: to.Ptr(int32(3))}
	resp2, err := queueClient.PeekMessages(context.Background(), opts)
	handleError(err)
	// check 3 messages retrieved
	fmt.Println(len(resp2.Messages))

	// delete the queue
	_, err = client.DeleteQueue(context.TODO(), "testqueue", nil)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_Update_Message() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.CreateQueue(context.TODO(), "testqueue", &azqueue.CreateOptions{
		Metadata: map[string]*string{"hello": to.Ptr("world")},
	})
	handleError(err)
	fmt.Println(resp)

	queueClient := client.NewQueueClient("testqueue")
	resp1, err := queueClient.EnqueueMessage(context.Background(), "test content", nil)
	handleError(err)
	fmt.Println(resp1)
	popReceipt := *resp1.Messages[0].PopReceipt
	messageID := *resp1.Messages[0].MessageID

	opts := &azqueue.UpdateMessageOptions{}
	_, err = queueClient.UpdateMessage(context.Background(), messageID, popReceipt, "new content", opts)
	handleError(err)

	resp3, err := queueClient.DequeueMessage(context.Background(), nil)
	handleError(err)
	// check message content has updated
	fmt.Println(resp3.Messages[0].MessageText)

	// delete the queue
	_, err = client.DeleteQueue(context.TODO(), "testqueue", nil)
	handleError(err)
	fmt.Println(resp)
}

func Example_client_Clear_Messages() {
	accountName, ok := os.LookupEnv("AZURE_STORAGE_ACCOUNT_NAME")
	if !ok {
		panic("AZURE_STORAGE_ACCOUNT_NAME could not be found")
	}
	serviceURL := fmt.Sprintf("https://%s.queue.core.windows.net/", accountName)

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	handleError(err)

	client, err := azqueue.NewServiceClient(serviceURL, cred, nil)
	handleError(err)

	resp, err := client.CreateQueue(context.TODO(), "testqueue", &azqueue.CreateOptions{
		Metadata: map[string]*string{"hello": to.Ptr("world")},
	})
	handleError(err)
	fmt.Println(resp)

	queueClient := client.NewQueueClient("testqueue")
	resp1, err := queueClient.EnqueueMessage(context.Background(), "test content", nil)
	handleError(err)
	fmt.Println(resp1)

	resp2, err := queueClient.ClearMessages(context.Background(), nil)
	handleError(err)
	fmt.Println(resp2)

	// delete the queue
	_, err = client.DeleteQueue(context.TODO(), "testqueue", nil)
	handleError(err)
	fmt.Println(resp)
}
