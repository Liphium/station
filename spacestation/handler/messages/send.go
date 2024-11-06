package message_handlers

import (
	"github.com/Liphium/station/main/integration"
	"github.com/Liphium/station/main/localization"
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipeshandler"
	"github.com/Liphium/station/spacestation/caching"
	"github.com/google/uuid"
)

// Action: msg_send
func sendMessage(c *pipeshandler.Context, action struct {
	Token string `json:"token"` // Timestamp token
	Data  string `json:"data"`
}) pipes.Event {

	// Validate request
	if len(action.Data) == 0 {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Verify the timestamp token
	timestamp, valid := VerifyTimestampToken(action.Token)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Get the connection (for getting the client id)
	connection, valid := caching.GetConnection(c.Client.ID)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}

	// Create the message and save it
	message := caching.Message{
		ID:           uuid.NewString(),
		Conversation: c.Client.Session + "@" + integration.Domain,
		Creation:     timestamp,
		Data:         action.Data,
		Edited:       false,
		Sender:       connection.ClientID,
	}

	// Add the message to the cache
	if err := caching.AddMessage(c.Client.Session, message); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	// TODO: Broadcast in an event

	return pipeshandler.SuccessResponse(c)
}
