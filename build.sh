#!/bin/bash

set -e

VERSION=0.0.1
docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.6 bash -c "go get -d && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bitesize-authz-webhook -v"
docker build -t geribatai/bitesize-authz-webhook:$VERSION .
docker push geribatai/bitesize-authz-webhook:$VERSION
