// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package azcosmos

// If you want to run and see the output of these tests
// 1. cd sdk/data/azcosmos
// 2. go test -v -run TestCloudChangeFeed_AIMHeader (or other test name)
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

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
		MaxItemCount: 1,
	}
	resp, err := container.getChangeFeed(context.Background(), nil, nil, options)
	if err != nil {
		t.Fatalf("getChangeFeed failed: %v", err)
	}

	fmt.Printf("TEST AIM Header: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	// fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
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

	// Use an ETag value of "5"
	continuation := "5"
	maxItemCount := int32(1)
	options := &ChangeFeedOptions{
		Continuation: &continuation,
		MaxItemCount: maxItemCount,
	}
	fmt.Printf("DEBUG: Requesting with MaxItemCount: %d\n", maxItemCount)

	resp, err := container.GetChangeFeedContainer(context.Background(), options)
	if err != nil {
		t.Fatalf("getChangeFeed failed: %v", err)
	}
	fmt.Printf("If-None-Match Header Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	// fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
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
	modifiedSince := time.Date(2025, 6, 14, 0, 0, 0, 0, time.UTC)
	fmt.Printf("Testing with If-Modified-Since: %s\n", modifiedSince.UTC().Format(time.RFC1123))

	options := &ChangeFeedOptions{
		ChangeFeedStartFrom: &modifiedSince, // Changed from IfModifiedSince
		MaxItemCount:        5,
	}

	// Note: toHeaders now requires partitionKeyRanges parameter
	headers := options.toHeaders(nil)
	if headers != nil {
		for k, v := range *headers {
			fmt.Printf("Header: %s: %s\n", k, v)
		}
	}
	resp, err := container.getChangeFeed(context.Background(), nil, nil, options)
	if err != nil {
		t.Fatalf("getChangeFeed failed: %v", err)
	}
	fmt.Printf("If-Modified-Since Header Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	// fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
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

	resp, err := container.getChangeFeed(context.Background(), nil, nil, options)
	if err != nil {
		t.Fatalf("getChangeFeed failed: %v", err)
	}

	fmt.Printf("Partition Key Header Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	// fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp.LSN)

	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}

// Testing Change Feed with FeedRange
func TestCloudChangeFeed_FeedRange(t *testing.T) {
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

	// First, get the partition key ranges to find a valid FeedRange
	pkrResp, err := container.getPartitionKeyRanges(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to get partition key ranges: %v", err)
	}

	if len(pkrResp.PartitionKeyRanges) == 0 {
		t.Skip("No partition key ranges found")
	}

	// Use the first partition key range as our FeedRange
	firstRange := pkrResp.PartitionKeyRanges[0]
	feedRange := &FeedRange{
		MinInclusive: firstRange.MinInclusive,
		MaxExclusive: firstRange.MaxExclusive,
	}

	options := &ChangeFeedOptions{
		MaxItemCount: 5,
		FeedRange:    feedRange,
	}

	resp, err := container.GetChangeFeedForEPKRange(context.Background(), feedRange, options)
	if err != nil {
		t.Fatalf("getChangeFeed failed: %v", err)
	}

	fmt.Printf("FeedRange Test: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)
	fmt.Printf("ETag header: %s\n", resp.ETag)
	// fmt.Printf("x-ms-continuation header: %s\n", resp.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp.LSN)
	fmt.Printf("Testing with FeedRange: MinInclusive=%s, MaxExclusive=%s\n", feedRange.MinInclusive, feedRange.MaxExclusive)

	for i, doc := range resp.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}
}

// Testing Composite Continuation Token functionality
func TestCloudChangeFeed_CompositeContinuationToken(t *testing.T) {
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

	// First, get the partition key ranges
	pkrResp, err := container.getPartitionKeyRanges(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to get partition key ranges: %v", err)
	}

	if len(pkrResp.PartitionKeyRanges) == 0 {
		t.Skip("No partition key ranges found")
	}

	// Use the first partition key range as our initial FeedRange
	firstRange := pkrResp.PartitionKeyRanges[0]
	feedRange := &FeedRange{
		MinInclusive: firstRange.MinInclusive,
		MaxExclusive: firstRange.MaxExclusive,
	}

	fmt.Printf("\n=== First Request with FeedRange ===\n")
	fmt.Printf("Using FeedRange: MinInclusive=%s, MaxExclusive=%s\n", feedRange.MinInclusive, feedRange.MaxExclusive)

	// First request with a feed range
	options := &ChangeFeedOptions{
		MaxItemCount: 1,
		FeedRange:    feedRange,
	}

	//resp1, err := container.getChangeFeed(context.Background(), nil, nil, options)
	resp1, err := container.GetChangeFeedForEPKRange(context.Background(), feedRange, options)
	if err != nil {
		t.Fatalf("First getChangeFeed failed: %v", err)
	}

	fmt.Printf("First response: ResourceID: %s, Documents count: %d\n", resp1.ResourceID, resp1.Count)
	// fmt.Printf("Documents from first response:\n")
	// for i, doc := range resp1.Documents {
	// 	fmt.Printf("Doc %d: %s\n", i, string(doc))
	// }
	fmt.Printf("ETag: %s\n", resp1.ETag)

	// Get the composite continuation token from the first response
	compositeToken, err := resp1.getCompositeContinuationToken()
	if err != nil {
		t.Fatalf("Failed to get composite continuation token: %v", err)
	}

	if compositeToken == "" {
		t.Skip("No composite continuation token generated (possibly no data or no valid continuation)")
	}

	fmt.Printf("\n=== Composite Continuation Token ===\n")
	fmt.Printf("Token: %s\n", compositeToken)

	// Parse and display the composite token structure
	var parsedToken compositeContinuationToken
	if err := json.Unmarshal([]byte(compositeToken), &parsedToken); err == nil {
		fmt.Printf("Parsed Token - ResourceID: %s\n", parsedToken.ResourceID)
		for i, cont := range parsedToken.Continuation {
			fmt.Printf("Continuation[%d]: Min=%s, Max=%s, Token=%s\n",
				i, cont.MinInclusive, cont.MaxExclusive, *cont.ContinuationToken)
		}
	}

	fmt.Printf("\n=== Second Request with Composite Token ===\n")

	// Second request using the composite continuation token
	options2 := &ChangeFeedOptions{
		MaxItemCount: 1,
		Continuation: &compositeToken,
		// Note: FeedRange is not set - it should be extracted from the composite token
	}

	resp2, err := container.GetChangeFeedForEPKRange(context.Background(), feedRange, options2)
	if err != nil {
		t.Fatalf("Second getChangeFeed failed: %v", err)
	}

	fmt.Printf("Second response: ResourceID: %s, Documents count: %d\n", resp2.ResourceID, resp2.Count)
	fmt.Printf("ETag: %s\n", resp2.ETag)
	fmt.Printf("Continuation Token: %s\n", resp2.ETag)

	// Verify that the FeedRange was set from the composite token
	if options2.FeedRange == nil {
		t.Error("Expected FeedRange to be set from composite token")
	} else {
		fmt.Printf("FeedRange extracted from composite token: Min=%s, Max=%s\n",
			options2.FeedRange.MinInclusive, options2.FeedRange.MaxExclusive)

		// Verify it matches the original feed range
		if options2.FeedRange.MinInclusive != feedRange.MinInclusive ||
			options2.FeedRange.MaxExclusive != feedRange.MaxExclusive {
			t.Errorf("FeedRange mismatch: expected Min=%s,Max=%s, got Min=%s,Max=%s",
				feedRange.MinInclusive, feedRange.MaxExclusive,
				options2.FeedRange.MinInclusive, options2.FeedRange.MaxExclusive)
		}
	}

	// Display documents from second response
	fmt.Printf("\nDocuments from second response:\n")
	for i, doc := range resp2.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}

	// Test continued pagination with composite token
	if resp2.Count > 0 {
		compositeToken2, err := resp2.getCompositeContinuationToken()
		if err != nil {
			t.Fatalf("Failed to get second composite continuation token: %v", err)
		}

		if compositeToken2 != "" {
			fmt.Printf("\n=== Third Request with Updated Composite Token ===\n")

			options3 := &ChangeFeedOptions{
				MaxItemCount: 1,
				Continuation: &compositeToken2,
			}

			resp3, err := container.GetChangeFeedForEPKRange(context.Background(), feedRange, options3)
			if err != nil {
				t.Fatalf("Third getChangeFeed failed: %v", err)
			}

			fmt.Printf("Third response: ResourceID: %s, Documents count: %d\n", resp3.ResourceID, resp3.Count)
			fmt.Printf("Continuation demonstrates token chaining works correctly\n")
		}
	}
}

// Testing Composite Continuation Token with multiple partitions
func TestCloudChangeFeed_CompositeContinuationTokenMultipleRanges(t *testing.T) {
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

	// Get all partition key ranges
	pkrResp, err := container.getPartitionKeyRanges(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to get partition key ranges: %v", err)
	}

	if len(pkrResp.PartitionKeyRanges) < 2 {
		t.Skip("Need at least 2 partition key ranges for this test")
	}

	fmt.Printf("\n=== Testing with Wide FeedRange (Multiple Partitions) ===\n")
	fmt.Printf("Total partition key ranges: %d\n", len(pkrResp.PartitionKeyRanges))

	// Create a feed range that spans multiple partitions
	// Using the full range from the first to the last partition
	wideFeedRange := &FeedRange{
		MinInclusive: pkrResp.PartitionKeyRanges[0].MinInclusive,
		MaxExclusive: pkrResp.PartitionKeyRanges[len(pkrResp.PartitionKeyRanges)-1].MaxExclusive,
	}

	fmt.Printf("Wide FeedRange: Min=%s, Max=%s\n", wideFeedRange.MinInclusive, wideFeedRange.MaxExclusive)

	// First request with wide feed range
	options := &ChangeFeedOptions{
		MaxItemCount: 5,
		FeedRange:    wideFeedRange,
	}

	resp, err := container.getChangeFeed(context.Background(), nil, nil, options)
	if err != nil {
		t.Fatalf("getChangeFeed with wide range failed: %v", err)
	}

	fmt.Printf("Response: ResourceID: %s, Documents count: %d\n", resp.ResourceID, resp.Count)

	// Get composite token
	compositeToken, err := resp.getCompositeContinuationToken()
	if err != nil {
		t.Fatalf("Failed to get composite continuation token: %v", err)
	}

	if compositeToken != "" {
		fmt.Printf("\nComposite token generated for wide range\n")

		// Use the composite token in a follow-up request
		options2 := &ChangeFeedOptions{
			MaxItemCount: 5,
			Continuation: &compositeToken,
		}

		resp2, err := container.getChangeFeed(context.Background(), nil, nil, options2)
		if err != nil {
			t.Fatalf("Second getChangeFeed with composite token failed: %v", err)
		}

		fmt.Printf("Follow-up response: Documents count: %d\n", resp2.Count)
		fmt.Printf("Successfully used composite token for wide range continuation\n")
	}
}

// Testing Composite Continuation Token with ETag only (no FeedRange)
func TestCloudChangeFeed_CompositeContinuationTokenWithETag(t *testing.T) {
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

	// First, get the partition key ranges
	pkrResp, err := container.getPartitionKeyRanges(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to get partition key ranges: %v", err)
	}

	if len(pkrResp.PartitionKeyRanges) == 0 {
		t.Skip("No partition key ranges found")
	}

	// Use the first partition key range as our initial FeedRange
	firstRange := pkrResp.PartitionKeyRanges[0]
	feedRange := &FeedRange{
		MinInclusive: firstRange.MinInclusive,
		MaxExclusive: firstRange.MaxExclusive,
	}

	fmt.Printf("\n=== First Request with ETag '5' ===\n")

	// First request with just an ETag (no FeedRange)
	etag := "5"
	options := &ChangeFeedOptions{
		MaxItemCount: 1,
		Continuation: &etag,
		FeedRange:    feedRange,
	}

	resp1, err := container.getChangeFeed(context.Background(), nil, nil, options)
	if err != nil {
		t.Fatalf("First getChangeFeed failed: %v", err)
	}

	fmt.Printf("First response: ResourceID: %s, Documents count: %d\n", resp1.ResourceID, resp1.Count)
	fmt.Printf("ETag from response: %s\n", resp1.ETag)
	fmt.Printf("Continuation Token: %s\n", resp1.ContinuationToken)
	fmt.Printf("LSN: %s\n", resp1.LSN)

	// Display documents from first response
	fmt.Printf("\nDocuments from first response:\n")
	for i, doc := range resp1.Documents {
		fmt.Printf("Doc %d: %s\n", i, string(doc))
	}

	// Try to get the composite continuation token from the first response
	fmt.Printf("\n=== Attempting to Create Composite Continuation Token ===\n")

	compositeToken, err := resp1.getCompositeContinuationToken()
	if err != nil {
		fmt.Printf("Error getting composite continuation token: %v\n", err)
	} else if compositeToken == "" {
		fmt.Printf("No composite continuation token generated\n")
		fmt.Printf("This is expected when there's no FeedRange in the continuation token\n")
	} else {
		fmt.Printf("Composite Token: %s\n", compositeToken)

		// Parse and display the composite token structure
		var parsedToken compositeContinuationToken
		if err := json.Unmarshal([]byte(compositeToken), &parsedToken); err == nil {
			fmt.Printf("Parsed Token - ResourceID: %s\n", parsedToken.ResourceID)
			for i, cont := range parsedToken.Continuation {
				fmt.Printf("Continuation[%d]: Min=%s, Max=%s, Token=%v\n",
					i, cont.MinInclusive, cont.MaxExclusive, cont.ContinuationToken)
			}
		}

		// Try using the composite token in a second request
		fmt.Printf("\n=== Second Request with Composite Token ===\n")

		options2 := &ChangeFeedOptions{
			MaxItemCount: 2,
			Continuation: &compositeToken,
			FeedRange:    feedRange,
		}

		resp2, err := container.getChangeFeed(context.Background(), nil, nil, options2)
		if err != nil {
			t.Fatalf("Second getChangeFeed failed: %v", err)
		}

		fmt.Printf("Second response: ResourceID: %s, Documents count: %d\n", resp2.ResourceID, resp2.Count)
		fmt.Printf("ETag: %s\n", resp2.ETag)
		fmt.Printf("Continuation Token: %s\n", resp2.ContinuationToken)
	}

	// Also test what happens when we manually check the continuation token structure
	fmt.Printf("\n=== Analyzing Continuation Token Structure ===\n")
	if resp1.ContinuationToken != "" {
		var contToken struct {
			Token *string `json:"token"`
			Range *struct {
				Min string `json:"min"`
				Max string `json:"max"`
			} `json:"range"`
		}

		if err := json.Unmarshal([]byte(resp1.ContinuationToken), &contToken); err == nil {
			fmt.Printf("Continuation token parsed successfully:\n")
			if contToken.Token != nil {
				fmt.Printf("  Token: %s\n", *contToken.Token)
			} else {
				fmt.Printf("  Token: null\n")
			}
			if contToken.Range != nil {
				fmt.Printf("  Range: Min=%s, Max=%s\n", contToken.Range.Min, contToken.Range.Max)
			} else {
				fmt.Printf("  Range: null\n")
			}
		} else {
			fmt.Printf("Failed to parse continuation token: %v\n", err)
		}
	}
}
