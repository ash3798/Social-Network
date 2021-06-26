#!/bin/bash
export GOPATH=/

cd /Social-Network

echo "Building the Binary"
go build

echo "Starting Social Network"
./Social-Network