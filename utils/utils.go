package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"gorm.io/gorm/logger"
)

func GetLogLevel() logger.LogLevel {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return logger.Error
	}

	switch logLevel {
	case "DEBUG":
		return logger.Info
	case "WARN":
		return logger.Warn
	case "INFO":
		return logger.Info
	case "SILENT":
		return logger.Silent
	default:
		return logger.Error
	}
}

func ContentTypeContains(list []string, str string) bool {
	for _, s := range list {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}

func ArrayContains(permissions []string, permission string) bool {
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}

	return false
}

func ExtractKeyFromURL(s3URL string) (string, error) {
	u, err := url.Parse(s3URL)
	if err != nil {
		return "", err
	}
	return strings.TrimPrefix(u.Path, "/"), nil
}

func DownloadFileFromS3(sess *session.Session, bucket, key string) (string, error) {
	downloader := s3manager.NewDownloader(sess)

	dir := filepath.Dir(key)
	if err := os.MkdirAll(filepath.Join(os.TempDir(), dir), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create directory %q, %v", dir, err)
	}

	tempFile, err := os.Create(filepath.Join(os.TempDir(), key))
	if err != nil {
		return "", fmt.Errorf("failed to create file %q, %v", key, err)
	}

	_, err = downloader.Download(tempFile,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})
	if err != nil {
		return "", fmt.Errorf("failed to download file, %v", err)
	}

	return tempFile.Name(), nil
}

func DeleteFileFromS3(sess *session.Session, bucket, key string) error {
	svc := s3.New(sess)
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

// IsAfterMidnight checks if the reservation time is after midnight
func IsAfterMidnight(reservationTime time.Time) bool {
	return reservationTime.Hour() >= 0 && reservationTime.Hour() < 12
}

// CompareTimes compares two time.Time objects
func CompareTimes(t1, t2 time.Time) bool {
	return t1.Hour()*60+t1.Minute() <= t2.Hour()*60+t2.Minute()
}

// RemovePlusSign removes the plus sign from the phone number
func RemovePlusSign(phoneNumber string) string {
	return strings.Replace(phoneNumber, "+", "", -1)
}
