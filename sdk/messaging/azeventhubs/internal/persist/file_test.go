// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package persist

import (
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestFilePersister_Read(t *testing.T) {
	namespace := "namespace"
	name := "name"
	group := "$Default"
	partitionID := "0"
	dir := path.Join(os.TempDir(), RandomName("read", 4))
	persister, err := NewFilePersister(dir)
	assert.Nil(t, err)
	ckp, err := persister.Read(namespace, name, group, partitionID)
	assert.NotNil(t, err)
	assert.Equal(t, NewCheckpointFromStartOfStream(), ckp)
}

func TestFilePersister_Write(t *testing.T) {
	namespace := "namespace"
	name := "name"
	group := "$Default"
	partitionID := "0"
	dir := path.Join(os.TempDir(), RandomName("write", 4))
	persister, err := NewFilePersister(dir)
	assert.Nil(t, err)
	ckp := NewCheckpoint("120", 22, time.Now())
	err = persister.Write(namespace, name, group, partitionID, ckp)
	assert.Nil(t, err)
	ckp2, err := persister.Read(namespace, name, group, partitionID)
	assert.Nil(t, err)
	assert.Equal(t, ckp.Offset, ckp2.Offset)
	assert.Equal(t, ckp.SequenceNumber, ckp2.SequenceNumber)
}

// RandomName generates a random Event Hub name
func RandomName(prefix string, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + "-" + string(b)
}
