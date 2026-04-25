# Gemini Context: mineserver-manager

`mineserver-manager` is a Go-based CLI tool designed to automate the installation, configuration, and management of Minecraft server instances.

## Project Overview

- **Purpose**: Streamline Minecraft server management, including version discovery, installation, Java Runtime Environment (JRE) management, and backups.
- **Architecture**:
  - **CLI Layer (`cmd/cli`)**: Built using [Cobra](https://github.com/spf13/cobra) for commands and [Viper](https://github.com/spf13/viper) for configuration.
  - **Core Logic (`minecraft/`)**:
    - `install_service.go`: Orchestrates the full installation process (downloading JARs, setting up JRE, generating scripts).
    - `backup_service.go`: Manages ZIP-based backups with automated rollover.
  - **Ecosystem Integration**:
    - `mojang/`: Client for interacting with official Mojang APIs for version information.
    - `java/`: Manages automated downloading and installation of appropriate Microsoft OpenJDK versions based on the Minecraft version requirements.
  - **Utilities (`utils/`)**: Common helpers for networking, file compression (ZIP/TarGZ), and system operations.

## Technologies

- **Language**: Go (version 1.26.2+)
- **CLI Framework**: Cobra
- **Configuration**: Viper & `initial-config-go`
- **Logging**: `log/slog`
- **Testing**: `testify`, `testcontainers-go`, `gock`
- **Release**: `goreleaser`

## Key Commands

### Development
- **Build**: `go build -o mineserver ./cmd/cli`
- **Test**: `make test` (Runs tests with coverage)
- **Lint**: `make lint` (Requires `golangci-lint`)
- **Vulnerability Check**: `make vulncheck` (Requires `govulncheck`)
- **Local Snapshot**: `make snapshot-local` (Uses `goreleaser`)

### Application Usage
- **Install Server**:
  ```bash
  mineserver install --version 1.21.3 --dest ./my-server --motd "My Awesome Server" --memory-limit 2g
  ```
- **List Versions**: `mineserver install --list`
- **Backup Instance**:
  ```bash
  mineserver backup save --instance-folder ./my-server --backup-folder ./backups --max-backup-files 5
  ```
- **Restore Backup**:
  ```bash
  mineserver backup restore --instance-folder ./restored-server --backup-file ./backups/my-server_2024-12-31_12-00-00_backup.zip
  ```

## Development Conventions

- **Configuration**: The application looks for a configuration file in `$HOME/.mineserver/config.yaml` or the current directory.
- **Logging**: Structured logging using `slog` is preferred.
- **Error Handling**: Use `fmt.Errorf("...: %w", err)` for error wrapping to maintain context.
- **Service Pattern**: Business logic is typically encapsulated in services (e.g., `InstallService`, `BackupService`) using the functional options pattern for configuration.
- **JRE Management**: The tool automatically maps Minecraft versions to required Java versions (11, 17, 21) and downloads the appropriate Microsoft OpenJDK build for Linux (amd64/arm64).
- **Start/Stop Scripts**: Every installation generates `start.sh` and `stop.sh` scripts in the destination folder, pre-configured with the correct Java path and memory limits.
