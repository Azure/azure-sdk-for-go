package storage

import (
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

	resp, err := cli.GetBlobService().ListContainers()
	if err != nil {
		t.Error(err)
	} else if resp == nil {
		t.Error("Got nil response")
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(bytes))
}

func getClient() (*StorageClient, error) {
	name := os.Getenv("ACCOUNT_NAME")
	if name == "" {
		return nil, errors.New("ACCOUNT_NAME not set")
	}
	key := os.Getenv("ACCOUNT_KEY")
	if key == "" {
		return nil, errors.New("ACCOUNT_KEY not set")
	}
	return NewBasicClient(name, key)
}
