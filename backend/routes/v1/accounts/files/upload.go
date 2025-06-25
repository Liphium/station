package files

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/settings"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/files/upload
func uploadFile(c *fiber.Ctx) error {

	if disabled {
		return integration.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	// Form data
	key := c.FormValue("key", "-")
	name := c.FormValue("name", "-")
	extension := c.FormValue("extension", "-")
	tag := c.FormValue("tag", "")
	if key == "-" || name == "-" || extension == "-" {
		return integration.InvalidRequest(c, "invalid form data")
	}
	file, err := c.FormFile("file")
	if err != nil {
		return integration.InvalidRequest(c, "no file")
	}
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	fileType := file.Header.Get("Content-Type")
	if fileType == "" {
		return integration.InvalidRequest(c, "invalid headers")
	}

	// Check if the file name is valid
	if strings.Contains(file.Filename, "/") || strings.Contains(file.Filename, "\\") {
		return integration.InvalidRequest(c, "invalid filename")
	}

	// Check if the tag is valid
	if len(tag) > 100 {
		return integration.InvalidRequest(c, "invalid tag")
	}

	// Get max upload size and max total storage
	storageLimit, err := settings.FilesMaxTotalStorage.GetValue()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	maxUploadSize, err := settings.FilesMaxUploadSize.GetValue()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Check file size
	if file.Size > maxUploadSize {
		return integration.FailedRequest(c, localization.ErrorFileTooLarge(maxUploadSize), nil)
	}

	// Check total storage
	totalStorage, err := CountTotalStorage(accId)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	if totalStorage+file.Size > storageLimit {
		return integration.FailedRequest(c, localization.ErrorFileStorageLimit(storageLimit), nil)
	}

	// Generate file name Format: a-[timestamp]-[objectIdentifier].[extension]
	fileId := "a-" + fmt.Sprintf("%d", time.Now().UnixMilli()) + "-" + auth.GenerateToken(16) + "." + extension
	if err := database.DBConn.Create(&database.CloudFile{
		Id:      fileId,
		Name:    name,
		Type:    file.Header.Get("Content-Type"),
		Key:     key,
		Account: accId,
		Tag:     tag,
		System:  false,
		Size:    file.Size,
	}).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	// Save the file to whatever repository is selected
	if fileRepoType == repoTypeR2 {

		// Open the file
		f, err := file.Open()
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}

		// Save the file to R2
		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(fileId),
			Body:   f,
		})
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	} else if fileRepoType == repoTypeLocal {

		// Save the file to local file storage
		err = c.SaveFile(file, saveLocation+fileId)
		if err != nil {
			return integration.FailedRequest(c, localization.ErrorServer, err)
		}
	} else {
		return integration.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	// Not encrypted cause this doesn't matter (and it's an unencrypted route)
	return c.JSON(fiber.Map{
		"success": true,
		"id":      fileId,
		"url":     urlPath,
	})
}
