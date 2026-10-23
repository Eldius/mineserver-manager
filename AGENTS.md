# AGENTS.md

Go CLI (`mineserver` binary) to install and manage Minecraft servers: version discovery, JRE provisioning, instance generation, backups. Module `github.com/eldius/mineserver-manager`; entrypoint `./cmd/cli`.

## Commands

- Build: `go build -o mineserver ./cmd/cli`
- Test: `go test ./...`
- Lint: `golangci-lint run` (no config file — runs with defaults)
- Vulncheck: `go tool govulncheck ./...`
- Full CI check: `go tool task validate` (= test + vulncheck + lint; what `.github/workflows/go.yml` runs)

### Gotchas

- **`go test ./...` fails without Docker**: `internal/ssh/ssh_test.go` uses testcontainers (pulls `lscr.io/linuxserver/openssh-server`). Without Docker: `go test $(go list ./... | grep -v /internal/ssh)`. This also makes `task test` / `task validate` fail.
- Dev tools are Go tool directives in `go.mod` — invoke via `go tool task`, `go tool goreleaser`, `go tool govulncheck`, `go tool dlv`, `go tool deadcode`. Exception: `golangci-lint` is expected at `$(go env GOPATH)/bin/golangci-lint` (the taskfile hardcodes that path; CI installs it via the golangci-lint.run curl script).
- Requires Go 1.27.3 per `go.mod` (CI pins 1.26.2 and relies on toolchain auto-download).
- The `Makefile` holds only the maintainer's personal workflows (manual local runs, remote deploys against a hardcoded test host) — CI checks and releases live in `taskfile.yaml`. `task release-tag` creates and pushes git tags; don't run it casually.

## Architecture

- `cmd/cli/cmd` — Cobra commands, lean-command pattern: `install.go` defines flags, `install_runner.go` holds the logic (same for `backup*.go`). Keep new commands structured this way.
- `internal/minecraft` — `InstallService` / `BackupService`, constructed with functional options.
- `internal/installer` — `ServerFlavor` interface; only Vanilla implemented (`vanilla.go`). `--flavor purpur` passes flag parsing but errors "not yet implemented".
- `internal/java` — JRE download/management per Minecraft version, including `RuntimeManager` (the JRE install port used by `InstallService`).
- `internal/provisioner` — renders the instance filesystem; templates embedded with `//go:embed all:templates` (`start.sh`, `stop.sh`, `log4j2.xml`) under `internal/provisioner/templates`, rendered via `text/template`.
- `internal/mojang` — Mojang API client.
- `internal/repository` — Storm (BoltDB) persistence. Constructed in the cmd layer and injected via `minecraft.WithRepository`; persistence is best-effort (install proceeds without it).
- `internal/config` — viper property keys and runtime config getters.
- `internal/buildinfo` — version vars (`Version`, `Commit`, ...) injected by goreleaser ldflags; shown by `mineserver version`.
- `internal/httpx` / `internal/archive` / `internal/ssh` — HTTP download + checksum helpers, zip/tar.gz packing, SSH exec; split out of the old grab-bag `internal/utils`, which now keeps only small generic helpers (`Must`, `PasswordPrompt`, path helpers).

## Conventions

- Logging: `log/slog` via `internal/logger`.
- Runtime config: `$HOME/.mineserver/config.yaml` (fallback `./config.yaml`), overridable with `--config`; defaults are registered in `cmd/cli/cmd/root.go`.
- Tests: testify + `gock` for HTTP mocking (Mojang hosts `launchermeta`/`piston-meta`/`piston-data.mojang.com`); fixtures live in `testdata/` dirs (`internal/mojang/testdata`, `internal/minecraft/testdata`, `internal/ssh/testdata`).
- Errors: per-package sentinel errors (`ErrXxx` in the returning package's `errors.go` — no shared errs package), always wrapped with `%w`. Custom types (`InstallError`, `DownloadError`, `APIError`) carry context and wrap a sentinel for the category. Check with `errors.Is`/`errors.As` only at the cmd boundary; services return errors, never panic. Don't export a sentinel that has no `errors.Is` consumer.
- Flag definitions live in `cmd/cli/cmd/*.go` — check them there if CLI behavior is unclear.

## Release

- Pushing a `v*` tag triggers goreleaser via `.github/workflows/releases.yml`; version info is injected by ldflags into `internal/buildinfo`.
- Local snapshot: `go tool goreleaser release --snapshot --clean`. Note its before-hooks run `go mod tidy`, tests, govulncheck, and lint — all must pass (Docker needed for tests).
- Builds are linux-only (amd64, arm64, armv6/7), `CGO_ENABLED=0`, static; binary name `mineserver`.
