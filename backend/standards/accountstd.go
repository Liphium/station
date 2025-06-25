package standards

import (
	"fmt"
	"os"
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

// Standards for the account address

const AddressAllowedCharacters = "^[\\p{Ll}\\p{Lu}\\p{N}_\\-]+$"

// Get the Liphium address for this town and an account id.
func LiphiumAddress(accountId string) string {
	return fmt.Sprintf("%s@%s", accountId, CurrentTown())
}

// Get the current town address, same as the one in the LPH address of everyone in this town.
func CurrentTown() string {
	protocol := os.Getenv("PROTOCOL")
	if protocol == "http://" {
		return protocol + os.Getenv("BASE_PATH")
	}
	return os.Getenv("BASE_PATH")
}

// Splits the Liphium address using the @, returns false for the boolean when invalid.
//
// Checks the first part of the address to make sure it's valid.
// Second part (town) should be verified using a request.
func SplitLiphiumAddress(address string) (string, string, bool) {
	accountId, townUrl, valid := strings.Cut(address, "@")
	if !valid {
		return "", "", false
	}
	matched, err := regexp.Match(AddressAllowedCharacters, []byte(accountId))
	return accountId, townUrl, err == nil && matched
}
