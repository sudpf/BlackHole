# BlackHole

BlackHole is a Go-based network operations project. It contains two runnable components:

- **VoidEngine**: an HTTP API service for users, network traffic data, unified error responses, request context, and API documentation.
- **Stash**: a configuration-driven ingestion pipeline that reads events from inputs such as syslog or Kafka, applies filters, and writes to outputs such as ClickHouse, Elasticsearch, syslog, or stdout.

The project favors explicit configuration, small internal packages, structured logs, and a clear separation between API contracts, services, and storage models.

## Features

- REST API service built on Gin.
- Unified API response format and localized error messages.
- Request-scoped trace ID support through `X-Trace-ID`.
- Separate access logs and application logs.
- MySQL-backed user storage.
- ClickHouse-backed network traffic storage.
- Swagger/OpenAPI documentation for VoidEngine.
- Stash pipeline with input, filter, handler, and output stages.
- Graceful shutdown support for long-running services.

## Requirements

- Go 1.25 or newer.
- MySQL, if VoidEngine user APIs are enabled.
- ClickHouse, if VoidEngine traffic APIs or Stash ClickHouse output are enabled.

## Quick Start

Clone the repository and run tests:

```bash
git clone git@github.com:sudpf/BlackHole.git
cd BlackHole
go test ./...
```

Build the binaries:

```bash
mkdir -p bin
go build -o bin/voidengine ./cmd/voidengine
go build -o bin/stash ./cmd/stash
```

Run VoidEngine:

```bash
./bin/voidengine -config-file conf/voidengine.toml
```

Run Stash:

```bash
./bin/stash -f conf/stash.yaml
```

## Configuration

VoidEngine uses TOML:

```toml
[app]
listen = "http://0.0.0.0:80"
request_timeout = "8s"
shutdown_timeout = "10s"
```

Stash uses YAML:

```yaml
app:
  listen: http://127.0.0.1:8002
  request_timeout: 30s
  shutdown_timeout: 10s
```

`listen` is expressed as a URL so the protocol and address stay in one field. The current runtime supports the `http` scheme.

The sample files under `conf/` are development examples. Update database addresses, usernames, passwords, log paths, and output destinations before running in another environment.

## API Documentation

VoidEngine serves Swagger UI at:

```text
http://<voidengine-host>/swagger/index.html
```

Generated OpenAPI files are stored under:

```text
docs/api/voidengine/
```

Regenerate the documentation with:

```bash
make swagger-generator
```

## Project Layout

```text
api/                         HTTP middleware, routing, response helpers, Swagger wiring
api/voidengine/openapi/       VoidEngine OpenAPI server and v1 handlers
cmd/                         Service entry points
conf/                        Example runtime configuration
docs/api/voidengine/          Generated OpenAPI files
internal/runtime/             Shared runtime helpers
internal/stash/               Stash configuration, app bootstrap, pipeline services
internal/voidengine/          VoidEngine config, contract, service, model, error codes
pkg/                         Shared packages: logging, errors, db, auth, request context
```

## Development

Run all tests:

```bash
go test ./...
```

Format Go code:

```bash
gofmt -w $(find . -name '*.go' -not -path './vendor/*')
```

The `GNUmakefile` also provides build targets such as `voidengine`, `stash`, and `swagger-generator`. Be aware that some targets run formatting, `go mod tidy`, vendoring, and generated document updates as part of the build flow.

## Logging

VoidEngine writes:

- access logs to `api.log`
- application logs to `<binary-name>.log`

Logs are JSON formatted and include request trace IDs when available.

## License

No license has been declared yet.
