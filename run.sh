#!/bin/bash

dir=week_02
#fname=visibility/main.go
fname=goroutines.go
#fname=chan_1.go

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd ${__dir}/${dir}

#export GOROOT=/Users/${USER}/go
#export GOPATH=${__dir}/${dir}

gofmt -w $fname || exit
#go run $fname || exit
GO111MODULE=off go run $fname || exit

#GO111MODULE=off cat data_map.txt | go run $fname || exit
#GO111MODULE=off cat data_map.txt | sort | go run $fname || exit
#GO111MODULE=off go test -v ./unique || exit

echo "Test run OK"
popd
