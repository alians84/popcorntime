package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type B2Client struct {
	Client *s3.Client
	Bucket string
}

func NewB2Client() (*B2Client, error) {
	// Получаем credentials из переменных окружения
	accountID := os.Getenv("B2_ACCOUNT_ID")
	applicationKey := os.Getenv("B2_APPLICATION_KEY")
	bucketName := os.Getenv("B2_BUCKET_NAME")
	b2Region := os.Getenv("B2_REGION")

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               fmt.Sprintf("https://s3.%s.backblazeb2.com", b2Region),
			SigningRegion:     b2Region, // Используем реальный регион B2
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accountID,
			applicationKey,
			"",
		)),
		config.WithRegion(b2Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &B2Client{
		Client: client,
		Bucket: bucketName,
	}, nil
}
