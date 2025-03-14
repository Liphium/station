package localization

var (
	ErrorTableNotFound = Translations{
		englishUS: "The table wasn't found for some reason, maybe try rejoining the Space?",
	}
	ErrorTableAlreadyJoined = Translations{
		englishUS: "You are already in tabletop. You can't join again.",
	}
	ErrorTableCreationIssue = Translations{
		englishUS: "There was an issue during the table creation. Please report this to the developers.",
	}
	ErrorObjectNotFound = Translations{
		englishUS: "This object doesn't exist anymore, maybe it has already been deleted?",
	}
	ErrorObjectAlreadyHeld = Translations{
		englishUS: "This object is already being held by someone else. Please try to modify it again later.",
	}
	ErrorObjectNotInQueue = Translations{
		englishUS: "You didn't ask to modify this object before the actual modification. This is an issue with the app, please contact the developers.",
	}
	ErrorInvalidAction = Translations{
		englishUS: "You can't do that right now. Please try again later.",
	}
	ErrorTableClientNotFound = Translations{
		englishUS: "This member wasn't found, maybe he's already left?",
	}
	ErrorRoomNotFound = Translations{
		englishUS: "This room wasn't found, maybe it's already been deleted?",
	}

	// Studio errors
	ErrorDidntJoinStudio = Translations{
		englishUS: "You haven't joined studio yet.",
	}
)
