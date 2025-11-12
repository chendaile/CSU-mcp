#!/usr/bin/env bash
set -euo pipefail

LOG_DIR="/app/var/logs/api"
mkdir -p /app/logs "${LOG_DIR}"

NGINX_PID_FILE="/run/nginx.pid"
NGINX_PID=""
NGINX_RUNNING=false

start_nginx() {
  nginx
  for _ in $(seq 1 10); do
    if [[ -f "${NGINX_PID_FILE}" ]]; then
      NGINX_PID="$(cat "${NGINX_PID_FILE}")"
      return 0
    fi
    sleep 0.2
  done
  echo "failed to start nginx" >&2
  exit 1
}

if command -v nginx >/dev/null 2>&1; then
  start_nginx
  NGINX_RUNNING=true
else
  echo "warning: nginx binary not found, skipping reverse proxy startup" >&2
fi

/app/bin/api-server &
CSUGO_PID=$!

cleanup() {
  kill "${CSUGO_PID}" "${MCP_PID}" 2>/dev/null || true
  if [[ "${NGINX_RUNNING}" == true && -n "${NGINX_PID}" ]]; then
    kill "${NGINX_PID}" 2>/dev/null || true
  fi
  wait "${CSUGO_PID}" 2>/dev/null || true
  wait "${MCP_PID}" 2>/dev/null || true
  if [[ "${NGINX_RUNNING}" == true && -n "${NGINX_PID}" ]]; then
    wait "${NGINX_PID}" 2>/dev/null || true
  fi
}
trap cleanup SIGINT SIGTERM

/app/bin/mcp-proxy &
MCP_PID=$!

if [[ "${NGINX_RUNNING}" == true && -n "${NGINX_PID}" ]]; then
  wait -n "${CSUGO_PID}" "${MCP_PID}" "${NGINX_PID}"
else
  wait -n "${CSUGO_PID}" "${MCP_PID}"
fi
