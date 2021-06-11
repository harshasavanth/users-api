package aws

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/harshasavanth/bookstore_users-api/logger"
	"github.com/harshasavanth/users-api/utils/rest_errors"
	"mime/multipart"
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

func Upload(file multipart.File, fileHeader *multipart.FileHeader, userid string) (filename string, error *rest_errors.RestErr) {
	s := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})))

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)
	_, err := s.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("profilepicsupload"),
		Key:    aws.String(userid + ".jpg"),
		Body:   bytes.NewReader(buffer),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", rest_errors.NewInternalServerError("error while uplpoading image")
		logger.Info(err.Error())
	}
	return fmt.Sprintf("https://profilepicsupload.s3.ap-south-1.amazonaws.com/%s.jpg", userid), nil

}
