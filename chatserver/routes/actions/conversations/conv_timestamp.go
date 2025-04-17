package conversation_actions

import (
	"time"

	"github.com/Liphium/station/chatserver/database"
	"github.com/Liphium/station/chatserver/util"
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

// Action: conv_timestamp
func HandleTimestamp(c *fiber.Ctx, token database.ConversationToken, data interface{}) error {

	// Okay, I will put an explanation here, so we all understand why this is here.
	// Time is hard and every single fucking time I boot into my Windows after coming from my
	// Linux installation the time is desynchronized by 2 hours. For this reason, we need time-
	// stamps to come from the fucking server and we also don't want people to troll, so we're
	// going to make jwt tokens that are just for time.
	// Yeah, over-engineering is my second name.

	// Sign a jwt token with the time
	time := time.Now().UnixMilli()
	tk, err := util.TimestampToken(time)
	if err != nil {
		return integration.FailedRequest(c, localization.ErrorServer, err)
	}

	return integration.ReturnJSON(c, fiber.Map{
		"success": true,
		"token":   tk,
		"stamp":   time,
	})
}
