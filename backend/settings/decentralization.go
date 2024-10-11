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
