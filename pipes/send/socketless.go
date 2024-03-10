package send

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/receive"
	"github.com/Liphium/station/pipes/util"
)

func Socketless(nodeEntity pipes.Node, message pipes.Message) error {

	if pipes.DebugLogs {
		pipes.Log.Printf("sent on [socketless] %s: %s", message.Channel.Channel, message.Event.Name)
	}

	if nodeEntity.ID == pipes.CurrentNode.ID {

		receive.HandleMessage("ws", message)
		return nil
	}

	err := util.PostRaw(nodeEntity.SL, map[string]interface{}{
		"token":   nodeEntity.Token,
		"message": message,
	})

	if err != nil {
		return err
	}

	return nil
}
