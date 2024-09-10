package files

import (
	"context"
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

// Route: /account/files/download/:id
func downloadFile(c *fiber.Ctx) error {

	if disabled {
		return util.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}

	id := c.Params("id")
	if id == "" {
		return util.InvalidRequest(c)
	}

	// Check for potentially malicious requests
	if strings.Contains(id, "/") {
		return util.InvalidRequest(c)
	}

	// Get the file from the database
	var file account.CloudFile
	if err := database.DBConn.Where("id = ?", id).Take(&file).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorFileNotFound, err)
	}

	// Send the file from the right location
	if fileRepoType == repoTypeR2 {
		// Retrieve file from R2
		obj, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(id),
		})
		if err != nil {
			return util.FailedRequest(c, localization.ErrorFileNotFound, err)
		}

		// Set headers for file download
		c.Set(fiber.HeaderContentType, *obj.ContentType)
		c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+id+`"`)

		// Stream the file to the client
		return c.SendStream(obj.Body)
	} else if fileRepoType == repoTypeLocal {
		// Send the file (it's encrypted so there is no checking of permissions required)
		return c.SendFile(saveLocation+id, true)
	} else {
		return util.FailedRequest(c, localization.ErrorFileDisabled, nil)
	}
}
