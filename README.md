
# mineserver-manager #

![Go](https://github.com/eldius/mineserver-manager/actions/workflows/go.yml/badge.svg)
![Latest release status](https://github.com/eldius/mineserver-manager/actions/workflows/releases.yml/badge.svg)

A CLI tool to help manage Minecraft servers: version discovery, automated installation (including the correct Java runtime), and backups.

## features ##

- Installs Vanilla servers (`--flavor purpur` is accepted but not implemented yet)
- Resolves `latest` to the newest stable version, or lists all available versions
- Downloads the correct Java runtime for the selected Minecraft version automatically
- Generates start/stop scripts and server configuration files
- Backup and restore, with optional rotation (`--max-backup-files`)

## install ##

Download the latest release for your arch (Linux: amd64, arm64, armv6/7):

```shell
## get latest version files
GET https://api.github.com/repos/eldius/mineserver-manager/releases/latest
```

```shell
## download latest version by arch
wget $( curl https://api.github.com/repos/eldius/mineserver-manager/releases/latest | jq -r '.assets[].browser_download_url' | grep <my-arch> )
```

Or build from source (requires Go 1.27.3+):

```shell
go build -o mineserver ./cmd/cli
```

## usage ##

### install a server ###

```shell
mineserver \
    install \
      --version 1.21.3 \
      --motd "A new test server" \
      --level-name "My test world" \
      --rcon-enabled \
      --headless \
      --dest ./my-server \
      --seed "5516949179205280665" \
      --memory-limit 2g \
      --whitelist-user Eldius
```

Notes:

- `--version` defaults to `latest`; `--dest` defaults to the current directory.
- When `--rcon-enabled` is set without `--rcon-passwd`, the password is asked interactively.

### list available versions ###

```shell
mineserver install --list
```

### backup ###

```shell
## save a backup (keeps at most 5 files)
mineserver backup save \
    --instance-folder ./my-server \
    --backup-folder ./backups \
    --max-backup-files 5

## restore a backup
mineserver backup restore \
    --instance-folder ./restored-server \
    --backup-file ./backups/my-server_2024-12-31_12-00-00_backup.zip
```

### other commands ###

- `mineserver version` — shows version/commit/build info.
- Global flags: `--debug` (verbose logs), `--home` (app home folder, defaults to `~/.mineserver`), `--config` (config file path).

### configuration ###

Settings are read from `$HOME/.mineserver/config.yaml` (fallback: `./config.yaml`), overridable with `--config`. Available keys and defaults are registered in `cmd/cli/cmd/root.go` (e.g. `minecraft.api.timeout`, `install.download.timeout`).

## development ##

```shell
go test ./...              # tests
golangci-lint run          # lint
go tool govulncheck ./...  # vulnerability check
go tool task validate      # all of the above — what CI runs
```

### web api / ui ###

A second binary, `mineserver-api` (`./cmd/api`), serves an embedded web UI
and a JSON API (`/api/*`) to install, list, start/stop and back up server
instances:

```shell
go run ./cmd/api --port 8080   # or: make api
```

Configuration (via `./config.yaml` or `~/.mineserver/config.yaml`):

```yaml
api:
  home:
    path: ~/.mineserver-api   # sandbox root: installs, backups and temp
                              # files all live (and must stay) under here
  auth:
    password: change-me       # login password for the UI; unset = API
                              # answers 503 on every endpoint
    token:
      ttl: 24h
```

Instances installed by the CLI outside the API home are listed but cannot
be started/stopped/backed up through the API (path access is confined to
the API home folder).

The UI source is a Vite + React + TypeScript app in `internal/web/ui`:

```shell
npm --prefix internal/web/ui run build   # builds into internal/web/static
npm --prefix internal/web/ui run dev     # dev server, proxies /api to :8080
```

The build output in `internal/web/static` is embedded via `go:embed` and
committed, so Go builds work without node — but rebuild and recommit it
whenever the UI changes (releases package whatever `static/` contains).

## reference links ##

- [APIs Mojang - Minecraft Wiki](https://minecraft.wiki/w/Mojang_API)
- [Mojang API - Wiki VG](https://wiki.vg/Mojang_API)
