// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azbatch_test

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/batch/azbatch"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

const recordingDir = "sdk/batch/azbatch/testdata"

var endpoint = "https://batch.local"

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	if ep, ok := os.LookupEnv("AZBATCH_ENDPOINT"); ok {
		endpoint = "https://" + ep
	}
	if recording.GetRecordMode() != recording.LiveMode {
		if proxy, err := recording.StartTestProxy(recordingDir, nil); err == nil {
			defer func() {
				if err := recording.StopTestProxy(proxy); err != nil {
					log.Fatal(err)
				}
			}()
		} else {
			return 1, err
		}
		if err := recording.RemoveRegisteredSanitizers([]string{
			"AZSDK3430", // $..id
			"AZSDK3493", // $..name
			"AZSDK4001", // default host replacement which doesn't replace region; adding a more robust one below
		}, nil); err != nil {
			return 1, err
		}
		u, err := url.Parse(endpoint)
		if err != nil {
			return 1, err
		}
		if err = recording.AddGeneralRegexSanitizer("batch.local", u.Host, nil); err != nil {
			return 1, err
		}
		if err = recording.AddBodyKeySanitizer("$.startTime", "42", "", nil); err != nil {
			return 1, err
		}
	}
	return m.Run(), nil
}

func createDefaultPool(t *testing.T) (*azbatch.Client, string) {
	client := record(t)
	pool := defaultPoolContent(t)
	_, err := client.CreatePool(ctx, pool, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _, _ = client.DeletePool(ctx, *pool.ID, nil) })
	return client, *pool.ID
}

func defaultPoolContent(t *testing.T) azbatch.CreatePoolContent {
	return azbatch.CreatePoolContent{
		ID:                   to.Ptr(randomString(t)),
		TargetDedicatedNodes: to.Ptr(int32(1)),
		TaskSchedulingPolicy: &azbatch.TaskSchedulingPolicy{
			NodeFillType: to.Ptr(azbatch.NodeFillTypePack),
		},
		VirtualMachineConfiguration: &azbatch.VirtualMachineConfiguration{
			DataDisks: []*azbatch.DataDisk{
				{
					DiskSizeGB:        to.Ptr(int32(1)),
					LogicalUnitNumber: to.Ptr(int32(1)),
				},
			},
			ImageReference: &azbatch.ImageReference{
				Offer:     to.Ptr("0001-com-ubuntu-server-jammy"),
				Publisher: to.Ptr("canonical"),
				SKU:       to.Ptr("22_04-lts"),
			},
			NodeAgentSKUID: to.Ptr("batch.node.ubuntu 22.04"),
		},
		VMSize: to.Ptr("Standard_A1_v2"),
	}
}

// firstReadyNode returns the first node in the pool that's ready to run tasks.
// It fails the test when no such node is found within 6 minutes.
func firstReadyNode(t *testing.T, client *azbatch.Client, poolID string) azbatch.Node {
	// note this assumes the pool has exactly one node, which is true for all test pools at time of writing
	node, err := poll(
		func() *azbatch.Node {
			var node *azbatch.Node
			for pgr := client.NewListNodesPager(poolID, nil); pgr.More(); {
				pg, err := pgr.NextPage(ctx)
				require.NoError(t, err)
				for _, node = range pg.Value {
					if node != nil {
						break
					}
				}
			}
			return node
		},
		func(n *azbatch.Node) bool {
			return n != nil && n.State != nil && (*n.State == azbatch.NodeStateIdle || *n.State == azbatch.NodeStateRunning)
		},
		6*time.Minute,
	)
	require.NoError(t, err)
	require.NotNil(t, node, "found no ready node")
	return *node
}

func poll[T any](get func() T, done func(T) bool, timeout time.Duration) (T, error) {
	const delay = 14 * time.Second
	ticks := int(timeout / delay)
	var t T
	for i := 0; i < ticks; i++ {
		t = get()
		if done(t) {
			return t, nil
		}
		if i < ticks-1 {
			recording.Sleep(delay)
		}
	}
	return t, errors.New("polling timed out")
}

func randomString(t *testing.T) string {
	id, err := recording.GenerateAlphaNumericID(t, t.Name(), 24, false)
	require.NoError(t, err)
	return strings.ReplaceAll(id, "/", "_")
}

func record(t *testing.T) *azbatch.Client {
	err := recording.Start(t, recordingDir, nil)
	require.NoError(t, err)
	t.Cleanup(func() {
		err := recording.Stop(t, nil)
		require.NoError(t, err)
	})
	transport, err := recording.NewRecordingHTTPClient(t, nil)
	require.NoError(t, err)
	cred, err := credential.New(nil)
	require.NoError(t, err)
	c, err := azbatch.NewClient(endpoint, cred, &azbatch.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Transport: transport,
		},
	})
	require.NoError(t, err)
	return c
}

// waitForTask creates a task with the given command line and waits 5 minutes for it to complete.
// It fails the test if the task doesn't complete within that time.
func waitForTask(t *testing.T, client *azbatch.Client, jobID, commandLine string) azbatch.Task {
	tid := randomString(t)
	_, err := client.CreateTask(ctx, jobID, azbatch.CreateTaskContent{
		CommandLine: to.Ptr(commandLine),
		ID:          to.Ptr(tid),
	}, nil)
	require.NoError(t, err)

	task, err := poll(
		func() azbatch.Task {
			gt, err := client.GetTask(ctx, jobID, tid, nil)
			require.NoError(t, err)
			return gt.Task
		},
		func(task azbatch.Task) bool {
			return task.State != nil && *task.State == azbatch.TaskStateCompleted
		},
		5*time.Minute,
	)
	require.NoError(t, err, "task isn't complete")
	return task
}
