# s3up - AWS S3 File Upload Tool

s3up is a command-line tool that simplifies uploading files to AWS S3 and generates pre-signed URLs for easy sharing.

## Features

- Upload files to AWS S3 buckets
- Generate pre-signed URLs with configurable expiration times
- Support for AWS CLI profiles
- Support for AWS credentials via environment variables
- Configurable default bucket, prefix, and URL expiration settings

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

### S3 Upload Configuration

Create a configuration file at `$XDG_CONFIG_HOME/s3up/config.yaml` (or `~/.config/s3up/config.yaml`):

```yaml
# Required: Default S3 bucket name
bucket: your-bucket-name

# Required: AWS region (default: ap-northeast-1)
region: ap-northeast-1

# Optional: Default prefix (folder) for uploads
prefix: uploads

# Optional: Default AWS profile name (default: "default")
profile: prod

# Optional: Default URL expiration time in minutes (default: 60)
expiration_mins: 180
```

### AWS Credentials

s3up supports two methods for AWS credentials configuration:

#### 1. Environment Variables

Set the following environment variables:
```bash
export AWS_ACCESS_KEY_ID=AKIAXXXXXXXXXXXXXXXX
export AWS_SECRET_ACCESS_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
export AWS_REGION=ap-northeast-1
```

When environment variables are set, they take precedence over AWS CLI configuration.

#### 2. AWS CLI Configuration

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
s3up [-e <expires-mins>] [-p <profile>] <file-path> [key-prefix]
```

Arguments:
- `file-path`: Path to the file to upload
- `key-prefix`: (Optional) Prefix to prepend to the S3 key (overrides config file prefix)

Options:
- `-e`: URL expiration time in minutes (overrides config file expiration_mins)
- `-p`: AWS profile name (overrides config file profile)

### Examples

1. Basic upload (uses all config defaults):
```bash
s3up document.pdf
```

2. Upload with custom key prefix (overrides config prefix):
```bash
s3up document.pdf custom/path/
```

3. Upload with custom expiration time (overrides config expiration_mins):
```bash
s3up -e 180 document.pdf
```

4. Upload with specific AWS profile (overrides config profile):
```bash
s3up -p prod document.pdf
```

5. Combine options:
```bash
s3up -e 180 -p prod document.pdf custom/path/
```

## Configuration Precedence

1. Command-line flags take highest precedence (-e, -p)
2. Command-line arguments override corresponding config values (key-prefix)
3. Configuration file values are used as defaults
4. Built-in defaults are used if neither flags nor config values are provided:
   - Region: ap-northeast-1
   - Expiration: 60 minutes
   - Profile: "default"
   - Prefix: "" (empty string)

## Error Handling

The tool handles various error scenarios:
- Missing or invalid configuration file
- Invalid AWS credentials
- Missing environment variables
- Invalid AWS CLI profile
- File access errors
- S3 upload failures
- Pre-signed URL generation failures

Error messages are displayed to stderr with appropriate exit codes.

## License

MIT License
