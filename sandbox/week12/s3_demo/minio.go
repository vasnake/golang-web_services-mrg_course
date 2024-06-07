package s3_demo

import (
	"crypto/md5"
	"errors"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go"
)

func MainMinioDemo() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "access_123"
	secretAccessKey := "secret_123"
	location := "us-east-1"
	useSSL := false

	minioClient, err := minio.NewWithRegion(endpoint,
		accessKeyID, secretAccessKey,
		useSSL, location)
	panicOnError("NewWithRegion failed", err)

	// list buckets

	buckets, err := minioClient.ListBuckets()
	panicOnError("ListBuckets failed", err)

	show("Existing buckets: ", buckets)
	for _, bucket := range buckets {
		show("bucket: ", bucket)
	}

	// create bucket

	bucketName := "photolist"
	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		var merr minio.ErrorResponse
		if errors.As(err, &merr) && merr.Code == "BucketAlreadyOwnedByYou" {
			show("bucket exists already: ", bucketName)
		} else {
			panicOnError("MakeBucket failed", err)
		}
	} else {
		show("bucket created: ", bucketName)
	}

	// upload file

	file, err := os.Open("/tmp/building_1.jpg")
	panicOnError("os.Open failed", err)

	objectName := "building_1.jpg"
	contentType := "image/jpeg"

	size, err := minioClient.PutObject(bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	panicOnError("PutObject failed", err)
	log.Printf("Successfully uploaded %s of size %d\n", objectName, size)

	// download file

	reader, err := minioClient.GetObject(bucketName, objectName, minio.GetObjectOptions{})
	panicOnError("GetObject failed", err)
	defer reader.Close()

	hasher := md5.New()
	io.Copy(hasher, reader)
	log.Printf("download file with md5sum: %x\n", hasher.Sum(nil))

	// set bucket policy

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
	panicOnError("SetBucketPolicy failed", err)
	show("policy is set: ", policy)
}
