package files

import (
	"context"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
)

type deleteRequest struct {
	Id string `json:"id"`
}

// Route: /account/files/delete
func deleteFile(c *fiber.Ctx) error {

	var req deleteRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	accId := util.GetAcc(c)

	// Get file
	var file account.CloudFile
	if database.DBConn.Where("account = ? AND id = ?", accId, req.Id).First(&file).Error != nil {
		return util.FailedRequest(c, "file.not_found", nil)
	}

	// Delete file from S3
	_, err := client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file.Id),
	})
	if err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	// Delete file from DB
	if err := database.DBConn.Delete(&file).Error; err != nil {
		return util.FailedRequest(c, "server.error", err)
	}

	return util.SuccessfulRequest(c)
}
