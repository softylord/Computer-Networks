#!/bin/bash
export GOPATH=`pwd`
go get github.com/gorilla/websocket
go install ./src/server
go install ./src/client
