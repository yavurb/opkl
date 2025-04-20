#!/usr/bin/env just --justfile

# recipe to display help information
help:
  @just --list

test:
  go test -v ./...

