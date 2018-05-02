#!/bin/bash

WORKDIR=github.com/edrans/cloudflarebeat

docker build -t cloudflarebeat-build:latest --build-arg WORKDIR=${WORKDIR} -f build/Dockerfile .
docker run -v "$PWD":/go/src/github.com/edrans/cloudflarebeat cloudflarebeat-build:latest
