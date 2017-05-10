#!/usr/bin/env sh
set -x
export GOPATH=$HOME/gocode
go build  -o br -ldflags "-X main.Build=`date -u +.%Y%m%d%.H%M%S` -X main.Revision=`git rev-parse HEAD` -X main.Version=`git describe --abbrev=0 --tags`" .
GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Build=`date -u +.%Y%m%d%.H%M%S` -X main.Revision=`git rev-parse HEAD` -X main.Version=`git describe --abbrev=0 --tags`" .
docker build -t www.stockpalm.com/breeze/core:`git describe --abbrev=0 --tags` .
docker push www.stockpalm.com/breeze/core:`git describe --abbrev=0 --tags`
docker tag www.stockpalm.com/breeze/core:`git describe --abbrev=0 --tags` www.stockpalm.com/breeze/core:latest
docker push  www.stockpalm.com/breeze/core:latest
