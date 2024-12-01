package settings

import "github.com/Liphium/station/main/localization"

var DecentralizationEnabled = &setting[bool]{
	Label:        localization.SettingDecentralizationEnabled,
	Name:         "decentralization.enabled",
	DefaultValue: true,
}

var DecentralizationAllowUnsafe = &setting[bool]{
	Label:        localization.SettingDecentralizationAllowUnsafe,
	Name:         "decentralization.allow_unsafe",
	DefaultValue: false,
}

var ChatMessagePullThreads = &intSetting{
	setting: setting[int64]{
		Label:        localization.SettingChatMessagePullThreads,
		Name:         "chat.message_pull_threads",
		DefaultValue: 5,
	},
	Min:     1,
	Max:     100,
	Devider: 1,
}
