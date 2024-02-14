// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionContainerSinglePartition(t *testing.T) {
	sessionContainer := newSessionContainer()

	expectedToken := "0:1#1234"
	sessionContainer.SetSessionToken("dbs/test/colls/test", "x1cRAMM0cZs=", expectedToken)

	actualToken := sessionContainer.GetSessionToken("dbs/test/colls/test")
	assert.Equal(t, expectedToken, actualToken)
}

func TestSessionContainerMultiplePartitions(t *testing.T) {
	sessionContainer := newSessionContainer()

	expectedToken := "0:1#1234,1:2#5678"
	sessionContainer.SetSessionToken("dbs/test/colls/test", "x1cRAMM0cZs=", expectedToken)

	actualToken := sessionContainer.GetSessionToken("dbs/test/colls/test")
	assert.Equal(t, expectedToken, actualToken)
}

func TestSessionContainerInvalidToken(t *testing.T) {
	sessionContainer := newSessionContainer()

	sessionContainer.SetSessionToken("dbs/test/colls/test", "x1cRAMM0cZs=", "invalid")

	actualToken := sessionContainer.GetSessionToken("dbs/test/colls/test")
	assert.Equal(t, "", actualToken)
}

func TestSessionContainerMultipleContainers(t *testing.T) {
	sessionContainer := newSessionContainer()

	expectedTokenColl1 := "0:1#1234"
	expectedTokenColl2 := "1:2#5678"
	sessionContainer.SetSessionToken("dbs/test/colls/test1", "x1cRAMM0cZs=", expectedTokenColl1)
	sessionContainer.SetSessionToken("dbs/test/colls/test2", "D80jAP5wbP4=", expectedTokenColl2)

	actualTokenColl1 := sessionContainer.GetSessionToken("dbs/test/colls/test1")
	assert.Equal(t, expectedTokenColl1, actualTokenColl1)

	actualTokenColl2 := sessionContainer.GetSessionToken("dbs/test/colls/test2")
	assert.Equal(t, expectedTokenColl2, actualTokenColl2)
}

func TestSessionContainerChangedRid(t *testing.T) {
	sessionContainer := newSessionContainer()

	expectedToken := "0:2#5678"
	sessionContainer.SetSessionToken("dbs/test/colls/test", "x1cRAMM0cZs=", "0:1#1234")
	sessionContainer.SetSessionToken("dbs/test/colls/test", "D80jAP5wbP4=", expectedToken)

	actualToken := sessionContainer.GetSessionToken("dbs/test/colls/test")
	assert.Equal(t, expectedToken, actualToken)
}

func TestSessionContainerClear(t *testing.T) {
	sessionContainer := newSessionContainer()

	sessionContainer.SetSessionToken("dbs/test/colls/test", "x1cRAMM0cZs=", "0:1#1234")
	sessionContainer.ClearSessionToken("dbs/test/colls/test")

	actualToken := sessionContainer.GetSessionToken("dbs/test/colls/test")
	assert.Equal(t, "", actualToken)
}
