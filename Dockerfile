FROM golang:1.17

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

ARG GIT_SHA="-"
ARG GIT_REF="-"

RUN go build -ldflags "-X github.com/webhookdb/webhookdb-cli/config.BuildTime=`date +"%Y-%m-%dT%H:%M:%SZ"` -X github.com/webhookdb/webhookdb-cli/config.BuildSha=${GIT_SHA} -X github.com/webhookdb/webhookdb-cli/config.Version=${GIT_REF}" -o webhookdb

ENTRYPOINT ["./webhookdb"]
