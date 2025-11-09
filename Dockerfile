# syntax=docker/dockerfile:1.7

############################
# Build unified binaries   #
############################
FROM golang:1.25 AS builder
WORKDIR /src

COPY go.mod go.sum ./
COPY internal ./internal
COPY cmd ./cmd
COPY third_party ./third_party

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -o /out/api-server ./cmd/api-server

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /out/mcp-proxy ./cmd/mcp-proxy

############################
# Runtime image            #
############################
FROM debian:bookworm-slim AS runtime

RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata bash && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV CSUGO_BASE_URL=http://127.0.0.1:12000 \
    CSUGO_HTTP_PORT=12000 \
    MCP_HTTP_ADDR=0.0.0.0:13000

RUN mkdir -p /app/bin /app/var/logs/api

COPY --from=builder /out/api-server /app/bin/api-server
COPY --from=builder /out/mcp-proxy /app/bin/mcp-proxy

COPY configs/api/conf ./configs/api/conf
COPY web/static ./web/static
COPY web/views ./web/views

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh && \
    mkdir -p /app/logs

EXPOSE 13000 12000
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
