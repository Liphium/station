package friends2_routes

import (
	"github.com/Liphium/station/backend/service"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/neogate"
	"github.com/gofiber/fiber/v2"
)

/*
# Architecture brainstorming

## Friends table
- Account is the column for your account.
- Target is the other person you're adding.
- Request boolean for if request or not.
- Token needed for Target or you to break the friendship (to verify on external servers).

## Friend adding process
1. Call /add to send request. In case not current server, call server/.../add_external, mirror database entry with returned result.

Target can then:
1. Call /remove to get rid of the request (also uses BreakToken to verify).
2. Call /add to accept request. Mirror database entries.

## Unauthorized
./friends/add_external: Send a friend request with their id, for external servers
./friends/remove_external: Remove a friend (BreakToken needed)

## Authorized
./friends/add: Send a friend request or accept, proxies to ./friends/add_external in case of different server
./friends/remove: Remove a friend, proxies to ./friend/remove_external ^

./friends/list: List all friends. With limit/offset.
TODO: See if other endpoints may be needed.

## Questions
- How do we handle account info? Names, profile, etc.
  - Cache could just work, store for 2 hours and otherwise re-get from database or external server
  - We don't want to also update the friends table when the profile changes
- How do we update account info on updates on the client?
  - Update event? Would be cool, we could just have an event that says the friend was updated and then on the client we sub to that event and make it so the client just quielty reloads that part of the UI, client should have an account map with all the account info anyway, make that reactive and update it when update event is called
  - External event channel needed ig, we need that anyway
- Should we handle the external event channel over the friends system?
  - Would make sense since that's the new main system that powers everything (before conv now friends)
  - Solves the question before this one with the update event
  - Needed for lots of things like ringing and stuff

*/

func Authorized(router fiber.Router) {

}

func Unauthorized(router fiber.Router) {

}

// Event for a new friend or request
func FriendEvent(request bool, account standards.LPHAddress, name string, displayName string, publicKey string, signatureKey string, profileKey string) neogate.Event {
	return neogate.Event{
		Name: "fr_rq",
		Data: requests.Map{
			"request":      request,
			"id":           account,
			"name":         name,
			"display_name": displayName,
			"pub":          publicKey,
			"sig":          signatureKey,
			"prf":          profileKey,
		},
	}
}

// Simple helper function using account info
func FriendEventFromAccountInfo(request bool, accInfo service.AccountInfo, profileKey string) neogate.Event {
	return FriendEvent(request, accInfo.Id, accInfo.Username, accInfo.DisplayName, accInfo.PublicKey, accInfo.SignatureKey, profileKey)
}
