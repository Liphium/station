package properties

//! ALL THE DATA IN THIS IS PUBLIC AND CAN BE ACCESED THROUGH /ACCOUNT/PROFILE/GET
type Profile struct {
	ID string `json:"id,omitempty" gorm:"primaryKey"` // Account ID

	//* Picture stuff
	Picture     string `json:"picture,omitempty"`      // File id of the picture
	Container   string `json:"container,omitempty"`    // Attachment container encrypted with profile key
	PictureData string `json:"picture_data,omitempty"` // Profile picture data (zoom, x/y offset) encrypted with profile key

	Data string `json:"data,omitempty"` // Encrypted data (if we need it in the future)
}
