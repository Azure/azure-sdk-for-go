//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type syncToken struct {
	id    string
	value string
	seqNo int64
}

type syncTokenPolicy struct {
	syncTokens map[string]syncToken
}

func newSyncTokenPolicy() *syncTokenPolicy {
	return &syncTokenPolicy{}
}

func parseToken(tok string) (syncToken, error) {
	if tok == "" {
		return syncToken{}, errors.New("empty token")
	}

	id := ""
	val := ""
	var seq int64

	for _, kv := range strings.Split(tok, ";") {
		if kv == "" {
			continue
		}

		eq := strings.Index(kv, "=")
		if eq == -1 || eq == len(kv)-1 {
			continue
		}

		n := strings.TrimSpace(kv[:eq])
		v := kv[eq+1:]

		if n == "sn" {
			sn, err := strconv.ParseInt(v, 0, 64)
			if err != nil {
				return syncToken{}, err
			}
			seq = sn
		} else {
			id = n
			val = v
		}
	}

	if id != "" && val != "" {
		return syncToken{
			id:    id,
			value: val,
			seqNo: seq,
		}, nil
	}

	return syncToken{}, errors.New("didn't parse all the required parts")
}

func (policy *syncTokenPolicy) addToken(tok string) {
	for _, t := range strings.Split(tok, ",") {
		if st, err := parseToken(t); err == nil {
			if existing := policy.syncTokens[st.id]; existing.seqNo < st.seqNo {
				policy.syncTokens[st.id] = st
			}
		}
	}
}

func (policy *syncTokenPolicy) Do(req *policy.Request) (*http.Response, error) {
	const syncTokenHeaderName = "Sync-Token"
	var tokens []string
	for _, st := range policy.syncTokens {
		tokens = append(tokens, st.id)
	}

	req.Raw().Header[syncTokenHeaderName] = tokens

	resp, err := req.Next()

	if err != nil {
		for _, st := range resp.Header[syncTokenHeaderName] {
			policy.addToken(st)
		}
	}

	return resp, err
}
