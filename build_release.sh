#!/bin/sh
set -ex
go fmt ./...
golint *.go
goimports -w *.go
go build -o cl.mac ./...
docker run -v /Users/ross/go/src/github.com/crewjam/cl:/go/src/github.com/crewjam/cl golang \
  bash -c 'cd /go/src/github.com/crewjam/cl && go get ./... && go build -o cl.linux ./...'
