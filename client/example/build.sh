#/bin/sh -e

BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
GitHash=`git rev-parse --short HEAD`
Version="0.1.0-dev"

go build -v -ldflags="-X main.GitHash=${GitHash} -X main.BuildTime=${BuildTime} -X main.Version=${dev}" -o client main.go