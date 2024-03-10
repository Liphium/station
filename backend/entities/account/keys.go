package account

type PublicKey struct {
	ID  string `json:"id" gorm:"primaryKey"` // Account id
	Key string `json:"key"`
}

type ProfileKey struct {
	ID  string `json:"id" gorm:"primaryKey"` // Account id
	Key string `json:"key"`                  // Encrypted with private key
}

type StoredActionKey struct {
	ID  string `json:"id" gorm:"primaryKey"` // Account id
	Key string `json:"key"`                  // Generated on the server
}

type SignatureKey struct {
	ID  string `json:"id" gorm:"primaryKey"` // Account id
	Key string `json:"key"`                  // Encrypted with private key
}
