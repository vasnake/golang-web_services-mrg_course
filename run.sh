#!/bin/bash

dir=week_01
fname=types.go

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
pushd ${__dir}/${dir}

gofmt -w $fname
go run $fname

popd
