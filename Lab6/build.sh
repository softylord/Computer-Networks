#!/bin/bash
export GOPATH=`pwd`
go get github.com/jlaffaye/ftp
go get github.com/goftp/server
go get github.com/goftp/file-driver
go get github.com/mmcdole/gofeed
go install ./src/client
go install ./src/server
