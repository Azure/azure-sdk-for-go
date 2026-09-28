// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package azcosmos

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNativeUserAgentOnWire(t *testing.T) {
	for _, tt := range []struct {
		name          string
		applicationID string
	}{
		{name: "empty"},
		{name: "application", applicationID: "orders"},
		{name: "24 characters", applicationID: strings.Repeat("a", 24)},
		{name: "25 characters", applicationID: strings.Repeat("a", 25)},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			var agents []string
			server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.UserAgent() == "" {
					// The capability probe is not a Cosmos service request.
					w.WriteHeader(http.StatusOK)
					return
				}
				mu.Lock()
				agents = append(agents, r.UserAgent())
				mu.Unlock()
				w.WriteHeader(http.StatusUnauthorized)
			}))
			// The native capability probe negotiates both protocols before fetching the account.
			server.Config.Protocols = new(http.Protocols)
			server.Config.Protocols.SetHTTP1(true)
			server.Config.Protocols.SetUnencryptedHTTP2(true)
			server.Start()
			defer server.Close()

			credential, err := NewKeyCredential(emulatorKey)
			require.NoError(t, err)
			client, err := NewClientWithKey(server.URL, credential, &ClientOptions{ApplicationID: tt.applicationID})
			require.NoError(t, err)
			defer func() { require.NoError(t, client.Close()) }()

			require.Error(t, client.Initialize(t.Context()))

			mu.Lock()
			defer mu.Unlock()
			require.NotEmpty(t, agents, "native driver sent no observed service request")
			suffixPattern := ` ` + regexp.QuoteMeta(tt.applicationID)
			if tt.applicationID == "" {
				suffixPattern = ` rustc/[^ |]+`
			}
			for _, agent := range agents {
				require.Contains(t, strings.Fields(agent), "azsdk-go-azcosmos/"+strings.TrimPrefix(serviceLibVersion, "v"))
				require.Regexp(t, `(?:^| )azsdk-rust-cosmos-driver/\S+(?: |$)`, agent)
				require.Regexp(t, suffixPattern+`\|F[0-9A-Fa-f]+$`, agent)
			}
		})
	}
}
