#!/bin/sh

if ! command -v electron >/dev/null 2>&1; then
	echo "pulse-remote-desktop: 'electron' not found on PATH." >&2
	echo "" >&2
	echo "The desktop launcher requires Electron to run." >&2
	echo "Install it from: https://github.com/electron/electron/releases" >&2
	echo "Or check your package manager (may be called 'electron' or 'electron38')." >&2
	exit 1
fi

if pkill -0 -f '/usr/lib/pulse-remote/desktop/app\.asar' 2>/dev/null; then
	pkill -f '/usr/lib/pulse-remote/desktop/app\.asar'
else
	exec electron /usr/lib/pulse-remote/desktop/app.asar "$@"
fi
