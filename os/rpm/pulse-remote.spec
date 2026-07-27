Name:           pulse-remote
Version:        @VERSION@
Release:        1%{?dist}
Summary:        Remote audio control for PulseAudio/PipeWire
License:        MIT
URL:            https://github.com/undg/pulse-remote
Source0:        %{name}-%{version}.tar.gz

BuildArch:      x86_64
Requires:       pulseaudio-libs, systemd
Recommends:     electron

%description
pulse-remote lets you control PulseAudio/PipeWire audio
from any device on your network via a web interface.

Includes a systemd user service (pulse-remote.service) and
an optional Electron desktop launcher (pulse-remote-desktop).

The desktop launcher requires Electron to be installed
separately. Install it from https://github.com/electron/electron
or your package manager.

%prep
%setup -q

%build
# Pre-built binary package — nothing to compile

%install
cp -a ./* %{buildroot}/

%files
/usr/bin/pulse-remote-server
/usr/bin/pulse-remote-desktop
/usr/lib/pulse-remote/desktop/app.asar
/usr/lib/pulse-remote/desktop/app.asar.unpacked/
/usr/share/applications/pulse-remote.desktop
/usr/share/icons/hicolor/256x256/apps/pulse-remote.png
/usr/lib/systemd/user/pulse-remote.service
/usr/share/man/man1/pulse-remote.1
/usr/share/licenses/pulse-remote/LICENSE

%post
if command -v systemctl >/dev/null 2>&1; then
	systemctl --user daemon-reload 2>/dev/null || true
fi

%preun
if command -v systemctl >/dev/null 2>&1; then
	systemctl --user stop pulse-remote.service 2>/dev/null || true
fi
