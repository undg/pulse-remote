# pulse-remote

![test Pipewire](https://github.com/undg/pulse-remote/actions/workflows/test-pipewire.yml/badge.svg)
![test PulseAudio](https://github.com/undg/pulse-remote/actions/workflows/test-pulseaudio.yml/badge.svg)
![audit](https://github.com/undg/pulse-remote/actions/workflows/audit.yml/badge.svg)
![tidy](https://github.com/undg/pulse-remote/actions/workflows/tidy.yml/badge.svg)

Audio mixer for PulseAudio and PipeWire — local desktop app and web UI accessible from any device on your network.

## Features

- **Universal Compatibility**: Works with both PulseAudio and PipeWire
- **Real-time Control**: Adjust volume, mute/unmute, and switch audio outputs instantly
- **Multiple Interfaces**: Web app (built-in), desktop app, or API access
- **WebSocket API**: Real-time updates and low-latency control
- **Multi-device Support**: Manage all audio sinks, sources, and applications

## Installation

### Option 1: Arch Linux (AUR)

Two packages are available:

- **`pulse-remote-bin`** — Pre-built binary, no compilation needed
- **`pulse-remote-git`** — Builds from latest git source

```bash
# Pre-built binary (fast install)
yay -S pulse-remote-bin
# or build from source
yay -S pulse-remote-git
```

The service is installed and auto-enabled. If needed, start manually with:

```bash
systemctl --user enable --now pulse-remote.service
```

### Option 2: Build from Source

```bash
git clone https://github.com/undg/pulse-remote
cd pulse-remote
make build
make install-local
```

### Option 3: Build from Source with Web Frontend

To bundle the latest built-in web UI into the binary:

```bash
git clone https://github.com/undg/pulse-remote
cd pulse-remote
make pull/web
make build
make install
```

`make pull/web` fetches the latest [pulse-remote-web](https://github.com/undg/pulse-remote-web), runs its tests, builds it, and commits the generated `_GUI/web` assets.

> Requires [pnpm](https://pnpm.io/) and Node.js installed.

## Usage

### Web Application

After starting the server, open your browser and navigate to:

```
http://localhost:8448
```

The built-in web interface provides full control over your audio system.

### Desktop Application

The desktop app is bundled with the server — launch it from your application launcher
or run:

```bash
pulse-remote-desktop
```

### API Access

The WebSocket API is available at:

```
ws://localhost:8448/api/v1/ws
```

REST endpoint for status:

```
http://localhost:8448/api/v1/status
```

For detailed API documentation, connect to the WebSocket endpoint and send a `GetSchema` action, or visit:

```
http://localhost:8448/api/v1/schema/status
http://localhost:8448/api/v1/schema/message
http://localhost:8448/api/v1/schema/response
```

## Configuration

### Debug Logging

Control log verbosity with the `DEBUG` environment variable:

```bash
# For systemd service
systemctl --user set-environment DEBUG=trace
systemctl --user restart pulse-remote

# For direct execution
DEBUG=debug ./build/bin/pulse-remote-server
```

Available levels:

- `TRACE` or `3` - Most verbose
- `DEBUG` or `2` - Debug information
- `INFO` or `1` - Default level
- `WARN` or `0` - Warnings only
- `ERR` or `-1` - Errors only

### View Logs

```bash
# For systemd service
journalctl --user -u pulse-remote.service -f

# Clean output
journalctl --user -u pulse-remote.service -f --output cat
```

## Development

### Prerequisites

- Go 1.25 or later (preferably installed with mise)
- PulseAudio or PipeWire
- Make

### Quick Start

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run with hot reload
make run/watch

# Format and tidy code
make tidy

# Run full quality checks
make audit
```

### Project Structure

```
pulse-remote/
├── .github/
│   └── workflows/         # CI/CD workflows (test, audit, tidy, release)
├── api/                   # Core API implementation
│   ├── buildinfo/         # Build metadata (version, commit, date)
│   ├── json/              # JSON schemas and REST endpoints
│   ├── logger/            # Zerolog logging setup
│   ├── pactl/             # PulseAudio/PipeWire control
│   │   └── generated/     # Auto-generated types from pactl JSON
│   ├── utils/             # Utility functions (network, etc.)
│   └── ws/                # WebSocket handlers and broadcasting
├── _GUI/web/              # Built-in web interface
│   ├── dist/              # Compiled web app assets
│   └── version            # Web interface version
├── os/                    # System integration files
│   ├── pulse-remote.1     # Man page
│   └── pulse-remote.service  # Systemd user service
├── scripts/               # Build and development scripts
│   ├── bump.sh            # Version bumping script
│   └── test-watch.sh      # Watch mode test runner
├── vendor/                # Vendored dependencies
├── .gitignore             # Git ignore patterns
├── .github/workflows/release.yml  # CI release workflow (make package + upload)
├── .mise.toml             # Mise tool version configuration
├── go.mod                 # Go module dependencies
├── go.sum                 # Go module checksums
├── LICENSE                # License file
├── main.go                # Application entry point
├── Makefile               # Build and development tasks
├── README.md              # This file
└── renovate.json          # Renovate dependency updates config
```

### Available Make Commands

- `make help` - List all available commands
- `make pull/web` - Pull and build the web frontend
- `make build` - Build the server binary
- `make test` - Run all tests with race detection
- `make test/watch` - Run tests in watch mode
- `make test/cover` - Run tests with coverage report
- `make tidy` - Format code and tidy dependencies
- `make audit` - Run quality checks (vet, staticcheck, govulncheck)
- `make run/watch` - Run with hot reload
- `make install` - Install files to $(DESTDIR)$(PREFIX) (for packaging)
- `make install-local` - Install to /usr and enable systemd service
- `make uninstall` - Remove installed files and disable service

### Testing

Run all tests:

```bash
make test
```

Run specific test:

```bash
go test -v -race -buildvcs ./api/pactl -run TestGetSinks
```

Run tests in watch mode:

```bash
make test/watch
```

### Code Style

See [AGENTS.md](AGENTS.md) for detailed coding guidelines.

## Related Projects

- [pulse-remote-web](https://github.com/undg/pulse-remote-web) - Web interface (included in this repo)
- [pulse-remote-desktop](https://github.com/undg/pulse-remote-desktop) - Standalone desktop application

## License

See [LICENSE](LICENSE) file for details.
