package townhall_settings

import (
	"github.com/Liphium/station/backend/settings"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /townhall/settings/set_int
func setBooleanSetting(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	if err := c.BodyParser(&req); err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Try to find the setting
	setting, valid := settings.SettingRegistryBoolean[req.Name]
	if !valid {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Try to decode the value
	val, err := setting.Decode(req.Value)
	if err != nil {
		return integration.InvalidRequest(c, "invalid request")
	}

	// Set the value in the database
	if err := setting.SetValue(val); err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.SuccessfulRequest(c)
}
