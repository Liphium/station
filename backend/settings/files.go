package settings

import "github.com/Liphium/station/main/localization"

var FilesMaxUploadSize = &setting[int64]{
	Label:        localization.SettingFilesMaxUploadSize,
	Name:         "files.max_upload_size",
	DefaultValue: 10 * 1024 * 1024, // 10 MB
}

var FilesMaxTotalStorage = &setting[int64]{
	Label:        localization.SettingFilesMaxTotalStorage,
	Name:         "files.max_total_storage",
	DefaultValue: 1024 * 1024 * 1024, // 1 GB
}
