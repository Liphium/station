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

	// Get the connection of the client (for the encrypted id)
	connections, valid := caching.GetAllConnections(c.Client.Session)
	if !valid {
		return pipeshandler.ErrorResponse(c, localization.ErrorInvalidRequest, nil)
	}
	member, valid := connections[c.Client.ID]
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
		Sender:       member.Data,
	}

	// Add the message to the cache
	if err := caching.AddMessage(c.Client.Session, message); err != nil {
		return pipeshandler.ErrorResponse(c, localization.ErrorServer, nil)
	}

	// Send the message to all members of the Space using an event
	// We don't handle the error here to leave open the list_after endpoints as a backup solution
	SendEventToMembers(c.Client.Session, pipes.Event{
		Name: "msg",
		Data: map[string]interface{}{
			"msg": message,
		},
	})

	return pipeshandler.SuccessResponse(c)
}
