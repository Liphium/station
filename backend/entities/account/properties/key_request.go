package properties

type KeyRequest struct {
	Session   string `gorm:"primaryKey"`
	Key       string // Public key of the session requesting it
	Signature string // Signature of the session requesting it
	Payload   string // Encrypted payload (from the session sending it)
	CreatedAt int64  `gorm:"not null,autoCreateTime:milli"`
}
