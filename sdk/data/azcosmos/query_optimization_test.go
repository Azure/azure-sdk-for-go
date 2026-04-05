// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/queryengine"
	"github.com/stretchr/testify/require"
)

func TestIsSimpleQuery(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected bool
	}{
		{
			name:     "basic select",
			query:    "SELECT * FROM c",
			expected: true,
		},
		{
			name:     "select with where",
			query:    "SELECT * FROM c WHERE c.id = @id",
			expected: true,
		},
		{
			name:     "select with projection",
			query:    "SELECT c.id, c.name FROM c WHERE c.type = 'user'",
			expected: true,
		},
		{
			name:     "select with ORDER BY",
			query:    "SELECT * FROM c ORDER BY c.timestamp",
			expected: false,
		},
		{
			name:     "select with GROUP BY",
			query:    "SELECT c.type, COUNT(1) FROM c GROUP BY c.type",
			expected: false,
		},
		{
			name:     "select with DISTINCT",
			query:    "SELECT DISTINCT c.type FROM c",
			expected: false,
		},
		{
			name:     "select with TOP",
			query:    "SELECT TOP 10 * FROM c",
			expected: false,
		},
		{
			name:     "select with OFFSET",
			query:    "SELECT * FROM c OFFSET 10 LIMIT 10",
			expected: false,
		},
		{
			name:     "select with COUNT",
			query:    "SELECT COUNT(1) FROM c",
			expected: false,
		},
		{
			name:     "select with SUM",
			query:    "SELECT SUM(c.amount) FROM c",
			expected: false,
		},
		{
			name:     "select with AVG",
			query:    "SELECT AVG(c.score) FROM c",
			expected: false,
		},
		{
			name:     "select with MIN",
			query:    "SELECT MIN(c.price) FROM c",
			expected: false,
		},
		{
			name:     "select with MAX",
			query:    "SELECT MAX(c.price) FROM c",
			expected: false,
		},
		{
			name:     "lowercase keywords",
			query:    "select * from c order by c.id",
			expected: false,
		},
		{
			name:     "mixed case",
			query:    "SELECT * FROM c Order By c.id",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSimpleQuery(tt.query)
			require.Equal(t, tt.expected, result, "query: %s", tt.query)
		})
	}
}

func TestSelectQueryExecutionMode_NoDirectTransport(t *testing.T) {
	container := &ContainerClient{
		database: &DatabaseClient{
			client: &Client{
				directTransport: nil,
			},
		},
	}

	pk := NewPartitionKeyString("test")
	opts := &QueryOptions{}

	mode := container.selectQueryExecutionMode(pk, opts)
	require.Equal(t, queryModeGateway, mode)
}

func TestSelectQueryExecutionMode_WithQueryEngine(t *testing.T) {
	container := &ContainerClient{
		database: &DatabaseClient{
			client: &Client{
				directTransport: &directModeTransport{},
			},
		},
	}

	pk := NewPartitionKeyString("test")
	opts := &QueryOptions{
		QueryEngine: &mockQueryEngine{},
	}

	mode := container.selectQueryExecutionMode(pk, opts)
	require.Equal(t, queryModeEngine, mode)
}

func TestSelectQueryExecutionMode_ODEEnabled(t *testing.T) {
	container := &ContainerClient{
		database: &DatabaseClient{
			client: &Client{
				directTransport: &directModeTransport{},
			},
		},
	}

	pk := NewPartitionKeyString("test")
	opts := &QueryOptions{}

	mode := container.selectQueryExecutionMode(pk, opts)
	require.Equal(t, queryModeODE, mode)
}

func TestSelectQueryExecutionMode_ODEDisabled(t *testing.T) {
	container := &ContainerClient{
		database: &DatabaseClient{
			client: &Client{
				directTransport: &directModeTransport{},
			},
		},
	}

	pk := NewPartitionKeyString("test")
	disabled := false
	opts := &QueryOptions{
		EnableOptimisticDirectExecution: &disabled,
	}

	mode := container.selectQueryExecutionMode(pk, opts)
	require.Equal(t, queryModeGateway, mode)
}

func TestSelectQueryExecutionMode_NoPartitionKey(t *testing.T) {
	container := &ContainerClient{
		database: &DatabaseClient{
			client: &Client{
				directTransport: &directModeTransport{},
			},
		},
	}

	pk := NewPartitionKey()
	opts := &QueryOptions{}

	mode := container.selectQueryExecutionMode(pk, opts)
	require.Equal(t, queryModeGateway, mode)
}

type mockQueryEngine struct{}

func (m *mockQueryEngine) SupportedFeatures() string {
	return "None"
}

func (m *mockQueryEngine) CreateQueryPipeline(_, _, _ string) (queryengine.QueryPipeline, error) {
	return nil, nil
}

func (m *mockQueryEngine) CreateReadManyPipeline(_ []queryengine.ItemIdentity, _ string, _ string, _ uint8, _ []string) (queryengine.QueryPipeline, error) {
	return nil, nil
}
