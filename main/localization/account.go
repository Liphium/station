package localization

import "fmt"

var (
	ErrorFriendNotFound = Translations{
		englishUS: "This user wasn't found. You sure this person exists? Maybe you just met them in your dreams?",
	}
	ErrorKeyNotFound = Translations{
		englishUS: "This key wasn't found. Make sure to look under the doormat.",
	}
	ErrorKeyAlreadySet = Translations{
		englishUS: "This key is already set.",
	}
	ErrorEntryNotFound = Translations{
		englishUS: "This vault entry couldn't be found.",
	}
	ErrorAccountNotFound = Translations{
		englishUS: "This account couldn't be found.",
	}
	ErrorSessionNotFound = Translations{
		englishUS: "This session couldn't be found. Maybe you need to login again?",
	}
	ErrorInvitesEmpty = Translations{
		englishUS: "You can't generate any invites right now as you don't have any invites left. You can ask an admin for invites if you want to invite someone.",
	}
)

func ErrorFriendLimitReached(limit int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("You've reached the maximum amount of %d friends. You must have a lot of friends, I feel lonely now. HOW DO YOU HAVE %d FRIENDS?", limit, limit),
	}
}

func ErrorVaultLimitReached(limit int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("You've reached the maximum amount of %d vault entries. Please delete some elements from the vault to add more items.", limit),
	}
}

func ErrorStoredActionLimitReached(limit int) Translations {
	return Translations{
		englishUS: fmt.Sprintf("The maximum amount of %d stored actions has been reached.", limit),
	}
}
