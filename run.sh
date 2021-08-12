#!/bin/bash

dir=week_01
fname=methods.go

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd ${__dir}/${dir}

gofmt -w $fname || exit
go run $fname || exit

echo "Test run OK"
popd
