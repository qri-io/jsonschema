#!/usr/bin/env bash

set -e
echo "" > coverage.txt
go test -v -race -coverprofile=profile.out -covermode=atomic github.com/qri-io/varName
if [ -f profile.out ]; then
  cat profile.out >> coverage.txt
  rm profile.out
fi