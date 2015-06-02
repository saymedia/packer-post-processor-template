#!/bin/bash


# Cleanup
rm -f *_test_out*.tfvars
rm -f crash.log
rm -f packer-post-processor-template


# Install deps
go get


# Testing
go test ./...
if [[ $? != 0 ]]; then
    exit 1
fi


# Linting
go vet ./...
if [[ $? != 0 ]]; then
    exit 1
fi

golint ./...
if [[ $? != 0 ]]; then
    exit 1
fi

gocyclo -over 15 .
if [[ $? != 0 ]]; then
    exit 1
fi


# Formatting
gofmt -s -d -l */*.go
if [[ $? != 0 ]]; then
    exit 1
fi


# Build
go build
if [[ $? != 0 ]]; then
    exit 1
fi


# Run
PACKER_LOG=1 $GOPATH/bin/packer build -var-file=amazon_test_variables.json amazon_test.json
if [[ $? != 0 ]]; then
    exit 1
fi

PACKER_LOG=1 TMPDIR=~/tmp $GOPATH/bin/packer build docker_test.json
if [[ $? != 0 ]]; then
    exit 1
fi
