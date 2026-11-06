#!/usr/bin/env bash
set -Eeuo pipefail
exec /bin/bash "$(dirname "$0")/resume_campaign.sh" ffi
