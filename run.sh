#!/bin/bash

dir=week_01/visibility
fname=main.go

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd ${__dir}/${dir}

#export GOROOT=/Users/${USER}/go
export GOPATH=${__dir}/${dir}

gofmt -w $fname || exit
go run $fname || exit

echo "Test run OK"
popd
