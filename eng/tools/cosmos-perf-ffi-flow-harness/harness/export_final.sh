#!/usr/bin/env bash
set -Eeuo pipefail
python3 "$(dirname "$0")/export_final.py"
