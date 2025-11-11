#!/usr/bin/env bash
set -euo pipefail

LOG_DIR="/app/var/logs/api"
mkdir -p /app/logs "${LOG_DIR}"

/app/bin/api-server &
CSUGO_PID=$!

cleanup() {
  kill "${CSUGO_PID}" "${MCP_PID}" 2>/dev/null || true
  wait "${CSUGO_PID}" 2>/dev/null || true
  wait "${MCP_PID}" 2>/dev/null || true
}
trap cleanup SIGINT SIGTERM

/app/bin/mcp-proxy &
MCP_PID=$!

wait -n "${CSUGO_PID}" "${MCP_PID}"
