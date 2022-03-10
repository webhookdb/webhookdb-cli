BIN := ./webhookdb
ARGS := WEBHOOKDB_API_HOST=http://localhost:18001
BUILDFLAGS = "-X github.com/lithictech/webhookdb-cli/config.BuildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/lithictech/webhookdb-cli/config.BuildSha=`git rev-list -1 HEAD`"
WEBSITE = ../webhookdb-api/webhookdb-website

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
	@go vet && echo "go vet ok"

test:
	go test .
	@WEBHOOKDB_LOG_LEVEL=fatal ginkgo -r --trace --race --progress --skipMeasurements

test-watch:
	@WEBHOOKDB_LOG_LEVEL=fatal ginkgo watch ./...

update-test-snapshots:
	UPDATE_SNAPSHOTS=true make test

bench:
	@WEBHOOKDB_LOG_LEVEL=fatal ginkgo -r --focus=benchmarks

check: vet lint

build:
	@go build -ldflags $(BUILDFLAGS) -o webhookdb

build-arm64:
	@GOOS=darwin go build -ldflags $(BUILDFLAGS) -o webhookdb

build-wasm:
	@GOOS=js GOARCH=wasm go build -ldflags $(BUILDFLAGS) -o webhookdb.wasm

build-all: build-arm64 build build-wasm

copy-to-web: ## Copy the WASM and MANUAL.md to the website directory.
	@cp webhookdb.wasm $(WEBSITE)/static/webterm
	@go run bin/copy-manual/main.go

docs-write: build ## Write a new copy of MANUAL.md.
	@$(BIN) docs build | grep -v '^%!(' > MANUAL.md

build-and-copy-to-web: docs-write build-wasm copy-to-web

update-lithic-deps:
	go get github.com/rgalanakis/golangal@latest
	go get github.com/lithictech/go-aperitif@latest

help:
	@go run ./main.go help

itest-auth-login: build
	$(ARGS) $(BIN) auth login --username=alpha@lithic.tech

itest-auth-otp-%: build
	$(ARGS) $(BIN) auth login --username=alpha@lithic.tech --token=$(*)

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
