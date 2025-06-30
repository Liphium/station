package magic_accounts

import (
	"fmt"
	"log"
	"time"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/backend/database"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
	auth_routes "github.com/Liphium/station/backend/routes/v1/accounts/auth"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/main/integration"
)

// Create a test token for the account with that username
func GetTestToken(p *mconfig.Plan, username string) (refreshToken string, token string) {

	var account database.Account
	magic_util.DatabaseError(database.DBConn.Where("username = ?", username).Preload("Rank").Take(&account))

	var session database.Session
	magic_util.DatabaseError(database.DBConn.Where("account = ?", account.ID).FirstOrCreate(&session, database.Session{
		Token:           integration.GenerateToken(20),
		Verified:        true,
		Account:         account.ID,
		PermissionLevel: account.Rank.Level,
		Device:          "test",
		LastUsage:       time.Now(),
		LastConnection:  time.Now(),
	}))

	// Create a real token from the server (wanna test as much as possible :D)
	res, err := requests.PostRequestURL(requests.CurrentPath("/accounts/auth/refresh"), auth_routes.RefreshRequest{
		Session: session.ID.String(),
		Token:   session.Token,
	})
	if err != nil {
		log.Fatalln("couldn't do refresh request:", err)
	}
	if !requests.ValueOr(res, "success", false) {
		log.Fatalln("request error:", requests.ValueOr(res, "message", "?"))
	}
	if !requests.ValueOr(res, "verified", true) {
		log.Fatalln("not verified smh")
	}
	if !requests.ValueOr(res, "valid", true) {
		log.Fatalln("not valid smh")
	}

	fmt.Println()
	fmt.Println("Refresh token:", session.Token)
	fmt.Println("Token (for requests):", requests.ValueOr(res, "token", "HUH"))
	fmt.Println()

	return session.Token, requests.ValueOr(res, "token", "HUH")
}
