//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azappconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSyncTokenParseTokenEmpty(t *testing.T) {
	_, err := parseToken("")
	require.Error(t, err)
}

func TestSyncTokenParseTokenWithoutSeqNo(t *testing.T) {
	tok1, err1 := parseToken("n=v")
	require.NoError(t, err1)
	require.NotEmpty(t, tok1)
	require.Equal(t, "n", tok1.id)
	require.Equal(t, "v", tok1.value)
	require.Empty(t, tok1.seqNo)

	tok2, err2 := parseToken("n=v;")
	require.NoError(t, err2)
	require.NotEmpty(t, tok2)
	require.Equal(t, "n", tok2.id)
	require.Equal(t, "v", tok2.value)
	require.Empty(t, tok2.seqNo)
}

func TestSyncTokenParseTokenWithSeqNo(t *testing.T) {
	tok1, err1 := parseToken("n=v;sn=1")
	require.NoError(t, err1)
	require.NotEmpty(t, tok1)
	require.Equal(t, "n", tok1.id)
	require.Equal(t, "v", tok1.value)
	require.Equal(t, int64(1), tok1.seqNo)

	tok2, err2 := parseToken("n=v;sn=1;")
	require.NoError(t, err2)
	require.NotEmpty(t, tok2)
	require.Equal(t, "n", tok2.id)
	require.Equal(t, "v", tok2.value)
	require.Equal(t, int64(1), tok2.seqNo)
}

func TestSyncTokenParseTokenWithSeqNoMax(t *testing.T) {
	tok1, err1 := parseToken("n=v;sn=9223372036854775807")
	require.NoError(t, err1)
	require.NotEmpty(t, tok1)
	require.Equal(t, "n", tok1.id)
	require.Equal(t, "v", tok1.value)
	require.Equal(t, int64(9223372036854775807), tok1.seqNo)

	tok2, err2 := parseToken("n=v;sn=9223372036854775807;")
	require.NoError(t, err2)
	require.NotEmpty(t, tok2)
	require.Equal(t, "n", tok2.id)
	require.Equal(t, "v", tok2.value)
	require.Equal(t, int64(9223372036854775807), tok2.seqNo)
}

func TestSyncTokenParseTokenWithSeqNoParseError(t *testing.T) {
	_, err1 := parseToken("n=v;sn=x")
	require.Error(t, err1)

	_, err2 := parseToken("n=v;sn=x;")
	require.Error(t, err2)
}

func TestSyncTokenParseTokenWithoutIdValue(t *testing.T) {
	_, err1 := parseToken("n=")
	require.Error(t, err1)

	_, err2 := parseToken("n=;")
	require.Error(t, err2)

	_, err3 := parseToken("n=;sn=1")
	require.Error(t, err3)

	_, err4 := parseToken("n=;sn=1;")
	require.Error(t, err4)

	_, err5 := parseToken(";")
	require.Error(t, err5)

	_, err6 := parseToken("sn=1")
	require.Error(t, err6)

	_, err7 := parseToken("sn=1;")
	require.Error(t, err7)
}

func TestSyncTokenParseTokenNameTrim(t *testing.T) {
	tok1, err1 := parseToken("  n  =v")
	require.NoError(t, err1)
	require.NotEmpty(t, tok1)
	require.Equal(t, "n", tok1.id)
	require.Equal(t, "v", tok1.value)
	require.Empty(t, tok1.seqNo)

	tok2, err2 := parseToken("  n   =v;")
	require.NoError(t, err2)
	require.NotEmpty(t, tok2)
	require.Equal(t, "n", tok2.id)
	require.Equal(t, "v", tok2.value)
	require.Empty(t, tok2.seqNo)

	tok3, err3 := parseToken("  n   =v;   sn  =1")
	require.NoError(t, err3)
	require.NotEmpty(t, tok3)
	require.Equal(t, "n", tok3.id)
	require.Equal(t, "v", tok3.value)
	require.Equal(t, int64(1), tok3.seqNo)

	tok4, err4 := parseToken("  n   =v;   sn  =1;")
	require.NoError(t, err4)
	require.NotEmpty(t, tok4)
	require.Equal(t, "n", tok4.id)
	require.Equal(t, "v", tok4.value)
	require.Equal(t, int64(1), tok4.seqNo)
}

func TestSyncTokenParseTokenWithSeqNoReverseOrder(t *testing.T) {
	tok1, err1 := parseToken("sn=1;n=v")
	require.NoError(t, err1)
	require.NotEmpty(t, tok1)
	require.Equal(t, "n", tok1.id)
	require.Equal(t, "v", tok1.value)
	require.Equal(t, int64(1), tok1.seqNo)

	tok2, err2 := parseToken("sn=1;n=v;")
	require.NoError(t, err2)
	require.NotEmpty(t, tok2)
	require.Equal(t, "n", tok2.id)
	require.Equal(t, "v", tok2.value)
	require.Equal(t, int64(1), tok2.seqNo)
}

func TestSyncTokenParseTokenValueContainingEqualSigns(t *testing.T) {
	tok1, err1 := parseToken("jtqGc1I4=MDoyOA==")
	require.NoError(t, err1)
	require.NotEmpty(t, tok1)
	require.Equal(t, "jtqGc1I4", tok1.id)
	require.Equal(t, "MDoyOA==", tok1.value)
	require.Empty(t, tok1.seqNo)

	tok2, err2 := parseToken("jtqGc1I4=MDoyOA==;")
	require.NoError(t, err2)
	require.NotEmpty(t, tok2)
	require.Equal(t, "jtqGc1I4", tok2.id)
	require.Equal(t, "MDoyOA==", tok2.value)
	require.Empty(t, tok2.seqNo)

	tok3, err3 := parseToken("jtqGc1I4=MDoyOA==;sn=28")
	require.NoError(t, err3)
	require.NotEmpty(t, tok3)
	require.Equal(t, "jtqGc1I4", tok3.id)
	require.Equal(t, "MDoyOA==", tok3.value)
	require.Equal(t, int64(28), tok3.seqNo)

	tok4, err4 := parseToken("jtqGc1I4=MDoyOA==;sn=28;")
	require.NoError(t, err4)
	require.NotEmpty(t, tok4)
	require.Equal(t, "jtqGc1I4", tok4.id)
	require.Equal(t, "MDoyOA==", tok4.value)
	require.Equal(t, int64(28), tok4.seqNo)
}
