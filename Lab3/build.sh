#!/bin/bash
export GOPATH=`pwd`
go get github.com/go-sql-driver/mysql
go get github.com/mmcdole/gofeed
go get github.com/mmcdole/gofeed/rss
go install ./src/install
