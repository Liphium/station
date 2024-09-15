package connect

import (
	"time"

	"github.com/Liphium/station/backend/database"
	"github.com/Liphium/station/backend/entities/account"
	"github.com/Liphium/station/backend/entities/app"
	"github.com/Liphium/station/backend/entities/node"
	"github.com/Liphium/station/backend/util"
	"github.com/Liphium/station/backend/util/nodes"
	"github.com/Liphium/station/main/localization"
	"github.com/gofiber/fiber/v2"
)

type connectRequest struct {
	Tag   string `json:"tag"`
	Token string `json:"token"`
}

// Route: /node/connect
func Connect(c *fiber.Ctx) error {

	// Parse request
	var req connectRequest
	if err := util.BodyParser(c, &req); err != nil {
		return util.InvalidRequest(c)
	}

	if !util.Permission(c, util.PermissionUseServices) {
		return util.FailedRequest(c, localization.ErrorNoPermission, nil)
	}

	// Get account
	accId, valid := util.GetAcc(c)
	if !valid {
		return util.InvalidRequest(c)
	}
	currentSessionId, err := util.GetSession(c)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}
	tk := req.Token

	var acc account.Account
	if err := database.DBConn.Preload("Sessions").Where("id = ?", accId).Take(&acc).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorAccountNotFound, nil)
	}

	// Check if account has key set
	if database.DBConn.Where("id = ?", acc.ID).Find(&account.PublicKey{}).Error != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Get the most recent session
	var mostRecent account.Session = account.Session{
		LastConnection: time.Unix(0, 10),
	}
	var sessionIds []string
	for _, session := range acc.Sessions {
		sessionIds = append(sessionIds, session.ID.String())

		if session.LastConnection.After(mostRecent.LastConnection) {
			mostRecent = session
		}
	}

	var currentSession account.Session
	if err := database.DBConn.Where("id = ?", currentSessionId).Take(&currentSession).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	if currentSession.Token != tk {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, nil)
	}

	// Get the app
	var application app.App
	if err := database.DBConn.Where("tag = ?", req.Tag).Take(&application).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorInvalidRequest, err)
	}

	// Get lowest load node
	var lowest node.Node

	// Connect to the same node if possible
	if mostRecent.Node != 0 {
		if err := database.DBConn.Model(&node.Node{}).Where("app_id = ? AND status = ? AND id = ?", application.ID, node.StatusStarted, mostRecent.Node).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorNotSetup, nil)
		}
	} else {
		if err := database.DBConn.Model(&node.Node{}).Where("app_id = ? AND status = ?", application.ID, node.StatusStarted).Order("load DESC").Take(&lowest).Error; err != nil {
			return util.FailedRequest(c, localization.ErrorNotSetup, nil)
		}
	}

	// Ping node (to see if it's online)
	if err := lowest.SendPing(); err != nil {

		// Set the node to error
		nodes.TurnOff(&lowest, node.StatusError)
		return util.FailedRequest(c, localization.ErrorNode, err)
	}

	// Generate a jwt token for the node
	token, err := util.ConnectionToken(accId, currentSessionId.String(), lowest.ID)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Generate a jwt token with session information
	sessionInformationToken, err := util.SessionInformationToken(accId, sessionIds)
	if err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	currentSession.LastConnection = time.Now()
	currentSession.Node = lowest.ID
	currentSession.App = application.ID
	if err := database.DBConn.Save(&currentSession).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	// Save node
	if err := database.DBConn.Save(&lowest).Error; err != nil {
		return util.FailedRequest(c, localization.ErrorServer, err)
	}

	return util.ReturnJSON(c, fiber.Map{
		"success": true,
		"domain":  lowest.Domain,
		"id":      lowest.ID,
		"token":   token,
		"s_info":  sessionInformationToken,
	})
}
