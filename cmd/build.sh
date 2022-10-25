#!/usr/bin/env bash
set -xe

export GOARCH=amd64
export GOOS=linux
export GCCGO=gc

version=$1

if [ -z $version ]; then
    version=v0.1
fi

go build -o bandwagonhost_cloud_exporter bandwagonhost_cloud_exporter.go
chmod +x bandwagonhost_cloud_exporter

tar -zcvf bandwagonhost_cloud_exporter-linux-amd64-${version}.tar.gz \
  bandwagonhost_cloud_exporter config/bandwagonhost_cloud_exporter.yaml config/bandwagonhost_cloud_exporter.service README.md
