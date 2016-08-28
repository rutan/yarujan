package uploader

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

type Uploader struct {
	c *s3.S3
}

func New() Uploader {
	u := Uploader{
		c: createS3Client(),
	}
	return u
}

func (self *Uploader) GetURLList(bucket string) ([]string, error) {
	res, err := self.c.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, err
	}

	list := make([]string, len(res.Contents))
	for i, content := range res.Contents {
		list[i] = self.GenerateURL(bucket, *content.Key)
	}

	return list, nil
}

func (self *Uploader) UploadBlob(bucket string, key string, blob []byte) (string, error) {
	_, err := self.c.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(blob),
	})
	return self.GenerateURL(bucket, key), err
}

func (self *Uploader) GenerateURL(bucket string, key string) string {
	if isPathStyle() {
		return os.Getenv("IMAGE_HOST") + "/" + bucket + "/" + key
	} else {
		return os.Getenv("IMAGE_HOST") + "/" + key
	}
}

func createS3Client() *s3.S3 {
	cre := credentials.NewEnvCredentials()
	cli := s3.New(session.New(), &aws.Config{
		Credentials:      cre,
		Region:           aws.String(os.Getenv("AWS_S3_REGION")),
		Endpoint:         aws.String(os.Getenv("AWS_S3_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(isPathStyle()),
	})
	return cli
}

func isPathStyle() bool {
	return os.Getenv("AWS_S3_FORCE_PATH_STYLE") == "1"
}
