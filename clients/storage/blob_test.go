package storage

import (
	"crypto/rand"
	"errors"
	"log"
	"os"
	"reflect"
	"sort"
	"sync"
	"testing"
)

const testContainerPrefix = "zzzztest-"

func TestListContainers(t *testing.T) {
	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	out, err := cli.ListContainers(ListContainersParameters{})
	if err != nil {
		t.Error(err)
	}
}

func TestListContainersPagination(t *testing.T) {
	cli, err := getClient()
	if err != nil {
		t.Fatal(err)
	}

	err = deleteTestContainers(cli)
	if err != nil {
		t.Fatal(err)
	}

	// Create test containers
	created := []string{}
	const n = 10
	for i := 0; i < n; i++ {
		name := randContainer(n)
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
	const pageSize = 2
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

		containers := resp.Containers.Containers

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

func deleteTestContainers(cli *BlobStorageClient) error {
	for {
		resp, err := cli.ListContainers(ListContainersParameters{Prefix: testContainerPrefix})
		if err != nil {
			return err
		}
		if len(resp.Containers.Containers) == 0 {
			break
		}
		for _, c := range resp.Containers.Containers {
			_, err := cli.DeleteContainer(c.Name)
			if err != nil {
				return err
			}
			log.Printf("Cleaning up leftover test container %s", c.Name)
		}
	}
	return nil
}

func TestCreateGetDeleteContainer(t *testing.T) {
	cnt := randString(10)

	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypePrivate)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.GetContainer(cnt)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Error(err)
	}
}

func TestPutBlockBlob(t *testing.T) {
	cnt := randString(5)
	blob := randString(10)
	body := randString(1024 * 4)

	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	_, err = cli.CreateContainer(cnt, ContainerAccessTypePrivate)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.PutBlob(cnt, blob, BlobTypeBlock, []byte(body))
	if err != nil {
		t.Error(err)
	}

	_, err = cli.DeleteBlob(cnt, blob)
	if err != nil {
		t.Error(err)
	}

	_, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Error(err)
	}
}

func getClient() (*BlobStorageClient, error) {
	name := os.Getenv("ACCOUNT_NAME")
	if name == "" {
		return nil, errors.New("ACCOUNT_NAME not set")
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

func randContainer(n int) string {
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
