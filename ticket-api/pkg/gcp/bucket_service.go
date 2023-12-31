package gcp

import (
	"context"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"ticket-api/config"
	"ticket-api/pkg"

	"cloud.google.com/go/storage"
	"golang.org/x/oauth2/google"
)

// BucketService -> handles the file upload/download functions
type BucketService struct {
	logger pkg.Logger
	client *storage.Client
	env    config.Config
}

// NewBucketService -> initialization for the BucketService struct
func NewBucketService(logger pkg.Logger, client *storage.Client, env config.Config) BucketService {
	return BucketService{
		logger: logger,
		client: client,
		env:    env,
	}
}

// UploadFile -> uploads the file to the cloud storage
func (s *BucketService) UploadFile(
	ctx context.Context, file io.Reader, fileName string, originalFileName string,
) (string, error) {
	var bucketName = s.env.StorageBucketName

	if bucketName == "" {
		s.logger.Fatal("please check your env file for STORAGE_BUCKET_NAME")
	}

	_, err := s.client.Bucket(bucketName).Attrs(ctx)

	if err == storage.ErrBucketNotExist {
		s.logger.Fatalf("provided bucket %v doesn't exists", bucketName)
	}

	if err != nil {
		s.logger.Fatalf("cloud bucket error: %v", err.Error())
	}

	wc := s.client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	wc.ContentType = "application/octet-stream"
	wc.ContentDisposition = "attachment; filename=" + originalFileName

	_, err = io.Copy(wc, file)
	if err != nil {
		return "", err
	}

	err = wc.Close()
	if err != nil {
		return "", err
	}

	u, err := url.ParseRequestURI("/" + bucketName + "/" + wc.Attrs().Name)

	if err != nil {
		return "", err
	}

	path := u.EscapedPath()
	path = strings.Replace(path, "/"+bucketName, "", 1)
	path = strings.Replace(path, "/", "", 1)

	return path, nil
}

// GetObjectSignedURL -> gets the signed url for the stored object
func (s *BucketService) GetObjectSignedURL(object string) (string, error) {
	var bucketName = s.env.StorageBucketName

	jsonKey, err := os.ReadFile("./config/serviceAccountKey.json")
	if err != nil {
		return "", nil
	}

	conf, err := google.JWTConfigFromJSON(jsonKey)

	if err != nil {
		return "", err
	}

	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "GET",
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Expires:        time.Now().Add(15 * time.Minute),
	}

	u, err := storage.SignedURL(bucketName, object, opts)

	if err != nil {
		return "", err
	}

	return u, nil
}

// RemoveObject -> removes the file from the storage bucket
func (s *BucketService) RemoveObject(objectName string, bucketName string, thumb bool) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	objectToDelete := s.client.Bucket(bucketName).Object(objectName)
	err := objectToDelete.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
