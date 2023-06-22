package middlewares

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"ticket-api/internal/constants"
	"ticket-api/internal/errors"
	"ticket-api/internal/responses"
	"ticket-api/internal/upload"
	"ticket-api/pkg"
	"ticket-api/pkg/gcp"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Extension string

const (
	JPEGFile Extension = ".jpeg"
	JPGFile  Extension = ".jpg"
	PNGFile  Extension = ".png"
)

type UploadConfig struct {
	// FieldName where to pull multipart files from
	FieldName string

	// BucketFolder where to put the uploaded files to
	BucketFolder string

	// Extensions array of extensions
	Extensions []Extension

	// ThumbnailEnabled set whether thumbnail is enabled or nor
	ThumbnailEnabled bool

	// ThumbnailWidth set thumbnail width
	ThumbnailWidth uint

	// WebpEnabled set whether thumbnail is enabled or nor
	WebpEnabled bool

	// ThumbnailHeight set thumbnail height
	ThumbnailHeight uint

	// Compute by Thumbnail height
	ThumbnailHeightCheck bool

	// Multiple set whether to upload multiple files with same key name
	Multiple bool
}

type UploadMiddleware struct {
	logger        pkg.Logger
	bucket        gcp.BucketService
	config        UploadConfig
	uploadService upload.UploadService
}

func NewUploadMiddleware(logger pkg.Logger, bucket gcp.BucketService, up upload.UploadService) UploadMiddleware {
	m := UploadMiddleware{
		bucket:        bucket,
		logger:        logger,
		uploadService: up,
	}
	return m
}

func (u UploadMiddleware) Config() UploadConfig {
	return UploadConfig{
		FieldName:            "file",
		BucketFolder:         "",
		Extensions:           []Extension{JPEGFile, PNGFile, JPGFile},
		ThumbnailEnabled:     false,
		ThumbnailWidth:       200,
		ThumbnailHeightCheck: false,
		ThumbnailHeight:      600,
		Multiple:             false,
	}
}

// Field modify field of upload
func (cfg UploadConfig) Field(name string) UploadConfig {
	cfg.FieldName = name
	return cfg
}

// Folder modify folder of upload
func (cfg UploadConfig) Folder(folder string) UploadConfig {
	cfg.BucketFolder = folder
	return cfg
}

// Extension modify upload extension
func (cfg UploadConfig) Extension(ext ...Extension) UploadConfig {
	cfg.Extensions = ext
	return cfg
}

// ThumbEnable enable thumbnail generation
func (cfg UploadConfig) ThumbEnable(enable bool) UploadConfig {
	cfg.ThumbnailEnabled = enable
	return cfg
}

// ThumbEnableHeight enable thumbnail generation by height
func (cfg UploadConfig) ThumbEnableHeight(enable bool) UploadConfig {
	cfg.ThumbnailHeightCheck = enable
	return cfg
}

// ThumbNailWidth customize thumbnail width
func (cfg UploadConfig) ThumbNailWidth(width uint) UploadConfig {
	cfg.ThumbnailWidth = width
	return cfg
}

// WebpEnable enables thumbnail generation
func (cfg UploadConfig) WebpEnable(enable bool) UploadConfig {
	cfg.WebpEnabled = enable
	return cfg
}

// MultipleFilesUpload enable multiple files to be uploaded with same key name
func (cfg UploadConfig) MultipleFilesUpload(enable bool) UploadConfig {
	cfg.Multiple = enable
	return cfg
}

// Push adds file upload configuration
func (u *UploadMiddleware) Push(config UploadConfig) UploadMiddleware {
	u.config = config
	return *u
}

// Handle handles file upload
func (u UploadMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {

		file, fileHeader, err := c.Request.FormFile(u.config.FieldName)
		if err == nil {
			ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
			if !u.matchesExtension(u.config, ext) {
				responses.ErrorJSON(c, http.StatusInternalServerError, "file extension not supported")
				c.Abort()
				return
			}

			fileByte, err := io.ReadAll(file)
			if err != nil {
				responses.ErrorJSON(c, http.StatusInternalServerError, errors.ErrFileRead)
				c.Abort()
				return
			}

			uploadedFiles, err := u.uploadFile(fileHeader, ext, u.config, fileByte)
			if err != nil {
				u.logger.Error("file-upload-error: ", err.Error())
				responses.ErrorJSON(c, http.StatusInternalServerError, err.Error())
				c.Abort()
				return
			}

			c.Set(constants.File, pkg.UploadedFiles(uploadedFiles))
			c.Next()
		} else {
			path := strings.Split(strings.Split(c.Request.URL.Path, "/")[1], "/")[0]

			if path == "event" && c.Request.Method == "PUT" {
				c.Next()
				return
			} else {
				responses.ErrorJSON(c, http.StatusInternalServerError, errors.ErrNoFile)
				c.Abort()
				return
			}
		}
	}
}

func (u *UploadMiddleware) uploadFile(fH *multipart.FileHeader, ext string, conf UploadConfig, fB []byte) ([]pkg.UploadMetadata, error) {
	ctx := context.Background()

	uploadFileName, fileUID := u.randomFileName(conf, ext)

	urlResponse, err := u.uploadService.UploadFile(ext, uploadFileName, fileUID, fB, ctx)
	if err != nil {
		return nil, err
	}

	var uploadedFiles []pkg.UploadMetadata

	uploadedFiles = append(
		uploadedFiles, pkg.UploadMetadata{
			FieldName: conf.FieldName,
			FileName:  fH.Filename,
			URL:       urlResponse,
			FileUID:   fileUID,
			Size:      fH.Size,
		},
	)

	return uploadedFiles, nil
}

func (u *UploadMiddleware) matchesExtension(c UploadConfig, ext string) bool {
	for _, e := range c.Extensions {
		if e == Extension(ext) {
			return true
		}
	}
	return false
}

func (u *UploadMiddleware) randomFileName(c UploadConfig, ext string) (randomName, uid string) {
	randUUID, _ := uuid.NewRandom()
	fileName := randUUID.String() + ext
	return u.bucketPath(c, fileName), randUUID.String()
}

func (u *UploadMiddleware) bucketPath(c UploadConfig, fileName string) string {
	if c.BucketFolder != "" {
		return fmt.Sprintf("%s/%s", c.BucketFolder, fileName)
	}
	return fileName
}
