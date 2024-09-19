package files

import (
	"context"
	"os"
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/delete
func deleteFile(c *fiber.Ctx) error {

	if disabled {
		return util.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	var req deleteRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}

	// Get file
	var file account.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return util.FailedRequest(c, localization.ErrorFileNotFound, nil)
	}

	// Check for potential malicious requests
	if strings.Contains(req.Id, "/") {
		return util.InvalidRequest(c)
	}

	// Check where the file should be deleted
	if fileRepoType == repoTypeR2 {

		// Delete the object from R2
		_, err := s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(file.Id),
		})
		if err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}
	} else if fileRepoType == repoTypeLocal {

		// Delete file from local file system
		err := os.Remove(saveLocation + req.Id)
		if err != nil {
			return util.FailedRequest(c, localization.ErrorServer, err)
		}
	} else {
		return util.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	// Delete file from DB
	if err := database.DBConn.Delete(&file).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
