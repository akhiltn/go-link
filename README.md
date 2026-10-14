# go-link

## Overview
**go-link** is a URL shortening service built with Go, BoltDB, and Fiber. This project aims to provide a simple and efficient way to shorten URLs and manage them.

## Features
- URL shortening
- URL redirection
- URL management
- RESTful API
- Swagger API documentation

## Technologies Used
- **Go**: The main programming language used for the project.
- **BoltDB**: A high-performance NoSQL database for Go.
- **Fiber**: An Express-inspired web framework for Go.
- **Swagger**: For API documentation.

## Installation

1. **Clone the repository:**
    ```sh
    git clone https://github.com/akhiltn/go-link.git
    cd go-link
    ```

2. **Install dependencies:**
    ```sh
    go mod tidy
    ```

3. **Run the application:**
    ```sh
    go run ./cmd/go-link
    ```

## Running with Podman

The `Dockerfile` is multi-stage: `base` -> `test` -> `build` -> `runtime` (small, non-root image).

### macOS: set up the Podman machine

On macOS, Podman runs containers in a Linux VM. Create and start it once:
```sh
podman machine init
podman machine start
```
On later sessions, `podman machine start` is enough. If you see `unable to connect to Podman socket ... connection refused`, the machine isn't running.

### Corporate network: `x509: certificate signed by unknown authority`

If image pulls fail with this error, your network inspects TLS traffic with a corporate CA that the Podman VM doesn't trust. Add your CA bundle (PEM) to the VM, then restart it:
```sh
podman machine ssh "sudo tee /etc/pki/ca-trust/source/anchors/corp-ca.crt >/dev/null && sudo update-ca-trust" < /path/to/ca_bundle.pem
podman machine stop
podman machine start
```
Repeat this if you recreate the machine with `podman machine rm`.

### Build and run
```sh
# Run tests in a clean container (optional)
make container-test

# Build the final image
make image

# Run it (data is persisted in the golink-data volume)
make run
```

The Makefile defaults to Podman. Override `CONTAINER_ENGINE=docker` to use Docker, or
`DOMAIN`, `PORT`, and `DATA_VOLUME` to change the runtime settings. The application
environment variables are `APP_PORT` (default `3000`), `DB_PATH` (default
`/data/go-link.db`), and `DOMAIN`.

### Hot reload (development)

The `dev` stage runs [air](https://github.com/air-verse/air), which rebuilds and restarts the app when source files change:
```sh
make dev
```
- `make dev` builds a static Linux Air binary and mounts it into the development
  image. It assumes the container image architecture matches the host architecture.
- `-v "$PWD":/src` lets air see your edits.
- `-v golink-data:/data` keeps the BoltDB file between restarts. Only one container can use it at a time (file lock), so stop other `go-link` containers first.
- The Air binary is mounted because the development image does not contain it.
- air polls for changes (`-build.poll true`) because file events are unreliable on macOS bind mounts.

## Usage

### API Endpoints

- **POST /
    - Request: `{ "url": "https://example.com", "short": "example" }`

- **GET /:shortened_url**
    - Redirects to the original URL.

- **DELETE /:shortened_url**
    - Deletes a short URL.

- **GET /api/allkv**
    - Lists all stored short URLs.

The short keys `api` and `swagger` are reserved (case-insensitive); creating one returns `400`.

### Example

To shorten a URL:
```sh
curl -X POST http://localhost:3000/ -d '{"url":"https://example.com", "short": "example" }' -H "Content-Type: application/json"
