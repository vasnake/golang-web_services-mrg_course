#!/bin/bash
# alias gr='bash -vxe /mnt/c/Users/valik/data/github/golang-web_services-mrg_course/run.sh'
PRJ_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# PATH=${PATH}:/mnt/c/bin/protoc-26.1-linux-x86_64/bin:${HOME}/go/bin

APP_SELECTOR=${GO_APP_SELECTOR:-week09}

go_run() {
    local selector="${1}"
    case $selector in
        hello_world | helloworld)   go_run_sandbox hello_world;;
        spec)                       go_run_sandbox spec;;
        week01)                     go_run_sandbox week01;;
        week01_test)                go_run_sandbox_test week01;;
        week01_tree_test)           go_run_sandbox_week01_tree_test;;
        week02)                     go_run_sandbox week02;;
        week02_signer_test)         go_run_sandbox_week02_signer_test;;
        week03)                     go_run_sandbox week03;;
        week03_finder_test)         go_run_sandbox_week03_finder_test;;
        week04)                     go_run_sandbox week04;;
        week04_homework)            go_run_sandbox_week04_search_test;;
        week05)                     go_run_sandbox week05;;
        week05_homework)            go_run_sandbox_week05_codegen_test;;
        week06)                     go_run_sandbox week06;;
        week06_homework)            go_run_sandbox_week06_db_explorer_test;;
        week07)                     go_run_sandbox week07;;
        week07_homework)            go_run_sandbox_week07_async_logger_test;;
        week08)                     go_run_sandbox week08;;
        week08_homework)            go_run_sandbox_week08_i2s_test;;
        week09)                     go_run_sandbox week09;;
        *)                          errorExit "Unknown program: ${selector}";;
    esac
}

go_run_module() {
    echo "####################################################################################################"
    go run \
        -ldflags="-X 'main.Version=$(git rev-parse HEAD)' -X 'main.Branch=$(git rev-parse --abbrev-ref HEAD)'" \
        ${1} \
        --comments=true --servers="127.0.0.1:8081,127.0.0.1:8082"
    exit_code=$?
    echo "####################################################################################################"
    return $exit_code
}

go_test_module() {
    echo "####################################################################################################"
    go test -v ${1}
    exit_code=$?
    echo "####################################################################################################"
    return $exit_code
}

go_run_sandbox_week08_i2s_test(){
    create_project(){ # need this only once
        pushd ${PRJ_DIR}/sandbox
        mkdir -p week08_homework/i2s
        pushd week08_homework/i2s

        if [ -f ./main.go ]; then
            echo "project created already"; exit 42
        else
            echo "creating project ..."
        fi

        go mod init i2s

        # makes go vet happy
cat > main.go << EOT
        package main
        func main() { panic("not yet") }
EOT

        go mod tidy

        popd # workspace
        # go work init
        # go work edit -dropuse=./week05_homework/code_gen
        go work use ./week08_homework/i2s
    }
    # create_project

    local module="i2s"
    local moduleDir=${PRJ_DIR}/sandbox/week08_homework/${module}
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week08_homework

    pushd ${moduleDir} && go mod tidy && popd
    gofmt -w $module || exit
    go vet $module
    # go vet -stringintconv=false $module # go doc cmd/vet
    # golangci-lint run $module # slow

    echo "####################################################################################################"
    # go test -v --failfast $module # during development: failfast
    go test -v $module
    # go test -v -race $module # final check

    # go run -race $module
    # go run $module
    # go_run_module $module

    # https://pkg.go.dev/cmd/go#hdr-Testing_flags
    # go test -bench . $module
    # go test -bench . -benchmem $module
    # go test -bench '.*Mem.*' -benchmem $module
    # go test -bench '.*Xml.*' -benchmem $module

    # go test -v -cover $module
    # go test -coverprofile=cover.out $module
    # go tool cover -html=cover.out -o cover.html

    exit_code=$?
    echo "####################################################################################################"
    popd    
    return $exit_code  
}
go_run_sandbox_week07_async_logger_test(){
    local module="async_logger"
    local moduleDir=${PRJ_DIR}/sandbox/week07_homework/${module}
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week07_homework

    create_project(){ # need this only once
        create project
        pushd ${PRJ_DIR}/sandbox # workspace
        mkdir -p week07_homework/async_logger
        pushd week07_homework/async_logger
        go mod init async_logger
        cat > main.go << EOT
package main
func main() { panic("not yet") }
EOT
        go mod tidy
        popd # workspace
        go work use ./week07_homework/async_logger
    }

    pushd ${moduleDir} && go mod tidy && popd
    gofmt -w $module || exit
    go vet $module
    # go vet -stringintconv=false $module # go doc cmd/vet

    echo "####################################################################################################"
    # go test -v --failfast $module # during development: failfast
    go test -v -race $module # final check

    # go run -race $module
    # go run $module
    # go_run_module $module

    # https://pkg.go.dev/cmd/go#hdr-Testing_flags
    # go test -bench . $module
    # go test -bench . -benchmem $module
    # go test -bench '.*Mem.*' -benchmem $module
    # go test -bench '.*Xml.*' -benchmem $module

    # go test -v -cover $module
    # go test -coverprofile=cover.out $module
    # go tool cover -html=cover.out -o cover.html

    exit_code=$?
    echo "####################################################################################################"
    popd    
    return $exit_code
}

go_run_sandbox_week06_db_explorer_test() {
    local module="db_explorer"
    local moduleDir=${PRJ_DIR}/sandbox/week06_homework/${module}
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week06_homework

    gofmt -w $module || exit
    go vet $module
    # go vet -stringintconv=false $module # go doc cmd/vet

    echo "####################################################################################################"
    go test -v -race $module
    # go test -v $module

    # go run -race $module
    # go run $module
    # go_run_module $module

    # https://pkg.go.dev/cmd/go#hdr-Testing_flags
    # go test -bench . $module
    # go test -bench . -benchmem $module
    # go test -bench '.*Mem.*' -benchmem $module
    # go test -bench '.*Xml.*' -benchmem $module

    # go test -v -cover $module
    # go test -coverprofile=cover.out $module
    # go tool cover -html=cover.out -o cover.html

    exit_code=$?
    echo "####################################################################################################"
    popd    
    return $exit_code    
}

go_run_sandbox_week05_codegen_test() {
    local module="codegen"
    local moduleDir=${PRJ_DIR}/sandbox/week05_homework/${module}
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week05_homework
    
    # rm codegen/api_handlers.go || echo $? # drop generated code

    gofmt -w $module || exit
    go vet $module
    # go vet -stringintconv=false $module # go doc cmd/vet

    # pushd ${moduleDir}/handlers_gen && gofmt -w ./ && go vet -printf=false && popd
    pushd ${moduleDir}/handlers_gen && gofmt -w ./ && go vet && popd
    pushd ${moduleDir} && go build handlers_gen/* && ./codegen api.go api_handlers.go && popd
    rm ${moduleDir}/codegen # drop binary file

    go test -v $module

    # go run -race $module
    # go run $module
    # go_run_module $module

    # https://pkg.go.dev/cmd/go#hdr-Testing_flags
    # go test -bench . $module
    # go test -bench . -benchmem $module
    # go test -bench '.*Mem.*' -benchmem $module
    # go test -bench '.*Xml.*' -benchmem $module

    # go test -v -cover $module
    # go test -coverprofile=cover.out $module
    # go tool cover -html=cover.out -o cover.html

    exit_code=$?
    echo "####################################################################################################"
    popd    
    return $exit_code
}

go_run_sandbox_week04_search_test() {
    local module="search"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week04_homework

    gofmt -w $module || exit
    go vet $module
    # go vet -stringintconv=false $module

    # go run -race $module
    # go run $module
    # go_run_module $module

    # https://pkg.go.dev/cmd/go#hdr-Testing_flags
    # go test -bench . $module
    # go test -bench . -benchmem $module
    # go test -bench '.*Mem.*' -benchmem $module
    # go test -bench '.*Xml.*' -benchmem $module
    # go test -v $module

    go test -v -cover $module
    # go test -coverprofile=cover.out $module
    # go tool cover -html=cover.out -o cover.html

    exit_code=$?
    echo "####################################################################################################"
    popd
    return $exit_code
}

go_run_sandbox_week03_finder_test() {
    local module="finder"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week03_homework

    gofmt -w $module || exit
    go vet $module

    go test -v $module
    echo "####################################################################################################"
    go test -bench . -benchmem $module
    exit_code=$?
    echo "####################################################################################################"

    # go tool pprof -http=:8083 /path/to/bin /path/to/out

    popd
    return $exit_code    
}

go_run_sandbox_week02_signer_test() {
    local module="signer"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox/week02_homework

    gofmt -w $module || exit
    go vet $module

    # go run $module
    # go run -race $module
    # exit_code=$?

    # go test -v $module -run TestByIlia
    # go test -v $module -run "TestMultiHash"
    # go test -v $module -run "TestSingleHash"
    # go test -v $module -run 'Test.*Results'
    # go test -v $module -run TestSigner
    # go test -v $module
    go_test_module $module
    # time go test -v -race $module
    # time go test -v -race -parallel 8 -failfast $module
    exit_code=$?

    popd
    return $exit_code    
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

    # go test -v $module
    go_test_module $module
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

    # go test -v $module
    go_test_module $module
    exit_code=$?

    popd
    return $exit_code
}

go_run_sandbox() {
    local module="${1}"
    local exit_code=0
    pushd ${PRJ_DIR}/sandbox

    create_module(){
        pushd ${PRJ_DIR}/sandbox
        mkdir -p ${module} && pushd ./${module}

        if [ -f ./main.go ]; then
            echo "project created already"; exit 42
        else
            echo "creating project ..."
        fi

        go mod init ${module}
        cat > main.go << EOT
package main
func main() { panic("not yet") }
EOT
        go mod tidy
        popd # sandbox
        go work use ./${module}        
        gofmt -w ${module}
        # go vet ${module}
        go test -v ${module}
        go run ${module}
        exit 42
    }
    # create_module

    # pushd ${module} && docker compose up&; popd

    gofmt -w $module || exit
    go vet $module
    # go vet -stringintconv=false $module
    # golangci-lint run $module # slow

    # go build -o /tmp/$module $module

    # go run -race $module
    # go run $module
    go_run_module $module

    # https://pkg.go.dev/cmd/go#hdr-Testing_flags
    # go test -bench . $module
    # go test -bench . -benchmem $module
    # go test -bench '.*Mem.*' -benchmem $module
    # go test -bench '.*Xml.*' -benchmem $module
    # go test -v -cover $module
    # go test -v $module
    # go test -bench . -gcflags '-l' $module # отключаем инлайнинг

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
