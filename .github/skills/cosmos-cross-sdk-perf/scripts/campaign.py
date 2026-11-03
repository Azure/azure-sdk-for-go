#!/usr/bin/env python3
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
# cSpell:words uprobes
# cSpell:words azsdk elif isinstance rglob rollup rustc shutil smaps vcpu

"""Lifecycle, validation, accounting, and reporting for the Cosmos perf campaign."""

from __future__ import annotations

import argparse
import csv
import hashlib
import json
import os
import pathlib
import re
import shutil
import subprocess
import sys
import tarfile
import tempfile
from dataclasses import dataclass
from typing import Any


SKILL_ROOT = pathlib.Path(__file__).resolve().parents[1]
CLASSIFIER_PATH = SKILL_ROOT / "classifiers" / "cpu-categories.json"
LEGS = ("go-v1", "go-v2-native", "rust")
RUN_CLASSES = ("clean-score", "perf", "go-diagnostic")
DIAGNOSTICS_FIX = "7201f4b70fac1184e66379278566f1ac8d299109"
ROLE_ID = "00000000-0000-0000-0000-000000000002"
OWNERSHIP_TAGS = {
    "azsdk-cosmos-perf": "true",
    "azsdk-cosmos-perf-owner": "campaign",
}
SENSITIVE_KEYS = re.compile(
    r"(token|secret|password|key|connection|string|endpoint|subscription|tenant|client|"
    r"resource.?id|account|database|container|public.?ip|private.?ip)",
    re.IGNORECASE,
)


class CampaignError(RuntimeError):
    """A user-actionable campaign failure."""


def load_json(path: pathlib.Path) -> dict[str, Any]:
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except FileNotFoundError as exc:
        raise CampaignError(f"missing file: {path}") from exc
    except json.JSONDecodeError as exc:
        raise CampaignError(f"invalid JSON in {path}: {exc}") from exc
    if not isinstance(value, dict):
        raise CampaignError(f"{path} must contain a JSON object")
    return value


def write_json(path: pathlib.Path, value: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def required(value: dict[str, Any], path: str) -> Any:
    current: Any = value
    for component in path.split("."):
        if not isinstance(current, dict) or component not in current:
            raise CampaignError(f"missing required configuration field: {path}")
        current = current[component]
    return current


def validate_config(config: dict[str, Any], allow_incomplete_azure: bool = False) -> None:
    azure = required(config, "azure")
    workload = required(config, "workload")
    sources = required(config, "sources")
    execution = required(config, "execution")
    profiling = required(config, "profiling")

    expected = {
        "partitionKeyPath": "/pk",
        "seedCount": 1024,
        "documentBytes": 1024,
        "operation": "point-read",
        "workers": [8, 16],
        "warmupSeconds": 30,
        "measurementSeconds": 300,
        "closedLoop": True,
        "keySequence": "round-robin-0-through-1023",
        "rawResponseMode": True,
        "validationOutsideTimedPath": True,
        "gatewayProtocol": "GatewayV1HTTP2",
        "binaryEncoding": False,
        "authentication": "account-key",
        "requiredRetries": 0,
        "requiredThrottles": 0,
    }
    for key, expected_value in expected.items():
        if workload.get(key) != expected_value:
            raise CampaignError(
                f"workload.{key} must be {expected_value!r} for the governed comparison; "
                f"got {workload.get(key)!r}"
            )

    if profiling.get("perfFrequencyHz") != 49:
        raise CampaignError("profiling.perfFrequencyHz must be 49")
    if profiling.get("capacityCores") != 4:
        raise CampaignError("profiling.capacityCores must be 4")
    if profiling.get("calibrationModes") != ["frame-pointer", "dwarf,16384"]:
        raise CampaignError(
            "profiling.calibrationModes must compare frame-pointer and dwarf,16384"
        )
    if profiling.get("selectedCallGraph") not in ("", "frame-pointer", "dwarf,16384"):
        raise CampaignError("profiling.selectedCallGraph has an unsupported value")
    tolerance = profiling.get("accountingTolerancePercent")
    if not isinstance(tolerance, (int, float)) or tolerance <= 0 or tolerance > 1:
        raise CampaignError("profiling.accountingTolerancePercent must be in (0, 1]")
    if execution.get("cleanRepeats", 0) < 2:
        raise CampaignError("execution.cleanRepeats must be at least 2 for A/A variance")
    if execution.get("runClasses") != list(RUN_CLASSES):
        raise CampaignError(f"execution.runClasses must be {list(RUN_CLASSES)!r}")
    if execution.get("authoritativeClass") != "clean-score":
        raise CampaignError("execution.authoritativeClass must be clean-score")
    if execution.get("stopAfterPilot") is not True:
        raise CampaignError("execution.stopAfterPilot must be true")
    expected_order = {
        "8": ["go-v1", "go-v2-native", "rust"],
        "16": ["rust", "go-v2-native", "go-v1"],
    }
    if execution.get("order") != expected_order:
        raise CampaignError(f"execution.order must be {expected_order!r}")

    if sources.get("goV1Version") not in ("", "v1.5.0"):
        raise CampaignError("sources.goV1Version must be v1.5.0 for the baseline")
    if sources.get("goV2Commit") != DIAGNOSTICS_FIX:
        raise CampaignError(
            "sources.goV2Commit must pin the verified diagnostics fix "
            f"{DIAGNOSTICS_FIX}, or this skill must be updated after verifying a containing commit"
        )
    for name in (
        "goV2Commit",
        "nativeDriverSourceCommit",
        "nativeDriverDistributionCommit",
        "rustCommit",
        "flameGraphCommit",
    ):
        ref = str(sources.get(name, ""))
        if not re.fullmatch(r"[0-9a-f]{40}", ref):
            raise CampaignError(f"sources.{name} must be a full 40-character commit SHA")
    if sources.get("flameGraphRepository") != "https://github.com/brendangregg/FlameGraph.git":
        raise CampaignError(
            "sources.flameGraphRepository must use Brendan Gregg's canonical FlameGraph repository"
        )

    if azure.get("cleanup") not in ("retain-vm", "all-campaign-owned"):
        raise CampaignError("azure.cleanup must be retain-vm or all-campaign-owned")
    provision_account = bool(azure.get("provisionCosmosAccount"))
    existing_account = str(azure.get("existingCosmosAccountResourceId", ""))
    if provision_account == bool(existing_account):
        raise CampaignError(
            "choose exactly one of azure.provisionCosmosAccount=true or "
            "azure.existingCosmosAccountResourceId"
        )
    if provision_account and not azure.get("cosmosAccountName"):
        raise CampaignError("azure.cosmosAccountName is required when provisioning an account")
    if not allow_incomplete_azure:
        for key in ("subscriptionId", "resourceGroup"):
            if not azure.get(key):
                raise CampaignError(f"azure.{key} is required")
        if provision_account and not azure.get("location"):
            raise CampaignError("azure.location is required when provisioning an account")


def run(command: list[str], *, capture: bool = True, check: bool = True) -> str:
    result = subprocess.run(
        command,
        check=False,
        text=True,
        stdout=subprocess.PIPE if capture else None,
        stderr=subprocess.PIPE if capture else None,
    )
    if check and result.returncode != 0:
        detail = (result.stderr or result.stdout or "").strip()
        raise CampaignError(f"command failed ({' '.join(command)}): {detail}")
    return (result.stdout or "").strip()


def az_json(arguments: list[str]) -> Any:
    output = run(["az", *arguments, "--only-show-errors", "--output", "json"])
    try:
        return json.loads(output)
    except json.JSONDecodeError as exc:
        raise CampaignError(f"Azure CLI returned invalid JSON for {' '.join(arguments)}") from exc


def prerequisites(phase: str) -> None:
    tools = ("az", "git", "python3")
    if phase == "vm":
        tools += ("perf", "readelf", "sha256sum", "tar", "rustc", "cargo", "go")
    missing = [tool for tool in tools if shutil.which(tool) is None]
    if missing:
        raise CampaignError("missing prerequisites: " + ", ".join(missing))
    account = az_json(["account", "show"])
    if not account.get("id"):
        raise CampaignError("Azure CLI is not signed in to a subscription")
    print(f"{phase} prerequisites satisfied.")


def normalize_location(value: str) -> str:
    return re.sub(r"[^a-z0-9]", "", value.lower())


def account_from_id(resource_id: str) -> tuple[str, str]:
    match = re.fullmatch(
        r"/subscriptions/[^/]+/resourceGroups/([^/]+)/providers/"
        r"Microsoft\.DocumentDB/databaseAccounts/([^/]+)",
        resource_id,
        re.IGNORECASE,
    )
    if not match:
        raise CampaignError("existing Cosmos account resource ID has an unexpected format")
    return match.group(1), match.group(2)


def choose_four_vcpu_sku(location: str, requested: str) -> str:
    candidates = az_json(
        [
            "vm",
            "list-skus",
            "--location",
            location,
            "--resource-type",
            "virtualMachines",
            "--all",
        ]
    )
    available: list[str] = []
    for sku in candidates:
        restrictions = sku.get("restrictions") or []
        if restrictions:
            continue
        capabilities = {
            capability.get("name"): capability.get("value")
            for capability in sku.get("capabilities") or []
        }
        if capabilities.get("vCPUs") == "4":
            available.append(str(sku.get("name")))
    if requested:
        if requested not in available:
            raise CampaignError(f"requested four-vCPU SKU {requested!r} is unavailable in {location}")
        return requested
    for preferred in ("Standard_D4s_v5", "Standard_D4as_v5", "Standard_D4s_v4"):
        if preferred in available:
            return preferred
    if not available:
        raise CampaignError(f"no unrestricted four-vCPU VM SKU is available in {location}")
    return sorted(available)[0]


def provision_plan(config: dict[str, Any]) -> dict[str, Any]:
    azure = config["azure"]
    run(["az", "account", "set", "--subscription", azure["subscriptionId"]])
    group_exists = run(
        [
            "az",
            "group",
            "exists",
            "--name",
            azure["resourceGroup"],
            "--output",
            "tsv",
        ]
    ).lower()
    if group_exists not in ("true", "false"):
        raise CampaignError("could not determine whether the campaign resource group exists")
    if group_exists == "true":
        raise CampaignError(
            "azure.resourceGroup must be a new, dedicated campaign resource group so teardown "
            "cannot delete pre-existing resources"
        )
    if azure["provisionCosmosAccount"]:
        location = azure["location"]
        account_rg = azure["resourceGroup"]
        account_name = azure["cosmosAccountName"]
        account_id = (
            f"/subscriptions/{azure['subscriptionId']}/resourceGroups/{account_rg}/providers/"
            f"Microsoft.DocumentDB/databaseAccounts/{account_name}"
        )
        account_owned = True
    else:
        account_id = azure["existingCosmosAccountResourceId"]
        account_rg, account_name = account_from_id(account_id)
        account = az_json(["cosmosdb", "show", "--ids", account_id])
        locations = account.get("locations") or []
        writable = [
            entry.get("locationName")
            for entry in locations
            if entry.get("failoverPriority") == 0
        ]
        location = writable[0] if writable else account.get("location")
        if not location:
            raise CampaignError("could not determine the Cosmos account write region")
        account_owned = False
    if azure.get("location") and normalize_location(azure["location"]) != normalize_location(location):
        raise CampaignError(
            f"configured location {azure['location']!r} does not match Cosmos write region {location!r}"
        )
    sku = choose_four_vcpu_sku(location, azure.get("vmSku", ""))
    return {
        "subscriptionId": azure["subscriptionId"],
        "location": location,
        "resourceGroup": azure["resourceGroup"],
        "vmName": azure["vmName"],
        "vmSku": sku,
        "adminUsername": azure["adminUsername"],
        "cosmosAccount": {
            "id": account_id,
            "resourceGroup": account_rg,
            "name": account_name,
            "owned": account_owned,
        },
        "database": azure["database"],
        "container": azure["container"],
        "provisionDatabaseContainer": bool(azure["provisionDatabaseContainer"]),
        "cleanup": azure["cleanup"],
        "resourceGroupOwned": True,
        "tags": OWNERSHIP_TAGS,
    }


def provision(config: dict[str, Any], state_path: pathlib.Path, dry_run: bool) -> None:
    validate_config(config)
    plan = provision_plan(config)
    print(json.dumps(plan, indent=2, sort_keys=True))
    if dry_run:
        return
    azure = config["azure"]
    tag_args = [f"{key}={value}" for key, value in OWNERSHIP_TAGS.items()]
    az_json(
        [
            "group",
            "create",
            "--name",
            plan["resourceGroup"],
            "--location",
            plan["location"],
            "--tags",
            *tag_args,
        ]
    )
    if plan["cosmosAccount"]["owned"]:
        az_json(
            [
                "cosmosdb",
                "create",
                "--name",
                plan["cosmosAccount"]["name"],
                "--resource-group",
                plan["resourceGroup"],
                "--locations",
                f"regionName={plan['location']}",
                "--default-consistency-level",
                "Session",
                "--tags",
                *tag_args,
            ]
        )
    vm = az_json(
        [
            "vm",
            "create",
            "--resource-group",
            plan["resourceGroup"],
            "--name",
            plan["vmName"],
            "--image",
            "Ubuntu2204",
            "--size",
            plan["vmSku"],
            "--admin-username",
            plan["adminUsername"],
            "--authentication-type",
            "ssh",
            "--generate-ssh-keys",
            "--assign-identity",
            "--public-ip-sku",
            "Standard",
            "--tags",
            *tag_args,
        ]
    )
    principal_id = required(vm, "identity.systemAssignedIdentity")
    az_json(
        [
            "cosmosdb",
            "sql",
            "role",
            "assignment",
            "create",
            "--account-name",
            plan["cosmosAccount"]["name"],
            "--resource-group",
            plan["cosmosAccount"]["resourceGroup"],
            "--scope",
            "/",
            "--principal-id",
            principal_id,
            "--role-definition-id",
            ROLE_ID,
        ]
    )
    if plan["provisionDatabaseContainer"]:
        az_json(
            [
                "cosmosdb",
                "sql",
                "database",
                "create",
                "--account-name",
                plan["cosmosAccount"]["name"],
                "--resource-group",
                plan["cosmosAccount"]["resourceGroup"],
                "--name",
                plan["database"],
            ]
        )
        az_json(
            [
                "cosmosdb",
                "sql",
                "container",
                "create",
                "--account-name",
                plan["cosmosAccount"]["name"],
                "--resource-group",
                plan["cosmosAccount"]["resourceGroup"],
                "--database-name",
                plan["database"],
                "--name",
                plan["container"],
                "--partition-key-path",
                "/pk",
                "--throughput",
                str(azure["throughput"]),
            ]
        )
    state = {
        **plan,
        "version": 1,
        "vmOwned": True,
        "resourceGroupOwned": plan["resourceGroupOwned"],
        "databaseContainerOwned": plan["provisionDatabaseContainer"],
        "principalId": principal_id,
    }
    write_json(state_path, state)
    print(f"Wrote ownership state to {state_path}")


def teardown(state: dict[str, Any], dry_run: bool) -> None:
    if state.get("tags") != OWNERSHIP_TAGS:
        raise CampaignError("ownership state tags do not match this skill")
    actions: list[list[str]] = []
    if state.get("cleanup") == "retain-vm":
        print(json.dumps({"actions": [], "reason": "cleanup preference retains the VM"}, indent=2))
        return
    if state.get("databaseContainerOwned") and not state["cosmosAccount"]["owned"]:
        actions.extend(
            [
                [
                    "az",
                    "cosmosdb",
                    "sql",
                    "container",
                    "delete",
                    "--yes",
                    "--account-name",
                    state["cosmosAccount"]["name"],
                    "--resource-group",
                    state["cosmosAccount"]["resourceGroup"],
                    "--database-name",
                    state["database"],
                    "--name",
                    state["container"],
                ],
                [
                    "az",
                    "cosmosdb",
                    "sql",
                    "database",
                    "delete",
                    "--yes",
                    "--account-name",
                    state["cosmosAccount"]["name"],
                    "--resource-group",
                    state["cosmosAccount"]["resourceGroup"],
                    "--name",
                    state["database"],
                ],
            ]
        )
    if state.get("resourceGroupOwned"):
        actions.append(
            [
                "az",
                "group",
                "delete",
                "--yes",
                "--no-wait",
                "--name",
                state["resourceGroup"],
            ]
        )
    print(json.dumps({"actions": actions}, indent=2))
    if dry_run:
        return
    for action in actions:
        run(action, capture=False)


@dataclass(frozen=True)
class PerfRow:
    symbol: str
    dso: str
    self_percent: float


def load_perf_rows(path: pathlib.Path) -> list[PerfRow]:
    rows: list[PerfRow] = []
    with path.open(newline="", encoding="utf-8") as handle:
        reader = csv.DictReader(handle)
        expected = {"symbol", "dso", "self_percent"}
        if set(reader.fieldnames or ()) != expected:
            raise CampaignError(f"{path} must have columns {sorted(expected)}")
        for index, row in enumerate(reader, start=2):
            try:
                percent = float(row["self_percent"])
            except ValueError as exc:
                raise CampaignError(f"{path}:{index}: invalid self_percent") from exc
            if percent < 0:
                raise CampaignError(f"{path}:{index}: self_percent cannot be negative")
            rows.append(PerfRow(row["symbol"], row["dso"], percent))
    if not rows:
        raise CampaignError(f"{path} has no perf rows")
    return rows


def classify(rows: list[PerfRow]) -> tuple[dict[str, float], list[str]]:
    document = load_json(CLASSIFIER_PATH)
    rules = [
        (entry["name"], [re.compile(pattern, re.IGNORECASE) for pattern in entry["patterns"]])
        for entry in document["categories"]
    ]
    totals = {name: 0.0 for name, _ in rules}
    unmapped: list[str] = []
    for row in rows:
        subject = f"{row.dso} {row.symbol}"
        for name, patterns in rules:
            if any(pattern.search(subject) for pattern in patterns):
                totals[name] += row.self_percent
                break
        else:
            unmapped.append(subject)
    return totals, unmapped


def compare_workload(metrics: dict[str, Any], config: dict[str, Any]) -> None:
    expected = config["workload"]
    actual = required(metrics, "workload")
    for key in (
        "partitionKeyPath",
        "seedCount",
        "documentBytes",
        "operation",
        "warmupSeconds",
        "measurementSeconds",
        "closedLoop",
        "keySequence",
        "rawResponseMode",
        "validationOutsideTimedPath",
        "gatewayProtocol",
        "binaryEncoding",
        "authentication",
    ):
        if actual.get(key) != expected.get(key):
            raise CampaignError(
                f"workload mismatch for {key}: expected {expected.get(key)!r}, "
                f"got {actual.get(key)!r}"
            )
    if actual.get("workers") != metrics.get("workers"):
        raise CampaignError("workload workers must match metrics workers")
    if actual.get("workers") not in expected.get("workers", []):
        raise CampaignError("workload workers are outside the configured pilot matrix")


def validate_metrics(metrics: dict[str, Any], leg: str) -> None:
    if metrics.get("leg") != leg:
        raise CampaignError(f"metrics leg must be {leg!r}")
    if metrics.get("runClass") not in RUN_CLASSES:
        raise CampaignError(f"metrics runClass must be one of {RUN_CLASSES}")
    if metrics.get("workers") not in (8, 16):
        raise CampaignError("metrics workers must be 8 or 16")
    if not isinstance(metrics.get("repeat"), int) or metrics["repeat"] < 1:
        raise CampaignError("metrics repeat must be a positive integer")
    proof = required(metrics, "transportProof")
    if proof.get("transportKind") != "gateway-v1":
        raise CampaignError("transport preflight did not prove classic Gateway V1")
    if proof.get("negotiatedHTTPVersion") != "HTTP/2":
        raise CampaignError("transport preflight did not prove HTTP/2")
    if proof.get("binaryEncoding") is not False:
        raise CampaignError("transport preflight did not prove binary encoding disabled")
    counts = required(metrics, "counts")
    for key in ("success", "error", "retry", "throttle"):
        if not isinstance(counts.get(key), int) or counts[key] < 0:
            raise CampaignError(f"counts.{key} must be a nonnegative integer")
    timing = required(metrics, "timing")
    if timing.get("wallSeconds", 0) <= 0 or timing.get("processCPUSeconds", -1) < 0:
        raise CampaignError("timing wall/process CPU seconds are invalid")
    for section, keys in (
        ("latencyMilliseconds", ("p50", "p90", "p95", "p99", "p999")),
        (
            "memoryBytes",
            ("rss", "hwm", "anonymous", "file", "cgroupCurrent", "cgroupPeak"),
        ),
        ("gauges", ("queueHighWater", "inFlightHighWater", "completionHighWater")),
    ):
        value = required(metrics, section)
        for key in keys:
            if key not in value:
                raise CampaignError(f"{section}.{key} is required")


def validate_provenance(
    provenance: dict[str, Any], leg: str, config: dict[str, Any], binary_path: pathlib.Path
) -> None:
    if provenance.get("leg") != leg:
        raise CampaignError(f"{leg}: provenance leg is missing or incorrect")
    source = required(provenance, "source")
    for key in ("repository", "commit", "tree"):
        if not source.get(key):
            raise CampaignError(f"{leg}: provenance source.{key} is required")
    if not re.fullmatch(r"[0-9a-f]{40}", source["commit"]):
        raise CampaignError(f"{leg}: source.commit must be a full commit SHA")
    if not re.fullmatch(r"[0-9a-f]{40}", source["tree"]):
        raise CampaignError(f"{leg}: source.tree must be a full tree SHA")
    if source.get("dirty") is not False:
        raise CampaignError(f"{leg}: source must be recorded as clean")
    for key in ("compiler", "buildCommand", "dependencyVersions"):
        if not provenance.get(key):
            raise CampaignError(f"{leg}: provenance.{key} is required")
    binary = required(provenance, "binary")
    if not re.fullmatch(r"[0-9a-f]{64}", str(binary.get("sha256", ""))):
        raise CampaignError(f"{leg}: binary.sha256 must be a SHA-256 digest")
    if not binary.get("buildId"):
        raise CampaignError(f"{leg}: binary.buildId is required")
    if binary.get("debugSymbols") not in ("embedded", "separate"):
        raise CampaignError(f"{leg}: binary.debugSymbols must be embedded or separate")
    actual_sha = hashlib.sha256(binary_path.read_bytes()).hexdigest()
    if actual_sha != binary["sha256"]:
        raise CampaignError(f"{leg}: recorded binary SHA-256 does not match raw/binary")

    expected_source = {
        "go-v2-native": config["sources"]["goV2Commit"],
        "rust": config["sources"]["rustCommit"],
    }.get(leg)
    if expected_source and source["commit"] != expected_source:
        raise CampaignError(f"{leg}: source commit does not match the effective configuration")
    if leg == "go-v1":
        versions = provenance["dependencyVersions"]
        if versions.get(config["sources"]["goV1Module"]) != config["sources"]["goV1Version"]:
            raise CampaignError(f"{leg}: Go v1 module version does not match the configuration")
    if leg == "go-v2-native":
        native = required(provenance, "nativeDriver")
        if native.get("sourceCommit") != config["sources"]["nativeDriverSourceCommit"]:
            raise CampaignError(f"{leg}: native-driver source commit does not match configuration")
        if (
            native.get("distributionCommit")
            != config["sources"]["nativeDriverDistributionCommit"]
        ):
            raise CampaignError(
                f"{leg}: native-driver distribution commit does not match configuration"
            )


def validate_leg(
    leg_dir: pathlib.Path, config: dict[str, Any], write_derived: bool = True
) -> dict[str, Any]:
    leg = leg_dir.name
    if leg not in LEGS:
        raise CampaignError(f"leg directory must be one of {LEGS}")
    metrics = load_json(leg_dir / "metrics.json")
    validate_metrics(metrics, leg)
    compare_workload(metrics, config)
    counts = metrics["counts"]
    if counts["success"] <= 0:
        raise CampaignError(f"{leg}: no successful operations")
    if counts["error"] != 0:
        raise CampaignError(f"{leg}: errors must be zero")
    if counts["retry"] != config["workload"]["requiredRetries"]:
        raise CampaignError(f"{leg}: retries must be zero")
    if counts["throttle"] != config["workload"]["requiredThrottles"]:
        raise CampaignError(f"{leg}: throttles must be zero")
    common_evidence = ["raw/binary", "provenance.json"]
    if metrics["runClass"] in ("clean-score", "perf"):
        common_evidence.extend(
            [
                "raw/proc-status.txt",
                "raw/smaps-rollup.txt",
                "raw/smaps.txt",
                "raw/cgroup-memory.txt",
                "raw/thread-stacks.txt",
            ]
        )
    for evidence in common_evidence:
        if not (leg_dir / evidence).is_file():
            raise CampaignError(f"{leg}: missing evidence {evidence}")
    run_class = metrics["runClass"]
    forbidden = {
        "clean-score": (
            "raw/perf.data",
            "raw/cpu.pprof",
            "raw/heap.pprof",
            "raw/runtime.trace",
        ),
        "perf": ("raw/cpu.pprof", "raw/heap.pprof", "raw/runtime.trace"),
        "go-diagnostic": ("raw/perf.data",),
    }[run_class]
    for evidence in forbidden:
        if (leg_dir / evidence).exists():
            raise CampaignError(f"{leg}: {run_class} run contains forbidden {evidence}")
    if run_class == "perf":
        for evidence in ("raw/perf.data", "raw/perf-quality.json", "perf-self.csv"):
            if not (leg_dir / evidence).is_file():
                raise CampaignError(f"{leg}: perf run missing {evidence}")
    if run_class == "go-diagnostic":
        if not leg.startswith("go-"):
            raise CampaignError("go-diagnostic runs are valid only for Go legs")
        for evidence in (
            "raw/cpu.pprof",
            "raw/heap.pprof",
            "raw/runtime.trace",
            "raw/block.pprof",
            "raw/mutex.pprof",
        ):
            if not (leg_dir / evidence).is_file():
                raise CampaignError(f"{leg}: Go diagnostic run missing {evidence}")
    validate_provenance(
        load_json(leg_dir / "provenance.json"), leg, config, leg_dir / "raw" / "binary"
    )

    wall = metrics["timing"]["wallSeconds"]
    process_cpu = metrics["timing"]["processCPUSeconds"]
    capacity = config["profiling"]["capacityCores"]
    busy_cores = process_cpu / wall
    capacity_percent = busy_cores / capacity * 100.0
    cpu_per_success_us = process_cpu * 1_000_000.0 / counts["success"]
    result = {
        "valid": True,
        "runClass": run_class,
        "processCPUSeconds": process_cpu,
        "busyCores": busy_cores,
        "fourVCPUCapacityPercent": capacity_percent,
        "cpuMicrosecondsPerSuccess": cpu_per_success_us,
    }
    if run_class != "perf":
        if write_derived:
            write_json(leg_dir / "derived" / "validation.json", result)
        return result

    rows = load_perf_rows(leg_dir / "perf-self.csv")
    categories, unmapped = classify(rows)
    if unmapped:
        raise CampaignError(f"{leg}: unmapped perf rows: {unmapped[:5]}")
    perf_total = sum(row.self_percent for row in rows)
    tolerance = config["profiling"]["accountingTolerancePercent"]
    if abs(perf_total - 100.0) > tolerance:
        raise CampaignError(
            f"{leg}: perf Self total {perf_total:.6f}% is outside 100+/-{tolerance}%"
        )
    category_total = sum(categories.values())
    if abs(category_total - perf_total) > 1e-9:
        raise CampaignError(f"{leg}: classified category total drifted from perf total")

    quality = load_json(leg_dir / "raw" / "perf-quality.json")
    if quality.get("lostEventPercent", 100) > config["profiling"]["maxLostEventPercent"]:
        raise CampaignError(f"{leg}: perf lost-event quality gate failed")
    if quality.get("resolvedSamplePercent", 0) < config["profiling"]["minResolvedSamplePercent"]:
        raise CampaignError(f"{leg}: perf resolved-sample quality gate failed")
    if quality.get("callchainSamplePercent", 0) < config["profiling"]["minCallchainSamplePercent"]:
        raise CampaignError(f"{leg}: perf callchain quality gate failed")
    if quality.get("callGraphMode") != config["profiling"]["selectedCallGraph"]:
        raise CampaignError(f"{leg}: perf call-graph mode differs from calibrated selection")
    bridged = {
        name: {
            "perfSelfPercent": percent,
            "processCPUSeconds": process_cpu * percent / 100.0,
            "fourVCPUCapacityPercent": capacity_percent * percent / 100.0,
        }
        for name, percent in categories.items()
    }
    result.update({
        "perfSelfTotalPercent": perf_total,
        "perfQuality": quality,
        "categories": bridged,
    })
    if write_derived:
        write_json(leg_dir / "derived" / "accounting.json", result)
    return result


def svg_bar_chart(
    title: str, labels: list[str], values: list[float], unit: str, path: pathlib.Path
) -> None:
    width, height = 900, 120 + len(labels) * 58
    maximum = max(values) if values else 1.0
    lines = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}">',
        '<rect width="100%" height="100%" fill="white"/>',
        f'<text x="20" y="32" font-family="sans-serif" font-size="22">{title}</text>',
    ]
    for index, (label, value) in enumerate(zip(labels, values)):
        y = 68 + index * 58
        bar = 650 * value / maximum if maximum else 0
        lines.extend(
            [
                f'<text x="20" y="{y + 18}" font-family="sans-serif" font-size="15">{label}</text>',
                f'<rect x="190" y="{y}" width="{bar:.2f}" height="24" fill="#0078d4"/>',
                f'<text x="{200 + bar:.2f}" y="{y + 18}" font-family="sans-serif" '
                f'font-size="14">{value:.3f} {unit}</text>',
            ]
        )
    lines.append("</svg>")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def svg_flywheel(accounting: dict[str, dict[str, Any]], path: pathlib.Path) -> None:
    colors = (
        "#0078d4",
        "#50e6ff",
        "#7fba00",
        "#ffb900",
        "#f25022",
        "#5c2d91",
        "#00b7c3",
        "#bad80a",
        "#e3008c",
        "#8a8886",
        "#4f6bed",
        "#a4262c",
    )
    categories = list(next(iter(accounting.values()))["categories"])
    width, height = 1200, 140 + len(accounting) * 80 + len(categories) * 24
    lines = [
        f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}">',
        '<rect width="100%" height="100%" fill="white"/>',
        '<text x="20" y="32" font-family="sans-serif" font-size="22">'
        "CPU attribution (perf Self composition; separate from Go pprof)</text>",
    ]
    y = 65
    for leg, result in accounting.items():
        lines.append(
            f'<text x="20" y="{y + 19}" font-family="sans-serif" font-size="15">{leg}</text>'
        )
        x = 180.0
        for index, category in enumerate(categories):
            value = result["categories"][category]["perfSelfPercent"]
            bar = 8.5 * value
            lines.append(
                f'<rect x="{x:.2f}" y="{y}" width="{bar:.2f}" height="28" '
                f'fill="{colors[index % len(colors)]}"/>'
            )
            x += bar
        y += 70
    y += 10
    for index, category in enumerate(categories):
        lines.extend(
            [
                f'<rect x="20" y="{y + index * 24}" width="16" height="16" '
                f'fill="{colors[index % len(colors)]}"/>',
                f'<text x="44" y="{y + 14 + index * 24}" font-family="sans-serif" '
                f'font-size="13">{category}</text>',
            ]
        )
    lines.append("</svg>")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def run_path(
    run_dir: pathlib.Path, workers: int, run_class: str, repeat: int, leg: str
) -> pathlib.Path:
    return run_dir / f"w{workers}" / run_class / f"repeat-{repeat}" / leg


def summarize_metrics(metrics: dict[str, Any], validation: dict[str, Any]) -> dict[str, Any]:
    success = metrics["counts"]["success"]
    wall = metrics["timing"]["wallSeconds"]
    return {
        "workers": metrics["workers"],
        "leg": metrics["leg"],
        "runClass": metrics["runClass"],
        "repeat": metrics["repeat"],
        "valid": True,
        "throughputOpsPerSecond": success / wall,
        "success": success,
        "error": metrics["counts"]["error"],
        "retry": metrics["counts"]["retry"],
        "throttle": metrics["counts"]["throttle"],
        "processCPUSeconds": metrics["timing"]["processCPUSeconds"],
        "busyCores": validation["busyCores"],
        "cpuMicrosecondsPerSuccess": validation["cpuMicrosecondsPerSuccess"],
        "latencyMilliseconds": metrics["latencyMilliseconds"],
        "memoryBytes": metrics["memoryBytes"],
        "transportProof": metrics["transportProof"],
    }


def arithmetic_mean(values: list[float]) -> float:
    return sum(values) / len(values)


def relative_range_percent(values: list[float]) -> float:
    mean = arithmetic_mean(values)
    return 0.0 if mean == 0 else (max(values) - min(values)) / mean * 100.0


def report(run_dir: pathlib.Path, config: dict[str, Any]) -> None:
    comparison_dir = run_dir / "comparison"
    comparison_dir.mkdir(parents=True, exist_ok=True)
    all_rows: list[dict[str, Any]] = []
    clean_groups: dict[tuple[int, str], list[dict[str, Any]]] = {}
    perf_accounting: dict[str, dict[str, Any]] = {}
    setup_rows: list[list[Any]] = []

    for workers in config["workload"]["workers"]:
        for repeat in range(1, config["execution"]["cleanRepeats"] + 1):
            for leg in config["execution"]["order"][str(workers)]:
                leg_dir = run_path(run_dir, workers, "clean-score", repeat, leg)
                validation = validate_leg(leg_dir, config)
                metrics = load_json(leg_dir / "metrics.json")
                row = summarize_metrics(metrics, validation)
                all_rows.append(row)
                clean_groups.setdefault((workers, leg), []).append(row)
                provenance = load_json(leg_dir / "provenance.json")
                setup_rows.append(
                    [
                        workers,
                        leg,
                        repeat,
                        provenance["source"]["commit"],
                        provenance["binary"]["sha256"],
                        metrics["transportProof"]["transportKind"],
                        metrics["transportProof"]["negotiatedHTTPVersion"],
                        metrics["workload"]["documentBytes"],
                        metrics["workload"]["seedCount"],
                        metrics["workload"]["keySequence"],
                    ]
                )
        for leg in config["execution"]["order"][str(workers)]:
            perf_dir = run_path(run_dir, workers, "perf", 1, leg)
            validation = validate_leg(perf_dir, config)
            metrics = load_json(perf_dir / "metrics.json")
            all_rows.append(summarize_metrics(metrics, validation))
            perf_accounting[f"w{workers}-{leg}"] = validation
        for leg in ("go-v1", "go-v2-native"):
            diagnostic_dir = run_path(run_dir, workers, "go-diagnostic", 1, leg)
            validation = validate_leg(diagnostic_dir, config)
            metrics = load_json(diagnostic_dir / "metrics.json")
            all_rows.append(summarize_metrics(metrics, validation))

    authoritative: list[dict[str, Any]] = []
    for (workers, leg), rows in sorted(clean_groups.items()):
        throughput = [row["throughputOpsPerSecond"] for row in rows]
        p99 = [row["latencyMilliseconds"]["p99"] for row in rows]
        cpu_per_read = [row["cpuMicrosecondsPerSuccess"] for row in rows]
        authoritative.append(
            {
                "workers": workers,
                "leg": leg,
                "runClass": "clean-score",
                "repeats": len(rows),
                "throughputOpsPerSecond": arithmetic_mean(throughput),
                "throughputAARelativeRangePercent": relative_range_percent(throughput),
                "latencyMilliseconds": {
                    key: arithmetic_mean(
                        [row["latencyMilliseconds"][key] for row in rows]
                    )
                    for key in ("p50", "p90", "p95", "p99", "p999")
                },
                "p99AARelativeRangePercent": relative_range_percent(p99),
                "processCPUSeconds": arithmetic_mean(
                    [row["processCPUSeconds"] for row in rows]
                ),
                "busyCores": arithmetic_mean([row["busyCores"] for row in rows]),
                "cpuMicrosecondsPerSuccess": arithmetic_mean(cpu_per_read),
                "cpuPerReadAARelativeRangePercent": relative_range_percent(cpu_per_read),
                "memoryBytes": {
                    key: arithmetic_mean([row["memoryBytes"][key] for row in rows])
                    for key in ("rss", "hwm", "anonymous", "file", "cgroupCurrent", "cgroupPeak")
                },
                "retry": 0,
                "throttle": 0,
                "valid": True,
            }
        )

    observer_effect: list[dict[str, Any]] = []
    for row in all_rows:
        if row["runClass"] == "clean-score":
            continue
        clean = next(
            item
            for item in authoritative
            if item["workers"] == row["workers"] and item["leg"] == row["leg"]
        )
        row["observerEffectPercent"] = {
            "throughput": (
                row["throughputOpsPerSecond"] / clean["throughputOpsPerSecond"] - 1
            )
            * 100.0,
            "cpuPerRead": (
                row["cpuMicrosecondsPerSuccess"] / clean["cpuMicrosecondsPerSuccess"] - 1
            )
            * 100.0,
            "p99": (
                row["latencyMilliseconds"]["p99"]
                / clean["latencyMilliseconds"]["p99"]
                - 1
            )
            * 100.0,
        }
        observer_effect.append(
            {
                "workers": row["workers"],
                "leg": row["leg"],
                "runClass": row["runClass"],
                **row["observerEffectPercent"],
            }
        )

    write_json(
        comparison_dir / "summary.json",
        {
            "authoritativeClass": "clean-score",
            "authoritative": authoritative,
            "observerEffect": observer_effect,
            "allRuns": all_rows,
        },
    )
    with (comparison_dir / "summary.csv").open("w", newline="", encoding="utf-8") as handle:
        writer = csv.writer(handle)
        writer.writerow(
            [
                "workers",
                "leg",
                "throughput_ops_s",
                "throughput_aa_range_pct",
                "p50_ms",
                "p95_ms",
                "p99_ms",
                "p99_aa_range_pct",
                "cpu_seconds",
                "busy_cores",
                "cpu_us_read",
                "cpu_aa_range_pct",
                "rss_bytes",
                "hwm_bytes",
                "valid",
            ]
        )
        for row in authoritative:
            writer.writerow(
                [
                    row["workers"],
                    row["leg"],
                    f"{row['throughputOpsPerSecond']:.6f}",
                    f"{row['throughputAARelativeRangePercent']:.6f}",
                    row["latencyMilliseconds"]["p50"],
                    row["latencyMilliseconds"]["p95"],
                    row["latencyMilliseconds"]["p99"],
                    f"{row['p99AARelativeRangePercent']:.6f}",
                    row["processCPUSeconds"],
                    f"{row['busyCores']:.6f}",
                    f"{row['cpuMicrosecondsPerSuccess']:.6f}",
                    f"{row['cpuPerReadAARelativeRangePercent']:.6f}",
                    row["memoryBytes"]["rss"],
                    row["memoryBytes"]["hwm"],
                    "true",
                ]
            )
    with (comparison_dir / "setup-equivalence.csv").open(
        "w", newline="", encoding="utf-8"
    ) as handle:
        writer = csv.writer(handle)
        writer.writerow(
            [
                "workers",
                "leg",
                "repeat",
                "source_commit",
                "binary_sha256",
                "transport",
                "http_version",
                "document_bytes",
                "seed_count",
                "key_sequence",
            ]
        )
        writer.writerows(setup_rows)
    with (comparison_dir / "observer-effect.csv").open(
        "w", newline="", encoding="utf-8"
    ) as handle:
        writer = csv.writer(handle)
        writer.writerow(
            [
                "workers",
                "leg",
                "run_class",
                "throughput_percent",
                "cpu_per_read_percent",
                "p99_percent",
            ]
        )
        for row in observer_effect:
            writer.writerow(
                [
                    row["workers"],
                    row["leg"],
                    row["runClass"],
                    f"{row['throughput']:.6f}",
                    f"{row['cpuPerRead']:.6f}",
                    f"{row['p99']:.6f}",
                ]
            )

    labels = [f"w{row['workers']} {row['leg']}" for row in authoritative]
    svg_bar_chart(
        "Authoritative clean-score throughput",
        labels,
        [row["throughputOpsPerSecond"] for row in authoritative],
        "reads/s",
        comparison_dir / "throughput.svg",
    )
    svg_bar_chart(
        "Authoritative clean-score P99 latency",
        labels,
        [row["latencyMilliseconds"]["p99"] for row in authoritative],
        "ms",
        comparison_dir / "latency.svg",
    )
    svg_flywheel(perf_accounting, comparison_dir / "cpu-flywheel.svg")

    lines = [
        "# Cosmos cross-SDK point-read pilot",
        "",
        "Only clean-score runs are authoritative for latency and throughput. Perf and Go "
        "diagnostic runs are observer-effect and attribution evidence only.",
        "",
        "| Workers | Leg | Reads/s | A/A range | P50 ms | P95 ms | P99 ms | CPU us/read | RSS MiB | HWM MiB |",
        "| ---: | --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |",
    ]
    for row in authoritative:
        lines.append(
            f"| {row['workers']} | {row['leg']} | {row['throughputOpsPerSecond']:.3f} | "
            f"{row['throughputAARelativeRangePercent']:.2f}% | "
            f"{row['latencyMilliseconds']['p50']:.3f} | "
            f"{row['latencyMilliseconds']['p95']:.3f} | "
            f"{row['latencyMilliseconds']['p99']:.3f} | "
            f"{row['cpuMicrosecondsPerSuccess']:.3f} | "
            f"{row['memoryBytes']['rss'] / 1048576:.3f} | "
            f"{row['memoryBytes']['hwm'] / 1048576:.3f} |"
        )
    lines.extend(
        [
            "",
            "## Observer effect",
            "",
            "Percent change from the matching clean-score mean. Positive throughput is faster; "
            "positive CPU/read or P99 means additional observer cost.",
            "",
            "| Workers | Leg | Run class | Throughput | CPU/read | P99 |",
            "| ---: | --- | --- | ---: | ---: | ---: |",
        ]
    )
    for row in observer_effect:
        lines.append(
            f"| {row['workers']} | {row['leg']} | {row['runClass']} | "
            f"{row['throughput']:+.2f}% | {row['cpuPerRead']:+.2f}% | "
            f"{row['p99']:+.2f}% |"
        )
    lines.extend(
        [
            "",
            "## Gates",
            "",
            "- Classic Gateway V1 and negotiated HTTP/2 were proven before every measured run.",
            "- Binary encoding was disabled; all timed operations were raw point reads.",
            "- Errors, retries, and throttles were zero.",
            "- Clean runs contain no perf, pprof, runtime trace, uprobes, or hot-path counters.",
            "- Perf Self composition normalizes independently to 100%; process CPU is a separate denominator.",
            "- The pilot stops at workers 8 and 16. Ashley-specific counters are not included.",
        ]
    )
    (comparison_dir / "report.md").write_text("\n".join(lines) + "\n", encoding="utf-8")

    flamegraph_links = []
    for path in sorted(run_dir.rglob("*.html")):
        if path.name == "dashboard.html":
            continue
        flamegraph_links.append(
            f'<li><a href="{path.relative_to(run_dir).as_posix()}">{path.relative_to(run_dir)}</a></li>'
        )
    dashboard = [
        "<!doctype html>",
        '<html><head><meta charset="utf-8"><title>Cosmos point-read pilot</title></head><body>',
        "<h1>Cosmos cross-SDK point-read pilot</h1>",
        "<p>Authoritative metrics come only from clean-score runs.</p>",
        '<ul><li><a href="comparison/report.md">Report</a></li>',
        '<li><a href="comparison/summary.csv">Summary CSV</a></li>',
        '<li><a href="comparison/summary.json">Summary JSON</a></li>',
        '<li><a href="comparison/observer-effect.csv">Observer effect CSV</a></li>',
        '<li><a href="comparison/setup-equivalence.csv">Setup equivalence</a></li>',
        '<li><a href="comparison/throughput.svg">Throughput chart</a></li>',
        '<li><a href="comparison/latency.svg">Latency chart</a></li>',
        '<li><a href="comparison/cpu-flywheel.svg">CPU composition</a></li></ul>',
        "<h2>Interactive flamegraphs</h2><ul>",
        *flamegraph_links,
        "</ul></body></html>",
    ]
    (run_dir / "dashboard.html").write_text("\n".join(dashboard) + "\n", encoding="utf-8")
    write_manifest(run_dir)
    print(f"Wrote pilot report to {comparison_dir / 'report.md'}")


def write_manifest(root: pathlib.Path) -> None:
    manifest_path = root / "SHA256SUMS"
    lines: list[str] = []
    for path in sorted(root.rglob("*")):
        if not path.is_file() or path == manifest_path:
            continue
        digest = hashlib.sha256(path.read_bytes()).hexdigest()
        lines.append(f"{digest}  {path.relative_to(root).as_posix()}")
    manifest_path.write_text("\n".join(lines) + "\n", encoding="utf-8")


def redact(value: Any) -> Any:
    if isinstance(value, dict):
        result: dict[str, Any] = {}
        for key, child in value.items():
            if SENSITIVE_KEYS.search(key):
                digest = hashlib.sha256(str(child).encode("utf-8")).hexdigest()[:16]
                result[key] = f"sha256:{digest}"
            else:
                result[key] = redact(child)
        return result
    if isinstance(value, list):
        return [redact(child) for child in value]
    return value


def sanitize(run_dir: pathlib.Path, output: pathlib.Path) -> None:
    allowed_names = {
        "metrics.json",
        "provenance.json",
        "perf-self.csv",
        "accounting.json",
        "summary.json",
        "summary.csv",
        "observer-effect.csv",
        "setup-equivalence.csv",
        "report.md",
        "dashboard.html",
        "throughput.svg",
        "latency.svg",
        "cpu-flywheel.svg",
        "SHA256SUMS",
    }
    with tempfile.TemporaryDirectory(prefix="cosmos-perf-sanitize-") as temp:
        root = pathlib.Path(temp) / "cosmos-perf-result"
        for source in sorted(run_dir.rglob("*")):
            is_derived_visual = (
                source.is_file()
                and "derived" in source.parts
                and source.suffix in (".svg", ".html", ".folded")
            )
            if not source.is_file() or (
                source.name not in allowed_names and not is_derived_visual
            ):
                continue
            relative = source.relative_to(run_dir)
            target = root / relative
            target.parent.mkdir(parents=True, exist_ok=True)
            if source.suffix == ".json":
                write_json(target, redact(load_json(source)))
            else:
                shutil.copy2(source, target)
        write_manifest(root)
        output.parent.mkdir(parents=True, exist_ok=True)
        with tarfile.open(output, "w:gz", format=tarfile.PAX_FORMAT) as archive:
            for path in sorted(root.rglob("*")):
                archive.add(path, arcname=path.relative_to(root.parent), recursive=False)
    print(f"Wrote sanitized bundle to {output}")


def calibrate_call_graph(
    config_path: pathlib.Path, frame_pointer_path: pathlib.Path, dwarf_path: pathlib.Path
) -> None:
    config = load_json(config_path)
    candidates = {
        "frame-pointer": load_json(frame_pointer_path),
        "dwarf,16384": load_json(dwarf_path),
    }
    passing: list[tuple[float, str]] = []
    for mode, quality in candidates.items():
        if quality.get("callGraphMode") != mode:
            raise CampaignError(f"{mode}: quality file has the wrong callGraphMode")
        if quality.get("lostEventPercent", 100) > config["profiling"]["maxLostEventPercent"]:
            continue
        if (
            quality.get("resolvedSamplePercent", 0)
            < config["profiling"]["minResolvedSamplePercent"]
        ):
            continue
        if (
            quality.get("callchainSamplePercent", 0)
            < config["profiling"]["minCallchainSamplePercent"]
        ):
            continue
        score = (
            quality["resolvedSamplePercent"]
            + quality["callchainSamplePercent"]
            - quality["lostEventPercent"]
        )
        passing.append((score, mode))
    if not passing:
        raise CampaignError("neither frame-pointer nor DWARF passed the perf quality gates")
    passing.sort(reverse=True)
    selected = passing[0][1]
    config["profiling"]["selectedCallGraph"] = selected
    write_json(config_path, config)
    print(f"Selected {selected} call graphs.")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    commands = parser.add_subparsers(dest="command", required=True)
    prerequisites_parser = commands.add_parser("prerequisites")
    prerequisites_parser.add_argument("--phase", choices=("control", "vm"), default="control")
    validate = commands.add_parser("validate")
    validate.add_argument("--config", type=pathlib.Path, required=True)
    provision_parser = commands.add_parser("provision")
    provision_parser.add_argument("--config", type=pathlib.Path, required=True)
    provision_parser.add_argument("--state", type=pathlib.Path, required=True)
    provision_parser.add_argument("--dry-run", action="store_true")
    teardown_parser = commands.add_parser("teardown")
    teardown_parser.add_argument("--state", type=pathlib.Path, required=True)
    teardown_parser.add_argument("--dry-run", action="store_true")
    leg_parser = commands.add_parser("validate-leg")
    leg_parser.add_argument("--leg-dir", type=pathlib.Path, required=True)
    leg_parser.add_argument("--config", type=pathlib.Path, required=True)
    report_parser = commands.add_parser("report")
    report_parser.add_argument("--run-dir", type=pathlib.Path, required=True)
    report_parser.add_argument("--config", type=pathlib.Path, required=True)
    sanitize_parser = commands.add_parser("sanitize")
    sanitize_parser.add_argument("--run-dir", type=pathlib.Path, required=True)
    sanitize_parser.add_argument("--output", type=pathlib.Path, required=True)
    calibrate_parser = commands.add_parser("calibrate-call-graph")
    calibrate_parser.add_argument("--config", type=pathlib.Path, required=True)
    calibrate_parser.add_argument("--frame-pointer", type=pathlib.Path, required=True)
    calibrate_parser.add_argument("--dwarf", type=pathlib.Path, required=True)
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    try:
        if args.command == "prerequisites":
            prerequisites(args.phase)
        elif args.command == "validate":
            validate_config(load_json(args.config))
            print("Configuration is valid.")
        elif args.command == "provision":
            provision(load_json(args.config), args.state, args.dry_run)
        elif args.command == "teardown":
            teardown(load_json(args.state), args.dry_run)
        elif args.command == "validate-leg":
            validate_leg(args.leg_dir, load_json(args.config))
            print(f"{args.leg_dir.name} is valid.")
        elif args.command == "report":
            report(args.run_dir, load_json(args.config))
        elif args.command == "sanitize":
            sanitize(args.run_dir, args.output)
        elif args.command == "calibrate-call-graph":
            calibrate_call_graph(args.config, args.frame_pointer, args.dwarf)
        return 0
    except CampaignError as exc:
        print(f"error: {exc}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
