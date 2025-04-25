package localization

import "fmt"

var (
	// General chat stuff
	ErrorKickNoPermission = Translations{
		englishUS: "You don't have permission to kick this person from the conversation.",
	}
	ErrorMessageAlreadyDeleted = Translations{
		englishUS: "This message has already been deleted. If this issue occurs, try restarting your app.",
	}
	ErrorMessageTooLong = Translations{
		englishUS: "Your message is too long. Please make sure it fits the requirements.",
	}
	ErrorMessageDeleteNoPermission = Translations{
		englishUS: "You don't have permission to delete this message.",
	}
	ErrorDecentralizationDisabled = Translations{
		englishUS: "Decentralization is currently disabled. Please contact the admin of your town if this bothers you.",
	}
	ErrorMemberNoPermission = Translations{
		englishUS: "You don't have any permission to perform this action.",
	}
)

func ErrorTooManySharedSpaces(limit int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("A conversation can only have %d active shared Spaces.", limit),
	}
}

func ErrorGroupMemberLimit(limit int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("You can't have more than %d group members.", limit),
	}
}
