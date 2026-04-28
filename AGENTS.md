# Gemini Context: mineserver-manager

`mineserver-manager` is a Go-based CLI tool designed to automate the installation, configuration, and management of Minecraft server instances.

## Project Overview

- **Purpose**: Streamline Minecraft server management, including version discovery, installation, Java Runtime Environment (JRE) management, and backups.

## Detailed Design & Code Organization

The project follows a clean architecture approach, prioritizing decoupling and testability through interfaces and dependency injection.

- **CLI Layer (`cmd/cli`)**: Built using [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper). It follows a **Lean Command** pattern:
  - Command files define flags and handle initial validation.
  - Action execution is delegated to "Runner" functions or direct Service calls.
- **Service Layer (`internal/minecraft/`)**: Orchestrates high-level domain operations.
  - `InstallService`: Coordinates server fetching, JRE setup, and instance provisioning.
  - `BackupService`: Manages instance archiving and rotation.
- **Installer & Strategy Pattern (`internal/installer/`)**:
  - `ServerFlavor` interface allows supporting multiple server types (Vanilla, Purpur, etc.) polymorphically.
  - `RuntimeManager` handles the lifecycle of Java Runtimes, ensuring the correct JDK is available for the selected Minecraft version.
- **Provisioner (`internal/provisioner/`)**:
  - Responsible for generating the filesystem layout.
  - Uses Go's `embed` package to bundle templates (`start.sh`, `stop.sh`, `log4j2.xml`).
  - Employs `text/template` for rendering configuration files.
- **Persistence (`internal/repository/`)**:
  - Implements the **Repository Pattern** to abstract data access.
  - Currently uses [Storm](https://github.com/asdine/storm) as a BoltDB-based embedded storage provider.
- **Domain Model (`internal/model/`)**: Contains pure data structures, keeping them decoupled from persistence or external API concerns.
- **External Clients (`internal/mojang/`)**: Specialized clients for interacting with third-party APIs (e.g., official Mojang metadata).

## Technologies & Strategy

- **Language**: Go (version 1.26.2+)
- **CLI Framework**: Cobra & Viper for robust CLI interactions and configuration management.
- **Persistence**: Storm (BoltDB based) for a zero-dependency, single-file embedded database.
- **Logging**: `log/slog` for structured, high-performance logging.
- **Testing Strategy**:
  - **Unit Testing**: Extensive use of `testify` and `gock` for mocking HTTP responses.
  - **Integration Testing**: Using `testcontainers-go` for environment-accurate tests (e.g., SSH validation).
- **Template-driven Configuration**: Leveraging Go's `text/template` ensures that startup scripts and configurations are consistent yet customizable.
- **Functional Options**: Services and clients use the functional options pattern for flexible, type-safe configuration.
- **Release**: `goreleaser` for automated cross-platform builds and artifact signing.

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

- **Internal Package**: Domain logic resides in `internal/` to maintain a clean API and prevent external imports.
- **Dependency Injection**: Heavy use of interface-based components to facilitate testing and modularity.
- **Surgical Updates**: Prefer targeted changes to existing files over large refactors unless explicitly requested.

## Improvement Points

- **Error Handling**: Enhance error wrapping with consistent use of `%w` and potentially introduce domain-specific error types.
- **Expanded Flavor Support**: Implement support for additional flavors like Forge, Fabric, and Quilt.
- **Plugin/Mod Management**: Integrate with SpigotMC/Modrinth/CurseForge APIs for automated plugin and mod installation.
- **Remote Orchestration**: Expand the existing SSH utilities to support remote instance management and deployment.
- **Concurrency**: Parallelize artifact downloads and backup processes to improve performance.
- **User Experience (UX)**:
  - Add interactive progress bars for long-running operations (downloads, backups).
  - Implement an interactive "wizard" mode for server installation.
- **Cloud Integration**: Add support for S3-compatible storage for backups.
- **API Mode**: Provide an option to run as a background service with a REST or gRPC API for remote control or Web UI integration.
