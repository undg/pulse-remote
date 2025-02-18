# Change these variables as necessary.
MAIN_PACKAGE_PATH := .
BINARY_NAME := pulse-remote-server
PKG_NAME := pulse-remote
SERVICE_NAME := pulse-remote.service

BUILD_TIME=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short=7 HEAD)
GIT_VERSION=$(shell git describe --tags --abbrev=0 | tr -d '\n')

BUILD_PKG_PATH=github.com/undg/go-prapi/buildinfo

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

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: tidy/ci
tidy/ci:
	tidy no-dirty

## audit: run quality control checks
.PHONY: audit/ci
audit/ci:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: audit
audit/full:
	tidy audit/ci test

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

## build/fe: get latest frontend from github and build in build/pr-web/dist
.PHONY: build/fe
build/fe:
	mkdir -p build

	# cd web/ && \
	# pnpm install && \
	# pnpm build
	# cd -

	cp -r web/dist build/web

## build/be: build the application
.PHONY: build/be
build/be: 
	go build -ldflags=${LDFLAGS} -o=build/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## build/clear: remove build/ directory for fresh start
.PHONY: build/clear
build/clear:
	# Delete build folder
	rm -rf build/

## build: build the application together with frontend
.PHONY: build
build:
	make build/clear
	make build/fe
	make build/be

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

.PHONY: install
install:

	make build
	@systemctl --user is-active ${SERVICE_NAME} >/dev/null 2>&1 && systemctl --user stop ${SERVICE_NAME} || true

	sudo install -Dm755 build/bin/${BINARY_NAME} /usr/bin/${BINARY_NAME}
	sudo install -Dm644 ${SERVICE_NAME} /usr/lib/systemd/user/${SERVICE_NAME}
	sudo install -Dm644 "LICENSE" "/usr/share/licenses/${PKG_NAME}/LICENSE"
	sudo install -Dm644 "pulse-remote.1" "/usr/share/man/man1/pulse-remote.1"

	sudo systemctl daemon-reload

	systemctl --user enable pulse-remote
	systemctl --user start pulse-remote

.PHONY: uninstall
uninstall:
	@systemctl --user is-active ${SERVICE_NAME} >/dev/null 2>&1 && systemctl --user stop ${SERVICE_NAME} || true
	systemctl --user disable ${SERVICE_NAME} 

	sudo rm /usr/bin/${BINARY_NAME}
	sudo rm /usr/lib/systemd/user/${SERVICE_NAME}
	sudo rm /usr/share/licenses/${PKG_NAME}/LICENSE
	sudo rm /usr/share/man/man1/pulse-remote.1

	systemctl --user daemon-reload


