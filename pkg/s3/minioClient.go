package s3

import (
	"log"

	"github.com/minio/minio-go"
)

// MinioConnection func for opening minio connection.
func MinioConnection(bucketName string, location string) (*minio.Client, error) {
	// Connect to MinIO
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}

	// Check if the bucket exists
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		log.Printf("Error checking if bucket exists: %v", err)
		return nil, err
	}

	if !exists {
		// Create bucket if it doesn't exist
		err = minioClient.MakeBucket(bucketName, location)
		if err != nil {
			log.Printf("Error creating bucket: %v", err)
			return nil, err
		}
		log.Printf("Successfully created bucket: %s", bucketName)
	} else {
		log.Printf("Bucket %s already exists", bucketName)
	}

	return minioClient, nil
}
