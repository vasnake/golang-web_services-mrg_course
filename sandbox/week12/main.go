package main

import (
	"week12/s3_demo"

	s3_photolist_main "week12/s3_photolist/cmd/photolist"

	// s3_acl "week12/s3_images_nginx_acl_photolist/cmd/photolist"
	// panic: proto: message session.AuthSession is already registered
	// See https://protobuf.dev/reference/go/faq#namespace-conflict

	photolist_tracing_ctx "week12/photolist_tracing_request_id/cmd/photolist"
)

func main() {
	// demo1()
	// demo2()
	// demo_103_images()
	// demo_104_acl()

	// demo_105_ctx()

	// # last version of photolist app
	demo_106_tracing_jaeger()
}

func demo1() {
	// s3 provider: minio s3 api
	s3_demo.MainS3Demo() // sandbox/week12/s3_demo$  docker compose -f ./docker-compose.yaml up&
	/*
	   Existing buckets:
	   2024-06-07T09:38:15.732Z: bucket: types.Bucket{CreationDate:time.Date(2024, time.June, 7, 9, 26, 44, 192000000, time.UTC), Name:(*string)(0xc000119e90), noSmithyDocumentSerde:document.NoSerde{}};
	   2024-06-07T09:38:15.734Z: smithy.APIError: "BucketAlreadyOwnedByYou"; "";
	   2024/06/07 12:38:15 Successfully uploaded building_1.jpg, res &{0 0 0 0 0 824636039984 0 %!d(types.RequestCharged=) 0 0 0 0 %!d(types.ServerSideEncryption=) 0 {map[{}:-783728287 {}:824634876576 {}:%!d(string=17D6AF56EB31E36B) {}:{13947880188658761375 59194494 12617440} {}:{0 63853349895 0} {}:{[{<nil> %!d(bool=false) %!d(bool=false) {map[{}:-783728287 {}:824634876576 {}:%!d(string=17D6AF56EB31E36B) {}:{13947880188658761375 59194494 12617440} {}:{0 63853349895 0} {}:%!d(string=dd9025bab4ad464b049177c95eb6ebf374d3b3fd1af9251148b658df7ac2e3e8)]}}]} {}:%!d(string=dd9025bab4ad464b049177c95eb6ebf374d3b3fd1af9251148b658df7ac2e3e8)]} {}}
	   2024/06/07 12:38:15 download file with md5sum: 93aaabaf6c9afc54965d721f108474df
	   see sandbox\week12\s3_demo\minio_data\photolist\building_1.jpg\2f65fa59-4166-4634-a340-6cef9cd87e0d\part.1
	*/
}

func demo2() {
	// s3 provider: minio api
	s3_demo.MainMinioDemo() // sandbox/week12/s3_demo$  docker compose -f ./docker-compose.yaml up&
	/*
	   2024-06-07T10:10:24.529Z: Existing buckets: []minio.BucketInfo{minio.BucketInfo{Name:"photolist", CreationDate:time.Date(2024, time.June, 7, 9, 26, 44, 192000000, time.UTC)}};
	   2024-06-07T10:10:24.529Z: bucket: minio.BucketInfo{Name:"photolist", CreationDate:time.Date(2024, time.June, 7, 9, 26, 44, 192000000, time.UTC)};
	   2024-06-07T10:10:24.531Z: bucket exists already: "photolist";
	   2024/06/07 13:10:24 Successfully uploaded building_1.jpg of size 204361
	   2024/06/07 13:10:24 download file with md5sum: 93aaabaf6c9afc54965d721f108474df
	*/
}

func demo_103_images() {
	// images storage: s3
	s3_photolist_main.MainDemo()
}

func demo_104_acl() {
	// config: viper
	// images from s3 via nginx + custom auth (images ACL)
	// add 2 auth services: photoauth (for nginx images acl); auth nano-service "user-sessions-db" (grpc)
	/*
	   работать это должно так:

	   в докере запускаются nginx и хранилища (mysql, s3)
	   sandbox\week12\s3_images_nginx_acl_photolist\deployments\docker-compose.yml

	   на хосте в трех разных терминалах запускаются три процесса (см конфиги sandbox\week12\s3_images_nginx_acl_photolist\configs\):
	   - sandbox\week12\s3_images_nginx_acl_photolist\cmd\auth\main.go
	   - sandbox\week12\s3_images_nginx_acl_photolist\cmd\photoauth\main.go
	   - sandbox\week12\s3_images_nginx_acl_photolist\cmd\photolist\main.go

	   браузером ходить надо на localhost:8080 - это nginx as reverse proxy, container:80
	   sandbox\week12\s3_images_nginx_acl_photolist\configs\nginx\nginx.conf
	   https://stackoverflow.com/questions/31324981/how-to-access-host-port-from-docker-container
	*/

	// pushd sandbox
	// docker compose -f ./week12/s3_images_nginx_acl_photolist/deployments/docker-compose.yml up

	// set -a && source week12/s3_images_nginx_acl_photolist/configs/common.env && set +a

	// set -a && source week12/s3_images_nginx_acl_photolist/configs/auth.env && set +a
	// go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" ./week12/s3_images_nginx_acl_photolist/cmd/auth/main.go

	// set -a && source week12/s3_images_nginx_acl_photolist/configs/photoauth.env && set +a
	// go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" ./week12/s3_images_nginx_acl_photolist/cmd/photoauth/main.go

	// export OAUTH_APP_SECRET=a***0
	// export OAUTH_APP_ID=O***F
	// GO_APP_SELECTOR=week12 gr # former: go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" ./week12/s3_images_nginx_acl_photolist/cmd/photolist/main.go -appid ${OAUTH_APP_ID:-foo} -appsecret ${OAUTH_APP_SECRET:-bar}

	// s3_acl.MainDemo()
	panic("uncomment: s3_acl.MainDemo()")
	/*
	   go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" ./week12/s3_images_nginx_acl_photolist/cmd/auth/main.go
	   2024/06/14 11:47:46 [startup 1] service 'auth', -ldflags info: buildHash '9c81fe0', buildTime '2024-06-14_08:47:43'
	   2024/06/14 11:47:46 [startup 2] service 'auth', 'go build -buildvcs' info: Version 'unknown', Revision 'unknown', DirtyBuild '%!s(bool=true)', LastCommit '0001-01-01 00:00:00 +0000 UTC', ShortInfo 'devel'
	   2024-06-14T08:47:46.154Z: service.port from config: "localhost:10000";
	   2024-06-14T08:47:46.154Z: sql.Open mysql DSN: "root:@tcp(localhost:3306)/photolist?charset=utf8&interpolateParams=true";
	   2024/06/14 11:47:46 [startup 6] grpc serve at tcp port 'localhost:10000'

	   go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" ./week12/s3_images_nginx_acl_photolist/cmd/photoauth/main.go
	   2024/06/14 11:48:08 [startup] photoauth, commit 9c81fe0, build 2024-06-14_07:48:39
	   2024/06/14 11:48:08 [startup] cfg.HTTP.Port "localhost:8081", example.env1 "", example.env2 "env config value"
	   2024-06-14T08:48:08.333Z: downstream svc session.grpc_addr from config: "localhost:10000";
	   2024-06-14T08:48:08.333Z: listen http.port from config: "localhost:8081";
	   2024-06-14T08:48:08.333Z: sql.Open mysql DSN: "root:@tcp(localhost:3306)/photolist?charset=utf8&interpolateParams=true";
	   2024/06/14 11:48:08 [startup] listening server at localhost:8081
	   2024/06/14 11:57:07 call UserRepository.IsFollowed - maybe user dataloader? 1 2
	   2024/06/14 11:57:24 call UserRepository.IsFollowed - maybe user dataloader? 1 2

	   go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" ./week12/s3_images_nginx_acl_photolist/cmd/photolist/main.go -appid ${OAUTH_APP_ID:-foo} -appsecret ${OAUTH_APP_SECRET:-bar}
	   2024/06/14 11:48:16 [startup] photolist, commit 9c81fe0, build 2024-06-14_07:55:02
	   2024-06-14T08:48:16.504Z: you must not show this! appid, appsecret: "O***F"; "a***0";
	   2024-06-14T08:48:16.506Z: listen http.port from config: "localhost:8082";
	   2024-06-14T08:48:16.506Z: sql.Open mysql DSN: "root:@tcp(localhost:3306)/photolist?charset=utf8&interpolateParams=true";
	   2024/06/14 11:48:16 [startup] listening server at localhost:8082
	   2024-06-14T08:56:20.096Z: oauth code: "a***9";
	   2024-06-14T08:56:20.675Z: oauth access token: &oauth2.Token{AccessToken:...
	   2024-06-14T08:56:21.270Z: api response: "{\"login\":...
	   2024-06-14T08:56:21.271Z: user email from oauth provider: ...
	   2024-06-14T08:56:21.271Z: user id from oauth provider: "2dv0h";
	   2024-06-14T08:56:21.271Z: creating app user ...
	   2024/06/14 11:56:21 call UserRepository.IsFollowed - maybe user dataloader? 1 2
	   2024/06/14 11:56:54 call UserRepository.IsFollowed - maybe user dataloader? 1 2
	   2024/06/14 11:57:07 call UserRepository.IsFollowed - maybe user dataloader? 1 2
	   2024/06/14 11:57:18 call UserRepository.IsFollowed - maybe user dataloader? 1 2
	   2024/06/14 11:57:24 call UserRepository.IsFollowed - maybe user dataloader? 1 2
	*/
}

func demo_105_ctx() {
	// tracing, request id, request context
	// log graphql operations (name, timing, path) via middleware

	_ = `
работать это должно так:

в докере запускаются nginx и хранилища (mysql, s3)
sandbox\week12\photolist_tracing_request_id\deployments\docker-compose.yml

на хосте в трех разных терминалах запускаются три процесса,
см. конфиги sandbox\week12\photolist_tracing_request_id\configs\
- sandbox\week12\photolist_tracing_request_id\cmd\auth\main.go
- sandbox\week12\photolist_tracing_request_id\cmd\photoauth\main.go
- sandbox\week12\photolist_tracing_request_id\cmd\photolist\main.go

браузером ходить надо на localhost:8080 - это nginx as reverse proxy
sandbox\week12\photolist_tracing_request_id\configs\nginx\nginx.conf
	`

	_ = `
pushd sandbox/week12/photolist_tracing_request_id
docker compose -f deployments/docker-compose.yml up

set -a && source configs/common.env && set +a
set -a && source configs/auth.env && set +a
go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" cmd/auth/main.go

set -a && source configs/common.env && set +a
set -a && source configs/photoauth.env && set +a
go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" cmd/photoauth/main.go

set -a && source configs/common.env && set +a
export OAUTH_APP_SECRET=a***0
export OAUTH_APP_ID=O***F
alias gr='bash -vxe golang-web_services/run.sh'
GO_APP_SELECTOR=week12 gr
# go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" cmd/photolist/main.go -appid ${OAUTH_APP_ID:-foo} -appsecret ${OAUTH_APP_SECRET:-bar}
	`

	photolist_tracing_ctx.MainDemo()
	_ = `
deployments-nginx-1    | 172.18.0.1 [17/Jun/2024:09:50:21 +0000] "GET /images/1/cf140d24-6984-4489-a53f-3a71fa623754_600.jpg HTTP/1.1" 304 516 "http://localhost:8080/photos/foo" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36" 0.046 a785ca0e76325b93042da8fa4b73acd8
deployments-nginx-1    | 172.18.0.1 [17/Jun/2024:09:50:24 +0000] "GET /photos/foo HTTP/1.1" 200 2361 "http://localhost:8080/photos/" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36" 0.009 f49b3c51fcf0797b7e1b850f4e4cb275
deployments-nginx-1    | 172.18.0.1 [17/Jun/2024:09:50:24 +0000] "GET /static/css/bootstrap/bootstrap.min.css.map HTTP/1.1" 404 223 "-" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36" 0.009 536127a450ae6c890991c5fc52a95985
...
deployments-nginx-1    | 172.18.0.1 [17/Jun/2024:09:52:04 +0000] "POST /graphql HTTP/1.1" 200 715 "http://localhost:8080/photos/" "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36" 0.052 de3622907e6d9ca26ff124a64addc468

go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" cmd/auth/main.go
2024/06/17 12:47:33 [startup] auth, commit 53cd5c7, build 2024-06-17_09:47:17
2024-06-17T09:47:33.902Z: service.port from config: "localhost:10000"; 
2024-06-17T09:47:33.902Z: sql.Open mysql DSN: "root:@tcp(localhost:3306)/photolist?charset=utf8&interpolateParams=true"; 
2024/06/17 12:47:33 [startup] listening server at localhost:10000
2024/06/17 12:50:21 [access] a785ca0e76325b93042da8fa4b73acd8 4.467128ms /session.Auth/Check '<nil>'
2024/06/17 12:52:04 [access] c830d6e230b5c61a59123b4efa6d8d95 1.886608ms /session.Auth/Check '<nil>'
2024/06/17 12:52:46 [access] ce0a56db60b6119d0975bde34942214e 1.297505ms /session.Auth/Check '<nil>'
2024/06/17 12:52:46 [access] 96cefb122b09b1aad7616341a22751eb 5.745257ms /session.Auth/Check '<nil>'

go run -ldflags "-X 'main.buildHash=${APP_COMMIT}' -X 'main.buildTime=${APP_BUILD_TIME}'" cmd/photoauth/main.go
2024/06/17 12:48:42 [startup] photoauth, commit 53cd5c7, build 2024-06-17_09:48:28
2024/06/17 12:48:42 [startup] cfg.HTTP.Port "localhost:8081", example.env1 "", example.env2 "example env2 config value"
2024-06-17T09:48:42.116Z: downstream svc session.grpc_addr from config: "localhost:10000"; 
2024-06-17T09:48:42.116Z: listen http.port from config: "localhost:8081"; 
2024-06-17T09:48:42.116Z: sql.Open mysql DSN: "root:@tcp(localhost:3306)/photolist?charset=utf8&interpolateParams=true"; 
2024/06/17 12:48:42 [startup] listening server at localhost:8081
2024/06/17 12:50:21 call UserRepository.IsFollowed - maybe user dataloader? 1 2
2024/06/17 12:50:21 [access] a785ca0e76325b93042da8fa4b73acd8 14.476252ms 127.0.0.1:52954 GET /api/v1/internal/images/auth
2024/06/17 12:52:04 [access] c830d6e230b5c61a59123b4efa6d8d95 4.169965ms 127.0.0.1:43846 GET /api/v1/internal/images/auth
2024/06/17 12:52:46 [access] ce0a56db60b6119d0975bde34942214e 2.810057ms 127.0.0.1:53306 GET /api/v1/internal/images/auth

go run week12 -appid Ov23lirslzXRwbCt2gJF -appsecret ada307c051a73e56fd7a8287b95e2f6f79aef860
2024/06/17 12:50:07 [startup] photolist, commit unknown, build unknown
2024-06-17T09:50:07.376Z: you must not show this! appid, appsecret: "O***F"; "a***0";
2024-06-17T09:50:07.380Z: listen http.port from config: "localhost:8082";
2024-06-17T09:50:07.380Z: sql.Open mysql DSN: "root:@tcp(localhost:3306)/photolist?charset=utf8&interpolateParams=true";
2024/06/17 12:50:07 [startup] listening server at localhost:8082
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 1.530794ms user '<nil>'
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 833ns user.name '<nil>'
...
2024/06/17 12:50:21 call UserRepository.IsFollowed - maybe user dataloader? 1 2
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 579.664µs me.followedUsers[0].followed '<nil>'
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 2.402446ms user.photos[0].user '<nil>'
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 295ns user.photos[0].user.name '<nil>'
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 408ns user.photos[0].user.avatar '<nil>'
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 1.327µs user.photos[0].user.id '<nil>'
2024/06/17 12:50:21 [resolver] a8ef85e77aeae6a391eed380ce0f8d45 7.889µs user.photos[0].user.followed '<nil>'
2024/06/17 12:50:21 [RequestMiddleware] a8ef85e77aeae6a391eed380ce0f8d45 7.297232ms  29
2024/06/17 12:50:21 [access] a8ef85e77aeae6a391eed380ce0f8d45 13.376275ms 127.0.0.1:40128 POST /graphql
...
2024/06/17 12:53:41 [resolver] 9227f31fd292b8789a5d1ae87d008fe7 6.427361ms ratePhoto '<nil>'
2024/06/17 12:53:41 [resolver] 9227f31fd292b8789a5d1ae87d008fe7 462ns ratePhoto.rating '<nil>'
2024/06/17 12:53:41 [resolver] 9227f31fd292b8789a5d1ae87d008fe7 1.196µs ratePhoto.id '<nil>'
2024/06/17 12:53:41 [RequestMiddleware] 9227f31fd292b8789a5d1ae87d008fe7 6.514879ms rateCommentToggle 3
2024/06/17 12:53:41 [access] 9227f31fd292b8789a5d1ae87d008fe7 7.679408ms 127.0.0.1:41406 POST /graphql

	`
}

func demo_106_tracing_jaeger() {
	// distributed tracing: photoauth <-grpc-> auth
	// open tracing, open telemetry, jaeger
	// request id, span, middleware
	// samplerconfig
	msg := `

Это последняя в курсе версия photolist,
поэтому деплой я сделал как в курсе: внутри контейнера (docker compose).
См. sandbox\week12_photolist_final\Makefile

pushd sandbox/week12_photolist_final/
go mod tidy

`
	show(msg)
}
