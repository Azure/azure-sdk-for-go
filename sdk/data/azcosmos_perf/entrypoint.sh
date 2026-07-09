#!/usr/bin/env bash
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.
set -euo pipefail

# Reserved for future Go profiling/debug hooks. Pyroscope profiles are collected
# externally via eBPF auto-instrumentation when PYROSCOPE_SERVER_URL is set.
exec cosmos-perf "$@"
