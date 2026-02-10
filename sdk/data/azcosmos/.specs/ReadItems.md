# ReadItems (ReadMany via Per-Partition Queries) — Design Specification

## 1. What the Python Code Does Today

The Python Cosmos DB SDK implements `read_items` (a.k.a. "read many") in two cooperating modules:

### Entry Point — `_cosmos_client_connection.py` → `read_items()`

- Accepts a list of `(item_id, partition_key_value)` tuples, request options, an optional `ThreadPoolExecutor`, and `max_concurrency`.
- Fetches the container's `PartitionKeyDefinition` (paths, kind, version).
- Delegates to `ReadItemsHelperSync`.

### Core Logic — `_read_items_helper.py` → `ReadItemsHelperSync`

#### 1.1 Partition Items by Physical Range

`_partition_items_by_range()`:

1. Groups input items by their **logical partition key value** (deduplicating range lookups).
2. For each unique PK value, computes the Effective Partition Key (EPK) range and resolves the overlapping **physical partition key range** via the routing map provider.
3. Produces `dict[range_id → list[(original_index, item_id, pk_value)]]`.

#### 1.2 Chunk by `max_items_per_query` (1 000)

`_create_query_chunks()`:

- Each physical range's item list is split into chunks of ≤ 1 000.
- Each chunk becomes a unit of work: `{range_id: [(idx, id, pk), …]}`.

#### 1.3 Execute Chunks Concurrently

`_execute_with_executor()`:

- Submits one `_execute_query_chunk_worker` per chunk to a `ThreadPoolExecutor`.
- Collects results preserving original input order (via the `original_index`).
- Sums `x-ms-request-charge` across all responses.

#### 1.4 Per-Chunk Worker

`_execute_query_chunk_worker()`:

| Chunk Size | Strategy |
|---|---|
| **1 item** | Single **point read** (`ReadItem`). 404 → `None`, item silently skipped in output. |
| **> 1 item** | Constructs a **parameterized SQL query** and runs it via `QueryItems`. |

#### 1.5 Query Construction (`_query_builder.py`)

Three query shapes are chosen depending on the items in the chunk:

| Condition | Query Shape | Example |
|---|---|---|
| PK path is `/id` **and** every pk value == item id | **ID-only IN** | `SELECT * FROM c WHERE c.id IN (@param_id0, @param_id1, …)` |
| All items share the **same logical PK value** | **PK = @pk AND id IN (…)** | `SELECT * FROM c WHERE c.pk = @pk AND c.id IN (@id0, @id1, …)` |
| Items span **multiple logical PK values** within the same physical range | **OR-of-conjunctions** | `SELECT * FROM c WHERE (c.id = @param_id0 AND c.pk = @param_pk00) OR (c.id = @param_id1 AND c.pk = @param_pk10) …` |

Special handling for undefined/empty/null partition key values → `IS_DEFINED(field) = false`.

All queries use **parameterized values** to avoid injection and enable plan caching.

The query is sent with the `x-ms-documentdb-partitionkeyrangeid` header set to the physical range, so the Gateway routes it to the correct replica.

#### 1.6 Result Assembly

- Returned items are re-sorted into the original input order.
- Missing items (404s and non-matching query results) are silently omitted.
- The aggregate `x-ms-request-charge` is returned.

---

## 2. What the Go Code Does Today

### 2.1 Public API

`ContainerClient.ReadManyItems(ctx, []ItemIdentity, *ReadManyOptions) → (ReadManyItemsResponse, error)`

Defined in `cosmos_container.go` (line 456).

### 2.2 Strategy Selection

| Condition | Path |
|---|---|
| `ReadManyOptions.QueryEngine != nil` | `executeReadManyWithEngine` — delegates to the external native query engine (CGo / WASM). |
| Otherwise (default) | `executeReadManyWithPointReads` — N individual point reads. |

### 2.3 `executeReadManyWithPointReads` (current default)

In `cosmos_container_read_many.go`:

- Creates a fixed-size goroutine pool (size = `min(concurrency, len(items))`).
- Each goroutine calls `ContainerClient.ReadItem` for a single item.
- 404 errors → item silently skipped.
- Any other error → cancels remaining work and returns the error.
- Results are collected in input order.
- Total `RequestCharge` is summed from individual responses.

**Limitation:** For N items this issues N separate HTTP round-trips. The RU cost is N × (point-read cost) and latency is dominated by the parallelism limit.

### 2.4 `executeReadManyWithEngine`

In `cosmos_container_read_many.go`:

- Fetches container properties (to get `PartitionKeyDefinition`) and raw partition key ranges.
- Builds item identities with JSON-serialized PK values.
- Creates a `ReadManyPipeline` via the external query engine.
- Runs the pipeline: the engine generates per-range `QueryRequest`s.
- `runEngineRequests` (in `cosmos_container_query_engine.go`) executes those requests concurrently via `sendQueryRequest`, piping results back through `ProvideData`.
- Returns merged items and total RU charge.

**Limitation:** Requires a native query engine binary (CGo / WASM). Not available in all environments.

### 2.5 Existing Types

| Type | File | Purpose |
|---|---|---|
| `ItemIdentity` | `cosmos_container.go` | `{ID string, PartitionKey PartitionKey}` |
| `ReadManyOptions` | `cosmos_read_many_request_options.go` | Session, consistency, dedicated gateway, query engine, max concurrency |
| `ReadManyItemsResponse` | `cosmos_query_response.go` | `{RequestCharge float32, Items [][]byte}` |
| `QueryParameter` | `cosmos_query.go` | `{Name string, Value any}` |
| `partitionKeyRange` | `cosmos_partition_key_range.go` | Physical range metadata |
| `PartitionKeyDefinition` | `partition_key_definition.go` | Paths, Kind, Version |

---

## 3. What the Go Code Should Do (Proposed: `executeReadManyWithQueries`)

**Replace** the point-read strategy (`executeReadManyWithPointReads`) with a new query-based strategy — `executeReadManyWithQueries` — that mirrors the Python logic. The old point-read function and all code exclusively supporting it are **deleted**.

### 3.1 High-Level Algorithm

```
ReadManyItems(ctx, items, opts)
  ├─ if QueryEngine != nil → executeReadManyWithEngine (unchanged)
  └─ else → executeReadManyWithQueries (NEW, replaces executeReadManyWithPointReads)
```

### 3.2 Steps

1. **Fetch container properties** → obtain `PartitionKeyDefinition` (paths, kind).
2. **Fetch partition key ranges** → `getPartitionKeyRanges(ctx, nil)`.
3. **Group items by physical partition range:**
   - For each unique logical PK value, compute the EPK hash and find the overlapping physical range using the `minInclusive`/`maxExclusive` intervals from step 2.
   - Build `map[rangeID] → []indexedItem` where `indexedItem = {originalIndex, itemID, pkValue}`.
4. **Chunk per-range lists** at `maxItemsPerQuery = 1000`.
5. **For each chunk, build a parameterized query** using the three-shape strategy from the Python `_QueryBuilder`:
   - **ID-only IN** — when PK path is `/id` and every PK value == item ID.
   - **PK + ID IN** — when all items in the chunk share a single logical PK.
   - **OR-of-conjunctions** — general multi-PK case.
6. **Execute chunks concurrently** using a goroutine pool (size = `determineConcurrency(opts.MaxConcurrency)`).
   - Each goroutine sends a `sendQueryRequest` with `x-ms-documentdb-partitionkeyrangeid` header targeting the physical range.
   - Paginate via continuation token until exhausted.
7. **Collect results:**
   - Match returned `id` fields back to original indices.
   - Sort by original index.
   - Sum `x-ms-request-charge`.
8. **Return `ReadManyItemsResponse`.**

### 3.3 EPK Hashing

The Python SDK computes the EPK via `_get_epk_range_for_partition_key`, which uses the V2 MurmurHash-based effective partition key algorithm. The Go SDK must implement the same hashing or use the Gateway's partition key range resolution.

**Recommended approach:** Use the existing `PartitionKey.toJsonString()` method and the Gateway-returned `partitionKeyRange.minInclusive`/`maxExclusive` boundaries. Each `PartitionKey` value maps to exactly one physical range. We can compute the effective partition key hash client-side (V2 hash) and binary-search the sorted range list, **or** we can simply send each query as a cross-partition query scoped to the known partition key value via the `x-ms-documentdb-partitionkey` header (letting the Gateway route it).

**Simplest correct approach for V1:** Set the `x-ms-documentdb-partitionkey` header (PK value as JSON array) on each per-partition query. The Gateway resolves the physical range automatically. This avoids implementing EPK hashing in Go and is what the existing `NewQueryItemsPager` already does. Chunking is still done by logical PK grouping (items with equal PK value go together).

If performance is a concern (too many round-trips when multiple logical PKs map to the same physical range), a follow-up can add client-side EPK hashing to coalesce by physical range. The Python SDK already does this optimisation.

### 3.4 Single-Item Optimisation

When a chunk contains exactly 1 item, issue a `SELECT * FROM c WHERE c.id = @id AND c.pk = @pk` query (or equivalent) rather than a point read. This keeps a single, uniform code path — all items flow through query construction and `sendQueryRequest`. No point-read logic is needed.

### 3.5 Error Handling

- Query returning fewer items than requested → silently omit missing items (mirrors Python behavior).
- Transient errors → bubble up immediately (cancel remaining work, consistent with current behavior).

---

## 4. Files to Create or Modify

### 4.1 New File: `cosmos_query_builder.go`

Contains the Go port of `_QueryBuilder`:

```go
// queryBuilder builds parameterized SQL queries for read-many operations.
type queryBuilder struct{}

func (qb queryBuilder) isIDPartitionKeyQuery(items []ItemIdentity, pkDef PartitionKeyDefinition) bool
func (qb queryBuilder) isSingleLogicalPartitionQuery(items []indexedItem) bool
func (qb queryBuilder) buildIDInQuery(items []indexedItem) (string, []QueryParameter)
func (qb queryBuilder) buildPKAndIDInQuery(items []indexedItem, pkDef PartitionKeyDefinition) (string, []QueryParameter)
func (qb queryBuilder) buildParameterizedQueryForItems(items []indexedItem, pkDef PartitionKeyDefinition) (string, []QueryParameter, error)
func (qb queryBuilder) getFieldExpression(path string) string
```

All functions are unexported (package-internal).

### 4.2 New File: `cosmos_query_builder_test.go`

Unit tests for every query builder function (no emulator needed):

- Various partition key shapes (single, hierarchical, `/id`-based).
- Edge cases: special characters in PK paths, null/undefined PK values, empty items list.
- Verify parameterized output matches expected SQL and parameter lists.

### 4.3 Modified File: `cosmos_container_read_many.go`

**Delete** `executeReadManyWithPointReads` entirely.

Add `executeReadManyWithQueries`:

```go
func (c *ContainerClient) executeReadManyWithQueries(
    ctx context.Context,
    items []ItemIdentity,
    readManyOptions *ReadManyOptions,
    operationContext pipelineRequestOptions,
) (ReadManyItemsResponse, error)
```

Implementation details:
- Calls `c.Read(ctx, nil)` to get `PartitionKeyDefinition`.
- Groups items by logical PK value into chunks.
- For each chunk, builds query via `queryBuilder`.
- Sends queries concurrently with `sendQueryRequest`, setting:
  - `x-ms-documentdb-partitionkey` header for PK-scoped routing.
  - `x-ms-documentdb-query-enablecrosspartition: true` when the query spans multiple PKs within a physical range (the OR-of-conjunctions case).
  - `x-ms-documentdb-partitionkeyrangeid` header for physical-range-scoped routing when EPK hashing is available.
- Handles continuation tokens (pagination) for each query.
- Re-sorts results into original input order.
- Returns aggregated `ReadManyItemsResponse`.

Internal helper types:

```go
type indexedItem struct {
    originalIndex int
    id            string
    pk            PartitionKey
}
```

### 4.4 Modified File: `cosmos_container.go`

In `ReadManyItems()` method, change the default strategy:

```go
// Current:
return c.executeReadManyWithPointReads(itemIdentities, readManyOptions, operationContext, ctx)

// Proposed:
return c.executeReadManyWithQueries(ctx, itemIdentities, readManyOptions, operationContext)
```

`executeReadManyWithPointReads` is deleted — there is no fallback to point reads.

### 4.5 Modified File: `cosmos_read_many_request_options.go`

No structural changes needed. The existing `ReadManyOptions` fields (SessionToken, ConsistencyLevel, DedicatedGatewayRequestOptions, MaxConcurrency) are already sufficient. The `QueryEngine` field continues to select the engine-based path.

### 4.6 Optional Future File: `cosmos_epk_hash.go`

If physical-range coalescing is desired (grouping multiple logical PKs that land on the same physical range into a single query), implement the V2 MurmurHash EPK computation. This is **not required for V1** but enables the full Python-equivalent optimisation.

---

## 5. Test Cases

All emulator tests follow the existing pattern in `emulator_cosmos_read_many_items_test.go`.

### 5.1 Unit Tests (`cosmos_query_builder_test.go`)

These do **not** require the emulator.

| # | Test Name | Description |
|---|---|---|
| U1 | `TestQueryBuilder_IDPartitionKey_Simple` | PK path is `/id`, 3 items where pk == id → produces `SELECT * FROM c WHERE c.id IN (…)` |
| U2 | `TestQueryBuilder_IDPartitionKey_Mismatch` | PK path is `/id` but pk ≠ id for one item → falls through to general query |
| U3 | `TestQueryBuilder_SingleLogicalPartition` | 5 items, all same PK value → produces `SELECT * FROM c WHERE c.myPk = @pk AND c.id IN (…)` |
| U4 | `TestQueryBuilder_MultipleLogicalPartitions` | 3 items, 2 different PKs → produces OR-of-conjunctions |
| U5 | `TestQueryBuilder_NestedPKPath` | PK path `/address/zipCode` → field expression `c["address"]["zipCode"]` |
| U6 | `TestQueryBuilder_NonIdentifierPKPath` | PK path `/my-pk` → field expression `c["my-pk"]` |
| U7 | `TestQueryBuilder_NullPartitionKey` | Item with `NullPartitionKey` → `IS_DEFINED(c.pk) = false` clause |
| U8 | `TestQueryBuilder_HierarchicalPK` | Multi-path PK definition → multiple conditions per item in OR clause |
| U9 | `TestQueryBuilder_EmptyItems` | Empty items list → `isIDPartitionKeyQuery` returns false, `isSingleLogicalPartitionQuery` returns false |
| U10 | `TestQueryBuilder_SingleItem` | Single item → `isSingleLogicalPartitionQuery` returns false (handled by point-read path) |
| U11 | `TestQueryBuilder_GetFieldExpression` | Various paths: `/pk`, `/a/b/c`, `/non-ident` |

### 5.2 Emulator Integration Tests (`emulator_cosmos_read_many_items_test.go`)

These require the Cosmos DB emulator. The existing tests (`TestReadMany_NilItemsSlice`, `TestReadMany_ReadSeveralItems`, `TestReadMany_NilIDReturnsError`, `TestReadMany_PartialFailure`, `TestReadMany_WithQueryEngine_*`) already cover the basic happy-path, empty-input, validation, and partial-miss scenarios. They serve as **regression tests** and must continue to pass after the strategy change.

One new integration test is added to exercise the fan-out behavior that distinguishes the query-based approach from point reads:

| # | Test Name | Setup | Action | Assertions |
|---|---|---|---|---|
| E1 | `TestReadManyWithQueries_MultipleLogicalPKs_SamePhysicalRange` | Container with `/pk` PK. Insert several items whose distinct logical PK values hash to the **same** physical partition key range (use the emulator's single-range default, or pick PK values known to collide). | ReadMany all items in a single call. | All items returned, correct content. Verifies: (a) in V1 (per-logical-PK routing), each logical PK produces its own query and all succeed; (b) in V2 (per-EPK-range routing), the items are coalesced into a single OR-of-conjunctions query targeting the physical range, and all still succeed. Positive `RequestCharge`. |

---

## 6. Implementation Notes

### 6.1 Query Routing Headers

For the per-partition query approach, each query request must include:

- `x-ms-documentdb-partitionkey`: The JSON-serialized PK value (e.g., `["myPkValue"]`). This lets the Gateway route the query to the correct physical partition without client-side EPK hashing.
- `x-ms-documentdb-query-enablecrosspartition: true`: Required for the OR-of-conjunctions case where items span multiple logical PKs.
- `x-ms-documentdb-partitionkeyrangeid`: Alternative to the PK header; directly specifies the physical range. Use this if/when EPK hashing is implemented.

### 6.2 Continuation Token Handling

Each `sendQueryRequest` may return a continuation token. The worker must loop, re-issuing the query with the continuation, until it is empty. This is especially important for large chunks near the 1000-item limit.

### 6.3 Concurrency Model

Reuse the existing `determineConcurrency` function and the goroutine-pool pattern from `runEngineRequests`.

### 6.4 Result Matching

Query results are matched back to the original input by the `id` field in the returned JSON document. The `originalIndex` is used to restore input order.

### 6.5 RU Cost Benefit

Per the Cosmos DB pricing model, a query returning N items costs less than N individual point reads when N > ~3, especially for small documents. The batched query approach also reduces HTTP round-trips, improving latency.

---

## 7. Open Questions

1. **V1 routing strategy:** Should V1 use `x-ms-documentdb-partitionkey` (simpler, no EPK hashing needed, but one query per logical PK) or invest in EPK hashing upfront (one query per physical range, matching Python)?
2. **Duplicate item IDs:** If the same `(id, pk)` appears twice in the input, should the output contain one or two copies?
3. **Order guarantee:** The Python SDK preserves input order. Should Go do the same? (Recommendation: yes.)
4. **Max items per query:** Is 1000 the correct limit for Go, or should it be configurable?
