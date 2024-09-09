package files

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
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
	tag := c.FormValue("tag", "")
	if key == "-" || name == "-" || extension == "-" {
		util.Log.Println("invalid form data")
		return util.InvalidRequest(c)
	}
	file, err := c.FormFile("file")
	if err != nil {
		util.Log.Println("no file")
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		util.Log.Println("invalid headers")
		return util.InvalidRequest(c)
	}

	// Check if the file name is valid
	if strings.Contains(file.Filename, "/") || strings.Contains(file.Filename, "\\") {
		return util.InvalidRequest(c)
	}

	// Check if the tag is valid
	if len(tag) > 100 {
		return util.InvalidRequest(c)
	}

	// Check file size
	if file.Size > maxUploadSize {
		return util.FailedRequest(c, fmt.Sprintf("file.too_large:%d", maxUploadSize/1_000_000), nil)
	}

	// Check total storage
	totalStorage, err := CountTotalStorage(accId)
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	if totalStorage+file.Size > maxTotalStorage {
		return util.FailedRequest(c, "file.storage_limit", nil)
	}

	// Generate file name Format: a-[timestamp]-[objectIdentifier].[extension]
	fileId := "a-" + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-" + auth.GenerateToken(16) + "." + extension
	if err := database.DBConn.Create(&account.CloudFile{
		Id:      fileId,
		Name:    name,
		Type:    file.Header.Get("Content-Type"),
		Key:     key,
		Account: accId,
		Tag:     tag,
		System:  false,
		Size:    file.Size,
	}).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Save the file to whatever repository is selected
	if fileRepoType == repoTypeR2 {

		// Open the file
		f, err := file.Open()
		if err != nil {
			return util.FailedRequest(c, "server.error", err)
		}

		// Save the file to R2
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileId),
			Body:   f,
		})
		if err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	} else if fileRepoType == repoTypeLocal {

		// Save the file to local file storage
		err = c.SaveFile(file, saveLocation+fileId)
		if err != nil {
			return util.FailedRequest(c, "server.error", err)
		}
	} else {
		return util.FailedRequest(c, "file.disabled", nil)
	}

	// Not encrypted cause this doesn't matter (and it's an unencrypted route)
	return c.JSON(fiber.Map{
		"success": true,
		"id":      fileId,
		"url":     urlPath,
	})
}
