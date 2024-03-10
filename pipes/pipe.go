package pipes

var DebugLogs = false

const ChannelP2P = "p"
const ChannelConversation = "c"
const ChannelBroadcast = "br"

type Event struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:"data"`
}

type Channel struct {
	Channel string   `json:"channel"` // Channel name
	Target  []string `json:"target"`  // User IDs to send to (node and user ID for p2p channel)
	Nodes   []string `json:"-"`       // Nodes to send to (only for conversation channel)
}

type Message struct {
	Channel Channel `json:"channel"`
	Event   Event   `json:"event"`
	Local   bool    `json:"-"` // Whether to only send to local clients (excluded from JSON)
}

func (c Channel) IsP2P() bool {
	return c.Channel == "p"
}

func (c Channel) IsConversation() bool {
	return c.Channel == "c"
}

func (c Channel) IsBroadcast() bool {
	return c.Channel == "br"
}

func P2PChannel(receiver string, receiverNode string) Channel {
	return Channel{
		Channel: ChannelP2P,
		Target:  []string{receiver, receiverNode},
	}
}

func Conversation(receivers []string, nodes []string) Channel {
	return Channel{
		Channel: ChannelConversation,
		Target:  receivers,
		Nodes:   nodes,
	}
}

func BroadcastChannel(receivers []string) Channel {
	return Channel{
		Channel: ChannelBroadcast,
		Target:  receivers,
	}
}
