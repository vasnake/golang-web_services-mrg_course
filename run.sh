#!/bin/bash

dir=week_01
#fname=visibility/main.go
fname=empty_2.go

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd ${__dir}/${dir}

#export GOROOT=/Users/${USER}/go
#export GOPATH=${__dir}/${dir}

gofmt -w $fname || exit
#go run $fname || exit
GO111MODULE=off go run $fname || exit

echo "Test run OK"
popd
