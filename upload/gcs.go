package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
)

func uploadToGCS(ctx context.Context, file io.Reader, filename string) (string, error) {
	logger := GetLogger(ctx, "helper::uploadToGCS")

	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Error("failed to create GCS client", "error", err)
		return "", fmt.Errorf("failed to create GCS client: %w", err)
	}
	defer client.Close()

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	objectName := fmt.Sprintf("uploads/%d-%s", time.Now().UnixNano(), filename)
	logger.Info("uploading to GCS", "bucket", bucketName, "object", objectName)

	bucket := client.Bucket(bucketName)
	object := bucket.Object(objectName)
	writer := object.NewWriter(ctx)

	if _, err := io.Copy(writer, file); err != nil {
		logger.Error("failed to write to GCS", "error", err)
		return "", fmt.Errorf("failed to write to GCS: %w", err)
	}

	if err := writer.Close(); err != nil {
		logger.Error("failed to close GCS writer", "error", err)
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}
	logger.Info("file written to GCS", "object", objectName)

	expiry := time.Now().Add(15 * time.Minute)
	signedURL, err := client.Bucket(bucketName).SignedURL(objectName, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: expiry,
	})
	if err != nil {
		logger.Error("failed to generate signed URL", "error", err)
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	logger.Info("signed URL generated", "object", objectName, "expires", expiry)

	return signedURL, nil
}