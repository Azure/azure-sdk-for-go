# Managed download flow: initial GetBlob, download hint, and Data Locality

This describes what `DownloadBuffer` and `DownloadFile` actually do in
`sdk/storage/azblob/blob/client.go` after the STG105 Data Locality work was integrated with
main's initial-GetBlob download architecture (PR #27438).

`DownloadStream` is the one-shot public read. It never fetches a layout on the caller's behalf —
it only honours an explicit `DownloadStreamOptions.LayoutEndpoint`. Layout discovery belongs to
the managed downloads described here.

## The flow

```
DownloadBuffer / DownloadFile
  │
  ├─ explicit Range.Count AND layout routing disabled
  │     └─ parallelDownload: chunked ranged reads, account endpoint
  │        (main's request pattern, unchanged)
  │
  └─ otherwise: initial ranged GetBlob  ──▶  account endpoint, always
        │
        │   captures: total size · ETag · first chunk of data · x-ms-download-hint
        │
        ├─ InvalidRange                    → empty blob, return 0
        ├─ no Content-Range / Content-Length → error (304 from caller conditions)
        │
        ├─ write the initial chunk to the destination
        │
        ├─ initial chunk covered the request?  → DONE. No GetLayout, no parallel reads.
        │
        └─ data remains:
              ├─ layout routing disabled        → parallel reads, account endpoint
              ├─ hint absent or not "layout"    → parallel reads, account endpoint
              └─ routing enabled AND hint == layout
                    → GetLayout for the remaining range only
                    → fallback / no ranges → account endpoint
                    → usable layout → per-chunk endpoint selection
```

The three conditions are `o.layoutAwareRoutingEnabled() && downloadHint != nil && *downloadHint
== generated.DownloadHintLayout`, evaluated only after the initial read has been consumed and
only when `initialChunkSize < count`.

## Why the initial read always uses the account endpoint

The hint arrives *on* the initial response, so there is nothing to route by yet. The first block
is therefore never layout-routed. This is a deliberate consequence of hint-gating: the previous
STG105 behaviour fetched the layout first and routed every chunk, at the cost of a GetLayout call
on every download including ones that finish in a single request.

## ETag consistency

The initial response is the consistency anchor. Its ETag is pinned as `If-Match` on:

- every remaining parallel chunk, and
- the `GetLayout` enumeration, because `resolveLayout` reads `o.AccessConditions` after the pin.

The pin is applied to a **copy** of the caller's `AccessConditions` and of the nested
`ModifiedAccessConditions`, so the caller's own structs come back untouched.

A caller-supplied `If-Match` is applied to the initial read as usual. If it matches, the ETag that
comes back is the same value, so pinning changes nothing. If it is stale, the initial read fails
with `ConditionNotMet` before any layout is fetched. The one case where behaviour differs is a
caller-supplied `ETagAny`, which is narrowed to the concrete version — that is intentional, since
`ETagAny` would otherwise allow chunks from different blob versions.

**Note for reviewers:** ETag pinning now also applies when layout routing is *disabled*, because
it comes from #27438 rather than from the Data Locality work. STG105 previously sent no `If-Match`
on that path.

## GetLayout

- **Range.** Only the remainder is enumerated: `Offset = o.Range.Offset + writerOffset`,
  `Count = remaining`. The part already delivered by the initial read is not described.
- **Pagination.** `GetLayoutPager` pages on `NextMarker`. The TypeSpec emitter exposes only
  `GetLayout(ctx, options)` — there is no generated pager and no `GetLayoutCreateRequest` /
  `GetLayoutHandleResponse` — so the pager is built on repeated single-options calls.
- **Cross-page consistency.** The first page's ETag is sent as `If-Match` on every later page, so
  one enumeration always describes one version of the blob. A caller-supplied `If-Match` wins.
- **Caching.** The layout lives in a `temporal.Resource` for the duration of one download, with
  `shouldRefreshLayout` refreshing ~30s before a 5-minute expiry. It is resolved **once**, before
  the parallel chunks start, so a fallback or no-layout answer costs one enumeration rather than
  one per chunk.
- **Fallback.** `getLayout` turns a 400 or 5xx into a cached fallback layout rather than an error,
  so the decision is made once. `resolveLayout` returns a nil resource for a fallback layout or a
  layout with no ranges, and the remainder is read from the account endpoint.

## Routing the remaining chunks

Each chunk's blob offset is `chunkStart + writerOffset + o.Range.Offset`. That offset is
binary-searched against the cached layout ranges (`getIdealEndpoint`) and the result is set as
`DownloadStreamOptions.LayoutEndpoint`.

`shared.LayoutPolicy` performs the rewrite. It is registered **per-call**, alongside main's
`shared.NewRangePolicy()`, in `base.GetAzClient`. Per-call matters: the policy moves the account
host into the `Host` header and replaces the URL host. Those mutations are made on the request
itself, so they survive retries. Running it per-retry would re-apply the rewrite to an
already-rewritten request and copy the layout host into the `Host` header.

If a mid-download layout refresh fails, the chunk falls back to the configured endpoint rather
than failing the download.

## Retries

`DownloadStreamResponse` carries the unexported `layoutEndpoint` that produced it, and
`NewRetryReader` passes it back into the retry's `DownloadStreamOptions`. Without this a retried
chunk silently reverted to the account endpoint — losing exactly the locality the routing exists
to provide.

Session authentication is re-applied on every retry, because the session policy sits in
`PerRetry`, inside the retry loop. The signature is built from the configured account name and
the URL path, never the host, so the layout rewrite cannot invalidate it.

## Edge cases

| Case | Behaviour |
|---|---|
| Empty blob | Initial read returns `InvalidRange`; download returns 0 with no error. |
| Blob smaller than one block | Completed by the initial read. No GetLayout, no parallel reads. |
| Explicit `Range.Count`, routing off | Skips the initial read entirely (main's request pattern). |
| Explicit `Range.Count`, routing on | Initial read still happens — it carries the hint and ETag. The count is not overwritten from `Content-Range`. |
| Non-zero `Range.Offset` | Offsets are relative to it throughout; the layout is requested for the remaining range only. |
| 304 Not Modified | No `Content-Range` and no `Content-Length` → explicit error rather than a size-parse failure. |
| GetLayout returns 204 / no ranges | Treated as "no layout"; remainder read from the account endpoint. |
| GetLayout returns 400 / 5xx | Cached fallback layout; one enumeration, then the account endpoint. |
| Layout cache expiry mid-download | Refreshed once via `temporal.Resource`; a failed refresh falls back per chunk. |
| Structured message | `Content-Length` is the encoded size, so decoded length comes from `Content-Range` for both the initial read and each chunk. |
| Progress | The initial read assigns the byte count; parallel chunks accumulate deltas under a shared lock. No double counting. |

## Scope

Data Locality applies to Get Blob and the managed downloads built on it, and nothing else.
Session authentication covers Get Blob, Get Blob Properties, Put Blob, Put Block and Put Block
List; those four write/metadata operations did **not** gain layout routing.
