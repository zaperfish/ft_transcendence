#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-https://localhost:7443}"
SIZE_MB="${SIZE_MB:-128}"
HOLD_SECONDS="${HOLD_SECONDS:-30}"

curl --fail --insecure -sS -X POST "${BASE_URL}/api/debug/gc-pressure?size_mb=${SIZE_MB}&hold_seconds=${HOLD_SECONDS}"
printf "\n"
