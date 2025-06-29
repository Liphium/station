package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/standards"
	"github.com/Liphium/station/backend/util/requests"
	"github.com/Liphium/station/main/localization"
	"github.com/google/uuid"
)

type AccountInfo struct {
	Id           standards.LPHAddress `json:"id"` // The address of the account, not just the id
	Username     string               `json:"name"`
	DisplayName  string               `json:"display_name"`
	PublicKey    string               `json:"pub"`
	SignatureKey string               `json:"sig"`
	CreatedAt    time.Time            `json:"-"`
}

// The accounts can be cached for a really long time as we have the update event, but for complete safety
// let's make it not too much. Database queries aren't *that* expensive.
const accountCacheDuration = time.Minute * 30

// Account cache for caching account info. (standards.LPHAddress -> AccountInfo)
//
// TODO: Invalidate when an update event arrives.
var accountCache *sync.Map = &sync.Map{}

// Get info about an account from the remote town or current town.
//
// Has a cache that automatically invalidates itself after “accountCacheDuration“ (currently 30 minutes) and
// also gets invalidated for every update event that's received containing the account's details.
//
// The duration of 30 minutes is there to make sure that, even if the update event gets lost somewhere, the
// account info is still correct after some time.
func LoadAccount(address standards.LPHAddress) (AccountInfo, localization.Translations, error) {

	// Check cache before doing anything else
	if obj, ok := accountCache.Load(address); ok && time.Since(obj.(AccountInfo).CreatedAt) < accountCacheDuration {
		return obj.(AccountInfo), nil, nil
	}

	// Make sure the address is valid
	accountId, origin, valid := address.Split()
	if !valid {
		return AccountInfo{}, nil, fmt.Errorf("invalid liphium address")
	}

	// Get the account info either from the database or the other server
	var accountInfo AccountInfo
	if origin != standards.CurrentTown() {
		res, err := requests.PostRequest(origin, "/accounts/get", requests.Map{
			"id": accountId,
		})
		if err != nil {
			return accountInfo, localization.ErrorOtherServer, fmt.Errorf("couldn't verify account on origin: %s", err)
		}
		if !requests.ValueOr(res, "success", false) {
			return accountInfo, nil, fmt.Errorf("account verification failed with: %s", requests.ValueOr(res, "message", "unknown error"))
		}
		if requests.ValueOr(res, "id", "-") != accountId {
			return accountInfo, localization.ErrorOtherServer, fmt.Errorf("sender account id is invalid")
		}
		accountInfo.Id = address
		accountInfo.Username = requests.ValueOr(res, "name", "")
		accountInfo.DisplayName = requests.ValueOr(res, "display_name", "")
		accountInfo.PublicKey = requests.ValueOr(res, "pub", "")
		accountInfo.SignatureKey = requests.ValueOr(res, "sg", "")
	} else {
		// Parse to UUID (as it's the standard at least on this server)
		accountUuid, err := uuid.Parse(accountId)
		if err != nil {
			return accountInfo, nil, fmt.Errorf("invalid account uuid: %s", err)
		}

		var account database.Account
		if err := database.DBConn.Where("id = ?", accountUuid).Take(&account).Error; err != nil {
			return accountInfo, nil, fmt.Errorf("couldn't get account from database: %s", err)
		}
		var publicKey database.PublicKey
		if err := database.DBConn.Where("id = ?", accountUuid).Take(&publicKey).Error; err != nil {
			return accountInfo, nil, fmt.Errorf("couldn't get public key from db: %s", err)
		}
		var signatureKey database.SignatureKey
		if err := database.DBConn.Where("id = ?", accountUuid).Take(&signatureKey).Error; err != nil {
			return accountInfo, nil, fmt.Errorf("couldn't get signature key from db: %s", err)
		}
		accountInfo.Id = address
		accountInfo.Username = account.Username
		accountInfo.DisplayName = account.DisplayName
		accountInfo.PublicKey = publicKey.Key
		accountInfo.SignatureKey = signatureKey.Key
	}

	if accountInfo.Username == "" || accountInfo.DisplayName == "" || accountInfo.PublicKey == "" || accountInfo.SignatureKey == "" {
		return accountInfo, nil, fmt.Errorf("invalid account info")
	}

	// Cache for future requests
	accountCache.Store(address, accountInfo)

	return accountInfo, nil, nil
}

func LoadAccounts() {
	// TODO: Implement
}
