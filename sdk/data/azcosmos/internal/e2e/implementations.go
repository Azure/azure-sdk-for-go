// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package e2e

// buildImplementations returns the Go implementation map for every scenario ID in the pinned
// catalog (see catalog/README.md for the pinned revision). It is kept in its own file because it
// is expected to grow and shrink as Go capabilities land; see issue #27654 for the tracking
// issue and #27690/#27655 for related follow-up work this map defers to.
func buildImplementations() []Implementation {
	const (
		noManagementCapabilitiesEndpoint = "the pinned hosted emulator now exposes a versioned GET /capabilities document, but Go has not yet ported the scenario's capability-document assertions; this remains a scenario-implementation gap, not a missing emulator endpoint."
		noChangeFeedAPI                  = "change feed is not yet implemented in azcosmos/v2; tracked by #27690."
		noFaultInjection                 = "requires deterministic hosted-emulator fault-injection controls; the pinned emulator exposes topology management actions but no fault-injection capability or endpoint."
		noTopologyControl                = "the pinned emulator exposes region/partition topology controls, but Go has no topology-control client and this implementation pass does not port topology scenarios."
		deferredScope                    = "in scope for this scenario catalog but deferred past this implementation pass for time; no missing Go capability is known."
	)

	return []Implementation{
		// --- Active: implemented against the public azcosmos/v2 API and the pinned emulator ---
		{ScenarioID: "bootstrap.primary-success", Status: StatusActive, Test: "TestBootstrapPrimarySuccess"},
		// item.lifecycle declares three profiles (smokeTests, lifecycleConsistencyMatrix,
		// readConsistencyOverrideMatrix); only the first (its resolved default) is certified
		// here. The other two need a multi-consistency-level account matrix this fixture cannot
		// build (same gap as consistency.feed-read-strategies below), so Active here means the
		// smokeTests-profile CRUD round trip is certified, not the full declared profile matrix.
		{ScenarioID: "item.lifecycle", Status: StatusActive, Test: "TestItemLifecycleSmoke"},
		{ScenarioID: "item.upsert-create-update", Status: StatusActive, Test: "TestItemUpsertCreateThenUpdate"},
		{ScenarioID: "item.create-conflict", Status: StatusActive, Test: "TestItemCreateConflictPreservesOriginal"},
		{ScenarioID: "item.not-found-wrong-partition-key", Status: StatusActive, Test: "TestItemNotFoundDoesNotCrossPartitionKeys"},
		{ScenarioID: "item.optimistic-concurrency", Status: StatusActive, Test: "TestItemOptimisticConcurrencyStaleETag"},
		{ScenarioID: "item.hierarchical-partition-key", Status: StatusActive, Test: "TestItemHierarchicalPartitionKeyPointAndPrefix"},
		{ScenarioID: "batch.atomicity", Status: StatusActive, Test: "TestBatchAtomicity"},
		{ScenarioID: "item.patch-state", Status: StatusActive, Test: "TestItemPatchState"},
		{ScenarioID: "item.validation-contracts", Status: StatusActive, Test: "TestItemValidationContracts"},
		{ScenarioID: "item.numeric-unique-key-equivalence", Status: StatusActive, Test: "TestItemNumericUniqueKeyEquivalence"},
		{ScenarioID: "query.parameterized-filter", Status: StatusActive, Test: "TestQueryParameterizedFilterAndOrder"},
		{ScenarioID: "query.invalid-syntax", Status: StatusActive, Test: "TestQueryInvalidSyntaxIsNotAnEmptyFeed"},
		{ScenarioID: "query.pagination-resume", Status: StatusActive, Test: "TestQueryPaginationResumesWithoutLossOrDuplication"},
		{ScenarioID: "diagnostics.success-and-error", Status: StatusActive, Test: "TestDiagnosticsCoverSuccessAndError"},

		// --- Blocked: missing Go product capability ---
		{ScenarioID: "item.scalar-partition-keys", Status: StatusBlocked, Reason: "the pinned in-memory emulator's header parser (azure_data_cosmos_driver::in_memory_emulator::epk::json_to_pk_component) rejects any JSON-object-shaped partition-key component with \"partition key components must be scalar\" - including the driver's own \"{}\" wire encoding for an Undefined partition key - so the catalog's required undefined-key case cannot be exercised through any client against this emulator build, not just Go. A TestItemScalarPartitionKeysRemainDistinct test still exists exercising the other four kinds (string/numeric/boolean/null) as supplementary regression coverage, but is not registered as certification evidence for this scenario; bumping the emulator pin is a separate decision (see catalog/README.md)."},
		{ScenarioID: "consistency.session-management", Status: StatusBlocked, Reason: "the catalog requires disabling automatic session-token capture and exercising continuation-based queries alongside the explicit token; azcosmos/v2 has no public option to disable automatic session capture (see OperationOptions/ClientOptions), so the full contract cannot be certified through the public API. A TestSessionTokenExplicitManagement test still exists verifying the narrower claim that an explicitly supplied session token is honored, but is not registered as certification evidence for this scenario."},
		{ScenarioID: "management.capabilities", Status: StatusBlocked, Reason: noManagementCapabilitiesEndpoint},
		// #27700 added azcosmos/v2's public database/container management API (Create/Read/Delete
		// for both), but not Replace or a List/Query API for either resource; the catalog's
		// requirement is "create, read, query, replace, and delete", so the full contract still
		// cannot be certified. TestDatabaseContainerManagementAvailableOperations exercises the
		// three operations that do exist as supplementary regression coverage, but is not
		// registered as certification evidence for this scenario.
		{ScenarioID: "management.resource-lifecycle", Status: StatusBlocked, Reason: "azcosmos/v2 (as of #27700) has Create/Read/Delete for databases and containers, but no Replace or List/Query API for either; the catalog requires create, read, query, replace, and delete, so the full contract is not yet certifiable through the public API."},
		{ScenarioID: "bootstrap.backup-fallback", Status: StatusBlocked, Reason: "NewClient/NewClientWithKey accept a single endpoint; azcosmos/v2 has no multi-endpoint/backup-bootstrap API yet."},
		{ScenarioID: "changefeed.pagination-resume", Status: StatusBlocked, Reason: noChangeFeedAPI},
		{ScenarioID: "changefeed.all-versions-start-validation", Status: StatusBlocked, Reason: noChangeFeedAPI},
		{ScenarioID: "diagnostics.handlers-telemetry", Status: StatusBlocked, Reason: "requires confirming azcosmos/v2 exposes diagnostics-handler/OpenTelemetry hooks equivalent to Rust's; not yet verified to exist for v2's native-driver diagnostics surface."},
		{ScenarioID: "query.feed-ranges", Status: StatusBlocked, Reason: "azcosmos/v2's ContainerClient has no public feed-range API (no equivalent of Rust's read_feed_ranges/feed_range_from_partition_key); only NewQueryItemsPager's partition-key/full-container scopes exist."},
		{ScenarioID: "consistency.response-token-capture", Status: StatusBlocked, Reason: "Rust's scenario pauses West US replication and routes with multi-region PreferredRegions to force an Eventual-then-Session switch. The pinned emulator and fixture can provide that account profile, but Go does not yet drive the emulator's replication controls for the scenario."},

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
		{ScenarioID: "item.quoted-partition-key-paths", Status: StatusBlocked, Reason: "requires a container whose partition-key path itself contains characters needing quoting; not yet provisioned in the emulator fixture; " + deferredScope},
		{ScenarioID: "configuration.binary-routing", Status: StatusBlocked, Reason: "requires verifying azcosmos/v2's binary-encoding override surface against the shared scenario's exact combinations; " + deferredScope},
		{ScenarioID: "consistency.feed-read-strategies", Status: StatusBlocked, Reason: "requires the lifecycleConsistencyMatrix account set across all five consistency levels; " + deferredScope},
		{ScenarioID: "consistency.session-staleness", Status: StatusBlocked, Reason: "requires disabling automatic session-token capture and observing a deliberately stale replica under controlled replication delay; " + deferredScope},
	}
}
