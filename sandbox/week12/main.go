package main

import "week12/s3_demo"

func main() {
	s3_demo.MainS3Demo() // sandbox/week12/s3_demo$  docker compose -f ./docker-compose.yaml up&
	/*
	   Existing buckets:
	   2024-06-07T09:38:15.732Z: bucket: types.Bucket{CreationDate:time.Date(2024, time.June, 7, 9, 26, 44, 192000000, time.UTC), Name:(*string)(0xc000119e90), noSmithyDocumentSerde:document.NoSerde{}};
	   2024-06-07T09:38:15.734Z: smithy.APIError: "BucketAlreadyOwnedByYou"; "";
	   2024/06/07 12:38:15 Successfully uploaded building_1.jpg, res &{0 0 0 0 0 824636039984 0 %!d(types.RequestCharged=) 0 0 0 0 %!d(types.ServerSideEncryption=) 0 {map[{}:-783728287 {}:824634876576 {}:%!d(string=17D6AF56EB31E36B) {}:{13947880188658761375 59194494 12617440} {}:{0 63853349895 0} {}:{[{<nil> %!d(bool=false) %!d(bool=false) {map[{}:-783728287 {}:824634876576 {}:%!d(string=17D6AF56EB31E36B) {}:{13947880188658761375 59194494 12617440} {}:{0 63853349895 0} {}:%!d(string=dd9025bab4ad464b049177c95eb6ebf374d3b3fd1af9251148b658df7ac2e3e8)]}}]} {}:%!d(string=dd9025bab4ad464b049177c95eb6ebf374d3b3fd1af9251148b658df7ac2e3e8)]} {}}
	   2024/06/07 12:38:15 download file with md5sum: 93aaabaf6c9afc54965d721f108474df
	   see sandbox\week12\s3_demo\minio_data\photolist\building_1.jpg\2f65fa59-4166-4634-a340-6cef9cd87e0d\part.1
	*/

	// s3_demo.MainMinioDemo() // sandbox/week12/s3_demo$  docker compose -f ./docker-compose.yaml up&
	/*
	   2024-06-07T10:10:24.529Z: Existing buckets: []minio.BucketInfo{minio.BucketInfo{Name:"photolist", CreationDate:time.Date(2024, time.June, 7, 9, 26, 44, 192000000, time.UTC)}};
	   2024-06-07T10:10:24.529Z: bucket: minio.BucketInfo{Name:"photolist", CreationDate:time.Date(2024, time.June, 7, 9, 26, 44, 192000000, time.UTC)};
	   2024-06-07T10:10:24.531Z: bucket exists already: "photolist";
	   2024/06/07 13:10:24 Successfully uploaded building_1.jpg of size 204361
	   2024/06/07 13:10:24 download file with md5sum: 93aaabaf6c9afc54965d721f108474df
	*/

}
