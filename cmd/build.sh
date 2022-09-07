#!/usr/bin/env bash
set -xe

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

go build -o bandwagonhost_cloud_exporter bandwagonhost_cloud_exporter.go

tar -zcvf bandwagonhost_cloud_exporter-linux-amd64.tar.gz bandwagonhost_cloud_exporter config/bandwagonhost_cloud_exporter.yaml README.md
