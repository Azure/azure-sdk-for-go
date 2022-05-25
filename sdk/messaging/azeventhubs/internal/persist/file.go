// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package persist

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path"
	"strings"
	"sync"
)

type (
	// FilePersister implements CheckpointPersister for saving to the file system
	FilePersister struct {
		directory string
		mu        sync.Mutex
	}
)

// NewFilePersister creates a FilePersister for saving to a given directory
func NewFilePersister(directory string) (*FilePersister, error) {
	err := os.MkdirAll(directory, 0777)
	return &FilePersister{
		directory: directory,
	}, err
}

func (fp *FilePersister) Write(namespace, name, consumerGroup, partitionID string, checkpoint Checkpoint) error {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	key := getFilePath(namespace, name, consumerGroup, partitionID)
	filePath := path.Join(fp.directory, key)
	bits, err := json.Marshal(checkpoint)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = file.Write(bits)
	if err != nil {
		return err
	}

	return file.Close()
}

func (fp *FilePersister) Read(namespace, name, consumerGroup, partitionID string) (Checkpoint, error) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	key := getFilePath(namespace, name, consumerGroup, partitionID)
	filePath := path.Join(fp.directory, key)

	f, err := os.Open(filePath)
	if err != nil {
		return NewCheckpointFromStartOfStream(), err
	}

	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, f)
	if err != nil {
		return NewCheckpointFromStartOfStream(), err
	}

	var checkpoint Checkpoint
	err = json.Unmarshal(buf.Bytes(), &checkpoint)
	return checkpoint, err
}

func getFilePath(namespace, name, consumerGroup, partitionID string) string {
	key := strings.Join([]string{namespace, name, consumerGroup, partitionID}, "_")
	return strings.Replace(key, "$", "", -1)
}
