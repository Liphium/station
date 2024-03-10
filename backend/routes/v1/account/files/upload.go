package files

import (
	"context"
	"fmt"
	"log"
	"node-backend/database"
	"node-backend/entities/account"
	"node-backend/util"
	"node-backend/util/auth"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/files/upload
func uploadFile(c *fiber.Ctx) error {

	if disabled {
		return util.FailedRequest(c, "file.disabled", nil)
	}

	// Form data
	key := c.FormValue("key", "-")
	name := c.FormValue("name", "-")
	extension := c.FormValue("extension", "-")
	favorite := c.FormValue("favorite", "false")
	if key == "-" || name == "-" || extension == "-" {
		log.Println("invalid form data")
		return util.InvalidRequest(c)
	}
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("no file")
		return util.InvalidRequest(c)
	}
	accId := util.GetAcc(c)
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		log.Println("invalid headers")
		return util.InvalidRequest(c)
	}

	// Check file size
	if file.Size > maxUploadSize {
		return util.FailedRequest(c, "file.too_large", nil)
	}

	// Check total storage
	totalStorage, err := CountTotalStorage(accId)
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	if totalStorage+file.Size > maxTotalStorage {
		return util.FailedRequest(c, "file.storage_limit", nil)
	}

	// Generate file name Format: a-[accountId]-[objectIdentifier].[extension]
	fileId := "a-" + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-" + accId + "-" + auth.GenerateToken(16) + "." + extension
	if err := database.DBConn.Create(&account.CloudFile{
		Id:       fileId,
		Name:     name,
		Type:     file.Header.Get("Content-Type"),
		Key:      key,
		Account:  accId,
		Favorite: favorite == "true",
		System:   false,
		Size:     file.Size,
	}).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	f, err := file.Open()
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Upload to R2
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileId),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}
	location := os.Getenv("R2_PUBLIC_URL") + fileId

	// Update file path
	if err := database.DBConn.Model(&account.CloudFile{}).Where("id = ?", fileId).Update("path", location).Error; err != nil {
		client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(file.Filename),
		})
		return util.FailedRequest(c, "server.error", err)
	}

	// Not encrypted cause this doesn't matter (and it's an unencrypted route)
	return c.JSON(fiber.Map{
		"success": true,
		"id":      fileId,
		"url":     location,
	})
}
