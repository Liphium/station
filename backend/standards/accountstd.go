package standards

import (
	"regexp"
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
)

// * Email standard
const EmailRegex = "^[a-zA-Z0-9]+(?:\\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:\\.[a-zA-Z0-9]+)*$"

func NormalizeEmail(email string) string {

	// Convert email to lowercase
	email = strings.ToLower(email)

	// Remove leading and trailing whitespaces
	email = strings.TrimSpace(email)

	// Remove dots (.) from the username part of the email
	parts := strings.Split(email, "@")
	username := parts[0]
	username = strings.ReplaceAll(username, ".", "")

	// Reconstruct the normalized email address
	normalizedEmail := username + "@" + parts[1]

	return normalizedEmail
}

func CheckEmail(email string) (bool, string) {

	// Check if email is valid
	match, err := regexp.Match(EmailRegex, []byte(email))
	if !match || err != nil {
		return false, ""
	}

	email = NormalizeEmail(email)
	if strings.Contains(email, " ") {
		return false, ""
	}

	return true, email
}

// * Account name#tag standard
const MinUsernameLength = 3
const MaxUsernameLength = 16
const AllowedCharactersRegex = "^[\\p{L}\\p{N}_\\-]+$"

const MinTagLength = 3
const MaxTagLength = 5

func CheckUsernameAndTag(username string, tag string) (bool, string) {

	// Check length of the username
	if len(username) < MinUsernameLength {
		return false, "username.invalid"
	}

	if len(username) > MaxUsernameLength {
		return false, "username.invalid"
	}

	// Check if the username is valid
	match, err := regexp.Match(AllowedCharactersRegex, []byte(username))
	if !match || err != nil {
		return false, "username.invalid"
	}

	// Check length of the tag
	if len(tag) < MinTagLength {
		return false, "tag.invalid"
	}

	if len(tag) > MaxTagLength {
		return false, "tag.invalid"
	}

	// Check if the username is valid
	match, err = regexp.Match(AllowedCharactersRegex, []byte(tag))
	if !match || err != nil {
		return false, "tag.invalid"
	}

	// Check if username and tag is available
	if database.DBConn.Where("username = ? AND tag = ?", username, tag).Take(&account.Account{}).RowsAffected > 0 {
		return false, "username.taken"
	}

	return true, ""
}
