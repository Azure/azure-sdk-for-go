// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

//go:build cgo && ((darwin && !ios && arm64) || (linux && !android && amd64))

package e2e

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	azcosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2"
	internale2e "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos/v2/internal/e2e"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/stretchr/testify/require"
)

// TestAzureLive is deliberately separate from emulator certification. The pipeline runner
// requires this suite and every selected scenario to execute, not skip.
func TestAzureLive(t *testing.T) {
	if os.Getenv("AZURE_COSMOS_E2E_BACKEND") != "azureLive" {
		t.Skip("set AZURE_COSMOS_E2E_BACKEND=azureLive for opt-in live tests")
	}
	require.Empty(t, os.Getenv("EMULATOR"), "live and emulator selection must not be combined")
	for _, axis := range []string{"AZURE_COSMOS_E2E_PROFILE", "AZURE_COSMOS_E2E_ACCOUNT",
		"AZURE_COSMOS_E2E_RUNTIME", "AZURE_COSMOS_E2E_CLIENT"} {
		require.Empty(t, os.Getenv(axis), "the bounded live lane selects its declared profiles; overrides are unsupported")
	}
	endpoint, databaseID, containerID := os.Getenv("AZCOSMOS_ENDPOINT"),
		os.Getenv("AZCOSMOS_DATABASE"), os.Getenv("AZCOSMOS_CONTAINER")
	target, err := url.Parse(endpoint)
	require.NoError(t, err)
	require.True(t, target.Scheme == "https" && target.Hostname() != "" &&
		target.User == nil && target.RawQuery == "" && target.Fragment == "",
		"an explicit HTTPS account endpoint without credentials or query is required")
	require.NotEmpty(t, databaseID)
	require.NotEmpty(t, containerID)
	require.Equal(t, "live", os.Getenv("AZURE_RECORD_MODE"), "live tests must never use playback credentials")
	if os.Getenv("TF_BUILD") != "" {
		for _, name := range []string{"SYSTEM_ACCESSTOKEN", "AZURESUBSCRIPTION_CLIENT_ID",
			"AZURESUBSCRIPTION_TENANT_ID", "AZURESUBSCRIPTION_SERVICE_CONNECTION_ID"} {
			require.NotEmpty(t, os.Getenv(name), "federated pipeline credential environment is incomplete: "+name)
		}
	}
	catalog, err := internale2e.LoadCatalog()
	require.NoError(t, err)
	region := os.Getenv("AZCOSMOS_REGION")
	require.NoError(t, internale2e.ValidateLiveSelection(catalog, internale2e.LiveScenarios, region))
	tests := map[string]func(*testing.T, *liveFixture){
		"bootstrap.primary-success":     liveBootstrap,
		"item.lifecycle":                liveLifecycle,
		"item.optimistic-concurrency":   liveConcurrency,
		"diagnostics.success-and-error": liveDiagnostics,
		"query.pagination-resume":       livePagination,
	}
	require.Len(t, tests, len(internale2e.LiveScenarios))
	if path := os.Getenv("AZCOSMOS_LIVE_PLAN_FILE"); path != "" {
		plan, err := json.Marshal(internale2e.LiveScenarios)
		liveNoError(t, "encode selected scenario manifest", err)
		liveNoError(t, "write selected scenario manifest", os.WriteFile(path, plan, 0600))
	}
	t.Logf("backend=azureLive auth=entra target=%s/%s go=%s", runtime.GOOS, runtime.GOARCH, runtime.Version())
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == "github.com/Azure/azure-cosmos-driver/linux/amd64" ||
				dep.Path == "github.com/Azure/azure-cosmos-driver/darwin/arm64" {
				t.Logf("native-driver=%s", dep.Version)
			}
		}
	}
	for _, scenario := range internale2e.LiveScenarios {
		test, exists := tests[scenario.ID]
		require.True(t, exists, "selected live scenario must have an implementation")
		t.Run(scenario.ID, func(t *testing.T) {
			t.Logf("scenario=%s profile=%s fidelity=full expected-region=%s catalog-region=%q",
				scenario.ID, scenario.Profile, region, internale2e.LiveCatalogRegion)
			test(t, newLiveFixture(t, endpoint, databaseID, containerID))
		})
	}
}

type liveFixture struct {
	ctx       context.Context
	client    *azcosmos.Client
	container *azcosmos.ContainerClient
	pk        azcosmos.PartitionKey
	prefix    string
}

func newLiveFixture(t *testing.T, endpoint, databaseID, containerID string) *liveFixture {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	t.Cleanup(cancel)
	cred, err := credential.New(nil)
	liveNoError(t, "create Entra credential", err)
	timeout := 20 * time.Second
	shared, err := azcosmos.NewRuntime(&azcosmos.RuntimeOptions{
		Operation: azcosmos.OperationOptions{EndToEndTimeout: &timeout},
	})
	liveNoError(t, "create isolated runtime", err)
	t.Cleanup(func() { liveNoError(t, "close runtime", shared.Close()) })
	// Newly deployed Cosmos RBAC can take time to propagate. Retry setup only, never scenarios.
	setupCtx, setupCancel := context.WithTimeout(ctx, 90*time.Second)
	defer setupCancel()
	var client *azcosmos.Client
	for {
		client, err = azcosmos.NewClient(setupCtx, endpoint, cred, &azcosmos.ClientOptions{Runtime: shared})
		var serviceErr *azcosmos.Error
		if !errors.As(err, &serviceErr) || serviceErr.StatusCode != 403 ||
			serviceErr.SubStatus != 5302 || setupCtx.Err() != nil {
			break
		}
		t.Log("setup: waiting for newly deployed data-plane RBAC (403/5302)")
		select {
		case <-setupCtx.Done():
			liveNoError(t, "wait for RBAC readiness", setupCtx.Err())
		case <-time.After(5 * time.Second):
		}
	}
	liveNoError(t, "initialize token-auth client", err)
	t.Cleanup(func() { liveNoError(t, "close client", client.Close()) })
	container, err := client.NewContainer(ctx, databaseID, containerID)
	liveNoError(t, "create container handle", err)
	metadata, err := container.Read(ctx, nil)
	liveNoError(t, "read container metadata", err)
	require.Equal(t, []string{"/pk"}, metadata.ContainerProperties.PartitionKeyDefinition.Paths)
	var random [16]byte
	_, err = rand.Read(random[:])
	liveNoError(t, "generate test partition", err)
	prefix := "go-live-" + hex.EncodeToString(random[:])
	return &liveFixture{ctx: ctx, client: client, container: container,
		prefix: prefix, pk: azcosmos.NewPartitionKeyString(prefix)}
}

// Never print raw credential/service errors: they can contain endpoint or authentication data.
func liveNoError(t *testing.T, operation string, err error) {
	t.Helper()
	if err == nil {
		return
	}
	t.Fatalf("%s: %s", operation, liveErrorSummary(err))
}

func liveErrorSummary(err error) string {
	var serviceErr *azcosmos.Error
	if errors.As(err, &serviceErr) {
		hints := "suppressed"
		if !serviceErr.FromWire && serviceErr.Code == azcosmos.CodeTransportFailure {
			hints = liveTransportCauseHints(serviceErr.Message)
		}
		return fmt.Sprintf("code=%s status=%d substatus=%d attempts=%d fromWire=%v causeHints=%s",
			serviceErr.Code, serviceErr.StatusCode, serviceErr.SubStatus,
			serviceErr.AttemptCount, serviceErr.FromWire, hints)
	}
	return fmt.Sprintf("%T (details suppressed)", err)
}

// Reconstruct hints from fixed labels; replacing known secrets in raw text can miss new formats.
// These are diagnostic hints, not a stable classification of the driver's unstructured message.
func liveTransportCauseHints(message string) string {
	message = strings.ToLower(message)
	var hints []string
	for _, rule := range []struct {
		label   string
		matches []string
	}{
		{"connect", []string{"error trying to connect", "tcp connect error", "client error (connect)"}},
		{"permission-denied", []string{"permission denied", "(os error 13)"}},
		{"operation-not-permitted", []string{"operation not permitted", "(os error 1)"}},
		{"connection-refused", []string{"connection refused", "(os error 111)", "(os error 61)"}},
		{"connection-reset", []string{"connection reset", "(os error 104)", "(os error 54)"}},
		{"network-unreachable", []string{"network is unreachable", "(os error 101)"}},
		{"host-unreachable", []string{"host is unreachable", "no route to host", "(os error 113)"}},
		{"dns", []string{"dns error", "dns lookup", "failed to lookup address", "name or service not known", "temporary failure in name resolution", "no such host"}},
		{"tls-certificate", []string{"invalid peer certificate", "invalid certificate", "certificate verify failed", "certificate verification failed", "unknownissuer", "certificatenotvalidforname"}},
		{"tls-handshake", []string{"tls handshake", "tls error", "received fatal alert"}},
		{"timeout", []string{"timed out", "timeout", "deadline has elapsed"}},
		{"unexpected-eof", []string{"unexpected eof", "unexpected end of file"}},
	} {
		for _, match := range rule.matches {
			if strings.Contains(message, match) {
				hints = append(hints, rule.label)
				break
			}
		}
	}
	if len(hints) == 0 {
		return "unclassified"
	}
	return strings.Join(hints, ",")
}

func TestLiveErrorSummaryTransportHints(t *testing.T) {
	const sensitive = `url=https://synthetic-user:synthetic-password@synthetic-account.documents.azure.com/dbs/private-db?sig=synthetic-signature
Authorization: Bearer synthetic-token-do-not-log
AccountKey=synthetic-key-do-not-log; client_secret=synthetic-client-secret
address=203.0.113.47 tenant=00000000-1111-2222-3333-444444444444
##vso[task.logissue type=error]synthetic-log-injection`
	for _, tc := range []struct {
		name    string
		message string
		want    string
	}{
		{"blocked connect", "error trying to connect: tcp connect error: Operation not permitted (os error 1)", "connect,operation-not-permitted"},
		{"permission", "client error (Connect): Permission denied (os error 13)", "connect,permission-denied"},
		{"refused", "Connection refused (os error 111)", "connection-refused"},
		{"reset", "Connection reset by peer (os error 104)", "connection-reset"},
		{"network", "Network is unreachable (os error 101)", "network-unreachable"},
		{"host", "No route to host (os error 113)", "host-unreachable"},
		{"dns", "DNS error: failed to lookup address: Temporary failure in name resolution", "dns"},
		{"certificate", `invalid peer certificate: NotValidForName; UnknownIssuer; subject="synthetic-subject"`, "tls-certificate"},
		{"tls", "TLS handshake failed: received fatal alert", "tls-handshake"},
		{"timeout", "operation timed out: deadline has elapsed", "timeout"},
		{"eof", "unexpected EOF", "unexpected-eof"},
		{"unknown", "unrecognized native failure", "unclassified"},
		{"empty", "", "unclassified"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			message := tc.message + "\n" + sensitive + "\nencoded=" + url.QueryEscape(sensitive)
			nativeErr := &azcosmos.Error{
				Code: azcosmos.CodeTransportFailure, StatusCode: 503, SubStatus: 20010,
				AttemptCount: 3, Message: message,
				Body: []byte(sensitive), SessionToken: azcosmos.SessionToken(sensitive),
				ActivityID: sensitive,
				Diagnostics: &azcosmos.Diagnostics{
					RegionsContacted: []string{sensitive},
					Attempts:         []azcosmos.DiagnosticAttempt{{Endpoint: sensitive, Region: sensitive}},
				},
			}
			err := fmt.Errorf("%s: %w", sensitive, nativeErr)
			require.Equal(t, "code=TransportFailure status=503 substatus=20010 attempts=3 fromWire=false causeHints="+tc.want,
				liveErrorSummary(err), "only fixed hints and existing status fields may be logged")
		})
	}
}

func TestLiveErrorSummarySuppressesOtherDetails(t *testing.T) {
	const message = "connection refused; Authorization: Bearer synthetic-secret-do-not-log"
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{"service", &azcosmos.Error{Code: azcosmos.CodeServiceUnavailable, StatusCode: 503, FromWire: true, Message: message},
			"code=ServiceUnavailable status=503 substatus=0 attempts=0 fromWire=true causeHints=suppressed"},
		{"wire transport classification", &azcosmos.Error{Code: azcosmos.CodeTransportFailure, StatusCode: 503, FromWire: true, Message: message},
			"code=TransportFailure status=503 substatus=0 attempts=0 fromWire=true causeHints=suppressed"},
		{"credential", &azcosmos.Error{Code: azcosmos.CodeAuthenticationFailed, StatusCode: 401, Message: message},
			"code=AuthenticationFailed status=401 substatus=0 attempts=0 fromWire=false causeHints=suppressed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, liveErrorSummary(tc.err))
		})
	}
	summary := liveErrorSummary(errors.New(message))
	require.Contains(t, summary, "(details suppressed)")
	require.NotContains(t, summary, message)
	require.NotContains(t, summary, "synthetic-secret")
}

func (f *liveFixture) create(t *testing.T, suffix string, value int64) azcosmos.ItemResponse {
	t.Helper()
	id := f.prefix + suffix
	// Register cleanup before submission: a failed/cancelled create may still commit.
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		_, err := f.container.DeleteItem(ctx, f.pk, id, nil)
		var serviceErr *azcosmos.Error
		if errors.As(err, &serviceErr) && serviceErr.Code == azcosmos.CodeNotFound {
			return
		}
		liveNoError(t, "clean up owned item", err)
	})
	body := item(id, f.prefix, value)
	response, err := f.container.CreateItem(f.ctx, f.pk, id, body, nil)
	liveNoError(t, "create item", err)
	require.Equal(t, 201, response.StatusCode)
	return response
}

func liveBootstrap(t *testing.T, f *liveFixture) {
	f.create(t, "-bootstrap", 1)
}

func liveLifecycle(t *testing.T, f *liveFixture) {
	created := f.create(t, "-lifecycle", 1)
	require.NotEmpty(t, created.SessionToken)
	response, err := f.container.ReadItem(f.ctx, f.pk, f.prefix+"-lifecycle", &azcosmos.ReadItemOptions{
		SessionToken: created.SessionToken,
	})
	liveNoError(t, "read with explicit session token", err)
	require.Equal(t, int64(1), mustUnmarshalItem(t, response.Value).Value)
	replaced, err := f.container.ReplaceItem(f.ctx, f.pk, f.prefix+"-lifecycle",
		item(f.prefix+"-lifecycle", f.prefix, 2), nil)
	liveNoError(t, "replace item", err)
	require.Equal(t, 200, replaced.StatusCode)
	read, err := f.container.ReadItem(f.ctx, f.pk, f.prefix+"-lifecycle", nil)
	liveNoError(t, "read replacement", err)
	require.Equal(t, int64(2), mustUnmarshalItem(t, read.Value).Value)
	deleted, err := f.container.DeleteItem(f.ctx, f.pk, f.prefix+"-lifecycle", nil)
	liveNoError(t, "delete item", err)
	require.Equal(t, 204, deleted.StatusCode)
	_, err = f.container.ReadItem(f.ctx, f.pk, f.prefix+"-lifecycle", nil)
	_ = liveCode(t, err, azcosmos.CodeNotFound, 404)
}

func liveCode(t *testing.T, err error, code azcosmos.Code, status int) *azcosmos.Error {
	t.Helper()
	var serviceErr *azcosmos.Error
	require.True(t, errors.As(err, &serviceErr), "expected typed Cosmos error")
	require.Equal(t, code, serviceErr.Code)
	require.Equal(t, status, serviceErr.StatusCode)
	require.True(t, serviceErr.FromWire)
	return serviceErr
}

func liveConcurrency(t *testing.T, f *liveFixture) {
	created := f.create(t, "-etag", 1)
	require.NotEmpty(t, created.ETag)
	condition, err := azcosmos.IfMatch(created.ETag)
	liveNoError(t, "create ETag precondition", err)
	body := item(f.prefix+"-etag", f.prefix, 2)
	updated, err := f.container.ReplaceItem(f.ctx, f.pk, f.prefix+"-etag", body,
		&azcosmos.ReplaceItemOptions{Precondition: condition})
	liveNoError(t, "conditional replacement", err)
	require.NotEqual(t, created.ETag, updated.ETag)
	_, err = f.container.ReplaceItem(f.ctx, f.pk, f.prefix+"-etag",
		item(f.prefix+"-etag", f.prefix, 3),
		&azcosmos.ReplaceItemOptions{Precondition: condition})
	_ = liveCode(t, err, azcosmos.CodePreconditionFailed, 412)
	read, err := f.container.ReadItem(f.ctx, f.pk, f.prefix+"-etag", nil)
	liveNoError(t, "read after rejected replacement", err)
	require.Equal(t, int64(2), mustUnmarshalItem(t, read.Value).Value)
}

func liveDiagnostics(t *testing.T, f *liveFixture) {
	f.create(t, "-diagnostics", 1)
	read, err := f.container.ReadItem(f.ctx, f.pk, f.prefix+"-diagnostics", nil)
	liveNoError(t, "read for success diagnostics", err)
	requireCriticalDiagnostics(t, read.Diagnostics, 200)
	require.NotEmpty(t, read.Diagnostics.RegionsContacted)
	require.NotEmpty(t, read.ActivityID)
	require.Positive(t, read.RequestCharge)
	_, err = f.container.ReadItem(f.ctx, f.pk, f.prefix+"-missing", nil)
	failure := liveCode(t, err, azcosmos.CodeNotFound, 404)
	requireCriticalDiagnostics(t, failure.Diagnostics, 404)
	require.NotEmpty(t, failure.ActivityID)
	require.NotEmpty(t, failure.Diagnostics.RegionsContacted)
	t.Logf("status=404 substatus=%d attempts=%d charge=%g", failure.SubStatus, failure.AttemptCount, failure.RequestCharge)
}

func livePagination(t *testing.T, f *liveFixture) {
	const count = 6
	var session azcosmos.SessionToken
	for i := 0; i < count; i++ {
		session = f.create(t, "-page-"+string(rune('a'+i)), int64(i)).SessionToken
	}
	query, err := azcosmos.NewQuery(`SELECT * FROM c WHERE c.pk = @pk ORDER BY c["value"] ASC`).WithParameter("@pk", f.prefix)
	liveNoError(t, "parameterize query", err)
	scope := azcosmos.NewFeedScopeForPartitionKey(f.pk)
	pager := f.container.NewQueryItemsPager[Item](query, scope, &azcosmos.QueryOptions{
		SessionToken: session, Feed: azcosmos.FeedOptions{PageSizeHint: 1},
	})
	t.Cleanup(func() { liveNoError(t, "close first pager", pager.Close()) })
	first, err := nextNonEmptyLivePage(f.ctx, pager)
	liveNoError(t, "fetch first nonempty page", err)
	require.NotEmpty(t, first.Items)
	require.Less(t, len(first.Items), count)
	token, err := pager.ContinuationToken(f.ctx)
	liveNoError(t, "checkpoint query", err)
	require.NotEmpty(t, token)
	liveNoError(t, "close first pager before resume", pager.Close())
	resumed := f.container.NewQueryItemsPager[Item](query, scope, &azcosmos.QueryOptions{
		SessionToken: session, Feed: azcosmos.FeedOptions{PageSizeHint: 1, ContinuationToken: token},
	})
	t.Cleanup(func() { liveNoError(t, "close resumed pager", resumed.Close()) })
	items := append([]Item(nil), first.Items...)
	for resumed.More() {
		page, err := resumed.NextPage(f.ctx)
		liveNoError(t, "resume query", err)
		items = append(items, page.Items...)
	}
	require.Len(t, items, count)
	for i, value := range items {
		require.Equal(t, int64(i), value.Value, "ordered resume must not lose or duplicate results")
	}
}

type liveItemPager interface {
	More() bool
	NextPage(context.Context) (azcosmos.QueryItemsResponse[Item], error)
}

func nextNonEmptyLivePage(ctx context.Context, pager liveItemPager) (azcosmos.QueryItemsResponse[Item], error) {
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil || len(page.Items) != 0 {
			return page, err
		}
	}
	return azcosmos.QueryItemsResponse[Item]{}, nil
}

type liveQueryPageSequence struct {
	pages []azcosmos.QueryItemsResponse[Item]
	err   error
	calls int
}

func (p *liveQueryPageSequence) More() bool {
	return len(p.pages) != 0 || p.err != nil
}

func (p *liveQueryPageSequence) NextPage(ctx context.Context) (azcosmos.QueryItemsResponse[Item], error) {
	p.calls++
	if err := ctx.Err(); err != nil {
		return azcosmos.QueryItemsResponse[Item]{}, err
	}
	if len(p.pages) == 0 {
		return azcosmos.QueryItemsResponse[Item]{}, p.err
	}
	page := p.pages[0]
	p.pages = p.pages[1:]
	return page, nil
}

func TestNextNonEmptyLivePage(t *testing.T) {
	first := azcosmos.QueryItemsResponse[Item]{Items: []Item{{ID: "first", Value: 1}}}
	second := azcosmos.QueryItemsResponse[Item]{Items: []Item{{ID: "second", Value: 2}}}
	failure := errors.New("query failed")
	for _, tc := range []struct {
		name      string
		pages     []azcosmos.QueryItemsResponse[Item]
		err       error
		want      azcosmos.QueryItemsResponse[Item]
		wantCalls int
	}{
		{"nonempty first page", []azcosmos.QueryItemsResponse[Item]{first, second}, nil, first, 1},
		{"empty nonterminal pages", []azcosmos.QueryItemsResponse[Item]{{}, {}, first, second}, nil, first, 3},
		{"empty exhausted feed", []azcosmos.QueryItemsResponse[Item]{{}, {}}, nil, azcosmos.QueryItemsResponse[Item]{}, 2},
		{"already exhausted", nil, nil, azcosmos.QueryItemsResponse[Item]{}, 0},
		{"error after empty page", []azcosmos.QueryItemsResponse[Item]{{}}, failure, azcosmos.QueryItemsResponse[Item]{}, 2},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pager := &liveQueryPageSequence{pages: tc.pages, err: tc.err}
			page, err := nextNonEmptyLivePage(t.Context(), pager)
			require.ErrorIs(t, err, tc.err)
			require.Equal(t, tc.want, page)
			require.Equal(t, tc.wantCalls, pager.calls, "stop at the first items, exhaustion, or error")
		})
	}
}
