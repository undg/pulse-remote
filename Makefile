# Change these variables as necessary.
MAIN_PACKAGE_PATH := .
BINARY_NAME := pulse-remote-server
PKG_NAME := pulse-remote
SERVICE_NAME := pulse-remote.service
MAN_NAME := pulse-remote.1

BUILD_TIME=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short=7 HEAD)
GIT_VERSION=$(shell git describe --tags --abbrev=0 | tr -d '\n')

BUILD_PKG_PATH=github.com/undg/go-prapi/api/buildinfo

LDFLAGS="-X '${BUILD_PKG_PATH}.GitVersion=${GIT_VERSION}' \
				-X '${BUILD_PKG_PATH}.BuildTime=${BUILD_TIME}' \
				-X '${BUILD_PKG_PATH}.GitCommit=${GIT_COMMIT}'"

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
no-dirty:
	git diff --exit-code

# generate_pactl_type: Generate Go struct from pactl JSON output
# $(1): pactl command (e.g., "list sinks", "list sources")
# $(2): type name (e.g., "sink", "source")
#
# Usage: $(call generate_pactl_type,<pactl_command>,<type_name>)
define generate_pactl_type
	# Run pactl, extract first item, generate Go struct
	pactl --format=json $(1) | jq '.[0]' | gojsonstruct \
		--package-name=pactl \
		--typename=Pactl$(shell echo '$(2)' | sed 's/./\U&/')JSON \
		--file-header="//lint:file-ignore ST1003 Ignore underscore naming in generated code" \
		--int-type=float64 \
		--o pactl/generated/$(2)-type.go
	@echo "Manual adjustment needed in pactl/generated/$(2)-type.go for accurate types"
endef

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## update-go: update Go version to latest available in mise
.PHONY: update-go
update-go:
	./scripts/update-go.sh

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: tidy/ci
tidy/ci:
	make tidy
	make no-dirty

## audit: run quality control checks
.PHONY: audit/ci
audit/ci:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: audit
audit:
	make tidy
	make audit/ci 
	make test

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/watch: run all tests in watch mode
.PHONY: test/watch
test/watch:
	./scripts/test-watch.sh

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

# ==================================================================================== #
# UTILS
# ==================================================================================== #

.PHONY: sink-type
sink-type:
	go install github.com/twpayne/go-jsonstruct/v3/cmd/gojsonstruct@latest
	$(call generate_pactl_type,list sinks,sink)

sink-item-type:
	go install github.com/twpayne/go-jsonstruct/v3/cmd/gojsonstruct@latest
	# ffplay -nodisp -autoexit -f lavfi -i "anullsrc=r=44100:cl=stereo" -loglevel quiet &
	$(call generate_pactl_type,list sink-inputs,apps)
	# killall ffplay

.PHONY source-type:
source-type:
	go install github.com/twpayne/go-jsonstruct/v3/cmd/gojsonstruct@latest
	$(call generate_pactl_type,list sources,source)

## typesgen: generate structs from json output
.PHONY: typesgen
typesgen:
	sink-type source-type tidy

## push: push changes to the remote Git repository
.PHONY: push
push:
	tidy audit no-dirty
	git push

.PHONY: bump/patch
bump/patch:
	./scripts/bump.sh patch

.PHONY: bump/minor
bump/minor:
	./scripts/bump.sh minor

.PHONY: bump/main
bump/main:
	./scripts/bump.sh main

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## pull/web: get latest frontend from github and build in _GUI/web/dist
.PHONY: pull/web 
pull/web:
	rm -rf /tmp/build/pulse-remote-web
	mkdir -p /tmp/build/pulse-remote-web
	git clone "https://github.com/undg/pulse-remote-web" /tmp/build/pulse-remote-web

	cd /tmp/build/pulse-remote-web/ && \
	pnpm install && \
	pnpm test:ci && \
	pnpm build

	cd -

	mkdir -p _GUI/web/

	# pulse-remote-web old version:
	cat _GUI/web/version

	git  -C /tmp/build/pulse-remote-web describe --long --abbrev=7 --tags | sed 's/^v//;s/\([^-]*-g\)/r\1/;s/-/./g' > _GUI/web/version

	cp -r /tmp/build/pulse-remote-web/dist _GUI/web

	git reset
	git add _GUI/web
	git add _GUI/web/version
	git commit -m "Update web to version $$(cat _GUI/web/version)"

	## pulse-remote-web new version:
	cat _GUI/web/version

## pull/desktop: get latest desktop frontend, build asar, and copy into _GUI/desktop/
.PHONY: pull/desktop 
pull/desktop:
	rm -rf /tmp/build/pulse-remote-desktop
	mkdir -p /tmp/build/pulse-remote-desktop
	git clone "https://github.com/undg/pulse-remote-desktop" /tmp/build/pulse-remote-desktop

	cd /tmp/build/pulse-remote-desktop/ && \
    mise trust && \
    NODE_OPTIONS="--max-old-space-size=8192" pnpm install && \
    pnpm build:unpacked

	cd -

	mkdir -p _GUI/desktop/
	cp /tmp/build/pulse-remote-desktop/dist/linux-unpacked/resources/app.asar _GUI/desktop/
	cp -r /tmp/build/pulse-remote-desktop/dist/linux-unpacked/resources/app.asar.unpacked _GUI/desktop/
	cp /tmp/build/pulse-remote-desktop/resources/icon.png _GUI/desktop/icon.png


## build: build the backend server
.PHONY: build
build: 
	rm -rf build/bin
	go build -ldflags=${LDFLAGS} -o=build/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## package: create release tarball with full FHS install tree (for aur-bin)
.PHONY: package
package: build
	make install DESTDIR=/tmp/pulse-pkg PREFIX=/usr
	PKG_VER=$$(echo "$(GIT_VERSION)" | sed 's/^v//'); \
	tar czf "pulse-remote_$${PKG_VER}_Linux_x86_64.tar.gz" -C /tmp/pulse-pkg .
	rm -rf /tmp/pulse-pkg; \
	sha256sum pulse-remote_*.tar.gz > checksums.txt; \
	echo "Release artifacts: pulse-remote_$${PKG_VER}_Linux_x86_64.tar.gz checksums.txt"

## deb: build .deb package
.PHONY: deb
deb: build
	PKG_VER=$$(echo "$(GIT_VERSION)" | sed 's/^v//'); \
	rm -rf /tmp/pulse-deb; \
	make install DESTDIR=/tmp/pulse-deb/pkg PREFIX=/usr; \
	mkdir -p /tmp/pulse-deb/pkg/DEBIAN; \
	sed "s/@VERSION@/$${PKG_VER}/" os/deb/control > /tmp/pulse-deb/pkg/DEBIAN/control; \
	cp os/deb/postinst /tmp/pulse-deb/pkg/DEBIAN/postinst; \
	cp os/deb/prerm /tmp/pulse-deb/pkg/DEBIAN/prerm; \
	chmod 755 /tmp/pulse-deb/pkg/DEBIAN/postinst /tmp/pulse-deb/pkg/DEBIAN/prerm; \
	dpkg-deb --build /tmp/pulse-deb/pkg "pulse-remote_$${PKG_VER}_amd64.deb"; \
	rm -rf /tmp/pulse-deb; \
	echo "Release artifact: pulse-remote_$${PKG_VER}_amd64.deb"

## rpm: build .rpm package (requires rpmbuild)
.PHONY: rpm
rpm: build
	PKG_VER=$$(echo "$(GIT_VERSION)" | sed 's/^v//'); \
	rm -rf /tmp/pulse-rpm; \
	mkdir -p /tmp/pulse-rpm/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}; \
	make install DESTDIR=/tmp/pulse-rpm/buildroot PREFIX=/usr; \
	mv /tmp/pulse-rpm/buildroot "/tmp/pulse-rpm/pulse-remote-$${PKG_VER}"; \
	tar czf "/tmp/pulse-rpm/rpmbuild/SOURCES/pulse-remote-$${PKG_VER}.tar.gz" \
		-C /tmp/pulse-rpm "pulse-remote-$${PKG_VER}"; \
	sed "s/@VERSION@/$${PKG_VER}/" os/rpm/pulse-remote.spec \
		> /tmp/pulse-rpm/rpmbuild/SPECS/pulse-remote.spec; \
	rpmbuild -bb --define "_topdir /tmp/pulse-rpm/rpmbuild" \
		/tmp/pulse-rpm/rpmbuild/SPECS/pulse-remote.spec; \
	mv /tmp/pulse-rpm/rpmbuild/RPMS/x86_64/*.rpm .; \
	rm -rf /tmp/pulse-rpm; \
	@echo "Release artifacts: pulse-remote-$${PKG_VER}-1.*.rpm"

## packages: build all release artifacts (tarball + deb + rpm)
.PHONY: packages
packages: package deb rpm
	sha256sum pulse-remote_*.tar.gz pulse-remote_*_amd64.deb pulse-remote-*.x86_64.rpm > checksums.txt

## run: build and run the application
.PHONY: run
run:
	make build
	while true; do build/bin/${BINARY_NAME};sleep 1; done

## run/watch: run the application with reloading on file changes
.PHONY: run/watch
run/watch:
	go run github.com/cosmtrek/air@v1.43.0 \
		--build.cmd "make build" --build.bin "build/bin/${BINARY_NAME}" --build.delay "100" \
		--build.exclude_dir "" \
		--build.include_ext "go, tpl, tmpl, html, css, scss, js, ts, sql, jpeg, jpg, gif, png, bmp, svg, webp, ico" \
		--misc.clean_on_exit "true"


# ==================================================================================== #
# INSTALL
# ==================================================================================== #

PREFIX = /usr/local
DESTDIR =

## install: install all files to $(DESTDIR)$(PREFIX). Used by packaging and CI.
.PHONY: install
install: build
	install -Dm755 "build/bin/${BINARY_NAME}" "$(DESTDIR)$(PREFIX)/bin/${BINARY_NAME}"
	install -Dm755 "os/launcher.sh" "$(DESTDIR)$(PREFIX)/bin/pulse-remote-desktop"
	install -Dm644 "_GUI/desktop/app.asar" "$(DESTDIR)$(PREFIX)/lib/pulse-remote/desktop/app.asar"
	cp -r "_GUI/desktop/app.asar.unpacked" "$(DESTDIR)$(PREFIX)/lib/pulse-remote/desktop/"
	install -Dm644 "os/pulse-remote.desktop" "$(DESTDIR)$(PREFIX)/share/applications/pulse-remote.desktop"
	install -Dm644 "_GUI/desktop/icon.png" "$(DESTDIR)$(PREFIX)/share/icons/hicolor/256x256/apps/pulse-remote.png"
	install -Dm644 "os/${SERVICE_NAME}" "$(DESTDIR)$(PREFIX)/lib/systemd/user/${SERVICE_NAME}"
	install -Dm644 "os/${MAN_NAME}" "$(DESTDIR)$(PREFIX)/share/man/man1/${MAN_NAME}"
	install -Dm644 "LICENSE" "$(DESTDIR)$(PREFIX)/share/licenses/${PKG_NAME}/LICENSE"

## install-local: install directly to /usr and enable systemd service (dev convenience)
.PHONY: install-local
install-local:
	@systemctl --user is-active ${SERVICE_NAME} >/dev/null 2>&1 && systemctl --user stop ${SERVICE_NAME} || true
	sudo make install PREFIX=/usr
	sudo systemctl daemon-reload
	systemctl --user enable --now ${SERVICE_NAME}

## uninstall: remove all installed files
.PHONY: uninstall
uninstall:
	@systemctl --user is-active ${SERVICE_NAME} >/dev/null 2>&1 && systemctl --user stop ${SERVICE_NAME} || true
	systemctl --user disable ${SERVICE_NAME}

	sudo rm -f "$(DESTDIR)$(PREFIX)/bin/${BINARY_NAME}"
	sudo rm -f "$(DESTDIR)$(PREFIX)/bin/pulse-remote-desktop"
	sudo rm -rf "$(DESTDIR)$(PREFIX)/lib/pulse-remote"
	sudo rm -f "$(DESTDIR)$(PREFIX)/share/applications/pulse-remote.desktop"
	sudo rm -f "$(DESTDIR)$(PREFIX)/share/icons/hicolor/256x256/apps/pulse-remote.png"
	sudo rm -f "$(DESTDIR)$(PREFIX)/lib/systemd/user/${SERVICE_NAME}"
	sudo rm -f "$(DESTDIR)$(PREFIX)/share/man/man1/${MAN_NAME}"
	sudo rm -rf "$(DESTDIR)$(PREFIX)/share/licenses/${PKG_NAME}"

	systemctl --user daemon-reload
