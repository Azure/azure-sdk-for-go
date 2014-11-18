package storage

import (
	"crypto/rand"
	"errors"
	"io/ioutil"
	"os"
	"testing"
)

func TestListContainers(t *testing.T) {
	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	resp, err := cli.ListContainers()
	if err != nil {
		t.Error(err)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))
}

func TestCreateGetDeleteContainer(t *testing.T) {
	cnt := randString(10)

	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	resp1, err := cli.CreateContainer(cnt)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Create: %v", resp1.Status)

	resp2, err := cli.GetContainer(cnt)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Get: %v", resp2.Status)

	resp3, err := cli.DeleteContainer(cnt)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Delete: %v", resp3.Status)
}

func TestPutBlockBlob(t *testing.T) {
	cnt := randString(5)
	blob := randString(10)
	body := randString(1024 * 4)

	cli, err := getClient()
	if err != nil {
		t.Error(err)
	}

	resp, err := cli.CreateContainer(cnt)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Create container: %v", resp.Status)

	resp, err = cli.PutBlob(cnt, blob, BlobTypeBlock, []byte(body))
	if err != nil {
		t.Error(err)
	}
	t.Logf("Put blob: %v", resp.Status)

	resp, err = cli.DeleteBlob(cnt, blob)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Delete blob: %v", resp.Status)

	resp, err = cli.DeleteContainer(cnt)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Delete container: %v", resp.Status)
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
