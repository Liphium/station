package settings

import "github.com/Liphium/station/main/localization"

var FilesMaxUploadSize = &intSetting{
	setting: setting[int64]{
		Label:        localization.SettingFilesMaxUploadSize,
		Name:         "files.max_upload_size",
		DefaultValue: 10 * 1024 * 1024, // 10 MB
	},
	Min:     0,
	Max:     200 * 1024 * 1024, // 200 MB (only reasonable value to allow on the current protocol)
	Devider: 1024 * 1024,       // So the thing shows the value in MB
}

var FilesMaxTotalStorage = &intSetting{
	setting: setting[int64]{
		Label:        localization.SettingFilesMaxTotalStorage,
		Name:         "files.max_total_storage",
		DefaultValue: 1024 * 1024 * 1024, // 1 GB
	},
	Min:     0,
	Max:     100 * 1024 * 1024 * 1024, // 100 GB
	Devider: 1024 * 1024,              // So the thing shows the value in MB
}
