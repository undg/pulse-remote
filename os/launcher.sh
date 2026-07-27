#!/bin/sh

if ! command -v electron >/dev/null 2>&1; then
	echo "pulse-remote-desktop: 'electron' not found on PATH." >&2
	echo "" >&2
	echo "The desktop launcher requires Electron to run." >&2
	echo "Install it from: https://github.com/electron/electron/releases" >&2
	echo "Or check your package manager (may be called 'electron' or 'electron38')." >&2
	exit 1
fi

# Always exec — the single-instance lock in the app handles toggle:
#   - Not running  → starts and shows the window
#   - Already running → toggles visibility of the existing instance

exec electron /usr/lib/pulse-remote/desktop/app.asar "$@"
