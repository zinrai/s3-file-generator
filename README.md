# S3 File Generator

This tool is designed to generate and upload random files to Amazon S3 or S3-compatible storage services. It's useful for populating S3 buckets with sample data.

## Note

This tool generates random file names using SHA512 hashes of the current timestamp. The content of each file is a simple string indicating the file number.

## Features

- Generates files with random content and uploads them to S3
- Supports both AWS S3 and S3-compatible services (like MinIO)
- Configurable number of files to generate
- Concurrent uploads using multiple workers
- Customizable S3 endpoint, region, and credentials

## Installation

Build the too:

```
$ go build
```

## Usage

Run the program with the following command-line flags:

```
$ s3-file-generator -bucket BUCKET_NAME -count FILE_COUNT -workers WORKER_COUNT -access-key ACCESS_KEY -secret-key SECRET_KEY [-endpoint CUSTOM_ENDPOINT] [-region REGION]
```

### Examples

Upload 1000 files to an AWS S3 bucket:

```
$ s3-file-generator -bucket my-test-bucket -count 1000 -workers 20 -access-key AKIAIOSFODNN7EXAMPLE -secret-key wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

Upload 500 files to a MinIO server:

```
$ kubectl apply -f minio-k8s.yaml.example
```

```
$ kubectl port-forward -n minio-ns svc/minio 9000:9000 9001:9001
```

```
$ s3-file-generator -bucket test-bucket -count 500 -workers 10 -access-key minio -secret-key minio123 -endpoint http://localhost:9000
```

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
