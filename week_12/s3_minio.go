package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go"
)

func main() {
	// connect

	endpoint := "127.0.0.1:9000"
	accessKeyID := "access_123"
	secretAccessKey := "secret_123"
	useSSL := false
	location := "us-east-1"

	minioClient, err := minio.NewWithRegion(endpoint,
		accessKeyID, secretAccessKey,
		useSSL, location)
	if err != nil {
		log.Fatalln(err)
	}

	// list

	buckets, err := minioClient.ListBuckets()
	if err != nil {
		log.Fatalln("ListBuckets err:", err)
	}

	fmt.Println("Existing buckets:")
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}

	// create

	bucketName := "photolist"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		var minErr minio.ErrorResponse
		if errors.As(err, &minErr) && minErr.Code == "BucketAlreadyOwnedByYou" {
			// log.Printf("%s already exists\n", bucketName)
		} else {
			log.Fatalln("cant create bucker", bucketName, err)
		}
	} else {
		log.Printf("bucket %s created\n", bucketName)
	}

	// upload

	objectName := "building_1.jpg"
	contentType := "image/jpeg"
	file, err := os.Open("../photo_samples/building_1.jpg")
	if err != nil {
		log.Fatalln("cant open file:", err)
	}

	n, err := minioClient.PutObject(bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("PutObject err:", err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
	// file.Seek(0, io.SeekStart)

	// download

	reader, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln("GetObject err:", err)
	}
	defer reader.Close()

	hasher := md5.New()
	io.Copy(hasher, reader)
	log.Printf("download file with md5sum: %x\n", hasher.Sum(nil))

	// -----
	// return

	// set bucket policy, allow read files for anonymous, allow read files directly from storage

	policy := `{
		"Version":"2012-10-17",
		"Statement":[
		   {
			  "Action":["s3:GetObject"],
			  "Effect":"Allow",
			  "Principal":{"AWS":["*"]},
			  "Resource":["arn:aws:s3:::` + bucketName + `/*"],
			  "Sid":""
		   }
		]
	}`
	err = minioClient.SetBucketPolicy(bucketName, policy)
	if err != nil {
		log.Fatalln("SetBucketPolicy err", err)
	}
}
