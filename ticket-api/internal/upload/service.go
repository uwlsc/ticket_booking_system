package upload

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"ticket-api/internal/errors"
	"ticket-api/pkg/gcp"
	"net/http"
	"strings"

	"github.com/chai2010/webp"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type UploadService struct {
	bucket gcp.BucketService
}

func NewUploadService(b gcp.BucketService) UploadService {
	return UploadService{bucket: b}
}

func (u *UploadService) UploadFile(ext, fileName, fileUID string, fB []byte, ctx context.Context) (string, error) {
	var fileToUpload io.Reader

	fileToUpload = bytes.NewReader(fB)

	if !(ext == ".pdf" || ext == ".mp4" || ext == ".svg" || ext == ".heic" || ext == ".heif") {
		origImage, err := u.getImage(bytes.NewReader(fB), ext)
		if err == nil {
			resizeImage := resize.Resize(1000, 0, origImage, resize.Lanczos3)

			buff := new(bytes.Buffer)
			if ext == ".jpg" || ext == ".jpeg" {
				if err := jpeg.Encode(buff, resizeImage, &jpeg.Options{Quality: 100}); err != nil {
					return "", err
				}
			}
			if ext == ".png" {
				if err := png.Encode(buff, resizeImage); err != nil {
					return "", err
				}
			}
			if ext == ".webp" {
				if err := webp.Encode(buff, resizeImage, &webp.Options{Lossless: true}); err != nil {
					return "", err
				}
			}

			fileToUpload = buff
		}
	}

	urlResponse, err := u.bucket.UploadFile(ctx, fileToUpload, fileName, fileUID)
	if err != nil {
		return "", err
	}

	return urlResponse, nil
}

func (u *UploadService) getImage(file io.Reader, ext string) (image.Image, error) {
	if ext == ".jpg" || ext == ".jpeg" {
		return jpeg.Decode(file)
	}
	if ext == ".png" {
		return png.Decode(file)
	}
	if ext == ".webp" {
		return webp.Decode(file)
	}
	return nil, errors.ErrExtensionMismatch
}

func (u *UploadService) UploadURL(url string, ctx context.Context) (string, error) {
	getImage, err := http.Get(url)
	if err != nil {
		return "", errors.ErrUrlRead
	}

	fileByte, err := io.ReadAll(getImage.Body)
	if err != nil {
		return "", errors.ErrFileRead
	}

	mimeType := http.DetectContentType(fileByte)

	imageNames := strings.Split(mimeType, "/")

	ext := "." + imageNames[1]

	if !(ext == ".jpeg" || ext == ".jpg" || ext == ".png" || ext == ".webp" || ext == ".heic") {
		return "", errors.ErrImageFormat
	}

	fileUID, _ := uuid.NewRandom()
	fileName := "user/" + fileUID.String() + ext

	urlResponse, err := u.UploadFile(ext, fileName, fileUID.String(), fileByte, ctx)
	if err != nil {
		return "", err
	}

	return urlResponse, nil
}
