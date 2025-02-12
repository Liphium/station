package manage

import (
	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/auth"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type newRequest struct {
	Token           string  `json:"token"`
	App             uint    `json:"app"` // App ID
	Domain          string  `json:"domain"`
	PeformanceLevel float32 `json:"performance_level"`
}

func newNode(c *fiber.Ctx) error {

	// Parse body to add request
	var req newRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	// Check if token is valid
	var ct database.NodeCreation
	if err := database.DBConn.Where("token = ?", req.Token).Take(&ct).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	if req.Domain == "" {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	if len(req.Domain) < 3 {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	var app database.App
	if err := database.DBConn.Take(&app, req.App).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Create node
	var created database.Node = database.Node{
		AppID:           req.App,
		Token:           auth.GenerateToken(300),
		Domain:          req.Domain,
		Load:            0,
		PeformanceLevel: req.PeformanceLevel,
		Status:          1,
	}

	if err := database.DBConn.Create(&created).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   created.Token,
		"app":     app.Name,
		"id":      created.ID,
	})
}
