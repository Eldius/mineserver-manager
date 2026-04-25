# Gemini Context: mineserver-manager

`mineserver-manager` is a Go-based CLI tool designed to automate the installation, configuration, and management of Minecraft server instances.

## Project Overview

- **Purpose**: Streamline Minecraft server management, including version discovery, installation, Java Runtime Environment (JRE) management, and backups.
- **Architecture**:
  - **CLI Layer (`cmd/cli`)**: Built using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper). Command logic is delegated to runner functions in `cmd/cli/cmd/*_runner.go`.
  - **Internal Logic (`internal/`)**:
    - **`minecraft/`**: Core orchestration logic for server management.
    - **`installer/`**: Specialized components for fetching artifacts (`Downloader`) and managing JDKs (`RuntimeManager`). Supports multiple server flavors (e.g., Vanilla, Purpur) via a strategy pattern.
    - **`provisioner/`**: Handles filesystem layout, template rendering (`start.sh`, `stop.sh`, `log4j2.xml`), and initial configuration.
    - **`repository/`**: Persistence layer using a repository pattern (currently implemented with [Storm](https://github.com/asdine/storm)).
    - **`model/`**: Pure domain data models (Instances, ServerProperties, etc.), decoupled from persistence and configuration logic.
    - **`mojang/`**: Client for interacting with official Mojang APIs.
    - **`utils/`**: Shared internal utilities for networking, compression, and system operations.

## Technologies

- **Language**: Go (version 1.26.2+)
- **CLI Framework**: Cobra & Viper
- **Persistence**: Storm (BoltDB based)
- **Logging**: `log/slog`
- **Testing**: `testify`, `testcontainers-go`, `gock`
- **Release**: `goreleaser`

## Key Commands

### Development
- **Build**: `go build -o mineserver ./cmd/cli`
- **Test**: `go test ./...`
- **Tidy**: `go mod tidy`
- **Lint**: `make lint`
- **Vulnerability Check**: `make vulncheck`

### Application Usage
- **Install Server**:
  ```bash
  mineserver install --flavor vanilla --version 1.21.3 --dest ./my-server --motd "My Awesome Server" --memory-limit 2g
  ```
- **List Versions**: `mineserver install --list` (defaults to vanilla flavor)
- **Backup Instance**:
  ```bash
  mineserver backup save --instance-folder ./my-server --backup-folder ./backups --max-backup-files 5
  ```
- **Restore Backup**:
  ```bash
  mineserver backup restore --instance-folder ./restored-server --backup-file ./backups/my-server_2024-12-31_12-00-00_backup.zip
  ```

## Development Conventions

- **Internal Package**: Almost all domain logic resides in `internal/` to prevent external imports and maintain a clean API.
- **Dependency Injection**: Services use functional options and interface-based components for better testability and modularity.
- **Lean Commands**: Cobra command files focus on flag definition and validation, delegating action execution to runner functions.
- **Service Pattern**: Orchestration logic is encapsulated in services (e.g., `InstallService`, `BackupService`).
- **Repository Pattern**: All database interactions must go through the `Repository` interface defined in `internal/repository`.
