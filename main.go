package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// Load configuration first to get defaults
	config, err := LoadS3Config()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// Define flags with config defaults
	expiresMins := flag.Int("e", config.ExpirationMins, "URL expiration time in minutes (overrides config)")
	profileName := flag.String("p", config.Profile, "AWS profile name (overrides config file)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-e <expires-mins>] [-p <profile>] <file-path> [key-prefix]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nArguments:")
		fmt.Fprintln(os.Stderr, "  file-path    Path to the file to upload")
		fmt.Fprintln(os.Stderr, "  key-prefix   Optional prefix to prepend to the S3 key (default: config file prefix)")
		fmt.Fprintln(os.Stderr, "\nOptions:")
		fmt.Fprintf(os.Stderr, "  -e int\n")
		fmt.Fprintf(os.Stderr, "        URL expiration time in minutes (default from config: %d)\n", config.ExpirationMins)
		fmt.Fprintf(os.Stderr, "  -p string\n")
		fmt.Fprintf(os.Stderr, "        AWS profile name (default from config: %s)\n", config.Profile)
		fmt.Fprintln(os.Stderr, "\nConfiguration:")
		fmt.Fprintf(os.Stderr, "  Region: %s (configured in config.yaml)\n", config.Region)
		fmt.Fprintf(os.Stderr, "  Bucket: %s (configured in config.yaml)\n", config.Bucket)
	}

	// Parse flags
	flag.Parse()

	// Get non-flag arguments
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}

	// Validate expiration time
	if *expiresMins <= 0 {
		fmt.Fprintf(os.Stderr, "Error: expiration time must be greater than 0 minutes\n")
		os.Exit(1)
	}

	// Create S3 client with profile from flag or config
	s3Client, err := NewS3Client(*profileName, config.Region)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating S3 client: %v\n", err)
		os.Exit(1)
	}

	filePath := args[0]

	// Determine the key prefix
	prefix := config.Prefix
	if len(args) > 1 {
		// If prefix is provided as argument, it overrides the config
		prefix = args[1]
	}

	// Clean up the prefix (remove leading/trailing slashes)
	prefix = strings.Trim(prefix, "/")

	// Build the full key
	key := filepath.Base(filePath)
	if prefix != "" {
		key = prefix + "/" + key
	}

	// Calculate expiration duration
	expires := time.Duration(*expiresMins) * time.Minute

	// Upload file and get URL
	ctx := context.Background()
	url, err := UploadAndGetURL(
		ctx,
		s3Client,
		filePath,
		config.Bucket,
		key,
		expires,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error uploading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(url)
}
