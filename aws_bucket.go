package filestores

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsBucket struct {
	config   aws.Config
	bucket   string
	client   *s3.Client
	uploader *manager.Uploader
}

func NewAwsBucket(bucket string, config aws.Config) *AwsBucket {
	var (
		client   = s3.NewFromConfig(config)
		uploader = manager.NewUploader(client)
	)
	return &AwsBucket{
		config:   config,
		bucket:   bucket,
		client:   client,
		uploader: uploader,
	}
}

func (c *AwsBucket) Save(input Storable) (string, error) {
	file := input.GetFile()
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileBytes := make([]byte, fileInfo.Size())
	if _, err := file.Read(fileBytes); err != nil {
		return "", err
	}
	filename := input.Filename()
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(fileBytes)

	_, err = c.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      &c.bucket,
		Key:         &filename,
		Body:        bytes.NewReader(fileBytes),
		ContentType: &contentType,
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", c.bucket, c.config.Region, filename), nil

}

func (c *AwsBucket) Delete(filepath string) error {
	if _, err := c.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: &c.bucket,
		Key:    &filepath,
	}); err != nil {
		return err
	}

	return nil
}
