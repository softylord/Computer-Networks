#!/bin/bash
export GOPATH=`pwd`
go get github.com/go-sql-driver/mysql
go get github.com/mmcdole/gofeed
go install ./src/client
go install ./src/spamer
