package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"
)

type S3Config struct {
	Bucket         string `yaml:"bucket"`          // Default bucket name
	Prefix         string `yaml:"prefix"`          // Optional default prefix (folder)
	Profile        string `yaml:"profile"`         // Optional AWS profile name
	Region         string `yaml:"region"`          // AWS region
	ExpirationMins int    `yaml:"expiration_mins"` // Default URL expiration time in minutes
}

// NewS3Client creates a new S3 client using credentials from either:
// 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION)
// 2. AWS CLI shared configuration (~/.aws/credentials or $XDG_CONFIG_HOME/aws/credentials)
func NewS3Client(profileName string, region string) (*s3.Client, error) {
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
		// Load AWS CLI shared configuration with specified profile and region
		opts := []func(*config.LoadOptions) error{
			config.WithSharedConfigProfile(profileName),
			config.WithRegion(region),
		}

		// Load the configuration
		cfg, err = config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config for profile '%s': %w", profileName, err)
		}
	}

	// Create and return the S3 client
	return s3.NewFromConfig(cfg), nil
}

// LoadS3Config loads the s3up configuration file
func LoadS3Config() (*S3Config, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	configPath := filepath.Join(configDir, "s3up", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config S3Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.Bucket == "" {
		return nil, fmt.Errorf("bucket name is required in config file")
	}
	if config.Region == "" {
		config.Region = "ap-northeast-1" // Default to Tokyo region if not specified
	}

	// Set default values
	if config.Profile == "" {
		config.Profile = "default"
	}
	if config.ExpirationMins <= 0 {
		config.ExpirationMins = 60 // Default to 60 minutes if not specified or invalid
	}

	return &config, nil
}
