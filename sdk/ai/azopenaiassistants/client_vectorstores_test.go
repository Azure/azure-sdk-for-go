//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestAssistantsWithVectorStores(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	fn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		uploadFileResp, err := client.UploadFile(context.Background(), strings.NewReader("hello world"), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
			Filename: getFileName(t, "txt"),
		})
		require.NoError(t, err)
		fileID := *uploadFileResp.ID

		defer mustDeleteFile(t, client, fileID)

		// you can associate a vector store with an assistant in one of two ways:
		// 1. Create the vector store separately, then associate it using FileSearch's option to pass in a vector store.
		// 2. As part of a CreateAssistant call, when added ot the FileSearch tool's options that let you pass in files and have it auto-create
		//     the vector store for you.

		// this is the "create vector store separately" path
		{
			createVectorStoreResp, err := client.CreateVectorStore(context.Background(), azopenaiassistants.VectorStoreBody{
				Name: to.Ptr("test vector store"),
			}, nil)
			require.NoError(t, err)

			vectorStoreID := *createVectorStoreResp.ID

			defer mustDeleteVectorStore(t, client, vectorStoreID)

			createVectorFileResp, err := client.CreateVectorStoreFile(context.Background(), vectorStoreID, fileID, nil)
			require.NoError(t, err)
			require.NotEmpty(t, createVectorFileResp)

			requireVectorStoreState(t, client, vectorStoreID, fileID)

			createAsstResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
				DeploymentName: &assistantsModel,
				Name:           to.Ptr("go - vector store test"),
				Tools: []azopenaiassistants.ToolDefinitionClassification{
					&azopenaiassistants.FileSearchToolDefinition{},
				},
				ToolResources: &azopenaiassistants.AssistantCreationOptionsToolResources{
					FileSearch: &azopenaiassistants.CreateFileSearchToolResourceOptions{
						VectorStoreIDs: []string{*createVectorStoreResp.ID},
					},
				},
			}, nil)
			requireNoErr(t, azure, err)

			defer mustDeleteAssistant(t, client, *createAsstResp.ID, azure)
		}

		// this is the "create vector store using CreateAssistant" path
		{
			createAsstResp, err := client.CreateAssistant(context.Background(), azopenaiassistants.CreateAssistantBody{
				DeploymentName: &assistantsModel,
				Name:           to.Ptr("go - vector store test"),
				Tools: []azopenaiassistants.ToolDefinitionClassification{
					&azopenaiassistants.FileSearchToolDefinition{},
				},
				ToolResources: &azopenaiassistants.AssistantCreationOptionsToolResources{
					FileSearch: &azopenaiassistants.CreateFileSearchToolResourceOptions{
						VectorStores: []azopenaiassistants.CreateFileSearchToolResourceVectorStoreOptions{
							{FileIDs: []string{fileID}},
						},
					},
				},
			}, nil)
			requireNoErr(t, azure, err)

			vectorStoreID := createAsstResp.ToolResources.FileSearch.VectorStoreIDs[0]
			requireVectorStoreState(t, client, vectorStoreID, fileID)

			mustDeleteAssistant(t, client, *createAsstResp.ID, azure)
			mustDeleteVectorStore(t, client, createAsstResp.ToolResources.FileSearch.VectorStoreIDs[0])
		}
	}

	t.Run("OpenAI", func(t *testing.T) {
		fn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		fn(t, true)
	})
}

func TestVectorStores(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	fn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		seed := time.Now().UnixNano()

		// create a couple of vector stores in the (unlikely) case, that we don't have enough to test with...
		resp1, err := client.CreateVectorStore(context.Background(), azopenaiassistants.VectorStoreBody{
			Name: to.Ptr(fmt.Sprintf("go-1-%X", seed)),
		}, nil)
		require.NoError(t, err)

		defer mustDeleteVectorStore(t, client, *resp1.ID)

		resp2, err := client.CreateVectorStore(context.Background(), azopenaiassistants.VectorStoreBody{
			Name: to.Ptr(fmt.Sprintf("go-2-%X", seed)),
		}, nil)
		require.NoError(t, err)

		defer mustDeleteVectorStore(t, client, *resp2.ID)

		// there's no scoping here (apart form the _entire_ resource) so it's possible our vector store isn't in here.
		listVectorStoresPager := client.NewListVectorStoresPager(&azopenaiassistants.ListVectorStoresOptions{
			Order: to.Ptr(azopenaiassistants.ListSortOrderDescending),
			Limit: to.Ptr[int32](1),
		})
		require.NoError(t, err)

		for i := 0; i < 2; i++ {
			require.True(t, listVectorStoresPager.More())

			resp, err := listVectorStoresPager.NextPage(context.Background())
			require.NoError(t, err)

			// it's highly likely that other people have (like me) polluted the global namespace with a ton of vector stores
			// so I'm not really looking for mine, just ensuring that the paging parameters work.
			require.Len(t, resp.Data, 1)
		}

		modifyResp, err := client.ModifyVectorStore(context.Background(), *resp1.ID, azopenaiassistants.VectorStoreUpdateBody{
			ExpiresAfter: &azopenaiassistants.VectorStoreUpdateOptionsExpiresAfter{
				Anchor: to.Ptr(azopenaiassistants.VectorStoreExpirationPolicyAnchorLastActiveAt),
				Days:   to.Ptr[int32](1),
			},
		}, nil)

		require.NoError(t, err)
		require.NotZero(t, modifyResp.ExpiresAfter)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		fn(t, true)
	})

	t.Run("OpenAI", func(t *testing.T) {
		fn(t, false)
	})
}

func TestVectorStoresWithBatch(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	fn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{
			Azure: azure,
		})

		seed := time.Now().UnixNano()

		createVectorResp, err := client.CreateVectorStore(context.Background(), azopenaiassistants.VectorStoreBody{
			Name: to.Ptr(fmt.Sprintf("go-vsb-%X", seed)),
		}, nil)
		require.NoError(t, err)

		defer mustDeleteVectorStore(t, client, *createVectorResp.ID)

		const numUploaded = 3
		fileIDMap := map[string]bool{}
		var fileIDs []string

		for i := 0; i < numUploaded; i++ {
			fileName := getFileName(t, "txt")
			uploadFileResp, err := client.UploadFile(context.Background(), strings.NewReader("hello world"), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
				Filename: fileName,
			})
			require.NoError(t, err)
			fileIDMap[*uploadFileResp.ID] = true
			fileIDs = append(fileIDs, *uploadFileResp.ID)
		}

		createBatchResp, err := client.CreateVectorStoreFileBatch(context.Background(), *createVectorResp.ID, azopenaiassistants.CreateVectorStoreFileBatchBody{
			FileIDs: fileIDs,
		}, nil)
		require.NoError(t, err)

		getResp, err := client.GetVectorStoreFileBatch(context.Background(), *createVectorResp.ID, *createBatchResp.ID, nil)
		require.NoError(t, err)
		require.Equal(t, *createBatchResp.ID, *getResp.ID)

		pager := client.NewListVectorStoreFileBatchFilesPager(*createVectorResp.ID, *createBatchResp.ID, &azopenaiassistants.ListVectorStoreFileBatchFilesOptions{
			Limit: to.Ptr[int32](2), // split all of our files across two pages
		})

		found := 0

		for pager.More() {
			resp, err := pager.NextPage(context.Background())
			require.NoError(t, err)

			for _, v := range resp.Data {
				require.True(t, fileIDMap[*v.ID])
				delete(fileIDMap, *v.ID)
				found++
			}

			require.LessOrEqual(t, len(resp.Data), 2)
		}

		require.Equal(t, numUploaded, found)
	}

	t.Run("AzureOpenAI", func(t *testing.T) {
		fn(t, true)
	})

	t.Run("OpenAI", func(t *testing.T) {
		fn(t, false)
	})
}

func requireVectorStoreState(t *testing.T, client *azopenaiassistants.Client, vectorStoreID string, fileID string) {
	// let's just do some simple vector store operations to make sure it all worked
	pager := client.NewListVectorStoreFilesPager(vectorStoreID, &azopenaiassistants.ListVectorStoreFilesOptions{
		Limit: to.Ptr[int32](10),
	})
	require.True(t, pager.More())

	listVectorFilesResp, err := pager.NextPage(context.Background())
	require.NoError(t, err)

	// we've only got one file in our vector store
	require.Equal(t, *(listVectorFilesResp.Data[0].ID), fileID)

	getVectorStoreFileResp, err := client.GetVectorStoreFile(context.Background(), vectorStoreID, fileID, nil)
	require.NoError(t, err)

	// sometimes the file is still uploading, even when I get here.
	require.Contains(t, []azopenaiassistants.VectorStoreFileStatus{
		azopenaiassistants.VectorStoreFileStatusCompleted,
		azopenaiassistants.VectorStoreFileStatusInProgress,
	}, *getVectorStoreFileResp.Status)
}

func mustDeleteFile(t *testing.T, client *azopenaiassistants.Client, fileID string) {
	_, err := client.DeleteFile(context.Background(), fileID, nil)
	require.NoError(t, err)
}

func mustDeleteVectorStore(t *testing.T, client *azopenaiassistants.Client, vectorStoreID string) {
	_, err := client.DeleteVectorStore(context.Background(), vectorStoreID, nil)
	require.NoError(t, err)
}

func mustDeleteAssistant(t *testing.T, client *azopenaiassistants.Client, assistantID string, azure bool) {
	_, err := client.DeleteAssistant(context.Background(), assistantID, nil)
	requireNoErr(t, azure, err)
}
