package services

import (
	"context"
	"github.com/jacksonopp/go-recipe/domain"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
	"log"
	"mime/multipart"
	"net/url"
	"time"
)

const (
	BUCKET_NAME = "go-recipe"
	URL_EXPIRY  = time.Hour * 24 * 7
)

type BucketService interface {
	UploadFile(userID uint, file *multipart.FileHeader) (*domain.File, error)
	GetFileByObjectName(objectName string) (*domain.File, error)
	GetFileByID(fileID uint) (*domain.File, error)
}

type bucketService struct {
	ctx   context.Context
	db    *gorm.DB
	minio *minio.Client
}

func NewBucketService(db *gorm.DB, minio *minio.Client) BucketService {
	ctx := context.Background()
	return &bucketService{ctx: ctx, db: db, minio: minio}
}

// UploadFile uploads a file to the bucket and returns the file object
// this method also creates a database entry for the file
//
// satisfying the BucketService interface
func (s *bucketService) UploadFile(userID uint, file *multipart.FileHeader) (*domain.File, error) {
	data, err := file.Open()
	defer func(data multipart.File) {
		err := data.Close()
		if err != nil {
			log.Println("failed to close file")
		}
	}(data)
	if err != nil {
		return nil, err
	}

	filename := url.QueryEscape(file.Filename)

	info, err := s.minio.PutObject(s.ctx, BUCKET_NAME, filename, data, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return nil, err
	}

	objUrl, err := s.createPresignedUrl(info.Key)
	if err != nil {
		return nil, err
	}

	dbFile := &domain.File{
		Name:      info.Key,
		Url:       objUrl,
		UrlExpiry: time.Now().Add(URL_EXPIRY),
		UserID:    userID,
	}
	err = s.db.Create(dbFile).Error
	if err != nil {
		return nil, err
	}

	return dbFile, nil
}

// GetFileByObjectName returns a file object by its object name
//
// satisfying the BucketService interface
func (s *bucketService) GetFileByObjectName(objectName string) (*domain.File, error) {
	file := &domain.File{}
	err := s.db.Where("name = ?", objectName).First(file).Error
	if err != nil {
		return nil, ErrFileNotFound
	}
	return file, nil
}

// GetFileByID returns a file object by its ID
//
// satisfying the BucketService interface
func (s *bucketService) GetFileByID(fileID uint) (*domain.File, error) {
	file := &domain.File{}
	err := s.db.First(file, fileID).Error
	if err != nil {
		return nil, ErrFileNotFound
	}

	// update the URL if it expires in less than 1 hour
	if file.UrlExpiry.Before(time.Now().Add(time.Hour)) {
		err = s.updateUrl(file)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}

// CreatePresignedUrl creates a presigned URL for an object in the bucket
func (s *bucketService) createPresignedUrl(objectName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+objectName)
	reqParams.Set("response-content-type", "application/octet-stream")
	reqParams.Set("response-expires", "Fri, 01 Jan 2100 00:00:00 GMT")

	presignedURL, err := s.minio.PresignedGetObject(s.ctx, BUCKET_NAME, objectName, URL_EXPIRY, reqParams)
	if err != nil {
		return "", err
	}

	log.Println("presigned URL:", presignedURL.String())

	return presignedURL.String(), nil
}

// getUrl returns the URL of a file by its object name
func (s *bucketService) getUrl(objectName string) (string, error) {
	file := &domain.File{}
	err := s.db.Where("name = ?", objectName).First(file).Error
	if err != nil {
		return "", ErrFileNotFound
	}
	return file.Url, nil
}

func (s *bucketService) updateUrl(file *domain.File) error {
	objUrl, err := s.createPresignedUrl(file.Name)
	if err != nil {
		return err
	}
	file.Url = objUrl
	file.UrlExpiry = time.Now().Add(URL_EXPIRY)
	return s.db.Save(file).Error
}
