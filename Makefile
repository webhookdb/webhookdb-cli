BIN := ./webhookdb
ARGS := WEBHOOKDB_API_HOST=http://localhost:18001
BUILDFLAGS = "-X github.com/webhookdb/webhookdb-cli/config.BuildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/webhookdb/webhookdb-cli/config.BuildSha=`git rev-list -1 HEAD`"
WEBSITE = ../webhookdb-api/webhookdb-website

ifdef GOROOT
GO := $(GOROOT)/bin/go
else
GO := go
endif

guardcmd-%:
	@hash $(*) > /dev/null 2>&1 || \
		(echo "ERROR: '$(*)' must be installed and available on your PATH."; exit 1)

guardenv-%:
	@if [ -z '${${*}}' ]; then echo 'ERROR: environment variable $* not set' && exit 1; fi

fmt:
	@go fmt ./...

lint: guardcmd-gofmt
	@test -z $$(gofmt -d -l . | tee /dev/stderr) && echo "gofmt ok"

vet:
	@$(GO) vet && echo "go vet ok"

test:
	$(GO) test .
	@WEBHOOKDB_LOG_LEVEL=fatal ginkgo -r --trace --race --progress --skipMeasurements

test-watch:
	@WEBHOOKDB_LOG_LEVEL=fatal ginkgo watch ./...

update-test-snapshots:
	UPDATE_SNAPSHOTS=true make test

bench:
	@WEBHOOKDB_LOG_LEVEL=fatal ginkgo -r --focus=benchmarks

check: vet lint

build:
	@$(GO) build -ldflags $(BUILDFLAGS) -o webhookdb

build-arm64:
	@GOOS=darwin $(GO) build -ldflags $(BUILDFLAGS) -o webhookdb

build-wasm:
	@GOOS=js GOARCH=wasm $(GO) build -ldflags $(BUILDFLAGS) -o webhookdb.wasm

_goreleaser-clean:
	rm -rf ./dist

goreleaser-build: _goreleaser-clean
	goreleaser build
goreleaser-local: _goreleaser-clean
	GITHUB_TOKEN=notreal goreleaser

wasm-server:
	$(GO) run bin/serve-wasm/main.go

build-all: build-arm64 build build-wasm

docs: build
	@DOCBUILD=true $(BIN) docs build
docs-write: build ## Write a new copy of MANUAL.md.
	@DOCBUILD=true $(BIN) docs build > MANUAL.md
docs-site: build ## Write MANUAL.md to the docs site.
	@DOCBUILD=true $(BIN) docs build --docsite > ../webhookdb/docs/docs/cli-reference.md

update-lithic-deps:
	$(GO) get github.com/rgalanakis/golangal@latest
	$(GO) get github.com/lithictech/go-aperitif@latest

help:
	@$(GO) run ./main.go help

itest-auth-login: build
	$(ARGS) $(BIN) auth login --username=alpha@webhookdb.com

itest-auth-otp-%: build
	$(ARGS) $(BIN) auth login --username=alpha@webhookdb.com --token=$(*)

itest-auth-logout: build
	$(ARGS) $(BIN) auth logout

#
#itest-integrations-create: build
#	$(ARGS) $(BIN) integrations create fake_v1

# SERVICES

itest-services-list: build
	$(ARGS) $(BIN) services list

# ORGS

itest-org-invite-%: build
	$(ARGS) $(BIN) org invite --username=$(*)

itest-org-join-%: build
	$(ARGS) $(BIN) org join $(*)

itest-org-list: build
	$(ARGS) $(BIN) org list

itest-org-members: build
	$(ARGS) $(BIN) org members

itest-org-members-%: build
	$(ARGS) $(BIN) org members --org=$(*)
