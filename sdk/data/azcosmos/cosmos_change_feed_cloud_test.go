package azcosmos

// If you want to run and see the output of these tests
// 1. cd sdk/data/azcosmos
// 2. go test -v -run TestCloudChangeFeed_AIMHeader (or other test name)
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

// Testing With only Incremental Feed Header and Item Count of 1
func TestCloudChangeFeed_AIMHeader(t *testing.T) {
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

	options := &ChangeFeedOptions{
		AIM:          "Incremental Feed",
		MaxItemCount: 1,
	}
	resp, err := container.GetChangeFeed(context.Background(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed failed: %v", err)
	}

	fmt.Printf("TEST AIM Header: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp.LSN)

	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}

// Testing If-None Match Header with ETag "5"
func TestCloudChangeFeed_IfNoneMatchHeader(t *testing.T) {
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

	// Use an ETag value of "15" (must be quoted)
	etag := azcore.ETag(`"5"`)

	options := &ChangeFeedOptions{
		IfNoneMatch:  &etag,
		MaxItemCount: 1,
	}
	resp, err := container.GetChangeFeed(context.Background(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed failed: %v", err)
	}
	fmt.Printf("If-None-Match Header Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp.LSN)

	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}

// Testing Date Modified Since Header w/ Jun 30th date
// Note since nothing has been modified, should get a 304 Not Modified response
func TestCloudChangeFeed_IfModifiedSinceHeader(t *testing.T) {
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

	// Use the fixed date: Mon, 30 Jun 2025 00:00:00 GMT
	modifiedSince := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Testing with If-Modified-Since: %s\n", modifiedSince.UTC().Format(time.RFC1123))
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}

	options := &ChangeFeedOptions{
		AIM:             "Incremental Feed",
		IfModifiedSince: &modifiedSince,
		MaxItemCount:    5,
	}

	headers := options.toHeaders()
	if headers != nil {
		for k, v := range *headers {
			fmt.Printf("Header: %s: %s\n", k, v)
		}
	}
	resp, err := container.GetChangeFeed(context.Background(), options)
	if err != nil {
		t.Fatalf("GetChangeFeed failed: %v", err)
	}
	fmt.Printf("If-Modified-Since Header Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp.LSN)

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

	fmt.Printf("Partition Key Header Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp.LSN)

	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}
