package storage

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"testing"
)

const testContainerPrefix = "zzzztest-"

func TestContainerExists(t *testing.T) {
	cnt := randContainer()

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	ok, err := cli.ContainerExists(cnt)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("Non-existing container returned as existing: %s", cnt)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypeBlob)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteContainer(cnt)

	ok, err = cli.ContainerExists(cnt)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("Existing container returned as non-existing: %s", cnt)
	}
}

func TestListContainersPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skip in short mode")
	}

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	err = deleteTestContainers(cli)
	if err != nil {
		t.Fatal(err)
	}

	const n = 5
	const pageSize = 2

	// Create test containers
	created := []string{}
	for i := 0; i < n; i++ {
		name := randContainer()
		_, err := cli.CreateContainer(name, ContainerAccessTypePrivate)
		if err != nil {
			t.Fatalf("Error creating test container: %s", err)
		}
		created = append(created, name)
	}
	sort.Strings(created)

	// Defer test container deletions
	defer func() {
		var wg sync.WaitGroup
		for _, cnt := range created {
			wg.Add(1)
			go func(name string) {
				_, err := cli.DeleteContainer(name)
				if err != nil {
					t.Logf("Error while deleting test container: %s", err)
				}
				wg.Done()
			}(cnt)
		}
		wg.Wait()
	}()

	// Paginate results
	seen := []string{}
	marker := ""
	for {
		resp, err := cli.ListContainers(ListContainersParameters{
			Prefix:     testContainerPrefix,
			MaxResults: pageSize,
			Marker:     marker})

		if err != nil {
			t.Fatal(err)
		}

		containers := resp.Containers

		if len(containers) > pageSize {
			t.Fatalf("Got a bigger page. Expected: %d, got: %d", pageSize, len(containers))
		}

		for _, c := range containers {
			seen = append(seen, c.Name)
		}

		marker = resp.NextMarker
		if marker == "" || len(containers) == 0 {
			break
		}
	}

	// Compare
	if !reflect.DeepEqual(created, seen) {
		t.Fatal("Wrong pagination results:\nExpected:\t\t%s\nGot:\t\t%s", created, seen)
	}
}

func TestCreateDeleteContainer(t *testing.T) {
	cnt := randContainer()

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypePrivate)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteContainer(cnt)

	_, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBlobExists(t *testing.T) {
	cnt := randContainer()
	blob := fmt.Sprintf("%s/%s", randString(5), randString(20))

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypeBlob)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteContainer(cnt)
	err = cli.PutBlockBlob(cnt, blob, strings.NewReader("Hello!"))
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteBlob(cnt, blob)

	ok, err := cli.BlobExists(cnt, blob+".foo")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("Non-existing blob returned as existing: %s/%s", cnt, blob)
	}

	ok, err = cli.BlobExists(cnt, blob)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("Existing blob returned as non-existing: %s/%s", cnt, blob)
	}
}

func TestListBlobsPagination(t *testing.T) {
	if testing.Short() {
		t.Skip("skip in short mode")
	}

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	cnt := randContainer()
	_, err = cli.CreateContainer(cnt, ContainerAccessTypePrivate)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteContainer(cnt)

	blobs := []string{}
	const n = 5
	const pageSize = 2
	for i := 0; i < n; i++ {
		name := randString(20)
		err := cli.PutBlockBlob(cnt, name, strings.NewReader("Hello, world!"))
		if err != nil {
			t.Fatal(err)
		}
		blobs = append(blobs, name)
	}
	sort.Strings(blobs)

	// Paginate
	seen := []string{}
	marker := ""
	for {
		resp, err := cli.ListBlobs(cnt, ListBlobsParameters{
			MaxResults: pageSize,
			Marker:     marker})
		if err != nil {
			t.Fatal(err)
		}

		for _, v := range resp.Blobs {
			seen = append(seen, v.Name)
		}

		marker = resp.NextMarker
		if marker == "" || len(resp.Blobs) == 0 {
			break
		}
	}

	// Compare
	if !reflect.DeepEqual(blobs, seen) {
		t.Fatalf("Got wrong list of blobs. Expected: %s, Got: %s", blobs, seen)
	}

	_, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPutSingleBlockBlob(t *testing.T) {
	cnt := randContainer()
	blob := randString(20)
	body := []byte(randString(1024 * 4))

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypeBlob)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteContainer(cnt)

	err = cli.PutBlockBlob(cnt, blob, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteBlob(cnt, blob)

	resp, err := cli.GetBlob(cnt, blob)
	if err != nil {
		t.Fatal(err)
	}

	// Verify contents
	if !reflect.DeepEqual(body, resp.body) {
		t.Fatalf("Wrong blob contents.\nExpected: %d bytes, Got: %d byes", len(body), len(resp.body))
	}

	_, err = cli.DeleteBlob(cnt, blob)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPutMultiBlockBlob(t *testing.T) {
	if testing.Short() {
		t.Skip("skip in short mode ")
	}

	var (
		cnt       = randContainer()
		blob      = randString(20)
		blockSize = 32 * 1024                                     // 32 KB
		body      = []byte(randString(blockSize*2 + blockSize/2)) // 3 blocks
	)

	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypeBlob)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteContainer(cnt)

	err = cli.putBlockBlob(cnt, blob, bytes.NewReader(body), blockSize)
	if err != nil {
		t.Fatal(err)
	}
	defer cli.DeleteBlob(cnt, blob)

	resp, err := cli.GetBlob(cnt, blob)
	if err != nil {
		t.Fatal(err)
	}

	// Verify contents
	if !reflect.DeepEqual(body, resp.body) {
		t.Fatalf("Wrong blob contents.\nExpected: %d bytes, Got: %d byes", len(body), len(resp.body))
	}

	_, err = cli.DeleteBlob(cnt, blob)
	if err != nil {
		t.Fatal(err)
	}

	_, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteTestContainers(cli *BlobStorageClient) error {
	for {
		resp, err := cli.ListContainers(ListContainersParameters{Prefix: testContainerPrefix})
		if err != nil {
			return err
		}
		if len(resp.Containers) == 0 {
			break
		}
		for _, c := range resp.Containers {
			_, err := cli.DeleteContainer(c.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getClient() (*BlobStorageClient, error) {
	name := os.Getenv("ACCOUNT_NAME")
	if name == "" {
		return nil, errors.New("ACCOUNT_NAME not set, need an empty storage account to test")
	}
	key := os.Getenv("ACCOUNT_KEY")
	if key == "" {
		return nil, errors.New("ACCOUNT_KEY not set")
	}
	cli, err := NewBasicClient(name, key)
	if err != nil {
		return nil, err
	}
	return cli.GetBlobService(), nil
}

func randContainer() string {
	return testContainerPrefix + randString(32-len(testContainerPrefix))
}

func randString(n int) string {
	if n <= 0 {
		panic("negative number")
	}
	const alphanum = "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
