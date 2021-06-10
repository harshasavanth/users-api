package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/harshasavanth/bookstore_users-api/logger"
	"github.com/harshasavanth/users-api/utils/rest_errors"
	"os"
)

var (
	s3session *s3.S3
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})))
}

func LinseningBuckets() (resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		logger.Info(err.Error())
		//rest_errors.NewInternalServerError("unable to up")
	}
	return resp
}

func Upload(file string, userid string) *rest_errors.RestErr {
	f, err := os.Open(file)
	if err != nil {
		logger.Info(err.Error())
		return rest_errors.NewBadRequestError(err.Error())
	}
	resp, err := s3session.PutObject(&s3.PutObjectInput{
		Body:   f,
		Bucket: aws.String("profilepicsupload"),
		Key:    aws.String(userid),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		logger.Info("error")
		logger.Info(err.Error())
		return rest_errors.NewInternalServerError(err.Error())
	}
	logger.Info("uploaded")
	logger.Info(fmt.Sprintf("%s", resp))
	return nil
}
