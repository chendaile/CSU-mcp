#!/usr/bin/env bash
set -euo pipefail

CSUGO_PORT="${CSUGO_HTTP_PORT:-12000}"
CONFIG_FILE="/app/configs/api/conf/app.conf"
LOG_DIR="/app/var/logs/api"

if [[ -n "${CSUGO_PORT}" ]]; then
  if grep -q "^httpport" "${CONFIG_FILE}"; then
    sed -ri "s/^httpport\s*=.*/httpport = ${CSUGO_PORT}/" "${CONFIG_FILE}"
  else
    echo "httpport = ${CSUGO_PORT}" >> "${CONFIG_FILE}"
  fi
fi

export CSUGO_BASE_URL="${CSUGO_BASE_URL:-http://127.0.0.1:${CSUGO_PORT}}"

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
