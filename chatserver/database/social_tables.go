package database

import "github.com/google/uuid"

type Post struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Creator     uuid.UUID
	Attachments string
	Content     string
	Key         string // Key of the post for the original creator (encrypted with vault key)
	Edited      bool
	Creation    int64 // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
}

// Types of visibility
const VisibilityPublic = "public" // Key will be stored encrypted with profile key
const VisibilityGroup = "group"   // Key will be encrypted with the conversation's key

// The keys for decrypting
type PostVisibility struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Post       uuid.UUID
	Creator    uuid.UUID // Creator of the post (this is cached here to make queries faster)
	Type       string    // The type of visibilty (look above)
	Identifier string    // Identifier of the visiblity group (for example the conversation id when the type is group visibility)
	Key        string    // Key of the post for the people in the visibility group
}
