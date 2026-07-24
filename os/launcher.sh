#!/bin/sh

PID=$(pgrep -f '/usr/lib/pulse-remote/desktop/app\.asar')

if [ -z "$PID" ]; then
	exec electron /usr/lib/pulse-remote/desktop/app.asar "$@"
else
	kill "$PID"
fi
