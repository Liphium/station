package conversations

import (
	"os"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
)

type Message struct {
	ID string `json:"id" gorm:"primaryKey"`

	Conversation string `json:"conversation" gorm:"not null"`
	Certificate  string `json:"certificate" gorm:"not null"`
	Creation     int64  `json:"creation"`               // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
	Data         string `json:"data" gorm:"not null"`   // Encrypted data
	Edited       bool   `json:"edited" gorm:"not null"` // Edited flag
	Sender       string `json:"sender" gorm:"not null"` // Sender ID (of conversation token)
}

func CheckSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*6
}

type CertificateClaims struct {
	Message      string `json:"mid"` // Message ID
	Conversation string `json:"c"`   // Conversation ID
	Sender       string `json:"sd"`  // Sender ID
	jwt.RegisteredClaims
}

// Generate a message certificate (used to verify message sender when editing/deleting messages)
func GenerateCertificate(id string, conversation string, sender string) (string, error) {

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, CertificateClaims{
		Message:      id,
		Conversation: conversation,
		Sender:       sender,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "chat-node",
		},
	})

	token, err := tk.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

// Get all claims in a message certificate
func GetCertificateClaims(certificate string) (*CertificateClaims, bool) {

	token, err := jwt.ParseWithClaims(certificate, &CertificateClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Minute))

	if err != nil {
		return &CertificateClaims{}, false
	}

	if claims, ok := token.Claims.(*CertificateClaims); ok && token.Valid {
		return claims, true
	}

	return &CertificateClaims{}, false
}

// Check if a message certificate is valid with some parameters
func (m *CertificateClaims) Valid(id string, conversation string, sender string) bool {
	return m.Message == id && m.Conversation == conversation && m.Sender == sender
}
