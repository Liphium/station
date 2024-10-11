package database

import "github.com/google/uuid"

// A regular post
type Post struct {

	// Data about the post
	ID      uuid.UUID `gorm:"primaryKey"`
	Creator uuid.UUID
	Parent  string // The parent post (in case this is a comment or sth, it's a string cause it can be empty)

	// The content of the post
	Attachments string
	Content     string
	Key         string // Key of the post for the original creator (encrypted with vault key)
	Edited      bool
	Creation    int64 // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
}

// A conversion helper to convert it to a sent post
func (p Post) ToSent(visibilities []PostVisibility) SentPost {
	// Convert all the visibilities to sendable versions
	var sentVisibilities []SentPostVisibility
	for _, visibility := range visibilities {
		sentVisibilities = append(sentVisibilities, SentPostVisibility{
			Type:       visibility.Type,
			Identifier: visibility.Identifier,
			Key:        visibility.Key,
		})
	}

	// Convert everything to the sent post
	return SentPost{
		ID:           p.ID.String(),
		Creator:      p.Creator.String(),
		Parent:       p.Parent,
		Attachments:  p.Attachments,
		Content:      p.Content,
		Edited:       p.Edited,
		Creation:     p.Creation,
		Visibilities: sentVisibilities,
	}
}

// Types of visibility
const VisibilityFriends = "friends"   // Key will be stored encrypted with profile key
const VisibilityConversation = "conv" // Key will be encrypted with the conversation's key

// The keys for decrypting
type PostVisibility struct {

	// Important data for visibility
	ID   uuid.UUID `gorm:"primaryKey"`
	Post uuid.UUID

	// Cached data from the post (this is cached here to make queries faster)
	Creator  uuid.UUID // Creator of the post
	Creation int64     // Time of creation from the post

	// Identifiers to get the visibilities
	Type       string // The type of visibilty (look above)
	Identifier string // Identifier of the visiblity target (for example the conversation id when the type is conversation)
	Key        string // Key of the post for the people in the visibility group
}

// A regular post (for sending it with json)
type SentPost struct {
	ID           string               `json:"id"`
	Creator      string               `json:"cr"`
	Parent       string               `json:"p,omitempty"` // Parent post, string to allow empty value
	Attachments  string               `json:"a"`
	Content      string               `json:"co"`
	Edited       bool                 `json:"e"`
	Creation     int64                `json:"crt"` // Unix timestamp (set by the client)
	Visibilities []SentPostVisibility `json:"v"`
}

type SentPostVisibility struct {
	Type       string `json:"t"` // The type of visibility (look above)
	Identifier string `json:"i"` // Identifier of the visibility target
	Key        string `json:"k"` // Key of the post for the people in the visibility group
}
