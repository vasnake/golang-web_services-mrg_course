пример с микросервисом авторизации

смысл:
* проверяем авторизацию в 1 месте из разных сервисов, можем горизонатльно масштабироваться (увеличивать количество серверов)
* скрываем детали реализации хранения - теперь это может быть мапка в памяти, мемкеш, таранутл, файлы, база, libastral

1. надо скачать protoc (https://github.com/google/protobuf/releases)
2. go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
3. go get -u google.golang.org/grpc
4. go get -u golang.org/x/net/context

Генерация кода:
* находясь в папке session сгенерируем код для го `protoc --go_out=plugins=grpc:. *.proto`
* подобной командой так же генерируется нужная обвязка для других поддерживаемых языков
* go_out означает что мы хотим сгенерировать код в этой папке для языка go
* plugins=grpc созначает что мы хотим использовать ещё плагин для генерации grpc-сервиса

дополнительная документация
* https://developers.google.com/protocol-buffers/docs/gotutorial
* https://github.com/grpc/grpc-go/tree/master/examples
* https://habrahabr.ru/company/infopulse/blog/265805/

windows:
Download protoc-win32.zip from https://developers.google.com/protocol-buffers/docs/downloads
Unzip and add location of the protoc.exe to your PATH environment variable
Run `protoc --version` from command prompt to verify
Verify the your GOPATH environment variable is set
Run `go get -u github.com/golang/protobuf/protoc-gen-go` from command prompt. This should install the binary to %GOPATH%/bin
Add `%GOPATH%/bin` to your PATH environment variable
Open a new command prompt, navigate to your .proto file, run `protoc --go_out=. *.proto` 

если ругается - надо читобы были доступны protoc.exe и protoc-gen-go.exe - прописать в PATH путь до них

# В 2024 году кое-что изменилось

лог моей борьбы за кодогенерацию
```s
// https://github.com/protocolbuffers/protobuf/releases/tag/v26.1
// https://github.com/protocolbuffers/protobuf/releases/download/v26.1/protoc-26.1-linux-x86_64.zip
// c:\bin\protoc-26.1-linux-x86_64\bin\protoc
// export PATH=${PATH}:/mnt/c/bin/protoc-26.1-linux-x86_64/bin:${HOME}/go/bin
// pushd sandbox/week07/grpc_1/
// go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
// go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
// protoc --go-grpc_out=. *.proto # service
// protoc --go_out=. --go-grpc_out=. *.proto # service + messages

https://grpc.io/docs/languages/go/quickstart/#prerequisites

$ protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto

```
