# AGENTS.md

Go CLI (`mineserver` binary) to install and manage Minecraft servers: version discovery, JRE provisioning, instance generation, backups. Module `github.com/eldius/mineserver-manager`; entrypoint `./cmd/cli`.

## Commands

- Build: `go build -o mineserver ./cmd/cli`
- Test: `go test ./...`
- Lint: `golangci-lint run` (no config file — runs with defaults)
- Vulncheck: `go tool govulncheck ./...`
- Full CI check: `go tool task validate` (= test + vulncheck + lint; what `.github/workflows/go.yml` runs)

### Gotchas

- **`go test ./...` fails without Docker**: `internal/utils/ssh_test.go` uses testcontainers (pulls `lscr.io/linuxserver/openssh-server`). Without Docker: `go test $(go list ./... | grep -v /internal/utils)`. This also makes `task test` / `task validate` fail.
- **Avoid `make test`**: it ends with `go tool cover -html=coverage.out`, which tries to open a browser. Use plain `go test`.
- Dev tools are Go tool directives in `go.mod` — invoke via `go tool task`, `go tool goreleaser`, `go tool govulncheck`, `go tool dlv`, `go tool deadcode`. Exception: `golangci-lint` is expected at `$(go env GOPATH)/bin/golangci-lint` (the taskfile hardcodes that path; CI installs it via the golangci-lint.run curl script).
- Requires Go 1.27.3 per `go.mod` (CI pins 1.26.2 and relies on toolchain auto-download).
- `Makefile` targets `put`, `install-remote`, `backup-remote`, `restore-remote` are the maintainer's personal workflows against a hardcoded test host — not general dev commands. `make release` / `task release-tag` create and push git tags; don't run them casually.

## Architecture

- `cmd/cli/cmd` — Cobra commands, lean-command pattern: `install.go` defines flags, `install_runner.go` holds the logic (same for `backup*.go`). Keep new commands structured this way.
- `internal/minecraft` — `InstallService` / `BackupService`, constructed with functional options.
- `internal/installer` — `ServerFlavor` interface; only Vanilla implemented (`vanilla.go`). `--flavor purpur` passes flag parsing but errors "not yet implemented"; the `purpur` cobra command is an empty stub.
- `internal/java` — JRE download/management per Minecraft version.
- `internal/provisioner` — renders the instance filesystem; templates embedded with `//go:embed all:templates` (`start.sh`, `stop.sh`, `log4j2.xml`) under `internal/provisioner/templates`, rendered via `text/template`.
- `internal/mojang` — Mojang API client.
- `internal/repository` — Storm (BoltDB) persistence.
- `internal/config` — viper property keys plus version vars (`Version`, `Commit`, ...) injected by goreleaser ldflags; shown by `mineserver version`.
- `internal/errs` — shared sentinel errors.

## Conventions

- Logging: `log/slog` via `internal/logger`.
- Runtime config: `$HOME/.mineserver/config.yaml` (fallback `./config.yaml`), overridable with `--config`; defaults are registered in `cmd/cli/cmd/root.go`.
- Tests: testify + `gock` for HTTP mocking (Mojang hosts `launchermeta`/`piston-meta`/`piston-data.mojang.com`); fixtures live in `internal/mojang/samples` and `internal/minecraft/test_samples`.
- Trust `cmd/cli/cmd/*.go` over README for flags — the README install snippet uses a stale name (`--enable-rcon`; the real flag is `--rcon-enabled`).

## Release

- Pushing a `v*` tag triggers goreleaser via `.github/workflows/releases.yml`; version info is injected by ldflags into `internal/config`.
- Local snapshot: `go tool goreleaser release --snapshot --clean`. Note its before-hooks run `go mod tidy`, tests, govulncheck, and lint — all must pass (Docker needed for tests).
- Builds are linux-only (amd64, arm64, armv6/7), `CGO_ENABLED=0`, static; binary name `mineserver`.
