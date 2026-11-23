# AGENTS.md

Go CLI (`mineserver` binary) to install and manage Minecraft servers: version discovery, JRE provisioning, instance generation, backups. Module `github.com/eldius/mineserver-manager`; entrypoint `./cmd/cli`.

## Commands

- Build: `go build -o mineserver ./cmd/cli`
- Test: `go test ./...`
- Lint: `golangci-lint run` (no config file — runs with defaults)
- Vulncheck: `go tool govulncheck ./...`
- Full CI check: `go tool task validate` (= test + vulncheck + lint; what `.github/workflows/go.yml` runs)

### Gotchas

- Dev tools are Go tool directives in `go.mod` — invoke via `go tool task`, `go tool goreleaser`, `go tool govulncheck`, `go tool dlv`, `go tool deadcode`. Exception: `golangci-lint` is expected at `$(go env GOPATH)/bin/golangci-lint` (the taskfile hardcodes that path; CI installs it via the golangci-lint.run curl script).
- Go 1.27.4 per `go.mod`; both CI workflows pin exactly that.
- The `Makefile` holds only the maintainer's personal workflows (manual local runs, remote deploys against a hardcoded test host) — CI checks and releases live in `taskfile.yaml`. `task release-tag` creates and pushes git tags; don't run it casually.

## Architecture

- `cmd/cli/cmd` — Cobra commands, lean-command pattern: `install.go` defines flags, `install_runner.go` holds the logic (same for `backup*.go`). Keep new commands structured this way.
- `cmd/api` + `internal/web` — second binary (`mineserver-api`): JSON API (`/api/*`) + embedded React SPA. Run locally with `make api`; goreleaser ships it in separate `mineserver-api_*` archives.
  - UI: Vite + React + TS app in `internal/web/ui`; `npm --prefix internal/web/ui run build` outputs to `internal/web/static`, which is embedded (`//go:embed all:static`) and committed — rebuild + recommit after UI changes (Go builds must not require node).
  - Auth: single password from `api.auth.password` exchanged at `POST /api/login` for an HMAC bearer token; unset password = all `/api/*` endpoints answer 503.
  - `Sandbox` (`internal/web/sandbox.go`) confines every API filesystem access to `api.home.path` (installs → `instances/`, backups → `backups/`, temp → `tmp/`); CLI-installed instances outside it are listed but not manageable. Use `Sandbox.Resolve`/`Contains` for any new path-touching endpoint.
- `internal/minecraft` — `InstallService` / `BackupService`, constructed with functional options.
- `internal/installer` — `ServerFlavor` interface; only Vanilla implemented (`vanilla.go`). `--flavor purpur` passes flag parsing but errors "not yet implemented".
- `internal/java` — JRE download/management per Minecraft version, including `RuntimeManager` (the JRE install port used by `InstallService`).
- `internal/provisioner` — renders the instance filesystem; templates embedded with `//go:embed all:templates` (`start.sh`, `stop.sh`, `log4j2.xml`) under `internal/provisioner/templates`, rendered via `text/template`.
- `internal/mojang` — Mojang API client.
- `internal/repository` — Storm (BoltDB) persistence. Constructed in the cmd layer and injected via `minecraft.WithRepository`; persistence is best-effort (install proceeds without it).
- `internal/config` — viper property keys and runtime config getters. Don't confuse it with `internal/minecraft/config`, which holds the per-instance functional options (`InstanceOpt`); runners import both (the former aliased `cfg`).
- `internal/buildinfo` — version vars (`Version`, `Commit`, ...) injected by goreleaser ldflags; shown by `mineserver version`.
- `internal/httpx` / `internal/archive` — HTTP download + checksum helpers, zip/tar.gz packing; split out of the old grab-bag `internal/utils`, which now keeps only small generic helpers (`Must`, `PasswordPrompt`, path helpers).

## Conventions

- Logging: `log/slog` via `internal/logger`.
- Bootstrap: both binaries wire config + OpenTelemetry through `cli-telemetry-helper/pkg/setup.Run` (env prefix `cli_telemetry`); defaults are registered in `cmd/cli/cmd/root.go`.
- Runtime config: file lookup checks `.` before `~/.mineserver`, overridable with `--config`. Gotcha: the committed repo-root `config.yaml` (debug JSON logs, LAN telemetry endpoints) is auto-loaded whenever you run the binary from the repo root.
- Tests: testify + `gock` for HTTP mocking (Mojang hosts `launchermeta`/`piston-meta`/`piston-data.mojang.com`); fixtures live in `testdata/` dirs (`internal/mojang/testdata`, `internal/minecraft/testdata`).
- Errors: per-package sentinel errors (`ErrXxx` in the returning package's `errors.go` — no shared errs package), always wrapped with `%w`. Custom types (`InstallError`, `DownloadError`, `APIError`) carry context and wrap a sentinel for the category. Check with `errors.Is`/`errors.As` only at the cmd boundary; services return errors, never panic. Don't export a sentinel that has no `errors.Is` consumer.
- Flag definitions live in `cmd/cli/cmd/*.go` — check them there if CLI behavior is unclear.

## Release

- Pushing a `v*` tag triggers goreleaser via `.github/workflows/releases.yml`; version info is injected by ldflags into `internal/buildinfo`.
- Local snapshot: `go tool goreleaser release --snapshot --clean`. Note its before-hooks run `go mod tidy`, tests, govulncheck, and lint — all must pass.
- Two linux-only builds (amd64, arm64, armv6/7; `CGO_ENABLED=0`, static): `mineserver` (CLI, `./cmd/cli`, `mineserver-manager_*` archives) and `mineserver-api` (`./cmd/api`, `mineserver-api_*` archives). The API build inherits the CLI's build settings via a YAML anchor (`&build-defaults`) — keep them in sync if the anchor is removed.
