#!/bin/bash
# alias gr='bash -vxe /mnt/c/Users/valik/data/github/golang-web_services-mrg_course/run.sh'
PRJ_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

APP_SELECTOR=${GO_APP_SELECTOR:-week02}

go_run() {
    local selector="${1}"
    case $selector in
        hello_world | helloworld)   go_run_sandbox hello_world;;
        spec)                       go_run_sandbox spec;;
        week01)                     go_run_sandbox week01;;
        week01_test)                go_run_sandbox_test week01;;
        week01_tree_test)           go_run_sandbox_week01_tree_test;;
        week02)                     go_run_sandbox week02;;
        *)                          errorExit "Unknown program: ${selector}";;
    esac
}

go_run_sandbox_week01_tree_test() {
    local module="tree"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week01_homework

    gofmt -w $module || exit
    go vet $module

    export RECURSIVE_TREE=no

    go run $module tree/testdata
    go run $module tree/testdata -f
    exit_code=$?

    go test -v $module
    exit_code=$?

    cd $module
    docker build -t mailgo_hw1 .

    popd
    return $exit_code
}

go_run_sandbox_test() {
    local module="${1}"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox

    gofmt -w $module || exit
    go vet $module

    go test -v $module
    exit_code=$?

    popd
    return $exit_code
}

go_run_sandbox() {
    local module="${1}"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox

    gofmt -w $module || exit
    go vet $module    
    # go vet -stringintconv=false spec

    go run $module    
    exit_code=$?

    popd
    return $exit_code
}

errorExit() {
    echo "$1" 1>&2
    exit 1
}

go_run ${APP_SELECTOR}
exit_code=$?
echo "Exit code: ${exit_code}"
exit $exit_code

installGo() {
    # https://go.dev/doc/install
    pushd /mnt/c/Users/valik/Downloads
    sudo rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
    popd
    # vim ~/.profile
}

# deprecated

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
