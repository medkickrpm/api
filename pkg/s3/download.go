package s3

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func DownloadFile(filename string) (*s3.GetObjectOutput, error) {
	downloader := s3.New(sess)

	resp, err := downloader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filename),
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
