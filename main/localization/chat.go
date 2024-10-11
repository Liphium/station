package localization

import "fmt"

var (
	ErrorAlreadyInSpace = Translations{
		englishUS: "You are already in a space. Please leave the current one to enter a new one.",
	}
	ErrorKickNoPermission = Translations{
		englishUS: "You don't have permission to kick this person from the conversation.",
	}
	ErrorGroupDataTooLong = Translations{
		englishUS: "The data of this conversation became too long. This shouldn't normally happen. You should probably contact the developers of this app.",
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
)

func ErrorGroupMemberLimit(limit int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("You can't have more than %d group members.", limit),
	}
}
