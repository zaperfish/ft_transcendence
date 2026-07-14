#!/usr/bin/env sh
set -eu

BASE_URL="${BASE_URL:-http://localhost:7080}"
MODE="${MODE:-sleep}"
GOROUTINES="${GOROUTINES:-1000}"
HOLD_SECONDS="${HOLD_SECONDS:-30}"

curl -sS -X POST "${BASE_URL}/api/debug/scheduler-pressure?mode=${MODE}&goroutines=${GOROUTINES}&hold_seconds=${HOLD_SECONDS}"
printf "\n"
