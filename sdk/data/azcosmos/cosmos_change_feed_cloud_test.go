package azcosmos

// If you want to run and see the output of these tests
// 1. cd sdk/data/azcosmos
// 2. go test -v -run TestIntegrationGetChangeFeed (or other test name)
import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/joho/godotenv"
)

func loadEnv(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}
}

func getEnvOrSkip(t *testing.T, key string) string {
	val := os.Getenv(key)
	if val == "" {
		t.Skipf("Missing required env var: %s", key)
	}
	return val
}

// Testing Change Feed with date modified since
func TestIntegrationGetChangeFeedDate(t *testing.T) {
	loadEnv(t)
	endpoint := getEnvOrSkip(t, "COSMOS_ENDPOINT")
	key := getEnvOrSkip(t, "COSMOS_KEY")
	dbName := getEnvOrSkip(t, "COSMOS_DATABASE")
	containerName := getEnvOrSkip(t, "COSMOS_CONTAINER")

	if endpoint == "" || key == "" || dbName == "" || containerName == "" {
		t.Skip("COSMOS_ENDPOINT, COSMOS_KEY, COSMOS_DATABASE, and COSMOS_CONTAINER must be set for integration test")
	}

	cred, err := NewKeyCredential(key)
	if err != nil {
		t.Fatalf("Failed to create KeyCredential: %v", err)
	}
	client, err := NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		t.Fatalf("Failed to create Cosmos client: %v", err)
	}

	db, err := client.NewDatabase(dbName)
	if err != nil {
		t.Fatalf("Failed to get database: %v", err)
	}

	container, err := db.NewContainer(containerName)
	if err != nil {
		t.Fatalf("Failed to get container: %v", err)
	}

	// Optionally, insert a document here for testing

	options := &ChangeFeedOptions{
		MaxItemCount: 5,
		// You can set IfModifiedSince, PartitionKey, etc. here
		IfModifiedSince: func() *time.Time { t := time.Now().Add(-24 * time.Hour); return &t }(),
	}

	resp, err := container.GetChangeFeed(context.Background(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed failed: %v", err)
	}

	fmt.Printf("ChangeFeed ResourceID: %s\n", resp.ResourceID)
	fmt.Printf("Documents count: %d\n", resp.Count)
	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}

// Testing Change Feed with Partition Key
func TestIntegrationGetChangeFeedPartitionKey(t *testing.T) {
	loadEnv(t)
	endpoint := getEnvOrSkip(t, "COSMOS_ENDPOINT")
	key := getEnvOrSkip(t, "COSMOS_KEY")
	dbName := getEnvOrSkip(t, "COSMOS_DATABASE")
	containerName := getEnvOrSkip(t, "COSMOS_CONTAINER")

	if endpoint == "" || key == "" || dbName == "" || containerName == "" {
		t.Skip("COSMOS_ENDPOINT, COSMOS_KEY, COSMOS_DATABASE, and COSMOS_CONTAINER must be set for integration test")
	}

	cred, err := NewKeyCredential(key)
	if err != nil {
		t.Fatalf("Failed to create KeyCredential: %v", err)
	}
	client, err := NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		t.Fatalf("Failed to create Cosmos client: %v", err)
	}

	db, err := client.NewDatabase(dbName)
	if err != nil {
		t.Fatalf("Failed to get database: %v", err)
	}

	container, err := db.NewContainer(containerName)
	if err != nil {
		t.Fatalf("Failed to get container: %v", err)
	}

	partitionKey := NewPartitionKeyString("11111")
	options := &ChangeFeedOptions{
		MaxItemCount: 5,
		PartitionKey: &partitionKey,
	}

	resp, err := container.GetChangeFeed(context.Background(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed failed: %v", err)
	}

	fmt.Printf("ChangeFeed ResourceID: %s\n", resp.ResourceID)
	fmt.Printf("Documents count: %d\n", resp.Count)
	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}

// Trying to grab all the changes from a container with paging
func TestIntegrationGetChangeFeedPaging(t *testing.T) {
	loadEnv(t)
	endpoint := getEnvOrSkip(t, "COSMOS_ENDPOINT")
	key := getEnvOrSkip(t, "COSMOS_KEY")
	dbName := getEnvOrSkip(t, "COSMOS_DATABASE")
	containerName := getEnvOrSkip(t, "COSMOS_CONTAINER")

	cred, err := NewKeyCredential(key)
	if err != nil {
		t.Fatalf("Failed to create KeyCredential: %v", err)
	}
	client, err := NewClientWithKey(endpoint, cred, nil)
	if err != nil {
		t.Fatalf("Failed to create Cosmos client: %v", err)
	}
	db, err := client.NewDatabase(dbName)
	if err != nil {
		t.Fatalf("Failed to get database: %v", err)
	}
	container, err := db.NewContainer(containerName)
	if err != nil {
		t.Fatalf("Failed to get container: %v", err)
	}

	var continuation *azcore.ETag
	page := 0
	for {
		options := &ChangeFeedOptions{
			MaxItemCount: 5, // or any page size you want
			IfNoneMatch:  continuation,
		}
		resp, err := container.GetChangeFeed(context.Background(), options)
		if err != nil {
			t.Fatalf("GetChangeFeed failed: %v", err)
		}
		fmt.Printf("Page %d: ResourceID: %s, Documents count: %d\n", page, resp.ResourceID, resp.Count)
		for i, doc := range resp.Documents {
			fmt.Printf("Page %d, Doc %d: %s\n", page, i, string(doc))
		}
		etag := resp.Response.ETag
		if resp.Count == 0 || etag == "" || (continuation != nil && *continuation == etag) {
			break // No more results or no new continuation token
		}
		continuation = &etag
		page++
	}
}
