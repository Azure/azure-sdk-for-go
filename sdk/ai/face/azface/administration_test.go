// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azface_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/face/azface"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestAdministrationLargeFaceList_CreateGetDelete(t *testing.T) {
	adminClient := newAdministrationClientForTest(t)
	largeFaceListClient := adminClient.NewAdministrationLargeFaceListClient()
	
	// Create a unique ID for the test
	timestamp := time.Now().Unix()
	largeFaceListID := fmt.Sprintf("test-large-face-list-%d", timestamp)
	name := fmt.Sprintf("Test Large Face List %d", timestamp)
	
	// Create large face list
	_, err := largeFaceListClient.Create(context.Background(), largeFaceListID, name, &azface.AdministrationLargeFaceListClientCreateOptions{
		UserData: to.Ptr("Test large face list for SDK testing"),
	})
	require.NoError(t, err)
	
	// Get the created large face list
	getResp, err := largeFaceListClient.Get(context.Background(), largeFaceListID, nil)
	require.NoError(t, err)
	require.NotNil(t, getResp.LargeFaceList)
	require.Equal(t, largeFaceListID, *getResp.LargeFaceList.LargeFaceListID)
	require.Equal(t, name, *getResp.LargeFaceList.Name)
	
	// List all large face lists (should include the one we just created)
	listResp, err := largeFaceListClient.GetLargeFaceLists(context.Background(), &azface.AdministrationLargeFaceListClientGetLargeFaceListsOptions{
		Start: to.Ptr(""),
		Top:   to.Ptr(int32(1000)),
	})
	require.NoError(t, err)
	require.NotNil(t, listResp.LargeFaceListArray)
	
	// Verify our list is in the results
	found := false
	for _, list := range listResp.LargeFaceListArray {
		if list.LargeFaceListID != nil && *list.LargeFaceListID == largeFaceListID {
			found = true
			break
		}
	}
	require.True(t, found, "Created large face list should be in the list")
	
	// Clean up - delete the large face list
	_, err = largeFaceListClient.Delete(context.Background(), largeFaceListID, nil)
	require.NoError(t, err)
}

func TestAdministrationLargeFaceList_AddFaceFromURL(t *testing.T) {
	adminClient := newAdministrationClientForTest(t)
	largeFaceListClient := adminClient.NewAdministrationLargeFaceListClient()
	
	// Create a unique ID for the test
	timestamp := time.Now().Unix()
	largeFaceListID := fmt.Sprintf("test-large-face-list-face-%d", timestamp)
	name := fmt.Sprintf("Test Large Face List with Face %d", timestamp)
	
	// Create large face list first
	_, err := largeFaceListClient.Create(context.Background(), largeFaceListID, name, nil)
	require.NoError(t, err)
	
	defer func() {
		// Clean up - delete the large face list
		largeFaceListClient.Delete(context.Background(), largeFaceListID, nil)
	}()
	
	// Add a face from URL
	addResp, err := largeFaceListClient.AddFaceFromURL(context.Background(), largeFaceListID, sampleImageURL, &azface.AdministrationLargeFaceListClientAddFaceFromURLOptions{
		UserData: to.Ptr("Test face added from URL"),
		TargetFace: []int32{}, // Empty slice for whole image
	})
	require.NoError(t, err)
	require.NotNil(t, addResp.PersistedFaceID)
	
	persistedFaceID := *addResp.PersistedFaceID
	
	// Get the added face
	getFaceResp, err := largeFaceListClient.GetFace(context.Background(), largeFaceListID, persistedFaceID, nil)
	require.NoError(t, err)
	require.NotNil(t, getFaceResp.PersistedFaceID)
	require.Equal(t, persistedFaceID, *getFaceResp.PersistedFaceID)
	
	// List faces in the large face list
	listFacesResp, err := largeFaceListClient.GetFaces(context.Background(), largeFaceListID, &azface.AdministrationLargeFaceListClientGetFacesOptions{
		Start: to.Ptr(""),
		Top:   to.Ptr(int32(1000)),
	})
	require.NoError(t, err)
	require.NotNil(t, listFacesResp.LargeFaceListFaceArray)
	require.GreaterOrEqual(t, len(listFacesResp.LargeFaceListFaceArray), 1)
	
	// Verify our face is in the results
	found := false
	for _, face := range listFacesResp.LargeFaceListFaceArray {
		if face.PersistedFaceID != nil && *face.PersistedFaceID == persistedFaceID {
			found = true
			break
		}
	}
	require.True(t, found, "Added face should be in the list")
	
	// Delete the face
	_, err = largeFaceListClient.DeleteFace(context.Background(), largeFaceListID, persistedFaceID, nil)
	require.NoError(t, err)
}