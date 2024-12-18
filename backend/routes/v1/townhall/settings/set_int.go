package townhall_settings

import (
	"log"

	"github.com/Liphium/station/backend/settings"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Route: /townhall/settings/set_int
func setIntegerSetting(c *fiber.Ctx) error {

	// Parse the request
	var req struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	if err := util.BodyParser(c, &req); err != nil {
		log.Println("invalid req", string(c.Locals("body").([]byte)))
		return util.InvalidRequest(c)
	}

	// Try to find the setting
	setting, valid := settings.SettingRegistryInteger[req.Name]
	if !valid {
		log.Println("setting not found")
		return util.InvalidRequest(c)
	}

	// Try to decode the value
	val, err := setting.Decode(req.Value)
	if err != nil {
		log.Println("can't decode")
		return util.InvalidRequest(c)
	}

	// Set the value in the database
	if err := setting.SetValue(val); err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.SuccessfulRequest(c)
}
