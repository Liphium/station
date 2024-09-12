package localization

import "fmt"

var (
	ErrorUsernameTaken = Translations{
		englishUS: "This username is taken. Please choose a different one.",
	}
	ErrorUsernameInvalid = Translations{
		englishUS: "Your username doesn't match the requirements.",
	}
	ErrorDisplayNameInvalid = Translations{
		englishUS: "Your display name doesn't match the requirements.",
	}
	ErrorEmailInvalid = Translations{
		englishUS: "Your email doesn't match the requirements.",
	}
	ErrorEmailNotFound = Translations{
		englishUS: "There is no account with this email address.",
	}
	ErrorSessionNotVerified = Translations{
		englishUS: "Your session isn't verified, please make sure to transfer the keys from your other devices first.",
	}
	ErrorPasswordIncorrect = Translations{
		englishUS: "Your password is incorrect. Please try again.",
	}
	ErrorAuthRatelimit = Translations{
		englishUS: "Please wait a few seconds before trying again.",
	}

	// Localization for general auth stuff
	AuthNextStepButton = Translations{
		englishUS: "Next step",
	}
	AuthSubmitButton = Translations{
		englishUS: "Submit",
	}

	// Localization for the auth start page
	AuthStartTitle = Translations{
		englishUS: "Your email, please.",
	}
	AuthStartEmailPlaceholder = Translations{
		englishUS: "you@email.com",
	}
	AuthStartCreateButton = Translations{
		englishUS: "Create an account",
	}

	// Localization for the password page
	LoginPasswordTitle = Translations{
		englishUS: "Your password, please.",
	}
	LoginPasswordPlaceholder = Translations{
		englishUS: "yourmum123 (don't use this)",
	}
)

func ErrorPasswordInvalid(minLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Please enter a password that is longer than %d characters.", minLength),
	}
}

func ErrorDisplayNameMinLength(minLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Your username has to be shorter than %d characters.", minLength),
	}
}

func ErrorDisplayNameMaxLength(maxLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Your display name has to be longer than %d characters.", maxLength),
	}
}

func ErrorUsernameMinLength(minLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Your username has to be longer than %d characters.", minLength),
	}
}

func ErrorUsernameMaxLength(maxLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Your username has to be shorter than %d characters.", maxLength),
	}
}
