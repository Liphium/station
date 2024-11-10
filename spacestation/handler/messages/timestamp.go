package message_handlers

import (
	"time"

	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
)

// Action: msg_timestamp
func generateTimestampToken(c *pipeshandler.Context, _ interface{}) pipes.Event {

	// Generate a new timestamp token
	time := time.Now().UnixMilli()
	tk, err := TimestampToken(time)
	if err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, err)
	}

	return pipeshandler.NormalResponse(c, map[string]interface{}{
		"success": true,
		"token":   tk,
		"time":    time,
	})
}
