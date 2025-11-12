#!/usr/bin/env bash
set -euo pipefail

LOG_DIR="/app/var/logs/api"
mkdir -p /app/logs "${LOG_DIR}"

TLS_CONF_PATH="/etc/nginx/conf.d/csu-mcp-https.conf"
TLS_CERT_DIR="${NGINX_CERT_DIR:-/etc/nginx/certs}"
TLS_ENABLED=false

NGINX_PID=""
NGINX_RUNNING=false

find_first_match() {
  local pattern=$1
  shopt -s nullglob
  local matches=($pattern)
  shopt -u nullglob
  if ((${#matches[@]} > 0)); then
    printf '%s\n' "${matches[0]}"
    return 0
  fi
  return 1
}

find_cert_file() {
  if [[ -d "${TLS_CERT_DIR}" ]]; then
    if find_first_match "${TLS_CERT_DIR}"/*.crt >/dev/null; then
      find_first_match "${TLS_CERT_DIR}"/*.crt
      return 0
    fi
    if find_first_match "${TLS_CERT_DIR}"/*.pem >/dev/null; then
      find_first_match "${TLS_CERT_DIR}"/*.pem
      return 0
    fi
  fi
  return 1
}

find_key_file() {
  if [[ -d "${TLS_CERT_DIR}" ]] && find_first_match "${TLS_CERT_DIR}"/*.key >/dev/null; then
    find_first_match "${TLS_CERT_DIR}"/*.key
    return 0
  fi
  return 1
}

configure_nginx_tls() {
  local cert="${NGINX_SSL_CERT_FILE:-}"
  local key="${NGINX_SSL_KEY_FILE:-}"

  if [[ -n "${cert}" && ! -f "${cert}" ]]; then
    echo "warning: specified TLS certificate ${cert} not found" >&2
    cert=""
  fi
  if [[ -n "${key}" && ! -f "${key}" ]]; then
    echo "warning: specified TLS key ${key} not found" >&2
    key=""
  fi

  if [[ -z "${cert}" ]]; then
    cert="$(find_cert_file || true)"
  fi
  if [[ -z "${key}" ]]; then
    key="$(find_key_file || true)"
  fi

  if [[ -n "${cert}" && -n "${key}" ]]; then
    cat >"${TLS_CONF_PATH}" <<EOF
server {
    listen 443 ssl http2 default_server;
    listen [::]:443 ssl http2 default_server;
    server_name _;

    ssl_certificate ${cert};
    ssl_certificate_key ${key};
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;

    location / {
        proxy_pass http://mcp_upstream;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
        proxy_http_version 1.1;
        proxy_set_header Connection \$connection_upgrade;
        proxy_set_header Upgrade \$http_upgrade;
    }
}
EOF
    TLS_ENABLED=true
    echo "nginx TLS enabled (cert: ${cert})" >&2
  else
    rm -f "${TLS_CONF_PATH}"
    TLS_ENABLED=false
    echo "nginx TLS disabled; serving HTTP only" >&2
  fi
}

start_nginx() {
  nginx -g "daemon off;" &
  NGINX_PID=$!
  for _ in $(seq 1 10); do
    if kill -0 "${NGINX_PID}" 2>/dev/null; then
      return 0
    fi
    sleep 0.2
  done
  echo "failed to start nginx" >&2
  exit 1
}

if command -v nginx >/dev/null 2>&1; then
  configure_nginx_tls
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
