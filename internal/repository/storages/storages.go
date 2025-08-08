package storages

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2/utils"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"path/filepath"
	"popcorntime-project/internal/config"
	"popcorntime-project/internal/models"
	"time"
)

type Service interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (*models.FilesS3, error)
	SaveFile(key string) (*models.FilesS3, error)
}

type service struct {
	db       *gorm.DB
	b2Client *config.B2Client
}

func NewService(db *gorm.DB, b2Client *config.B2Client) Service {
	return &service{
		db:       db,
		b2Client: b2Client,
	}
}

func (r *service) UploadFile(ctx context.Context, file *multipart.FileHeader) (*models.FilesS3, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 2. Генерируем уникальное имя файла (если нужно)
	fileKey := "images/" + utils.UUID() + filepath.Ext(file.Filename)

	// 3. Загружаем в S3/B2
	if err := r.UploadFileS3(ctx, fileKey, src); err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	saveFile, err := r.SaveFile(fileKey)
	if err != nil {
		return nil, fmt.Errorf("error save file: %w", err)
	}

	return saveFile, nil
}

func (r *service) UploadFileS3(ctx context.Context, key string, file io.Reader) error {
	_, err := r.b2Client.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.b2Client.Bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	return err
}

func (r *service) SaveFile(key string) (*models.FilesS3, error) {
	file := models.FilesS3{
		FileURL:   key,
		CreatedAt: time.Now(),
	}

	if err := r.db.Create(&file).Error; err != nil {
		return nil, err
	}

	return &file, nil
}
