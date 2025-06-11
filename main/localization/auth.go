package localization

import (
	"fmt"
)

var (
	ErrorUsernameTaken = Translations{
		englishUS: "This username is taken. Please choose a different one.",
	}
	ErrorUsernameInvalid = Translations{
		englishUS: "Your username can only contain lowercase letters, numbers together with _ and -.",
	}
	ErrorEmailInvalid = Translations{
		englishUS: "Please enter a valid email address.",
	}
	ErrorEmailAlreadyInUse = Translations{
		englishUS: "The email address you entered is already in use, please try to use a different one.",
	}
	ErrorEmailNotFound = Translations{
		englishUS: "There is no account with this email address.",
	}
	ErrorEmailCodeInvalid = Translations{
		englishUS: "The code you entered doesn't match the one in your email, try again.",
	}
	ErrorSessionNotVerified = Translations{
		englishUS: "Your session isn't verified, please make sure to transfer the keys from your other devices first.",
	}
	ErrorPasswordIncorrect = Translations{
		englishUS: "Your password is incorrect. Please try again.",
	}
	ErrorPasswordsDontMatch = Translations{
		englishUS: "Your passwords don't match. Please try again.",
	}
	ErrorAuthRatelimit = Translations{
		englishUS: "Please wait a few seconds before trying again.",
	}
	ErrorInviteInvalid = Translations{
		englishUS: "This invite isn't valid, maybe it's already been used by someone else?",
	}
	ErrorSSONotCompleted = Translations{
		englishUS: "SSO hasn't been completed yet.",
	}

	// Localization for general auth stuff
	AuthNextStepButton = Translations{
		englishUS: "Next step",
	}
	AuthSubmitButton = Translations{
		englishUS: "Submit",
	}
	AuthFinishButton = Translations{
		englishUS: "Finish",
	}
	AuthResendEmailButton = Translations{
		englishUS: "Resend email",
	}

	// Localization for the auth start page
	AuthStartTitle = Translations{
		englishUS: "Welcome to Liphium!",
	}
	AuthStartDescription = Translations{
		englishUS: "With an account or not, just click \"Next step\"  and we'll figure out the rest.",
	}
	AuthStartEmailLabel = Translations{
		englishUS: "Email",
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

	// Localization for register invite form
	RegisterInviteTitle = Translations{
		englishUS: "Your invite, please.",
	}
	RegisterInvitePlaceholder = Translations{
		englishUS: "your-invite-here",
	}

	// Localization for the email validate form
	RegisterCodeTitle = Translations{
		englishUS: "Verify your email.",
	}
	RegisterCodeDescription = Translations{
		englishUS: "We sent you an email containing a verification code in the title and text. In case you entered the wrong email address in the beginning, you can still change it by using the input field below. If the email didn't reach your inbox, try clicking the resend button below.",
	}
	RegisterCodePlaceholder = Translations{
		englishUS: "111111",
	}
	RegisterResendEmailDescription = Translations{
		englishUS: "We sent you another email. Let's hope this one reached your inbox!",
	}

	// Localization for username registration form
	RegisterUsernameTitle = Translations{
		englishUS: "Create your username.",
	}
	RegisterUsernameDescription = Translations{
		englishUS: "Your username is the name other people can use to add you as a friend. It can only contain lowercase characters, numbers together with _ or -.",
	}
	RegisterUsernamePlaceholder = Translations{
		englishUS: "test123",
	}

	// Localization for display name registration form
	RegisterDisplayNameTitle = Translations{
		englishUS: "Create your display name.",
	}
	RegisterDisplayNameDescription = Translations{
		englishUS: "Your display name is the name everyone sees. No special requirements.",
	}
	RegisterDisplayNamePlaceholder = Translations{
		englishUS: "Test 123",
	}

	// Localization for the password adding form
	RegisterPasswordTitle = Translations{
		englishUS: "Create your password.",
	}
	RegisterPasswordRequirements = Translations{
		englishUS: "No big requirements, just create a password that's longer than 8 characters. And please don't use a bad one!",
	}
	RegisterPasswordPlaceholder = Translations{
		englishUS: "Password",
	}
	RegisterPasswordConfirmPlaceholder = Translations{
		englishUS: "Confirm password",
	}

	// Localization for SSO start
	RegisterSSOTitle = Translations{
		englishUS: "Sign in with SSO.",
	}
	RegisterSSODescription = Translations{
		englishUS: "This Liphium town uses SSO for accounts. If you don't know what to do from here, please ask the owner of your town. We'll check if you finished SSO sign-in automatically.",
	}
	RegisterSSOButton = Translations{
		englishUS: "Open auth provider",
	}
	RegisterSSOStatus = Translations{
		englishUS: "SSO sign-in",
	}
	RegisterSSOComplete = Translations{
		englishUS: "SSO has been completed. You can now safely return to the app.",
	}
)

func ErrorPasswordInvalid(minLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Please enter a password that is longer than %d characters.", minLength),
	}
}

func ErrorDisplayNameMinLength(minLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Your display name has to be longer than %d characters.", minLength),
	}
}

func ErrorDisplayNameMaxLength(maxLength int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Your display name has to be shorter than %d characters.", maxLength),
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

func ErrorRegistrationFailed(msg Translations) Translations {
	return Translations{
		englishUS: fmt.Sprintf("We're sorry, but registration failed. You may have to try registration again from scratch. %s", msg[englishUS]),
	}
}

func AuthRegisterCodeEmailCooldown(seconds int64) Translations {
	return Translations{
		englishUS: fmt.Sprintf("Please wait %d seconds before requesting another email.", seconds),
	}
}
