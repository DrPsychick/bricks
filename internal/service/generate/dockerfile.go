// Copyright © 2018 by PACE Telematics GmbH. All rights reserved.
// Created at 2018/08/31 by Vincent Landgraf

package generate

import (
	"html/template"
	"log"
	"os"
)

// DockerfileOptions configure the output of the generated docker
// file
type DockerfileOptions struct {
	Commands CommandOptions
}

// Dockerfile generate a dockerfile using the given options
// for specified path
func Dockerfile(path string, options DockerfileOptions) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	err = dockerTemplate.Execute(f, options)
	if err != nil {
		log.Fatal(err)
	}
}

var dockerTemplate = template.Must(template.New("Dockerfile").Parse(
	`FROM golang:1.11 as builder
RUN go get gopkg.in/alecthomas/gometalinter.v2
RUN gometalinter.v2 --install
WORKDIR /tmp/service
ADD . .

# Lin, vet & test
# (many linters from gometalinter don't support go mod and therefore need to be disabled)
RUN gometalinter.v2 --disable-all --vendor -E gocyclo -E goconst -E golint -E ineffassign -E gotypex -E deadcode ./... && \
	go vet -mod vendor ./... && \
	go test -mod vendor -v -race -cover ./...

# Build
RUN go install -mod vendor ./cmd/{{ .Commands.DaemonName }} && \
	go install -mod vendor ./cmd/{{ .Commands.ControlName }}

FROM alpine
RUN apk update && apk add ca-certificates && apk add tzdata && rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/{{ .Commands.DaemonName }} /usr/local/bin/
COPY --from=builder /go/bin/{{ .Commands.ControlName }} /usr/local/bin/

EXPOSE 3000
ENV PORT 3000
ENTRYPOINT ["/usr/local/bin/{{ .Commands.DaemonName }}"]
`))
