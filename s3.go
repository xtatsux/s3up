package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/schollz/progressbar/v3"
)

// progressReader wraps an io.Reader and updates a progress bar
type progressReader struct {
	reader     io.Reader
	size       int64
	bar        *progressbar.ProgressBar
}

func (r *progressReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if n > 0 {
		r.bar.Add(n)
	}
	return n, err
}

func UploadAndGetURL(ctx context.Context, client *s3.Client, filePath, bucket, key string, expires time.Duration) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Get file size for progress bar
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Create progress bar
	bar := progressbar.DefaultBytes(
		fileInfo.Size(),
		"uploading",
	)

	// Create a custom progress reader
	reader := &progressReader{
		reader: file,
		size:   fileInfo.Size(),
		bar:    bar,
	}

	// Upload the file with progress bar
	size := fileInfo.Size()
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: &size, // Pass pointer to size
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate presigned URL
	presignClient := s3.NewPresignClient(client)
	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignResult.URL, nil
}
