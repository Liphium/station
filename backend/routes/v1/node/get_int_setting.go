package node

import (
	"github.com/Liphium/station/backend/settings"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /node/get_int_setting
func getIntSetting(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		ID      string `json:"id"`
		Token   string `json:"token"`
		Setting string `json:"setting"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Validate node token
	_, err := nodes.Node(nodeToU(req.ID), req.Token)
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Get the setting and the value of it
	setting, valid := settings.SettingRegistryInteger[req.Setting]
	if !valid {
		return integration.FailedRequest(c, localization.ErrorInvalidRequestContent, nil)
	}
	val, err := setting.GetValue()
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, nil)
	}

	return c.JSON(fiber.Map{
		"success": true,
		"value":   val,
	})
}
