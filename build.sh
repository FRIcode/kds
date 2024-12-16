#!/bin/bash

set -x

GOOS=linux GOARCH=amd64 go build -o target/kds-amd64-linux
GOOS=linux GOARCH=arm64 go build -o target/kds-arm64-linux
