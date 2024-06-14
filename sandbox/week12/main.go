package main

import (
	"week12/s3_demo"
	s3_acl "week12/s3_images_nginx_acl_photolist/cmd/photolist"
	s3_photolist_main "week12/s3_photolist/cmd/photolist"
)

func main() {
	// demo1()
	// demo2()
	// demo_103_images()

	// # I_AM_HERE
	demo_104_acl()

	// demo_105_ctx()
	// demo_106_tracing_jaeger()
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
	// add 2 auth services: photoauth (for nginx images acl); auth nano-service "user-sessions-db" (grpc)
	// images from s3 via nginx + custom auth (images ACL)
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

	s3_acl.MainDemo()
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
	// tracing, request id
	// log graphql operations (name, timing, path) via middleware
}

func demo_106_tracing_jaeger() {
	// distributed tracing: photoauth - auth grpc
	// open tracing, open telemetry, jaeger
	// request id, span, middleware
	// samplerconfig
}
