package magic_accounts

import (
	"fmt"
	"log"

	"github.com/Liphium/magic/mconfig"
	"github.com/Liphium/station/backend/database"
	magic_util "github.com/Liphium/station/backend/magic/scripts/util"
	"github.com/Liphium/station/backend/util/auth"
)

const DefaultPub = "some_pub"
const DefaultSig = "some_sig"

// Create a test account.
func TestAccount(p *mconfig.Plan, name string) {
	magic_util.PrepareEnvironment(p)
	database.Connect()

	// Create the actual account
	acc := &database.Account{
		Email:       name + "@liphium.app",
		DisplayName: name,
		Username:    name,
		RankID:      1, // Default
	}
	if err := database.DBConn.Create(&acc).Error; err != nil {
		log.Fatalln("couldn't create account:", err)

	}

	// Create the password
	hash, err := auth.HashPassword("yourmum123", acc.ID)
	if err != nil {
		log.Fatalln("couldn't hash yourmum123:", err)
	}
	if err := database.DBConn.Create(&database.Authentication{
		Account: acc.ID,
		Type:    database.AuthTypePassword,
		Secret:  hash,
	}).Error; err != nil {
		log.Fatalln("couldn't create password:", err)
	}

	// Create fake pub and signature key
	if err := database.DBConn.Create(&database.PublicKey{
		ID:  acc.ID,
		Key: DefaultPub,
	}).Error; err != nil {
		log.Fatalln("couldn't create public key:", err)
	}
	if err := database.DBConn.Create(&database.SignatureKey{
		ID:  acc.ID,
		Key: DefaultSig,
	}).Error; err != nil {
		log.Fatalln("couldn't create signature key:", err)
	}
	// Print in case not testing
	if p.Profile != "test" {
		fmt.Println("Created account")
		fmt.Println("E-Mail:", name+"@liphium.app")
		fmt.Println("Password: yourmum123")
	}
}
