#!/bin/bash
# alias gr='bash -x ~/github/vasnake/go/mrg_course/golang-web_services-mrg_course/run.sh'

# program to run

# app_name=hello_world
app_name=spec
# dir=week_01
# fname=hello_world.go
# fname=vars_2.go

# dir=week_02
# fname=visibility/main.go # package!
# fname=afterfunc.go

# run
__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# pushd ${__dir}/${dir}
pushd ${__dir}/sandbox/

# export GOROOT=/Users/${USER}/go
# export GOPATH=${__dir}/${dir}

# fix formatting
# gofmt -w $fname || exit
gofmt -w $app_name || exit

# go run $fname || exit
# GO111MODULE=off go run $fname || exit
go run $app_name

# experiments

#GO111MODULE=off cat data_map.txt | go run $fname || exit
#GO111MODULE=off cat data_map.txt | sort | go run $fname || exit
#GO111MODULE=off go test -v ./unique || exit

# echo "go run ${fname}: SUCCESS"
popd
