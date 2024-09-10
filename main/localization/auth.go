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
