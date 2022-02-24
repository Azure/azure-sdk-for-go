//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

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

	{
		pos := 0
		nm := ""
		vl := ""

		last := len(tok) - 1
		for i, ch := range []rune(tok) {
			if nm == "" && ch == '=' {
				nm = tok[0:i]
				pos = i + 1
				continue
			}

			if ch == ';' || i == last {
				frag := tok[pos:(i - pos)]
				if i == last {
					frag = frag[:len(frag)-1]
				}

				if nm == "" {
					nm = frag
				} else {
					vl = frag
				}

				nm = strings.TrimSpace(nm)

				if nm == "sn" {
					sn, err := strconv.ParseInt(vl, 0, 64)
					if err != nil {
						return syncToken{}, err
					}
					seq = sn
				} else if id == "" {
					id = nm
					val = vl
				}

				nm = ""
				vl = ""
				pos = i + 1
			}
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
