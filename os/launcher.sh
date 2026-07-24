#!/bin/sh

if pkill -0 -f '/usr/lib/pulse-remote/desktop/app\.asar' 2>/dev/null; then
	pkill -f '/usr/lib/pulse-remote/desktop/app\.asar'
else
	exec electron /usr/lib/pulse-remote/desktop/app.asar "$@"
fi
