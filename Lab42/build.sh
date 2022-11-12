#!/bin/bash
export GOPATH=`pwd`
go get github.com/mgutz/logxi/v1
go get github.com/gliderlabs/ssh
go get golang.org/x/crypto/ssh
go install ./src/server
go install ./src/client
