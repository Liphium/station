package localization

import "fmt"

var (
	ErrorFileNotFound = Translations{
		englishUS: "This file couldn't be found.",
	}
	ErrorFileDisabled = Translations{
		englishUS: "File uploading is currently not available. Please try again later or contact the owner of your town.",
	}
)

func ErrorFileTooLarge(maxSize int64) Translations {
	size := TranslateFileSize(float64(maxSize))
	return Translations{
		englishUS: fmt.Sprintf("The max file size is %s.", size[englishUS]),
	}
}

func ErrorFileStorageLimit(maxSize int64) Translations {
	size := TranslateFileSize(float64(maxSize))
	return Translations{
		englishUS: fmt.Sprintf("You reached your maximum storage limit of %s.", size[englishUS]),
	}
}

// Format amount of bytes to a translated string
func TranslateFileSize(fileSize float64) Translations {
	if fileSize < 1024 {
		return fileSizeBytes(fileSize)
	}

	if fileSize < 1024*1024 {
		return fileSizeKilobytes(fileSize / 1024)
	}

	if fileSize < 1024*1024*1024 {
		return fileSizeMegabytes(fileSize / (1024 * 1024))
	}

	return fileSizeGigabytes(fileSize / (1024 * 1024 * 1024))
}

func fileSizeBytes(fileSize float64) Translations {
	return Translations{
		englishUS: fmt.Sprintf("%.2f B", fileSize),
	}
}

func fileSizeKilobytes(fileSize float64) Translations {
	return Translations{
		englishUS: fmt.Sprintf("%.2f KB", fileSize),
	}
}

func fileSizeMegabytes(fileSize float64) Translations {
	return Translations{
		englishUS: fmt.Sprintf("%.2f MB", fileSize),
	}
}

func fileSizeGigabytes(fileSize float64) Translations {
	return Translations{
		englishUS: fmt.Sprintf("%.2f GB", fileSize),
	}
}
