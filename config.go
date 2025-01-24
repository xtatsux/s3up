package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client creates a new S3 client using credentials from either:
// 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION)
// 2. AWS CLI shared configuration (~/.aws/credentials or $XDG_CONFIG_HOME/aws/credentials)
func NewS3Client(profileName string) (*s3.Client, error) {
	ctx := context.Background()
	var cfg aws.Config
	var err error

	// Check for environment variables
	if os.Getenv("AWS_ACCESS_KEY_ID") != "" && os.Getenv("AWS_SECRET_ACCESS_KEY") != "" {
		// When environment variables are present, they take precedence
		cfg, err = config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config from environment: %w", err)
		}
	} else {
		// Fall back to AWS CLI shared configuration with specified profile
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithSharedConfigProfile(profileName),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config from shared configuration: %w", err)
		}
	}

	return s3.NewFromConfig(cfg), nil
}
