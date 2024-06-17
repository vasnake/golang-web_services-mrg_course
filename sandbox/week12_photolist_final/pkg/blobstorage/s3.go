package blobstorage

import (
	"errors"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Storage struct {
	session *session.Session
	client  *s3.S3
	bucket  *string
}

func NewS3Storage(host, access, secret, bucketName string) (*S3Storage, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(access, secret, ""),
		Endpoint:         aws.String(host),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}

	storage := &S3Storage{
		bucket: aws.String(bucketName),
	}

	storage.session = session.New(s3Config)
	storage.client = s3.New(storage.session)

	_, err := storage.client.CreateBucket(&s3.CreateBucketInput{
		Bucket: storage.bucket,
	})
	if err != nil {
		var awsErr awserr.Error
		if errors.As(err, &awsErr) && awsErr.Code() == "BucketAlreadyOwnedByYou" {
			// do nothing, pre-setup
		} else {
			return nil, err
		}
	}

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
	_, err = storage.client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: storage.bucket,
		Policy: aws.String(policy),
	})
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (storage *S3Storage) Put(data io.ReadSeeker, objectName, contentType string, userID uint32) error {
	_, err := storage.client.PutObject(&s3.PutObjectInput{
		Body:        data,
		Bucket:      storage.bucket,
		Key:         aws.String(objectName),
		ContentType: aws.String(contentType),
		Metadata: map[string]*string{
			"user-id": aws.String(strconv.Itoa(int(userID))),
		},
	})
	return err
}
