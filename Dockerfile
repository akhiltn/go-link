# syntax=docker/dockerfile:1

# ---- base: toolchain + sources ----
FROM golang:1.23-alpine AS base
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY docs ./docs

# ---- test: run with `podman build --target test .` ----
FROM base AS test
RUN go test ./...

# ---- dev: hot reload; mount source at /src and a linux air binary at /usr/local/bin/air ----
FROM base AS dev
ENV DB_PATH=/data/go-link.db
EXPOSE 3000
CMD ["air", "-build.cmd", "go build -o /tmp/go-link ./cmd/go-link", "-build.bin", "/tmp/go-link", "-build.poll", "true"]

# ---- build: static binary ----
FROM base AS build
ARG VERSION=dev
ENV CGO_ENABLED=0
RUN go build -trimpath \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/go-link ./cmd/go-link

# ---- runtime: minimal non-root image ----
FROM alpine:3.20 AS runtime
RUN addgroup -S app && adduser -S -G app app \
    && mkdir /data && chown app:app /data
COPY --from=build /out/go-link /usr/local/bin/go-link
USER app
ENV APP_PORT=3000 \
    DB_PATH=/data/go-link.db
VOLUME /data
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget -qO- http://127.0.0.1:${APP_PORT}/livez || exit 1
ENTRYPOINT ["go-link"]
