# pulse-remote

![test Pipewire](https://github.com/undg/go-prapi/actions/workflows/test-pipewire.yml/badge.svg)
![test PulseAudio](https://github.com/undg/go-prapi/actions/workflows/test-pulseaudio.yml/badge.svg)
![audit](https://github.com/undg/go-prapi/actions/workflows/audit.yml/badge.svg)
![tidy](https://github.com/undg/go-prapi/actions/workflows/tidy.yml/badge.svg)

## Pulse Remote Backend

A simple and powerful PulseAudio Remote API for Linux systems.

## What is this?

go-prapi is a backend implementation for [pulse-remote](https://github.com/undg/pulse-remote) written in Go. It provides a WebSocket-based API to control and gather information from PulseAudio devices and sinks.

## Features

- Works with Linux PulseAudio and PipeWire
- WebSocket communication for real-time updates
- Control volume, mute status, and audio outputs
- Retrieve information about audio cards and sinks

## Quick Start

1. Clone the repository
2. Run the server:
3. The server will start on `ws://localhost:8448/api/v1/ws`

## Frontend

An actively developed frontend for this API is available at [pr-web](https://github.com/undg/pr-web).

To use the frontend:

1. Build the pr-web project
2. Copy or symlink the build output to the `frontend` folder in this project

Example (if pr-web is in a sibling directory):
```bash
ln -s ../pr-web/dist frontend
```


## API

For detailed API documentation, connect to the WebSocket endpoint and send a `GetSchema` action.

## Development

Check the Makefile for available commands:

- `make test`: Run tests
- `make build`: Build the application
- `make run/watch`: Run with auto-reload on file changes

## Debugging

Use build it logger. You can set environmental variable `DEBUG` to filter out or show more logs.

By default it's set to `"INFO"` or `"1"`.

All available options:

* `"TRACE"` or `"3"`
* `"DEBUG"` or `"2"`
* `"INFO"` or `"1"`
* `"WARN"` or `"0"`
* `"ERR"` or `"-1"`


Example of logger in the code.

logger.Trace().Msg("from logger.Trace")
logger.Debug().Msg("from logger.Debug")
logger.Info().Msg("from logger.Info")
logger.Warn().Msg("from logger.Warn")
logger.Error().Msg("from logger.Error")
logger.Fatal().Msg("from logger.Fatal")
logger.Panic().Msg("from logger.Panic")

# CLI snippets

Few useful commands

@TODO (undg) 2025-02-17: clean doc

```bash
make install

make uninstall

systemctl --user start pulse-remote.service

systemctl --user set-environment DEBUG=trace # see available options in Debugging section

systemctl --user restart pulse-remote.service

systemctl --user unset-environment DEBUG

journalctl --user -u pulse-remote.service -f --output cat

```

