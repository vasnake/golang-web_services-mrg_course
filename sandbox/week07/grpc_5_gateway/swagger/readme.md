Изначально в этом фолдере ничего нет, кроме ридми.
Джейсон берем из предыдущего урока.

По джейсону генерим код клиента для сервиса.

По этому коду создаем (нет, уже для нас создали) код consumer.go для демонстрации.

Что делать в командной строке, для запуска демо
```s
pushd sandbox/week07/grpc_5_gateway/swagger/

go get -u github.com/go-swagger/go-swagger/cmd/swagger
go install github.com/go-swagger/go-swagger/cmd/swagger

swagger serve ../session_grpc5.swagger.json -p 8082 --no-open # http://localhost:8082/docs

rm -rf sess-client &&\
  mkdir -p sess-client &&\
  swagger generate client -f ../session_grpc5.swagger.json -A sess-client/ -t ./sess-client/

# don't need this, server exists already
rm -rf sess-server &&\
  mkdir -p sess-server &&\
  swagger generate server -f ../session_grpc5.swagger.json -A sess-server/ -t ./sess-server/
```
