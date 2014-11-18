package storage

import (
	"crypto/rand"
	"errors"
	"os"
	"testing"
)

func TestListContainers(t *testing.T) {
	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	out, err := cli.ListContainers()
	if err != nil {
		t.Error(err)
	}
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

func randString(n int) string {
	const alphanum = "0123456789abcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}
