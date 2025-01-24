# s3up - AWS S3 File Upload Tool

s3up is a command-line tool that simplifies uploading files to AWS S3 and generates pre-signed URLs for easy sharing.

## Features

- Upload files to AWS S3 buckets
- Generate pre-signed URLs with configurable expiration times
- Support for AWS CLI profiles
- Support for S3 bucket prefixes (folders)
- Support for AWS credentials via environment variables

## Prerequisites

- Go 1.21 or later
- AWS credentials (via environment variables or AWS CLI configuration)

## Installation

```bash
# Clone the repository
git clone https://github.com/xtatsux/s3up.git
cd s3up

# Build the binary
go build -o s3up
```

## Configuration

s3up supports two methods for AWS credentials configuration:

### 1. Environment Variables

Set the following environment variables:
```bash
export AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
export AWS_REGION=ap-northeast-1
```

When environment variables are set, they take precedence over AWS CLI configuration.

### 2. AWS CLI Configuration

Alternatively, use AWS CLI's shared configuration:

```bash
# Default location for AWS CLI credentials
$XDG_CONFIG_HOME/aws/credentials
# or
~/.aws/credentials
```

Example AWS credentials file:
```ini
[default]
aws_access_key_id = AKIAXXXXXXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
region = ap-northeast-1

[prod]
aws_access_key_id = AKIAXXXXXXXXXXXXXXXX
aws_secret_access_key = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
region = ap-northeast-1
```

## Usage

Basic command format:
```bash
s3up <file-path> <bucket-name>[/<prefix>] [-e <expires-mins>] [-p <profile>]
```

Options:
- `-e`: URL expiration time in minutes (default: 60)
- `-p`: AWS profile name (default: "default", ignored when using environment variables)

### Examples

1. Basic upload (uses environment variables or default profile, 60-minute expiration):
```bash
s3up document.pdf mybucket/uploads/
```

2. Upload with custom profile (when using AWS CLI configuration):
```bash
s3up document.pdf mybucket/uploads/ -p prod
```

3. Upload with custom expiration time (180 minutes):
```bash
s3up document.pdf mybucket/uploads/ -e 180
```

4. Upload to bucket root:
```bash
s3up document.pdf mybucket
```

## Error Handling

The tool handles various error scenarios:
- Invalid AWS credentials
- Missing environment variables
- Invalid AWS CLI profile
- File access errors
- S3 upload failures
- Pre-signed URL generation failures

Error messages are displayed to stderr with appropriate exit codes.

## License

MIT License
