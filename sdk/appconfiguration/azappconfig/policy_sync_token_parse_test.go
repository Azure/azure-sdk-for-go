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
	require.NotEmpty(t, err)
}

func TestSyncTokenParseTokenWithoutSeqNo(t *testing.T) {
	{
		tok1, err1 := parseToken("n=v")
		require.Empty(t, err1)
		require.NotEmpty(t, tok1)
		require.Equal(t, tok1.id, "n")
		require.Equal(t, tok1.value, "v")
		require.Empty(t, tok1.seqNo)
	}

	{
		tok2, err2 := parseToken("n=v;")
		require.Empty(t, err2)
		require.NotEmpty(t, tok2)
		require.Equal(t, tok2.id, "n")
		require.Equal(t, tok2.value, "v")
		require.Empty(t, tok2.seqNo)
	}
}

func TestSyncTokenParseTokenWithSeqNo(t *testing.T) {
	{
		tok1, err1 := parseToken("n=v;sn=1")
		require.Empty(t, err1)
		require.NotEmpty(t, tok1)
		require.Equal(t, tok1.id, "n")
		require.Equal(t, tok1.value, "v")
		require.Equal(t, tok1.seqNo, 1)
	}

	{
		tok2, err2 := parseToken("n=v;sn=1;")
		require.Empty(t, err2)
		require.NotEmpty(t, tok2)
		require.Equal(t, tok2.id, "n")
		require.Equal(t, tok2.value, "v")
		require.Equal(t, tok2.seqNo, 1)
	}
}

func TestSyncTokenParseTokenWithSeqNoMax(t *testing.T) {
	{
		tok1, err1 := parseToken("n=v;sn=9223372036854775807")
		require.Empty(t, err1)
		require.NotEmpty(t, tok1)
		require.Equal(t, tok1.id, "n")
		require.Equal(t, tok1.value, "v")
		require.Equal(t, tok1.seqNo, 9223372036854775807)
	}

	{
		tok2, err2 := parseToken("n=v;sn=9223372036854775807;")
		require.Empty(t, err2)
		require.NotEmpty(t, tok2)
		require.Equal(t, tok2.id, "n")
		require.Equal(t, tok2.value, "v")
		require.Equal(t, tok2.seqNo, 9223372036854775807)
	}
}

func TestSyncTokenParseTokenWithSeqNoParseError(t *testing.T) {
	{
		_, err1 := parseToken("n=v;sn=x")
		require.NotEmpty(t, err1)
	}

	{
		_, err2 := parseToken("n=v;sn=x;")
		require.NotEmpty(t, err2)
	}
}

func TestSyncTokenParseTokenWithoutIdValue(t *testing.T) {
	{
		_, err1 := parseToken("n=")
		require.NotEmpty(t, err1)
	}

	{
		_, err2 := parseToken("n=;")
		require.NotEmpty(t, err2)
	}

	{
		_, err3 := parseToken("n=;sn=1")
		require.NotEmpty(t, err3)
	}

	{
		_, err4 := parseToken("n=;sn=1;")
		require.NotEmpty(t, err4)
	}

	{
		_, err5 := parseToken(";")
		require.NotEmpty(t, err5)
	}

	{
		_, err6 := parseToken("sn=1")
		require.NotEmpty(t, err6)
	}

	{
		_, err7 := parseToken("sn=1;")
		require.NotEmpty(t, err7)
	}
}

func TestSyncTokenParseTokenNameTrim(t *testing.T) {
	{
		tok1, err1 := parseToken("  n  =v")
		require.Empty(t, err1)
		require.NotEmpty(t, tok1)
		require.Equal(t, tok1.id, "n")
		require.Equal(t, tok1.value, "v")
		require.Empty(t, tok1.seqNo)
	}

	{
		tok2, err2 := parseToken("  n   =v;")
		require.Empty(t, err2)
		require.NotEmpty(t, tok2)
		require.Equal(t, tok2.id, "n")
		require.Equal(t, tok2.value, "v")
		require.Empty(t, tok2.seqNo)
	}

	{
		tok3, err3 := parseToken("  n   =v;   sn  =1")
		require.Empty(t, err3)
		require.NotEmpty(t, tok3)
		require.Equal(t, tok3.id, "n")
		require.Equal(t, tok3.value, "v")
		require.Equal(t, tok3.seqNo, 1)
	}

	{
		tok4, err4 := parseToken("  n   =v;   sn  =1;")
		require.Empty(t, err4)
		require.NotEmpty(t, tok4)
		require.Equal(t, tok4.id, "n")
		require.Equal(t, tok4.value, "v")
		require.Equal(t, tok4.seqNo, 1)
	}
}

func TestSyncTokenParseTokenWithSeqNoReverseOrder(t *testing.T) {
	{
		tok1, err1 := parseToken("sn=1;n=v")
		require.Empty(t, err1)
		require.NotEmpty(t, tok1)
		require.Equal(t, tok1.id, "n")
		require.Equal(t, tok1.value, "v")
		require.Equal(t, tok1.seqNo, 1)
	}

	{
		tok2, err2 := parseToken("sn=1;n=v;")
		require.Empty(t, err2)
		require.NotEmpty(t, tok2)
		require.Equal(t, tok2.id, "n")
		require.Equal(t, tok2.value, "v")
		require.Equal(t, tok2.seqNo, 1)
	}
}

func TestSyncTokenParseTokenValueContaningEqualSigns(t *testing.T) {
	{
		tok1, err1 := parseToken("jtqGc1I4=MDoyOA==")
		require.Empty(t, err1)
		require.NotEmpty(t, tok1)
		require.Equal(t, tok1.id, "jtqGc1I4")
		require.Equal(t, tok1.value, "MDoyOA==")
		require.Empty(t, tok1.seqNo)
	}

	{
		tok2, err2 := parseToken("jtqGc1I4=MDoyOA==;")
		require.Empty(t, err2)
		require.NotEmpty(t, tok2)
		require.Equal(t, tok2.id, "jtqGc1I4")
		require.Equal(t, tok2.value, "MDoyOA==")
		require.Empty(t, tok2.seqNo)
	}

	{
		tok3, err3 := parseToken("jtqGc1I4=MDoyOA==;sn=28")
		require.Empty(t, err3)
		require.NotEmpty(t, tok3)
		require.Equal(t, tok3.id, "jtqGc1I4")
		require.Equal(t, tok3.value, "MDoyOA==")
		require.Equal(t, tok3.seqNo, 28)
	}

	{
		tok4, err4 := parseToken("jtqGc1I4=MDoyOA==;sn=28;")
		require.Empty(t, err4)
		require.NotEmpty(t, tok4)
		require.Equal(t, tok4.id, "jtqGc1I4")
		require.Equal(t, tok4.value, "MDoyOA==")
		require.Equal(t, tok4.seqNo, 28)
	}
}
