#!/bin/bash
export GOPATH=`pwd`
go get github.com/gliderlabs/ssh
go install ./src/ssh-server
go install ./src/http-server
