package main

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	bucketName := flag.String("bucket", "", "S3 bucket name")
	fileCount := flag.Int("count", 100, "Number of files to generate")
	workerCount := flag.Int("workers", 10, "Number of concurrent workers")
	endpoint := flag.String("endpoint", "", "S3 compatible endpoint (optional, e.g., http://localhost:9000 for Minio)")
	accessKey := flag.String("access-key", "", "Access key for S3 or S3-compatible service")
	secretKey := flag.String("secret-key", "", "Secret key for S3 or S3-compatible service")
	region := flag.String("region", "us-east-1", "AWS region or dummy region for S3-compatible service")
	flag.Parse()

	if *bucketName == "" {
		log.Fatal("Bucket name is required")
	}

	ctx := context.TODO()
	var cfg aws.Config
	var err error

	if *endpoint != "" {
		// S3-compatible service configuration
		if *accessKey == "" || *secretKey == "" {
			log.Fatal("Access key and secret key are required for S3-compatible services")
		}
		creds := credentials.NewStaticCredentialsProvider(*accessKey, *secretKey, "")
		cfg, err = awsconfig.LoadDefaultConfig(
			ctx,
			awsconfig.WithCredentialsProvider(creds),
			awsconfig.WithRegion(*region),
		)
		if err != nil {
			log.Fatalf("Failed to load config for S3-compatible service: %v", err)
		}
	} else {
		// AWS S3 configuration
		cfg, err = awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(*region))
		if err != nil {
			log.Fatalf("Failed to load config for AWS S3: %v", err)
		}
		// If access key and secret key are provided, use them
		if *accessKey != "" && *secretKey != "" {
			cfg.Credentials = credentials.NewStaticCredentialsProvider(*accessKey, *secretKey, "")
		}
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if *endpoint != "" {
			o.UsePathStyle = true
			o.BaseEndpoint = aws.String(*endpoint)
		}
	})

	jobs := make(chan int, *fileCount)
	var wg sync.WaitGroup

	// Adjust worker count if it's more than file count
	actualWorkerCount := min(*workerCount, *fileCount)

	for w := 1; w <= actualWorkerCount; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg, client, *bucketName)
	}

	for j := 1; j <= *fileCount; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()

	fmt.Printf("File generation complete. %d files uploaded.\n", *fileCount)
}

func worker(id int, jobs <-chan int, wg *sync.WaitGroup, client *s3.Client, bucketName string) {
	defer wg.Done()
	for j := range jobs {
		fileName := generateFileName()
		content := strings.NewReader(fmt.Sprintf("This is file %d", j))

		_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket:      &bucketName,
			Key:         &fileName,
			Body:        content,
			ContentType: aws.String("text/plain"),
		})

		if err != nil {
			log.Printf("Worker %d: Failed to upload file %s: %v", id, fileName, err)
		} else {
			fmt.Printf("Worker %d: Uploaded file: %s\n", id, fileName)
		}
	}
}

func generateFileName() string {
	nanoEpoch := time.Now().UnixNano()
	hash := sha512.Sum512([]byte(fmt.Sprintf("%d", nanoEpoch)))
	return hex.EncodeToString(hash[:])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
