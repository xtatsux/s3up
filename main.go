package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		profileName = flag.String("p", "default", "AWS profile name")
		expiresMins = flag.Int("e", 60, "URL expiration time in minutes")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <file-path> <bucket-name>[/<prefix>] [-e <expires-mins>] [-p <profile>]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	filePath := flag.Arg(0)
	bucketSpec := flag.Arg(1)

	// Split bucket and prefix
	parts := strings.SplitN(bucketSpec, "/", 2)
	bucket := parts[0]
	prefix := ""
	if len(parts) > 1 {
		prefix = parts[1]
	}

	// Create S3 client using AWS CLI credentials
	s3Client, err := NewS3Client(*profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating S3 client: %v\n", err)
		os.Exit(1)
	}

	// Upload file and get URL
	ctx := context.Background()
	url, err := UploadAndGetURL(
		ctx,
		s3Client,
		filePath,
		bucket,
		prefix,
		time.Duration(*expiresMins)*time.Minute,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error uploading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(url)
}
