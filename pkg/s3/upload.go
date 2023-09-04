package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"os"
)

func UploadFile(filename string, file multipart.File) error {
	uploader := s3.New(sess)

	_, err := uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
		Body:   file,
	})

	return err
}
