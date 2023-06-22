package gcp

import (
	"context"
	"ticket-api/config"
	"ticket-api/pkg"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// NewBucketStorage creates a new storage client
func NewBucketStorage(logger pkg.Logger, env config.Config) *storage.Client {
	bucketName := env.StorageBucketName
	ctx := context.Background()
	if bucketName == "" {
		logger.Error("Please check your env file for StorageBucketName")
	}
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./config/serviceAccountKey.json"))
	if err != nil {
		logger.Fatal(err.Error())
	}
	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err == storage.ErrBucketNotExist {
		logger.Fatalf("Provided bucket %v doesn't exists", bucketName)
	}
	if err != nil {
		logger.Fatalf("Cloud bucket error: %v", err.Error())
	}
	return client
}
