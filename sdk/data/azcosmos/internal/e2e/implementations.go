// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

// buildImplementations returns the Go implementation map for every scenario ID in the pinned
// catalog (see catalog/README.md for the pinned revision). It is kept in its own file because it
// is expected to grow and shrink as Go capabilities land; see issue #27654 for the tracking
// issue and #27690/#27655 for related follow-up work this map defers to.
func buildImplementations() []Implementation {
	const (
		noManagementAPI   = "azcosmos/v2 has no public database/container management API yet; NewDatabase/NewContainer only build local handles against pre-provisioned emulator resources."
		noChangeFeedAPI   = "change feed is not yet implemented in azcosmos/v2; tracked by #27690."
		noBatchAPI        = "transactional batch is not yet implemented in azcosmos/v2; tracked by #27690."
		noFaultInjection  = "requires the hosted emulator's deterministic fault-injection management API. The currently pinned EmulatorSourceRevision (a05f024d, 2026-09-02) predates the e2e_tests catalog (2026-09-11) and has no such API; bumping the pin is a separate decision (see catalog/README.md)."
		noTopologyControl = "requires the hosted emulator's region/partition topology-transition management API, which the currently pinned emulator revision does not expose (same gap as resilience.*). Dynamic topology is also explicitly out of this issue's initial-implementation scope."
		deferredScope     = "in scope for this scenario catalog but deferred past this implementation pass for time; no missing Go capability is known."
	)

	return []Implementation{
		// --- Active: implemented against the public azcosmos/v2 API and the pinned emulator ---
		{ScenarioID: "bootstrap.primary-success", Status: StatusActive, Test: "TestBootstrapPrimarySuccess"},
		{ScenarioID: "item.lifecycle", Status: StatusActive, Test: "TestItemLifecycleSmoke"},
		{ScenarioID: "item.upsert-create-update", Status: StatusActive, Test: "TestItemUpsertCreateThenUpdate"},
		{ScenarioID: "item.create-conflict", Status: StatusActive, Test: "TestItemCreateConflictPreservesOriginal"},
		{ScenarioID: "item.not-found-wrong-partition-key", Status: StatusActive, Test: "TestItemNotFoundDoesNotCrossPartitionKeys"},
		{ScenarioID: "item.optimistic-concurrency", Status: StatusActive, Test: "TestItemOptimisticConcurrencyStaleETag"},
		{ScenarioID: "item.scalar-partition-keys", Status: StatusActive, Test: "TestItemScalarPartitionKeysRemainDistinct"},
		{ScenarioID: "item.hierarchical-partition-key", Status: StatusActive, Test: "TestItemHierarchicalPartitionKeyPointAndPrefix"},
		{ScenarioID: "query.parameterized-filter", Status: StatusActive, Test: "TestQueryParameterizedFilterAndOrder"},
		{ScenarioID: "query.invalid-syntax", Status: StatusActive, Test: "TestQueryInvalidSyntaxIsNotAnEmptyFeed"},
		{ScenarioID: "query.pagination-resume", Status: StatusActive, Test: "TestQueryPaginationResumesWithoutLossOrDuplication"},
		{ScenarioID: "diagnostics.success-and-error", Status: StatusActive, Test: "TestDiagnosticsCoverSuccessAndError"},
		{ScenarioID: "consistency.session-management", Status: StatusActive, Test: "TestSessionTokenExplicitManagement"},

		// --- Blocked: missing Go product capability ---
		{ScenarioID: "management.capabilities", Status: StatusBlocked, Reason: noManagementAPI},
		{ScenarioID: "management.resource-lifecycle", Status: StatusBlocked, Reason: noManagementAPI},
		{ScenarioID: "bootstrap.backup-fallback", Status: StatusBlocked, Reason: "NewClient/NewClientWithKey accept a single endpoint; azcosmos/v2 has no multi-endpoint/backup-bootstrap API yet."},
		{ScenarioID: "batch.atomicity", Status: StatusBlocked, Reason: noBatchAPI},
		{ScenarioID: "changefeed.pagination-resume", Status: StatusBlocked, Reason: noChangeFeedAPI},
		{ScenarioID: "changefeed.all-versions-start-validation", Status: StatusBlocked, Reason: noChangeFeedAPI},
		{ScenarioID: "diagnostics.handlers-telemetry", Status: StatusBlocked, Reason: "requires confirming azcosmos/v2 exposes diagnostics-handler/OpenTelemetry hooks equivalent to Rust's; not yet verified to exist for v2's native-driver diagnostics surface."},
		{ScenarioID: "query.feed-ranges", Status: StatusBlocked, Reason: "azcosmos/v2's ContainerClient has no public feed-range API (no equivalent of Rust's read_feed_ranges/feed_range_from_partition_key); only NewQueryItemsPager's partition-key/full-container scopes exist."},
		{ScenarioID: "consistency.response-token-capture", Status: StatusBlocked, Reason: "Rust's scenario pauses West US replication and routes with multi-region PreferredRegions to force an Eventual-then-Session switch; the pinned emulator is single-region with no replication-pause management API, so this cannot be arranged honestly. " + noFaultInjection},

		// --- Blocked: requires hosted-emulator capability beyond the currently pinned revision ---
		{ScenarioID: "resilience.deadline", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.hedge-deadline", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.hedging-topology", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.hedging", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.partition-circuit-breaker", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.partition-topology-retry", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.request-timeout-retry", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.service-retry", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.throttling-retry", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "resilience.transport-retry", Status: StatusBlocked, Reason: noFaultInjection},
		{ScenarioID: "topology.continuation-transitions", Status: StatusBlocked, Reason: noTopologyControl},
		{ScenarioID: "topology.partition-split-merge", Status: StatusBlocked, Reason: noTopologyControl},
		{ScenarioID: "topology.region-lifecycle", Status: StatusBlocked, Reason: noTopologyControl},
		{ScenarioID: "topology.replication-pause-resume", Status: StatusBlocked, Reason: noTopologyControl},
		{ScenarioID: "topology.write-failover", Status: StatusBlocked, Reason: noTopologyControl},

		// --- Blocked: in scope, deferred this pass (no missing capability, just not yet ported) ---
		{ScenarioID: "item.patch-state", Status: StatusBlocked, Reason: "PatchItem exists publicly, but the shared scenario's exact post-image/patch-strategy assertions have not yet been verified against Go's patch semantics; " + deferredScope},
		{ScenarioID: "item.validation-contracts", Status: StatusBlocked, Reason: "requires auditing every shared validation-contract case against Go's argument-validation error codes one by one; " + deferredScope},
		{ScenarioID: "item.numeric-unique-key-equivalence", Status: StatusBlocked, Reason: "requires a container with a unique-key policy, which needs container management (see management.* above) or a dedicated emulator-config fixture; " + deferredScope},
		{ScenarioID: "item.quoted-partition-key-paths", Status: StatusBlocked, Reason: "requires a container whose partition-key path itself contains characters needing quoting; not yet provisioned in the emulator fixture; " + deferredScope},
		{ScenarioID: "configuration.binary-routing", Status: StatusBlocked, Reason: "requires verifying azcosmos/v2's binary-encoding override surface against the shared scenario's exact combinations; " + deferredScope},
		{ScenarioID: "consistency.feed-read-strategies", Status: StatusBlocked, Reason: "requires the lifecycleConsistencyMatrix account set across all five consistency levels; " + deferredScope},
		{ScenarioID: "consistency.session-staleness", Status: StatusBlocked, Reason: "requires disabling automatic session-token capture and observing a deliberately stale replica under controlled replication delay; " + deferredScope},
	}
}
