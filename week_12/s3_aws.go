package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// connect

	bucketName := "photolist"
	bucket := aws.String(bucketName) // NB datatype wrappers for convinient use of strings and other types

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("access_123", "secret_123", ""),
		Endpoint:         aws.String("http://127.0.0.1:9000"), // minio local service, S3 imitation
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	// list

	input := &s3.ListBucketsInput{}
	buckets, err := s3Client.ListBuckets(input)
	if err != nil {
		log.Fatalln("ListBuckets err:", err)
	}

	fmt.Println("Existing buckets:")
	for _, bucket := range buckets.Buckets {
		fmt.Println(bucket)
	}

	// create

	_, err = s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: bucket,
	})
	if err != nil {
		// NB check error type and content
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "BucketAlreadyOwnedByYou" {
			// log.Printf("%s already exists\n", bucketName)
		} else {
			log.Fatalln("cant create bucker", bucketName, err)
		}
	}

	// upload

	objectName := "building_1.jpg"
	contentType := "image/jpeg"
	file, err := os.Open("../photo_samples/building_1.jpg")
	if err != nil {
		log.Fatalln("cant open file:", err)
	}

	res, err := s3Client.PutObject(&s3.PutObjectInput{
		Body:        file,
		Bucket:      bucket,
		Key:         aws.String(objectName),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Fatalln("PutObject err:", err)
	}

	log.Printf("Successfully uploaded %s, res %d\n", objectName, res)

	// download

	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: bucket,
		Key:    aws.String(objectName),
	})
	if err != nil {
		log.Fatalln("GetObject ewr:", err)
	}
	defer result.Body.Close()

	hasher := md5.New()
	io.Copy(hasher, result.Body)
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
	_, err = s3Client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: bucket,
		Policy: aws.String(policy),
	})
	if err != nil {
		log.Fatalln("PutBucketPolicy err", err)
	}

}

func stuff() {

	// downloader := s3manager.NewDownloader(newSession)
	// numBytes, err := downloader.Download(file,
	// 	&s3.GetObjectInput{
	// 		Bucket: bucket,
	// 		Key:    aws.String(objectName),
	// 	})
	// if err != nil {
	// 	fmt.Println("Failed to download file", err)
	// 	return
	// }
	// fmt.Println("Downloaded file", file.Name(), numBytes, "bytes")

	// file.Seek(0, io.SeekStart)
	// objectName = "2_" + objectName
	// res, err = s3Client.PutObject(&s3.PutObjectInput{
	// 	Body:        file,
	// 	Bucket:      bucket,
	// 	Key:         aws.String(objectName),
	// 	ContentType: aws.String(contentType),
	// })
	// if err != nil {
	// 	log.Fatalln("PutObject err:", err)
	// }
	// log.Printf("Successfully uploaded %s, res %d\n", objectName, res)

	// result, err := s3Client.GetBucketPolicy(&s3.GetBucketPolicyInput{
	// 	Bucket: bucket,
	// })
	// fmt.Println(result, err)
}
