// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package armedgeactions_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/edgeactions/armedgeactions"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
)

const handlerCode = "function handler(event) { return event; }"

func versionCode(deploymentType armedgeactions.EdgeActionVersionDeploymentType) (string, error) {
	content := []byte(handlerCode)
	switch deploymentType {
	case armedgeactions.EdgeActionVersionDeploymentTypeFile:
	case armedgeactions.EdgeActionVersionDeploymentTypeZip:
		var buffer bytes.Buffer
		archive := zip.NewWriter(&buffer)
		// Fixed DOS timestamp and creator metadata keep request bodies identical across platforms.
		entry, err := archive.CreateHeader(&zip.FileHeader{
			Name: "handler.js", Method: zip.Store,
			ModifiedDate: 33, //nolint:staticcheck // Modified adds an extended timestamp, changing recorded ZIP payload bytes.
		})
		if err != nil {
			return "", err
		}
		if _, err = entry.Write(content); err != nil {
			return "", err
		}
		if err = archive.Close(); err != nil {
			return "", err
		}
		content = buffer.Bytes()
	default:
		return "", fmt.Errorf("unsupported deployment type %q", deploymentType)
	}
	return base64.StdEncoding.EncodeToString(content), nil
}

func hasStatus(err error, status int) bool {
	var responseError *azcore.ResponseError
	return errors.As(err, &responseError) && responseError.StatusCode == status
}

// A non-nil cleanup owns only the new action, including a partially failed create.
func createOwnedAction(ctx context.Context, factory *armedgeactions.ClientFactory, group, name string) (func(context.Context) error, error) {
	client := factory.NewClient()
	_, err := client.Get(ctx, group, name, nil)
	if err == nil {
		return nil, fmt.Errorf("refusing to overwrite existing EdgeAction %s/%s", group, name)
	}
	if !hasStatus(err, http.StatusNotFound) {
		return nil, err
	}
	cleanup := func(ctx context.Context) error { return cleanupAction(ctx, factory, group, name) }
	poller, err := client.BeginCreate(ctx, group, name, armedgeactions.EdgeAction{
		Location: to.Ptr("global"),
		SKU:      &armedgeactions.SKUType{Name: to.Ptr("Standard"), Tier: to.Ptr("Standard")},
		Tags:     map[string]*string{"purpose": to.Ptr("initial"), "remove-me": to.Ptr("initial")},
	}, nil)
	if err != nil {
		return cleanup, err
	}
	_, err = pollForScenario(ctx, poller)
	return cleanup, err
}

func cleanupAction(ctx context.Context, factory *armedgeactions.ClientFactory, group, name string) error {
	// A service contract failure must not prevent deletion of resources we own.
	ctx = context.WithValue(ctx, skipValidationKey{}, true)
	client := factory.NewClient()
	action, err := client.Get(ctx, group, name, nil)
	if err != nil {
		if hasStatus(err, http.StatusNotFound) {
			return nil
		}
		return err
	}
	if action.Properties == nil {
		return errors.New("cannot clean up an action without properties")
	}
	if err := waitForTerminal(ctx, action.Properties.ProvisioningState, func(ctx context.Context) (*armedgeactions.ProvisioningState, error) {
		response, err := client.Get(ctx, group, name, nil)
		if err != nil {
			return nil, err
		}
		if response.Properties == nil {
			return nil, errors.New("action GET returned no properties during cleanup")
		}
		return response.Properties.ProvisioningState, nil
	}); err != nil {
		return fmt.Errorf("wait for action before deletion: %w", err)
	}
	filters := factory.NewEdgeActionExecutionFiltersClient()
	filterPager := filters.NewListByEdgeActionPager(group, name, nil)
	var filterItems []*armedgeactions.EdgeActionExecutionFilter
	for filterPager.More() {
		page, err := filterPager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, filter := range page.Value {
			if filter == nil || filter.Name == nil {
				return errors.New("cannot clean up an execution filter without a name")
			}
			filterItems = append(filterItems, filter)
		}
	}
	versions := factory.NewEdgeActionVersionsClient()
	versionPager := versions.NewListByEdgeActionPager(group, name, nil)
	var items []*armedgeactions.EdgeActionVersion
	defaults := 0
	for versionPager.More() {
		page, err := versionPager.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, version := range page.Value {
			if version == nil || version.Name == nil || version.Properties == nil || version.Properties.IsDefaultVersion == nil {
				return errors.New("cannot clean up a version without a name and default state")
			}
			err := waitForTerminal(ctx, version.Properties.ProvisioningState, func(ctx context.Context) (*armedgeactions.ProvisioningState, error) {
				response, err := versions.Get(ctx, group, name, *version.Name, nil)
				if err != nil {
					return nil, err
				}
				if response.Properties == nil {
					return nil, errors.New("version GET returned no properties during cleanup")
				}
				version.Properties = response.Properties
				return response.Properties.ProvisioningState, nil
			})
			if err != nil {
				return fmt.Errorf("wait for version before deletion: %w", err)
			}
			if version.Properties.IsDefaultVersion == nil {
				return errors.New("version GET returned no default state during cleanup")
			}
			switch *version.Properties.IsDefaultVersion {
			case armedgeactions.EdgeActionIsDefaultVersionTrue:
				defaults++
			case armedgeactions.EdgeActionIsDefaultVersionFalse:
			default:
				return fmt.Errorf("invalid default state for version %s", *version.Name)
			}
			items = append(items, version)
		}
	}
	if defaults > 1 {
		return errors.New("cannot clean up versions with multiple defaults")
	}
	// Validate and collect all pages before deletion so paging cannot miss a removed child.
	sort.SliceStable(items, func(i, j int) bool {
		return *items[i].Properties.IsDefaultVersion == armedgeactions.EdgeActionIsDefaultVersionFalse &&
			*items[j].Properties.IsDefaultVersion == armedgeactions.EdgeActionIsDefaultVersionTrue
	})
	for _, filter := range filterItems {
		if filter.Properties == nil {
			return errors.New("cannot clean up a filter without properties")
		}
		err := waitForTerminal(ctx, filter.Properties.ProvisioningState, func(ctx context.Context) (*armedgeactions.ProvisioningState, error) {
			response, err := filters.Get(ctx, group, name, *filter.Name, nil)
			if err != nil {
				return nil, err
			}
			if response.Properties == nil {
				return nil, errors.New("filter GET returned no properties during cleanup")
			}
			return response.Properties.ProvisioningState, nil
		})
		if err != nil {
			return fmt.Errorf("wait for filter before deletion: %w", err)
		}
		poller, err := filters.BeginDelete(ctx, group, name, *filter.Name, nil)
		if err != nil {
			return err
		}
		if _, err = testutil.PollForTest(ctx, poller); err != nil {
			return err
		}
	}
	for _, version := range items {
		poller, err := versions.BeginDelete(ctx, group, name, *version.Name, nil)
		if err != nil {
			return err
		}

		if _, err = testutil.PollForTest(ctx, poller); err != nil {
			return err
		}
	}
	poller, err := client.BeginDelete(ctx, group, name, nil)
	if err != nil {
		return err
	}
	_, err = testutil.PollForTest(ctx, poller)
	return err
}

func waitForTerminal(ctx context.Context, state *armedgeactions.ProvisioningState, read func(context.Context) (*armedgeactions.ProvisioningState, error)) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if state == nil {
			return errors.New("cleanup requires a provisioning state")
		}
		if terminalStatus(string(*state)) {
			return nil
		}
		recording.Sleep(5 * time.Second)
		if err := ctx.Err(); err != nil {
			return err
		}
		var err error
		state, err = read(ctx)
		if err != nil {
			return err
		}
	}
}
