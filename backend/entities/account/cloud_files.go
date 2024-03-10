package account

type CloudFile struct {
	Id        string `json:"id,omitempty"`   // Format: a-[accountId]-[objectIdentifier]
	Name      string `json:"name,omitempty"` // File name (encrypted with file key)
	Path      string `json:"path,omitempty"` // Path to the CDN
	Type      string `json:"type,omitempty"` // Mime type
	Key       string `json:"key,omitempty"`  // Encryption key (encrypted with account public key)
	Account   string `json:"account,omitempty"`
	Size      int64  `json:"size,omitempty"` // In bytes
	Favorite  bool   `json:"favorite,omitempty"`
	System    bool   `json:"system,omitempty"` // If in use by system (won't be deleted)
	CreatedAt int64  `json:"created,omitempty" gorm:"not null,autoCreateTime:milli"`
}
