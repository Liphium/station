package connect

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/verify"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Tag   string `json:"tag"`
	Token string `json:"token"`
	Extra string `json:"extra"`
}

// Route: /node/connect
func Connect(c *fiber.Ctx) error {

	// Parse request
	var req connectRequest
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Make sure the extra parameter isn't abused
	if len(req.Extra) >= 64 {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	if !verify.InfoLocals(c).HasPermission(verify.PermissionUseServices) {
		return integration.FailedRequest(c, localization.ErrorNoPermission, nil)
	}

	// Get account
	accId, err := verify.InfoLocals(c).GetAccountUUID()
	if err != nil {
		return integration.InvalidRequest(c, "invalid account id")
	}
	currentSessionId, err := verify.InfoLocals(c).GetSessionUUID()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}
	tk := req.Token

	var acc database.Account
	if err := database.DBConn.Preload("Sessions").Where("id = ?", accId).Take(&acc).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorAccountNotFound, nil)
	}

	// Check if account has key set
	if database.DBConn.Where("id = ?", acc.ID).Find(&database.PublicKey{}).Error != nil {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	var currentSession database.Session
	if err := database.DBConn.Where("id = ?", currentSessionId).Take(&currentSession).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	if currentSession.Token != tk {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Get the app
	var application database.App
	if err := database.DBConn.Where("tag = ?", req.Tag).Take(&application).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	// Get the node
	var chosenNode database.Node
	if err := database.DBConn.Model(&database.Node{}).Where("app_id = ? AND status = ?", application.ID, database.StatusStarted).Take(&chosenNode).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorNotSetup, nil)
	}

	// Generate a jwt token for the node
	token, err := util.ConnectionToken(accId, currentSessionId.String(), req.Extra, chosenNode.ID)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	currentSession.LastConnection = time.Now()

	// Only chat receives stored actions (in the future this could be in the database)
	if application.Tag == "liphium_chat" {
		currentSession.Node = chosenNode.ID
		currentSession.App = application.ID
	}
	if err := database.DBConn.Save(&currentSession).Error; err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"domain":  chosenNode.Domain,
		"id":      chosenNode.ID,
		"token":   token,
	})
}
