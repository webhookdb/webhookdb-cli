BIN := ./webhookdb
ARGS := API_HOST=http://localhost:19001

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
	@go vet

test:
	go test .
	@LOG_LEVEL=fatal ginkgo -r --trace --race --progress --skipMeasurements

test-watch:
	@LOG_LEVEL=fatal ginkgo watch ./...

update-test-snapshots:
	UPDATE_SNAPSHOTS=true make test

bench:
	@LOG_LEVEL=fatal ginkgo -r --focus=benchmarks

build:
	@go build -ldflags \
		"-X github.com/lithictech/webhookdb-cli/config.BuildTime=`date -u +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/lithictech/webhookdb-cli/config.BuildSha=`git rev-list -1 HEAD`" \
		-o webhookdb

update-lithic-deps:
	go get github.com/rgalanakis/golangal@latest
	go get github.com/lithictech/go-aperitif@latest

help:
	@go run ./main.go help

itest-auth-register: build
	$(ARGS) $(BIN) auth register --username=x@y.com
itest-auth-login: build
	$(ARGS) $(BIN) auth login --username=natalie@lithic.tech

itest-auth-otp-%: build
	$(ARGS) $(BIN) auth otp --username=natalie@lithic.tech --token=$(*)

itest-auth-logout: build
	$(ARGS) $(BIN) auth logout


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
