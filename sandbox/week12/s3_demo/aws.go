package s3_demo

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func MainS3Demo() {
	// https://aws.github.io/aws-sdk-go-v2/docs/migrating/

	// s3Config := &aws.Config{
	// 	Credentials:      credentials.NewStaticCredentials("access_123", "secret_123", ""),
	// 	Endpoint:         aws.String("http://127.0.0.1:9000"),
	// 	Region:           aws.String("us-east-1"),
	// 	DisableSSL:       aws.Bool(true),
	// 	S3ForcePathStyle: aws.Bool(true),
	// }
	// newSession := session.New(s3Config)
	// cfg, err := config.LoadDefaultConfig(context.TODO())
	// if err != nil {	}
	// s3Client := s3.New(newSession)

	ctx := context.TODO()

	s3Config, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("access_123", "secret_123", "")),
		config.WithRegion("us-east-1"),
	)
	panicOnError("LoadDefaultConfig failed", err)

	s3Client := s3.NewFromConfig(s3Config, func(opts *s3.Options) {
		opts.BaseEndpoint = aws.String("http://127.0.0.1:9000")
		opts.EndpointOptions.DisableHTTPS = true
		opts.UsePathStyle = true
	})

	// list existing buckets

	input := &s3.ListBucketsInput{}
	buckets, err := s3Client.ListBuckets(ctx, input)
	panicOnError("ListBuckets", err)

	fmt.Println("Existing buckets:")
	for _, bucket := range buckets.Buckets {
		show("bucket: ", bucket)
	}

	// create bucket

	bucketName := "photolist"
	bucket := aws.String(bucketName)

	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: bucket})
	if err != nil {
		// panicOnError("CreateBucket error", err)
		var apierr smithy.APIError
		if errors.As(err, &apierr) {
			if apierr.ErrorCode() == "BucketAlreadyOwnedByYou" {
				show("bucket exists already: ", *bucket)
			} else {
				panicOnError("CreateBucket error, APIError unknown code", err)
			}
		} else {
			panicOnError("CreateBucket error, not APIError", err)
		}
		// var awsErr awserr.Error
		// if errors.As(err, &awsErr) && awsErr.Code() == "BucketAlreadyOwnedByYou" {
		// 	log.Printf("%s already exists\n", bucketName)
		// } else {
		// 	log.Fatalln("cant create bucker", bucketName, err)
		// }
	}

	// upload file

	file, err := os.Open("/tmp/building_1.jpg")
	panicOnError("os.Open failed", err)

	objectName := "building_1.jpg"
	contentType := "image/jpeg"
	res, err := s3Client.PutObject(ctx, &s3.PutObjectInput{
		Body:        file,
		Bucket:      bucket,
		Key:         aws.String(objectName),
		ContentType: aws.String(contentType),
	})
	panicOnError("PutObject failed", err)
	log.Printf("Successfully uploaded %s, res %d\n", objectName, res)

	// download file

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: bucket,
		Key:    aws.String(objectName),
	})
	panicOnError("GetObject failed", err)
	defer result.Body.Close()

	hasher := md5.New()
	io.Copy(hasher, result.Body)
	log.Printf("download file with md5sum: %x\n", hasher.Sum(nil))

	// set policy

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
	_, err = s3Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: bucket,
		Policy: aws.String(policy),
	})
	panicOnError("PutBucketPolicy failed", err)

	/*
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
	   _, err = s3Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
	   	Bucket: bucket,
	   	Policy: aws.String(policy),
	   })
	   if err != nil {
	   	log.Fatalln("PutBucketPolicy err", err)
	   }
	*/
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
