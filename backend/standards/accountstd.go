package standards

import (
	"regexp"
	"strings"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/main/localization"
)

// * Email standard
const EmailRegex = "^[a-zA-Z0-9]+(?:\\.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*(?:\\.[a-zA-Z0-9]+(?:-[a-zA-Z0-9]+)*)*$"

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

// * Account name standard
const MinUsernameLength = 3
const MaxUsernameLength = 20
const UsernameAllowedCharacters = "^[\\p{Ll}\\p{N}_\\-]+$"

// Check the requirements for a username
func CheckUsername(username string) localization.Translations {

	// Check length of the username
	if len(username) < MinUsernameLength {
		return localization.ErrorUsernameMinLength(MinUsernameLength)
	}

	if len(username) > MaxUsernameLength {
		return localization.ErrorUsernameMaxLength(MaxUsernameLength)
	}

	// Check if the username is valid
	match, err := regexp.Match(UsernameAllowedCharacters, []byte(username))
	if !match || err != nil {
		return localization.ErrorUsernameInvalid
	}

	// Check if username is available
	if database.DBConn.Where("username = ?", username).Take(&database.Account{}).RowsAffected > 0 {
		return localization.ErrorUsernameTaken
	}

	return nil
}

// * Account display name standard
const MaxDisplayNameLength = 20

// Check the requirements for a display name
func CheckDisplayName(username string) localization.Translations {

	// Check length of the username
	if len(username) < MinUsernameLength {
		return localization.ErrorDisplayNameMinLength(MinUsernameLength)
	}

	if len(username) > MaxDisplayNameLength {
		return localization.ErrorDisplayNameMaxLength(MaxDisplayNameLength)
	}

	return nil
}
